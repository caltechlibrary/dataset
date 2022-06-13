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
	records = map[string]map[string]interface{}{
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
)

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
	if res.StatusCode != 200 {
		t.Errorf("unexpected response, %q", res.Status)
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
		if res.StatusCode != 200 {
			t.Errorf("unexpected response, %q", res.Status)
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
		//fmt.Printf("DEBUG body -> %s\n", body)
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

	mimeType := "application/json"
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
			res, err := http.Post(u, mimeType, body)
			if err != nil {
				t.Errorf("http.Post(%q, %q, %s) error %s", u, mimeType, body, err)
				t.FailNow()
			}
			if res.StatusCode != http.StatusOK {
				t.Errorf("expected http (POST) status %q, got %q", http.StatusText(http.StatusOK), http.StatusText(res.StatusCode))
				t.FailNow()
			}
			// Read Back Object
			res, err = http.Get(u)
			if err != nil {
				t.Errorf("expected http (GET) status %q, got %q", http.StatusText(http.StatusOK), http.StatusText(res.StatusCode))
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
				req, _ := http.NewRequest(http.MethodPut, u, body)
				req.Header.Set("content-type", mimeType)
				client := new(http.Client)
				res, err := client.Do(req)
				if err != nil {
					t.Errorf("put failed, %s", err)
				}
				if res.StatusCode != http.StatusOK {
					t.Errorf("expected status code %q, got %q", http.StatusText(http.StatusOK), http.StatusText(res.StatusCode))
					t.FailNow()
				}
			} else if k == "Freda-L" {
				// Little Freda went on to a different plain of existence
				req, _ := http.NewRequest(http.MethodDelete, u, nil)
				req.Header.Set("content-type", mimeType)
				client := new(http.Client)
				res, err := client.Do(req)
				if err != nil {
					t.Errorf("delete failed, %s", err)
				}
				if res.StatusCode != http.StatusOK {
					t.Errorf("expected status code %q, got %q", http.StatusText(http.StatusOK), http.StatusText(res.StatusCode))
					t.FailNow()
				}
			}
		}
	}
}

func TestRunAPI(t *testing.T) {
	dName := "testout"
	if _, err := os.Stat(dName); os.IsNotExist(err) {
		os.MkdirAll(dName, 0775)
	}
	// Setup up a test collection
	cName := path.Join(dName, "apitest.ds")
	if _, err := os.Stat(cName); err == nil {
		os.RemoveAll(cName)
	}
	dbName := path.Join(dName, "collections.db")
	if _, err := os.Stat(dbName); err == nil {
		os.RemoveAll(dbName)
	}

	dsnURI := "sqlite://" + dbName
	if err := SetupTestCollection(cName, dsnURI, records); err != nil {
		t.Errorf("failed to setup %q %q, %s", cName, dsnURI, err)
		t.FailNow()
	}

	// Setup a test configuration
	fName := path.Join("testout", "settings.json")
	if _, err := os.Stat(fName); err == nil {
		os.RemoveAll(fName)
	}
	settings := new(config.Settings)
	settings.Host = "localhost:8585"
	cfg := new(config.Config)
	cfg.CName = cName
	cfg.DsnURI = dsnURI
	cfg.Keys = true
	cfg.Create = true
	cfg.Read = true
	cfg.Update = true
	cfg.Delete = true
	cfg.Attach = false
	cfg.Retrieve = false
	cfg.Prune = false
	settings.Collections = append(settings.Collections, cfg)
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
}
