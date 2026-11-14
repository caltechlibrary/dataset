package dataset

// api_status_test.go pins the HTTP status codes datasetd answers with.
//
// TestLegacyStatus records the answers of datasetd 2.5.1, including the ones
// workspace DR-0020 calls wrong. It is the compatibility pin for the v2.5.2
// work (agents/projects/dataset/design/datasetd_http_status_v2_5_2.md): with
// strict_status off, none of these may change.

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// statusHarness is a datasetd API served by httptest over scratch collections.
type statusHarness struct {
	t      *testing.T
	api    *API
	server *httptest.Server
	client *http.Client
}

const (
	// harnessOpenCollection has every method enabled and some queries.
	harnessOpenCollection = "objects.ds"
	// harnessReadOnlyCollection only allows keys and read, so POST is a disabled method.
	harnessReadOnlyCollection = "readonly.ds"
)

// newStatusHarness builds the scratch collections and starts the service. The
// customize func may change the settings before they are written.
//
// @param t {*testing.T} the test
// @param customize {func(*Settings)} edits the settings, or nil
// @returns {*statusHarness} the running harness, closed when the test ends
func newStatusHarness(t *testing.T, customize func(dir string, s *Settings)) *statusHarness {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	for _, cName := range []string{harnessOpenCollection, harnessReadOnlyCollection} {
		c, err := Init(cName, "sqlite://collection.db")
		if err != nil {
			t.Fatalf("Init(%q), %s", cName, err)
		}
		if err := c.Create("k1", map[string]interface{}{"id": "k1", "v": "original"}); err != nil {
			t.Fatalf("Create k1 in %q, %s", cName, err)
		}
		c.Close()
	}
	open := &Config{
		CName: harnessOpenCollection, DsnURI: "sqlite://collection.db",
		Keys: true, Create: true, Read: true, Update: true, Delete: true,
		QueryFn: map[string]string{
			"all":     "select src from objects",
			"by_key":  "select src from objects where _key = ?",
			"by_path": "select json_extract(src, ?) as v from objects",
		},
	}
	readOnly := &Config{
		CName: harnessReadOnlyCollection, DsnURI: "sqlite://collection.db",
		Keys: true, Read: true,
	}
	settings := &Settings{Host: "localhost:0", Collections: []*Config{open, readOnly}}
	if customize != nil {
		customize(dir, settings)
	}
	if err := settings.WriteFile("settings.yaml", 0664); err != nil {
		t.Fatalf("WriteFile settings.yaml, %s", err)
	}
	api := new(API)
	if err := api.Init("datasetd", "settings.yaml"); err != nil {
		t.Fatalf("api.Init, %s", err)
	}
	h := &statusHarness{t: t, api: api}
	h.server = httptest.NewServer(http.HandlerFunc(api.Router))
	h.client = &http.Client{
		// A redirect is an answer to inspect, not to follow.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	t.Cleanup(func() {
		h.server.Close()
		for _, c := range api.CMap {
			c.Close()
		}
	})
	return h
}

// reset puts the open collection back to one record, k1, so cases are independent.
func (h *statusHarness) reset() {
	h.t.Helper()
	c := h.api.CMap[harnessOpenCollection]
	keys, _ := c.Keys()
	for _, k := range keys {
		c.Delete(k)
	}
	if err := c.Create("k1", map[string]interface{}{"id": "k1", "v": "original"}); err != nil {
		h.t.Fatalf("reset, %s", err)
	}
}

