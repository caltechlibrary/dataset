package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
	"time"
)

var (
	contentType        = "application/json"
	dName              = "testout"
	query 			   = "query_test"
	testCollectionsSrc = []byte(`{
	"attachment_test.ds": {
		"one": {
			"one": 1,
			"two": "too",
			"three": true
		}
	},
	"objects_test.ds": {
		"Miller-A": {
			"id":     "Miller-A",
			"given":  "Arthor",
			"family": "Miller",
			"genre":  ["plays", "film"]
		},
		"Lopez-T": {
			"id":     "Lopez-T",
			"given":  "Tom",
			"family": "Lopez",
			"genre":  ["plays", "radio-theater"]
		}
	}
}`)
)

func setupApiTestCollection(cName string, dsnURI string, records map[string]map[string]interface{}) error {
	// Create collection.json using v1 structures
	if len(cName) == 0 {
		return fmt.Errorf("missing a collection name")
	}
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	c, err := Init(cName, dsnURI)
	if err != nil {
		return err
	}
	defer c.Close()
	// Now populate with some test records records.
	for key, obj := range records {
		if err := c.Create(key, obj); err != nil {
			return err
		}
	}
	return nil
}


// sameStrings compares one slice of strings to another by converting
// them to a JSON represention and compariing the representation.
func sameStrings(expected []string, got []string) bool {
	src1, _ := JSONMarshalIndent(expected, "", "    ")
	src2, _ := JSONMarshalIndent(got, "", "    ")
	return bytes.Compare(src1, src2) == 0
}

// sameMapping compares one map to another by converting
// them to a JSON representation and comparing the representation.
func sameMapping(expected map[string]interface{}, got map[string]interface{}) bool {
	src1, _ := JSONMarshalIndent(expected, "", "    ")
	src2, _ := JSONMarshalIndent(got, "", "    ")
	return bytes.Compare(src1, src2) == 0
}

// sameObjects compares one slice of objects to an another
// by converting them to JSON and comparing the representation.
func sameObjects(expected []map[string]interface{}, got []map[string]interface{}) bool {
	if len(expected) != len(got) {
		return false
	}
	for i := 0; i < len(expected); i++ {
		if !sameMapping(expected[i], got[i]) {
			return false
		}
	}
	return true
}

func makePayload(src []byte) (io.Reader, error) {
	return bytes.NewBuffer(src), nil
}

func makeObjectPayload(o map[string]interface{}) (io.Reader, error) {
	src, err := JSONMarshalIndent(o, "", "    ")
	if err != nil {
		return nil, err
	}
	return makePayload(src)
}

func makeRequest(u string, method string, payload io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, u, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", contentType)
	client := new(http.Client)
	return client.Do(req)
}

func assertHTTPStatus(expected int, got int) error {
	if got == http.StatusNotImplemented {
		return fmt.Errorf("expected status %s, got %s <- %s",
			http.StatusText(expected), http.StatusText(got), http.StatusText(http.StatusNotImplemented))
	}
	if expected != got {
		return fmt.Errorf("expected http status %q, got %q", http.StatusText(expected), http.StatusText(got))
	}
	return nil
}

func sameObjectSrc(expected []byte, got []byte) bool {
	expectedO := map[string]interface{}{}
	gotO := map[string]interface{}{}
	if err := json.Unmarshal(expected, &expectedO); err != nil {
		return false
	}
	if err := json.Unmarshal(got, &gotO); err != nil {
		return false
	}
	if len(expectedO) != len(gotO) {
		return false
	}
	for k, v := range expectedO {
		gotV, ok := gotO[k]
		if !ok {
			return false
		}
		//FIXME: assumes simple type, i.e. not object or array
		eT := fmt.Sprintf("%T", v)
		gT := fmt.Sprintf("%T", gotV)
		if eT != gT {
			return false
		}
		if v != gotV {
			return false
		}
	}
	return true
}

func clientTestVersion(t *testing.T, settings *Settings) {
	fmt.Printf("starting client test versions\n")
	// Run through a set of the
	u := fmt.Sprintf("http://%s/api/version", settings.Host)
	res, err := http.Get(u)
	if err != nil {
		t.Errorf("failed to get version of api, %s", err)
		t.FailNow()
	}
	defer res.Body.Close()
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if len(body) == 0 {
		t.Errorf("expected a response body, got %q", body)
		t.FailNow()
	}
}

