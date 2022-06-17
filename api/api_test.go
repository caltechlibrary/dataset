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
	"github.com/caltechlibrary/dataset/dotpath"
)

var (
	contentType        = "application/json"
	dName              = "testout"
	testCollectionsSrc = []byte(`{
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
	src1, _ := json.MarshalIndent(expected, "", "    ")
	src2, _ := json.MarshalIndent(got, "", "    ")
	return bytes.Compare(src1, src2) == 0
}

// sameMapping compares one map to another by converting
// them to a JSON representation and comparing the representation.
func sameMapping(expected map[string]interface{}, got map[string]interface{}) bool {
	src1, _ := json.MarshalIndent(expected, "", "    ")
	src2, _ := json.MarshalIndent(got, "", "    ")
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
}

func clientTestAttachments(t *testing.T, settings *config.Settings) {
	t.Errorf("clientTestAttachments() not implemented.")
}

func clientTestFrames(t *testing.T, settings *config.Settings) {
	cPath := path.Join(dName, "frames_test.ds")
	cName := path.Base(cPath)
	c, err := ds.Open(cPath)
	if err != nil {
		t.Errorf("Failed to open test collection %q, %s", cPath, err)
	}
	defer c.Close()

	// Double check loaded test data
	frmKeys := []string{"Sam-M", "Freda-L", "Flanders-J"}
	attributes := []string{"family", "given"}
	dotPaths := []string{".family", ".given"}
	labels := []string{"Family Name", "Given Name"}
	for _, key := range frmKeys {
		m := map[string]interface{}{}
		if err := c.Read(key, m); err != nil {
			t.Errorf("could not open %q, %s", key, err)
			t.FailNow()
		}
		for _, attr := range attributes {
			if given, ok := m[attr]; !ok {
				t.Errorf("failed to find .%s in %s, %+v", attr, key, m)
			} else if given == "" {
				t.Errorf("failed to find .%s is empty string %s, %+v", attr, key, m)
			}
			expr := fmt.Sprintf(".%s", attr)
			val, err := dotpath.Eval(expr, m)
			if err != nil {
				t.Errorf("dotpath.Eval(%q, %+v) failed, %s", expr, m, err)
			}
			if val == nil {
				t.Errorf("dotpath.Eval(%q, %+v) should not return nil value", expr, m)
				t.FailNow()
			}
		}
	}
	exFrameName := "expected"
	expectedFrm, err := c.FrameCreate(exFrameName, frmKeys, dotPaths, labels, true)
	if err != nil {
		t.Errorf("should be able to created the %q frame, %s", exFrameName, err)
	}
	if expectedFrm == nil {
		t.Errorf("something went really wrong, should have an %q frame", exFrameName)
	}

	// Check to make sure frame does not exist
	frameName := "test_names"
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

	// Create the frame "test_names"
	frameDef := map[string]interface{}{
		"keys":      frmKeys,
		"dot_paths": dotPaths,
		"labels":    labels,
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
	// Attempt to read the frame that was just created.
	frm, err := c.FrameRead(frameName)
	if err != nil {
		t.Errorf("frame failed to be created %q, %s", frameName, err)
		t.FailNow()
	}
	if len(frm.ObjectMap) == 0 {
		t.Errorf("something went wrong, frm.ObjectMap should have objects, %+v", frm)
		t.FailNow()
	}
	hasError := false
	for k, v := range frm.ObjectMap {
		if v == nil || len(v.(map[string]interface{})) == 0 {
			t.Errorf("something went wrong, frm.ObjectMap[%q] should have a populated map[string]interface{}, %+v", k, frm.ObjectMap)
			hasError = true
		}
	}
	if hasError {
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
	if len(frameList) != 2 {
		t.Errorf("expected two frames defined, got %d", len(frameList))
		t.FailNow()
	}
	foundFrame := false
	for _, name := range frameList {
		if name == frameName {
			foundFrame = true
			break
		}
	}
	if !foundFrame {
		t.Errorf("expected frame %q, it was missing from %+v", frameName, frameList)
		t.FailNow()
	}

	// Get the frame's def
	expectedDef, err := c.FrameDef(frameName)
	if err != nil {
		t.Errorf("expected to find frame def for %q, got %s", frameName, err)
		t.FailNow()
	}
	u = fmt.Sprintf("http://%s/api/%s/frame-def/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	src, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected to read the body of the request, %s", err)
	}
	res.Body.Close()
	def := map[string]interface{}{}
	if err := json.Unmarshal(src, &def); err != nil {
		t.Errorf("json.Unmarshal(%s, def) failed, %s", src, err)
		t.FailNow()
	}
	if !sameMapping(expectedDef, def) {
		exSrc, _ := json.MarshalIndent(expectedDef, "", "    ")
		t.Errorf("expected map %s, got %s", exSrc, src)
	}

	expectedKeys := c.FrameKeys(frameName)
	u = fmt.Sprintf("http://%s/api/%s/frame-keys/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	src, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected to read the body of the request, %s", err)
	}
	res.Body.Close()
	gotKeys := []string{}
	if err := json.Unmarshal(src, &gotKeys); err != nil {
		t.Errorf("json.Unmarshal(%s, gotKeys) failed, %s", src, err)
		t.FailNow()
	}
	if !sameStrings(expectedKeys, gotKeys) {
		exSrc, _ := json.MarshalIndent(expectedKeys, "", "    ")
		t.Errorf("expected map %s, got %s", exSrc, src)
	}

	expectedObjects, err := c.FrameObjects(frameName)
	if err != nil {
		t.Errorf("Could get get %q objects, %s", frameName, err)
		t.FailNow()
	}
	u = fmt.Sprintf("http://%s/api/%s/frame-objects/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodGet, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	src, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected to read the body of the request, %s", err)
	}
	res.Body.Close()
	gotObjects := []map[string]interface{}{}
	if err := json.Unmarshal(src, &gotObjects); err != nil {
		t.Errorf("json.Unmarshal(%s, gotObjects) failed, %s", src, err)
		t.FailNow()
	}
	if !sameObjects(expectedObjects, gotObjects) {
		exSrc, _ := json.MarshalIndent(expectedObjects, "", "    ")
		t.Errorf("expected map %s, got %s", exSrc, src)
	}

	expectedObjects, err = c.FrameObjects(frameName)
	if err != nil {
		t.Errorf("Could get get %q objects, %s", frameName, err)
		t.FailNow()
	}

	// save old object list
	oldObjects := gotObjects[:]
	// Update a record so we can then try frame-refresh
	key := "Flanders-J"
	buf := bytes.NewBuffer([]byte(`{
    "id":        "Flanders-J",
    "given":     "Capt. Marcelle",
    "family":    "Du Champ",
	"character":   true,
	"description": "Sky Pirate Captain"
}`))
	u = fmt.Sprintf("http://%s/api/%s/object/%s", settings.Host, cName, key)
	res, err = makeRequest(u, http.MethodPut, buf)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Make a refresh call
	u = fmt.Sprintf("http://%s/api/%s/frame-refresh/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodPut, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodGet, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// Verify records changed.
	gotObjects, err = c.FrameObjects(frameName)
	if err != nil {
		t.Errorf("Could get get %q objects, %s", frameName, err)
		t.FailNow()
	}
	if sameObjects(oldObjects, gotObjects) {
		t.Errorf("FrameRefresh failed to update objects")
		t.FailNow()
	}

	//FIXME: need tests for reframe
	newKeys, err := c.Keys()
	if err != nil {
		t.Errorf("failed to get keys from %q, %s", cName, err)
		t.FailNow()
	}
	src, err = json.MarshalIndent(newKeys, "", "    ")
	if err != nil {
		t.Errorf("failed to marshal key list, %s", err)
		t.FailNow()
	}
	u = fmt.Sprintf("http://%s/api/%s/frame-reframe/%s", settings.Host, cName, frameName)
	buf = bytes.NewBuffer(src)
	res, err = makeRequest(u, http.MethodPut, buf)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodPut, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	updatedKeys := c.FrameKeys(frameName)
	if !sameStrings(newKeys, updatedKeys) {
		t.Errorf("expected keys %+v, got %+v", newKeys, updatedKeys)
		t.FailNow()
	}
	objects, err := c.FrameObjects(frameName)
	if len(objects) != len(newKeys) {
		t.Errorf("expected %d objects, got %d", len(newKeys), len(objects))
		t.FailNow()
	}

	frameName = "expected"
	u = fmt.Sprintf("http://%s/api/%s/frame/%s", settings.Host, cName, frameName)
	res, err = makeRequest(u, http.MethodDelete, nil)
	if err != nil {
		t.Errorf("makeRequest(%q, %q, nil) %s", u, http.MethodDelete, err)
		t.FailNow()
	}
	if err := assertHTTPStatus(http.StatusOK, res.StatusCode); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if c.HasFrame(frameName) {
		t.Errorf("expected frame %q to be deleted, found it", frameName)
		t.FailNow()
	}
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
	testCollections := map[string]interface{}{}
	if err := json.Unmarshal(testCollectionsSrc, &testCollections); err != nil {
		t.Errorf("Failed to unpack text data to populate collections, %s", err)
		t.FailNow()
	}
	for name, testRecords := range testCollections {
		records := map[string]map[string]interface{}{}
		for k, v := range testRecords.(map[string]interface{}) {
			records[k] = v.(map[string]interface{})
		}
		pName := path.Join(dName, name)
		if _, err := os.Stat(pName); err == nil {
			os.RemoveAll(pName)
		}

		dbName := path.Join(pName, "collection.db")
		dsnURI := "sqlite://" + dbName
		if err := SetupTestCollection(pName, dsnURI, records); err != nil {
			t.Errorf("Failed to setup test collection %q, %s", pName, err)
			t.FailNow()
		}
		// Setup a test configuration
		cfg := new(config.Config)
		cfg.CName = pName
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