// do sends a request and returns the status, the Content-Type, the body and the Location header.
func (h *statusHarness) do(method, path, contentType, body string) (int, http.Header, string) {
	h.t.Helper()
	req, err := http.NewRequest(method, h.server.URL+path, strings.NewReader(body))
	if err != nil {
		h.t.Fatalf("NewRequest, %s", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	res, err := h.client.Do(req)
	if err != nil {
		h.t.Fatalf("%s %s, %s", method, path, err)
	}
	defer res.Body.Close()
	src, _ := io.ReadAll(res.Body)
	return res.StatusCode, res.Header, string(src)
}

// TestLegacyStatus pins every datasetd 2.5.1 answer in the design table.
// A row that disagrees with the table is a finding, not a test to loosen.
func TestLegacyStatus(t *testing.T) {
	h := newStatusHarness(t, nil)
	const j = "application/json"
	over := `{"id":"big","v":"` + strings.Repeat("x", jsonSizeLimit+1) + `"}`
	o := "/api/" + harnessOpenCollection
	cases := []struct {
		name        string
		method      string
		path        string
		contentType string
		body        string
		want        int
	}{
		{"read existing", "GET", o + "/object/k1", "", "", 200},
		{"read missing is 200 with {}", "GET", o + "/object/nope", "", "", 200},
		{"delete existing", "DELETE", o + "/object/k1", "", "", 200},
		{"delete missing is 200", "DELETE", o + "/object/nope", "", "", 200},
		{"create new", "POST", o + "/object/new", j, `{"id":"new"}`, 201},
		{"post to existing overwrites with 200", "POST", o + "/object/k1", j, `{"id":"k1","v":"changed"}`, 200},
		{"update existing", "PUT", o + "/object/k1", j, `{"id":"k1","v":"changed"}`, 200},
		{"put to missing is 200", "PUT", o + "/object/nope", j, `{"id":"nope"}`, 200},
		{"post invalid JSON is 406", "POST", o + "/object/new", j, `{not json`, 406},
		{"post empty body is 406", "POST", o + "/object/new", j, ``, 406},
		{"post text/plain is accepted", "POST", o + "/object/new", "text/plain", `hello`, 201},
		{"post over 1 MiB is 400", "POST", o + "/object/big", j, over, 400},
		{"query ok", "GET", o + "/query/all", "", "", 200},
		{"query missing parameter is 406", "GET", o + "/query/by_key/key", "", "", 406},
		{"query failing on user input is 406", "GET", o + "/query/by_path/path?path=$[", "", "", 406},
		{"query name not defined is 400", "GET", o + "/query/nope", "", "", 400},
		{"unknown collection is 403", "GET", "/api/nope.ds/keys", "", "", 403},
		{"unknown route is 403", "GET", o + "/nope", "", "", 403},
		{"disabled method is 405", "POST", "/api/" + harnessReadOnlyCollection + "/object/new", j, `{"id":"new"}`, 405},
		{"keys", "GET", o + "/keys", "", "", 200},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h.reset()
			got, _, body := h.do(tc.method, tc.path, tc.contentType, tc.body)
			if got != tc.want {
				t.Errorf("%s %s: want %d, got %d, body %q", tc.method, tc.path, tc.want, got, body)
			}
		})
	}

	t.Run("read missing body is {}", func(t *testing.T) {
		h.reset()
		_, _, body := h.do("GET", o+"/object/nope", "", "")
		if strings.TrimSpace(body) != "{}" {
			t.Errorf("want {}, got %q", body)
		}
	})
	t.Run("post to existing silently overwrites", func(t *testing.T) {
		h.reset()
		h.do("POST", o+"/object/k1", j, `{"id":"k1","v":"changed"}`)
		obj := map[string]interface{}{}
		if err := h.api.CMap[harnessOpenCollection].Read("k1", obj); err != nil {
			t.Fatal(err)
		}
		if obj["v"] != "changed" {
			t.Errorf("legacy POST of an existing key overwrites, got %+v", obj)
		}
	})
	t.Run("put to missing creates nothing", func(t *testing.T) {
		h.reset()
		h.do("PUT", o+"/object/nope", j, `{"id":"nope"}`)
		if h.api.CMap[harnessOpenCollection].HasKey("nope") {
			t.Errorf("legacy PUT of a missing key creates nothing, but nope exists")
		}
	})
	t.Run("delete missing leaves the others", func(t *testing.T) {
		h.reset()
		h.do("DELETE", o+"/object/nope", "", "")
		if !h.api.CMap[harnessOpenCollection].HasKey("k1") {
			t.Errorf("k1 should still exist")
		}
	})
}