func clientTestKeys(t *testing.T, settings *Settings) {
	fmt.Printf("starting client test keys\n")

	for _, cfg := range settings.Collections {
		c, err := Open(cfg.CName)
		if err != nil {
			t.Errorf("Open(%q) failed, %s", cfg.CName, err)
			t.FailNow()
		}
		expectedKeys, err := c.Keys()
		if err != nil {
			t.Errorf("c.Key() %s", err)
			t.FailNow()
		}
		cName := path.Base(cfg.CName)
		// Run through a set of the
		u := fmt.Sprintf("http://%s/api/%s/keys", settings.Host, cName)
		res, err := http.Get(u)
		if err != nil {
			t.Errorf("failed to get keys from api, %s", err)
			t.FailNow()
		}
		defer res.Body.Close()
		if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
			t.Error(err)
			t.FailNow()
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		if len(body) == 0 {
			t.Errorf("expected a response body, got %q", body)
			t.FailNow()
		}
		keys := []string{}
		err = json.Unmarshal(body, &keys)
		if err != nil {
			t.Errorf("failed to unmarshal\n%s\n%s", body, err)
			t.FailNow()
		}
		if len(expectedKeys) != len(keys) {
			t.Errorf("expected %d keys, got %d", len(expectedKeys), len(keys))
			t.FailNow()
		}
		sort.Strings(expectedKeys)
		sort.Strings(keys)
		for i := 0; i < len(expectedKeys); i++ {
			if expectedKeys[i] != keys[i] {
				t.Errorf("expected key %q, got %q", expectedKeys[i], keys[i])
				t.FailNow()
			}
		}
	}
}

func clientTestObjects(t *testing.T, settings *Settings) {
	fmt.Printf("starting client test add objects\n")

	newRecords := map[string]map[string]interface{}{
		"Flanders-J": {
			"name":        "Jack Flanders",
			"played-by":   "Robert Lorick",
			"character":   true,
			"description": "Metaphysical Detective",
		},
		"Freda-L": {
			"name":        "Little Freda",
			"played-by":   "P.J. O'Rorke",
			"character":   true,
			"description": "A wise Venusian",
		},
		"Sam-M": {
			"name":        "Mojo Sam",
			"played-by":   "Dave Adams",
			"character":   true,
			"description": "The wise You-do man",
		},
	}
	newPost := make(map[string][]byte)
	for k, o := range newRecords {
		src, _ := JSONMarshalIndent(o, "", "    ")
		newPost[k] = src
	}

	for _, cfg := range settings.Collections {
		cName := path.Base(cfg.CName)
		if cName == "objects_test.ds" {
			//FIXME: only work with object_tests.ds
			for k, v := range newPost {
				// Add an Object
				u := fmt.Sprintf("http://%s/api/%s/object/%s", settings.Host, cName, k)
				body := bytes.NewBuffer(v)
				res, err := makeRequest(u, http.MethodPost, body)
				if err != nil {
					t.Errorf("http.Post(%q, %q, %s) error %s", u, contentType, body, err)
					t.FailNow()
				}
				if err := assertHTTPStatus(http.StatusCreated, res.StatusCode); err != nil {
					t.Error(err)
					t.FailNow()
				}
				// Read Back Object
				res, err = makeRequest(u, http.MethodGet, nil)
				if err != nil {
					t.Errorf("makeRequest(%q, %q, nil) -> %s", u, http.MethodGet, err)
					t.FailNow()
				}
				if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
					t.Error(err)
					t.FailNow()
				}
				src, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				if !sameObjectSrc(v, src) {
					t.Errorf("expected %s, got %s", v, src)
				}

				// We lost Robert Lorick and Dave Adams.
				if k == "Flanders-J" || k == "Sam-M" {
					o := map[string]interface{}{}
					json.Unmarshal(v, &o)
					o["active"] = false
					src, _ := JSONMarshal(o)
					body = bytes.NewBuffer(src)
					res, err := makeRequest(u, http.MethodPut, body)
					if err != nil {
						t.Errorf("makeRequest(%q, %q, %s) %s", u, http.MethodPut, body, err)
					}
					if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
						t.Error(err)
						t.FailNow()
					}
				} else if k == "Freda-L" {
					// Little Freda went on to a different plain of existence
					res, _ := makeRequest(u, http.MethodDelete, nil)
					if err != nil {
						t.Errorf("delete failed, %s", err)
					}
					if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
						t.Error(err)
						t.FailNow()
					}
				}
				// Run test on query support
				if qStmt, ok := cfg.QueryFn[query]; ok {
					fmt.Fprintf(os.Stderr, "DEBUG cName: %q DsnURI: %q, qStmt: %s\n", cName, cfg.DsnURI, qStmt)
					u := fmt.Sprintf("http://%s/api/%s/query/%s", settings.Host, cName, query)
					res, err := makeRequest(u, http.MethodGet, nil)
					if err != nil {
						t.Errorf("http.Get(%q, %q, %+v) error %s", u, contentType, nil, err)
						t.FailNow()
					}
					fmt.Fprintf(os.Stderr, "DEBUG res -> %+v\n", res)
					if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
						t.Error(err)
						t.FailNow()
					}
				} else {
					fmt.Fprintf(os.Stderr, "DEBUG no query defined for %q -> %+v\n", cName, cfg.QueryFn)
				}
			}
		}
	}
}

