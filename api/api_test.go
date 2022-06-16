package api

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
	"testing"
	"time"

	// Caltech Library packages
	ds "github.com/caltechlibrary/dataset"
	"github.com/caltechlibrary/dataset/config"
)

var (
	contentType = "application/json"
	dName       = "testout"
	records     = map[string]map[string]interface{}{
		"Miller-A": {
			"id":     "Miller-A",
			"given":  "Arthor",
			"family": "Miller",
			"genre":  []string{"plays", "film"},
		},
		"Lopez-T": {
			"id":     "Lopez-T",
			"given":  "Tom",
			"family": "Lopez",
			"genre":  []string{"plays", "radio-theater"},
		},
	}

	frames = map[string]map[string]string{
		"name": {
			".given":  "Given Name",
			".family": "Family Name",
		},
	}
)

func makePayload(src []byte) (io.Reader, error) {
	return bytes.NewBuffer(src), nil
}

func makeObjectPayload(o map[string]interface{}) (io.Reader, error) {
	src, err := json.MarshalIndent(o, "", "    ")
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

func clientTestVersion(t *testing.T, settings *config.Settings) {
	fmt.Printf("starting client test verisons\n")
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

func clientTestKeys(t *testing.T, settings *config.Settings) {
	fmt.Printf("starting client test keys\n")

	for _, cfg := range settings.Collections {
		c, err := ds.Open(cfg.CName)
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

func clientTestObjects(t *testing.T, settings *config.Settings) {
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
		src, _ := json.MarshalIndent(o, "", "    ")
		newPost[k] = src
	}

	for _, cfg := range settings.Collections {
		cName := path.Base(cfg.CName)
		for k, v := range newPost {
			// Add an Object
			u := fmt.Sprintf("http://%s/api/%s/object/%s", settings.Host, cName, k)
			body := bytes.NewBuffer(v)
			res, err := makeRequest(u, http.MethodPost, body)
			if err != nil {
				t.Errorf("http.Post(%q, %q, %s) error %s", u, contentType, body, err)
				t.FailNow()
			}
			if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
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
				src, _ := json.Marshal(o)
				body = bytes.NewBuffer(src)
				res, err := makeRequest(u, http.MethodPut, body)
				if err != nil {
					t.Errorf("put failed, %s", err)
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
		}
	}
}

func clientTestAttachments(t *testing.T, settings *config.Settings) {
	t.Errorf("clientTestAttachments() not implemented.")
}

func clientTestFrames(t *testing.T, settings *config.Settings) {
	framedRecords := map[string]map[string]interface{}{
		"Miller-A": {
			"id":        "Miller-A",
			"given":     "Arthor",
			"family":    "Miller",
			"character": false,
			"vocations": []string{"playright", "writer", "author"},
		},
		"Lopez-T": {
			"id":        "Lopez-T",
			"given":     "Tom",
			"family":    "Lopez",
			"character": false,
			"vocations": []string{"playright", "producer", "director", "sound-engineer", "voice actor", "disc jockey"},
		},
		"Flanders-J": {
			"given":       "Jack",
			"family":      "Jack Flanders",
			"played-by":   "Robert Lorick",
			"character":   true,
			"description": "Metaphysical Detective",
		},
		"Freda-L": {
			"given":       "Little",
			"family":      "Freda",
			"played-by":   "P.J. O'Rorke",
			"character":   true,
			"description": "A wise Venusian",
		},
		"Sam-M": {
			"given":       "Mojo",
			"family":      "Sam",
			"played-by":   "Dave Adams",
			"character":   true,
			"description": "The wise You-do man",
		},
	}

	cPath := path.Join(dName, "frames_test.ds")
	dbName := path.Join(cPath, "collections.db")
	cName := path.Base(cPath)
	dsnURI := "sqlite://" + dbName
	if err := SetupTestCollection(cPath, dsnURI, framedRecords); err != nil {
		t.Errorf("SetupTestCollection(%q, %+v) -> %s", cName, records, err)
		t.FailNow()
	}

	// Check to make sure frame does not exist
	frameName := "names"
	u := fmt.Sprintf("http://%s/api/%s/has-frame/%s", settings.Host, cName, frameName)
	res, err := makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusNotFound, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if _, err := ioutil.ReadAll(res.Body); err != nil {
		t.Errorf("expected to read the body of reqest, %s", err)
		t.FailNow()
	}
	res.Body.Close()

	// Create the frame "names"
	frameDef := map[string]interface{}{
		"dot_paths": []string{".given", ".family"},
		"labels":    []string{"Given Name", "Family Name"},
		"keys":      []string{"Sam-M", "Flanders-J"},
	}
	payload, err := makeObjectPayload(frameDef)
	u = fmt.Sprintf("http://%s/api/%s/frame/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodPost, payload)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, %s) -> %s", u, http.MethodPost, payload, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Check if we have the "names" frame.
	u = fmt.Sprintf("http://%s/api/%s/has-frame/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if _, err := ioutil.ReadAll(res.Body); err != nil {
		t.Errorf("expected to read the body of request, %s", err)
	}
	res.Body.Close()

	// List the frames in a collection
	u = fmt.Sprintf("http://%s/api/%s/frames/", settings.Host, cName)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	src, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected to read the body of the request, %s", err)
	}
	res.Body.Close()
	frameList := []string{}
	if err := json.Unmarshal(src, &frameList); err != nil {
		t.Errorf("Unmarshal(%s, %+v) -> %s", src, frameList, err)
	}
	if len(frameList) != 1 {
		t.Errorf("expected on frame defined, got %d", len(frameList))
		t.FailNow()
	}
	if frameList[0] != frameName {
		t.Errorf("expected frame name %q, got %q", frameName, frameList[0])
		t.FailNow()
	}

	//	t.Errorf("clientTestFrames() incomplete")
}

func TestRunAPI(t *testing.T) {
	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.MkdirAll(dName, 0775)
	}
	// Setup up a test collection
	fName := path.Join(dName, "settings.json")
	if _, err := os.Stat(fName); err == nil {
		os.RemoveAll(fName)
	}
	settings := new(config.Settings)
	settings.Host = "localhost:8585"
	colNames := []string{"apitest.ds", "frames_test.ds"}
	for _, name := range colNames {
		cName := path.Join(dName, name)
		if _, err := os.Stat(cName); err == nil {
			os.RemoveAll(cName)
		}
		dbName := path.Join(cName, "collection.db")
		dsnURI := "sqlite://" + dbName
		if err := SetupTestCollection(cName, dsnURI, records); err != nil {
			t.Errorf("failed to setup %q %q, %s", cName, dsnURI, err)
			t.FailNow()
		}

		// Setup a test configuration
		cfg := new(config.Config)
		cfg.CName = cName
		cfg.DsnURI = dsnURI
		cfg.Keys = true
		cfg.Create = true
		cfg.Read = true
		cfg.Update = true
		cfg.Delete = true
		cfg.Attach = true
		cfg.Retrieve = true
		cfg.Prune = true
		cfg.FrameRead = true
		cfg.FrameWrite = true
		settings.Collections = append(settings.Collections, cfg)
	}
	if err := settings.WriteFile(fName, 0664); err != nil {
		t.Errorf("failed to save config %q, %s", fName, err)
		t.FailNow()
	}

	setupWait := "5s"
	wait, _ := time.ParseDuration(setupWait)
	fmt.Printf(`
Launching API Tests at http://%q

Press Ctr-C if tests hang, you should see requests log output

Client testings starts in %s (s = seconds)
`, settings.Host, setupWait)
	appName := os.Args[0]

	// Add some records to our test collection.

	// Send our web service into the back ground so we can run
	// a client test.
	go func() {
		if err := RunAPI(appName, fName); err != nil {
			t.Errorf("Expected API to setup and run, %s", err)
		}
	}()
	time.Sleep(wait)
	clientTestVersion(t, settings)
	clientTestKeys(t, settings)
	clientTestObjects(t, settings)
	clientTestKeys(t, settings)
	clientTestFrames(t, settings)
	clientTestAttachments(t, settings)
}