// TestReadMissingKeyStores finds where the 200 {} of a missing key comes from
// (plan step 0.3): the stores, or the handler.
func TestReadMissingKeyStores(t *testing.T) {
	dir := t.TempDir()
	for name, dsn := range map[string]string{"sql.ds": "sqlite://collection.db", "pairtree.ds": ""} {
		t.Run(name, func(t *testing.T) {
			c, err := Init(filepath.Join(dir, name), dsn)
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			o := map[string]interface{}{}
			err = c.Read("nope", o)
			t.Logf("%s: Read of a missing key returned err=%v obj=%v HasKey=%t", name, err, o, c.HasKey("nope"))
			// The finding is recorded by the log line above; the store that
			// reports no error is what makes the handler answer 200 {}.
			if name == "sql.ds" && err != nil {
				t.Errorf("expected the SQL store to report no error for a missing key (the source of 200 {}), got %v", err)
			}
		})
	}
}

// TestAbsolutePathCollection reproduces the A3 bug: a collection configured
// with an absolute dataset: path routes by its base name, but GetCfg compares
// the whole string, so a named query answers 400 and create_success is silently
// ignored. These fail until Phase 1 step 1.3.
func TestAbsolutePathCollection(t *testing.T) {
	h := newStatusHarness(t, func(dir string, s *Settings) {
		s.Collections[0].CName = filepath.Join(dir, harnessOpenCollection)
		s.Collections[0].CreateSuccess = "/thanks"
		s.Collections[0].CreateError = "/oops"
	})
	o := "/api/" + harnessOpenCollection
	t.Run("named query works", func(t *testing.T) {
		got, _, body := h.do("GET", o+"/query/all", "", "")
		if got != 200 {
			t.Errorf("A3: want 200, got %d, body %q", got, body)
		}
	})
	t.Run("create_success redirect is honoured", func(t *testing.T) {
		got, hdr, body := h.do("POST", o+"/object/formrec", "application/x-www-form-urlencoded", "id=formrec&v=1")
		if got != 302 || !strings.HasPrefix(hdr.Get("Location"), "/thanks") {
			t.Errorf("A3: want 302 to /thanks, got %d Location %q body %q", got, hdr.Get("Location"), body)
		}
	})
	t.Run("routes that do not use the config work either way", func(t *testing.T) {
		if got, _, _ := h.do("GET", o+"/object/k1", "", ""); got != 200 {
			t.Errorf("read: want 200, got %d", got)
		}
	})
}

// TestQueryFailureReasons is tier A1: the four query-routing failures keep 400
// but say which one happened, in the body. Unknown collection and no queries
// defined cannot be reached by a routed request (the route is not registered,
// so the router answers first), so Query is called directly for those.
func TestQueryFailureReasons(t *testing.T) {
	h := newStatusHarness(t, nil)
	o := "/api/" + harnessOpenCollection
	direct := func(cName string, options []string) (int, string) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/"+cName+"/query", nil)
		Query(w, r, h.api, cName, "query", options)
		return w.Code, w.Body.String()
	}
	routed := func(path string) (int, string) {
		code, _, body := h.do("GET", path, "", "")
		return code, body
	}
	cases := []struct {
		name   string
		get    func() (int, string)
		reason string
	}{
		{"missing query name", func() (int, string) { return routed(o + "/query") }, "missing query name"},
		{"undefined query name", func() (int, string) { return routed(o + "/query/nope") }, `"nope" is not defined`},
		{"unknown collection", func() (int, string) { return direct("unknown.ds", []string{"all"}) }, "unknown collection"},
		{"no queries defined", func() (int, string) { return direct(harnessReadOnlyCollection, []string{"all"}) }, "no queries defined"},
	}
	seen := map[string]string{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := tc.get()
			if status != 400 {
				t.Errorf("status stays 400, got %d", status)
			}
			if !strings.Contains(body, tc.reason) {
				t.Errorf("body should say %q, got %q", tc.reason, body)
			}
			if other, dup := seen[body]; dup {
				t.Errorf("body %q is identical to the one for %q", body, other)
			}
			seen[body] = tc.name
		})
	}
}