func clientTestAttachments(t *testing.T, settings *Settings) {
	cPath := path.Join(dName, "attachment_test.ds")
	cName := path.Base(cPath)
	/*
		c, err := Open(cPath)
		if err != nil {
			t.Errorf("Failed to open test collection %q, %s", cPath, err)
		}
		defer c.Close()
	*/

	// Write out a file to attach, attach it and then test "attachments"
	// route.
	key := "123"
	src := []byte(`one,two,three
1,3,2
4,7,6
10,100,50
131,313,113
`)
	aName := "numbers.csv"
	if _, err := os.Stat(aName); err == nil {
		os.RemoveAll(aName)
	}
	if err := ioutil.WriteFile(aName, src, 0664); err != nil {
		t.Errorf("failed to write %q, %s", aName, err)
		t.FailNow()
	}

	// Retrieve the attachments (should be no attachments first time)
	u := fmt.Sprintf("http://%s/api/%s/attachments/%s", settings.Host, cName, key)
	res, err := makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) -> %s", u, http.MethodGet, err)
		t.FailNow()
	}
	defer res.Body.Close()
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if len(body) == 0 {
		t.Errorf("expected a response body, got %q", body)
		t.FailNow()
	}
	l := []string{}
	if err := json.Unmarshal(body, &l); err != nil {
		t.Errorf("expected a list of attachments %q, %s", len(l), err)
		t.FailNow()
	}
	if len(l) > 0 {
		t.Errorf("expected a no attachments, got %+v", l)
		t.FailNow()
	}

	// Adding a new document
	fName := "numbers.csv"
	u = fmt.Sprintf("http://%s/api/%s/attachment/%s/%s", settings.Host, cName, key, fName)
	payload, err := makePayload(src)
	res, err = makeRequest(u, http.MethodPost, payload)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, %q) -> %s", u, http.MethodPost, payload, err)
		t.FailNow()
	}
	defer res.Body.Close()
	if err := assertHTTPStatus(http.StatusCreated, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if len(body) == 0 {
		t.Errorf("expected a response body, got %q", body)
		t.FailNow()
	}

	// Retrieve the attachments (should be one attachment second time)
	u = fmt.Sprintf("http://%s/api/%s/attachments/%s", settings.Host, cName, key)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) -> %s", u, http.MethodGet, err)
		t.FailNow()
	}
	defer res.Body.Close()
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if len(body) == 0 {
		t.Errorf("expected a response body, got %q", body)
		t.FailNow()
	}
	l = []string{}
	if err := json.Unmarshal(body, &l); err != nil {
		t.Errorf("expected a list of attachments %q, %s", len(l), err)
		t.FailNow()
	}
	if len(l) != 1 {
		t.Errorf("expected a no attachments, got %+v", l)
		t.FailNow()
	}
	if l[0] != path.Base(aName) {
		t.Errorf("expected %q, got, %q", path.Base(aName), l[0])
		t.FailNow()
	}

	// Get the file we added and make sure it looks OK
	for _, filename := range l {
		u = fmt.Sprintf("http://%s/api/%s/attachment/%s/%s", settings.Host, cName, key, filename)
		res, err = makeRequest(u, http.MethodGet, nil)
		if err != nil {
			t.Errorf("makeRequest(%q, %q, nil) -> %s", u, http.MethodGet, err)
			t.FailNow()
		}
		defer res.Body.Close()
		if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
			t.Error(err)
			t.FailNow()
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("%s", err)
			t.FailNow()
		}
		if len(body) == 0 {
			t.Errorf("expected a response body, got %q", body)
			t.FailNow()
		}
		outName := "download-"+filename
		if err := ioutil.WriteFile(outName, body, 0664); err != nil {
			t.Errorf("unable to write requested file %q, %s", filename, err)
		}
	}
	for _, filename := range l {
		u = fmt.Sprintf("http://%s/api/%s/attachment/%s/%s", settings.Host, cName, key, filename)
		res, err = makeRequest(u, http.MethodDelete, nil)
		if err != nil {
			t.Errorf("makeRequest(%q, %q, nil) -> %s", u, http.MethodGet, err)
			t.FailNow()
		}
		defer res.Body.Close()
		if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
	// Now make sure attachments were all removed.
	u = fmt.Sprintf("http://%s/api/%s/attachments/%s", settings.Host, cName, key)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) -> %s", u, http.MethodGet, err)
		t.FailNow()
	}
	defer res.Body.Close()
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if len(body) == 0 {
		t.Errorf("expected a response body, got %q", body)
		t.FailNow()
	}
	l = []string{}
	if err := json.Unmarshal(body, &l); err != nil {
		t.Errorf("expected a list of attachments %q, %s", len(l), err)
		t.FailNow()
	}
	if len(l) > 0 {
		t.Errorf("expected a no attachments, got %+v", l)
		t.FailNow()
	}
}

