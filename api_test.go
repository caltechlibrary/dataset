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
	"path/filepath"
	"sort"
	"testing"
	"time"

	// 3rd Party library
	"gopkg.in/yaml.v3"
)

var (
	contentType        = "application/json"
	dName              = "testout"
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
	},
	"frames_test.ds": {
		"Miller-A": {
			"id":        "Miller-A",
			"given":     "Arthor",
			"family":    "Miller",
			"character": false,
			"vocations": ["playright", "writer", "author"]
		},
		"Lopez-T": {
			"id":        "Lopez-T",
			"given":     "Tom",
			"family":    "Lopez",
			"character": false,
			"vocations": ["playright", "producer", "director", "sound-engineer", "voice actor", "disc jockey"]
		},
		"Flanders-J": {
			"id":          "Flanders-J",
			"given":       "Jack",
			"family":      "Flanders",
			"played-by":   "Robert Lorick",
			"character":   true,
			"description": "Metaphysical Detective"
		},
		"Freda-L": {
			"id":          "Freda-L",
			"given":       "Little",
			"family":      "Freda",
			"played-by":   "P.J. O'Rorke",
			"character":   true,
			"description": "A wise Venusian"
		},
		"Sam-M": {
			"id":          "Sam-M",
			"given":       "Mojo",
			"family":      "Sam",
			"played-by":   "Dave Adams",
			"character":   true,
			"description": "The wise You-do man"
		}
	}
}`)
)

// sameStrings compares one slice of strings to another by converting
// them to a JSON represention and compariing the representation.
func sameStrings(expected []string, got []string) bool {
	src1, _ := json.Marshal(expected)
	src2, _ := json.Marshal(got)
	return bytes.Compare(src1, src2) == 0
}

// sameMapping compares one map to another by converting
// them to a JSON representation and comparing the representation.
func sameMapping(expected map[string]interface{}, got map[string]interface{}) bool {
	src1, _ := json.Marshal(expected)
	src2, _ := json.Marshal(got)
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
	src, err := json.Marshal(o)
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
		cName := filepath.Base(cfg.CName)
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
		src, _ := json.Marshal(o)
		newPost[k] = src
	}

	for _, cfg := range settings.Collections {
		cName := filepath.Base(cfg.CName)
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
					src, _ := json.Marshal(o)
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
			}
		}
	}
}


func TestRunAPI(t *testing.T) {
	
	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.RemoveAll(dName)
	}
	os.MkdirAll(path.Join(dName, "htdocs"), 0775)
	// Setup up a test collection
	dbName := "recipes"
	cName := path.Join(dName, fmt.Sprintf("%s.ds", dbName))
	_, err := Init(cName, "sqlite://collection.db")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fName := path.Join(dName, fmt.Sprintf("%s_api.yaml", dbName))
	if _, err := os.Stat(fName); err == nil {
		os.RemoveAll(fName)
	}
	err = os.WriteFile(fName, []byte(`host: localhost:8010
htdocs: testout/htdocs
collections:
  - dataset: recipes.ds
    keys: true
    create: true
    read: true
    update: true
    codemeta: true
    query:
      list_objects: |
        select src
        from recipes
        order by _Key
      list_recent: |
        select src
        from recipes
        order by created desc
        limit ?,?
`), 0666)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	src, err := os.ReadFile(fName)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	settings := &Settings{}
	if err := yaml.Unmarshal(src, &settings); err != nil {
		t.Error(err)
		t.FailNow()
	}

	appName := filepath.Base(appName)

	setupWait := "5s"
	wait, _ := time.ParseDuration(setupWait)
	fmt.Printf(`
Launching API Tests at http://%s

Press Ctr-C if tests hang, you should see requests log output

Client testings starts in %s (s = seconds)
`, settings.Host, setupWait)

	// Add some records to our test collection.

	// Send our web service into the back ground so we can run
	// a client test.
	debug := false
	go func() {
		if err := RunAPI(appName, fName, debug); err != nil {
			t.Errorf("Expected API to setup and run, %s", err)
		}
	}()
	time.Sleep(wait)
	clientTestVersion(t, settings)
	clientTestKeys(t, settings)
	clientTestObjects(t, settings)
	clientTestKeys(t, settings)
}