// TestQueryBadRequestRedirects pins a bug found in Phase 1: Query passed an
// error string as statusIsError's errorRedirect, so a POST with an unsupported
// content type answered 303 to a Location that was the error text. It is a 400
// with the reason in the body.
func TestQueryBadRequestRedirects(t *testing.T) {
	h := newStatusHarness(t, nil)
	got, hdr, body := h.do("POST", "/api/"+harnessOpenCollection+"/query/all", "text/plain", "x")
	if got != 400 {
		t.Errorf("want 400, got %d Location %q body %q", got, hdr.Get("Location"), body)
	}
	if hdr.Get("Location") != "" {
		t.Errorf("an error must not redirect, Location %q", hdr.Get("Location"))
	}
	if !strings.Contains(body, "unsupported content type") {
		t.Errorf("body should say why, got %q", body)
	}
}

// TestAllowHeader is tier A2: a 405 lists the methods the route does support.
func TestAllowHeader(t *testing.T) {
	h := newStatusHarness(t, nil)
	cases := []struct{ method, path, want string }{
		{"PATCH", "/api/" + harnessOpenCollection + "/object/k1", "DELETE, GET, POST, PUT"},
		{"POST", "/api/" + harnessReadOnlyCollection + "/object/k1", "GET"},
		{"DELETE", "/api/" + harnessOpenCollection + "/keys", "GET"},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			got, hdr, _ := h.do(tc.method, tc.path, "", "")
			if got != 405 {
				t.Fatalf("want 405, got %d", got)
			}
			if hdr.Get("Allow") != tc.want {
				t.Errorf("Allow: want %q, got %q", tc.want, hdr.Get("Allow"))
			}
		})
	}
}

// TestDuplicateBaseNameRefused is tier A3: two collections whose paths share a
// base name are routed by the same name, so the second would silently replace
// the first. Init refuses them.
func TestDuplicateBaseNameRefused(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	var cfgs []*Config
	for _, sub := range []string{"a", "b"} {
		p := filepath.Join(dir, sub, "same.ds")
		if err := os.MkdirAll(filepath.Dir(p), 0775); err != nil {
			t.Fatal(err)
		}
		c, err := Init(p, "sqlite://collection.db")
		if err != nil {
			t.Fatal(err)
		}
		c.Close()
		cfgs = append(cfgs, &Config{CName: p, DsnURI: "sqlite://collection.db", Keys: true})
	}
	s := &Settings{Host: "localhost:0", Collections: cfgs}
	if err := s.WriteFile("settings.yaml", 0664); err != nil {
		t.Fatal(err)
	}
	api := new(API)
	err := api.Init("datasetd", "settings.yaml")
	if err == nil {
		t.Fatal("two collections named same.ds must be refused")
	}
	if !strings.Contains(err.Error(), "same.ds") {
		t.Errorf("the message should name the collection, got %q", err)
	}
}

// TestPutFormWithoutModel: a form (non-JSON) PUT to a collection that has no
// model used to panic in Update (c.Model.GetPrimaryId() ahead of the nil
// check), so the client saw a dropped connection instead of an answer.
func TestPutFormWithoutModel(t *testing.T) {
	h := newStatusHarness(t, nil)
	got, _, body := h.do("PUT", "/api/"+harnessOpenCollection+"/object/k1", "application/x-www-form-urlencoded", "id=k1&v=viaform")
	if got != 200 {
		t.Errorf("want 200, got %d %q", got, body)
	}
}