func TestRunAPI(t *testing.T) {
	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.MkdirAll(dName, 0775)
	}
	// Setup up a test collection
	fName := path.Join(dName, "settings.yaml")
	if _, err := os.Stat(fName); err == nil {
		os.RemoveAll(fName)
	}
	settings := new(Settings)
	settings.Host = "localhost:8585"
	testCollections := map[string]interface{}{}
	if err := json.Unmarshal(testCollectionsSrc, &testCollections); err != nil {
		t.Errorf("Failed to unpack text data to populate collections, %s", err)
		t.FailNow()
	}
	for cName, testRecords := range testCollections {
		records := map[string]map[string]interface{}{}
		for k, v := range testRecords.(map[string]interface{}) {
			records[k] = v.(map[string]interface{})
		}
		pName := path.Join(dName, cName)
		if _, err := os.Stat(pName); err == nil {
			os.RemoveAll(pName)
		}
		dsnURI := "sqlite://collection.db"
		if err := setupApiTestCollection(pName, dsnURI, records); err != nil {
			t.Errorf("Failed to setup test collection %q, %s", pName, err)
			t.FailNow()
		}
		// Setup a test configuration
		cfg := new(Config)
		cfg.CName = cName
		cfg.DsnURI = dsnURI
		cfg.Keys = true
		cfg.Create = true
		cfg.Read = true
		cfg.Update = true
		cfg.Delete = true
		cfg.Attachments = true
		cfg.Attach = true
		cfg.Retrieve = true
		cfg.Prune = true
		cfg.QueryFn = map[string]string{}
		tName := strings.TrimSuffix(cName, ".ds")
		cfg.QueryFn[query] = fmt.Sprintf(`select count(*) as src from %s`, tName)
		if settings.Collections == nil {
			settings.Collections = []*Config{}
		}
		settings.Collections = append(settings.Collections, cfg)
	}
	//fmt.Fprintf(os.Stderr, "DEBUG writing test YAML to %q\n", fName)
	//src, _ := YAMLMarshal(settings)
	//fmt.Fprintf(os.Stderr, "DEBUG yaml src ->\n%s\n", src)
	if err := settings.WriteFile(fName, 0664); err != nil {
		t.Errorf("failed to save config %q, %s", fName, err)
		t.FailNow()
	}

	setupWait := "5s"
	wait, _ := time.ParseDuration(setupWait)
	fmt.Printf(`
Launching API Tests at http://%s

Press Ctr-C if tests hang, you should see requests log output

Client testings starts in %s (s = seconds)
`, settings.Host, setupWait)
	appName := os.Args[0]

	// Add some records to our test collection.

	// Send our web service into the back ground so we can run
	// a client test.
	debug := false
	go func() {
		os.Chdir(dName)
		fName = "settings.yaml"
		wDir, _ := os.Getwd()
		fmt.Fprintf(os.Stderr, "INFO Running tests %q %q %t from %s\n", path.Base(appName), fName, debug, path.Base(wDir))
		if err := RunAPI(appName, fName, debug); err != nil {
			t.Errorf("Expected API to setup and run, %s", err)
		}
	}()
	time.Sleep(wait)
	clientTestVersion(t, settings)
	clientTestKeys(t, settings)
	clientTestObjects(t, settings)
	clientTestKeys(t, settings)
	clientTestAttachments(t, settings)
}
