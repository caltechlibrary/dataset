package dataset

// api_strict_test.go is the strict_status contract (plan Phase 2): one case per
// row of the design table, with strict_status on. TestLegacyStatus pins the
// same requests with it off.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/caltechlibrary/models"
)

func strictOn(dir string, s *Settings) { s.StrictStatus = true }

// TestStrictStatus is the table of agents/projects/dataset/design/datasetd_http_status_v2_5_2.md.
func TestStrictStatus(t *testing.T) {
	h := newStatusHarness(t, strictOn)
	const j = "application/json"
	const form = "application/x-www-form-urlencoded"
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
		{"read missing", "GET", o + "/object/nope", "", "", 404},
		{"delete existing", "DELETE", o + "/object/k1", "", "", 200},
		{"delete missing", "DELETE", o + "/object/nope", "", "", 404},
		{"create new", "POST", o + "/object/new", j, `{"id":"new"}`, 201},
		{"json post to existing", "POST", o + "/object/k1", j, `{"id":"k1","v":"changed"}`, 409},
		{"json post to existing with charset", "POST", o + "/object/k1", j + "; charset=utf-8", `{"id":"k1","v":"changed"}`, 409},
		{"update existing", "PUT", o + "/object/k1", j, `{"id":"k1","v":"changed"}`, 200},
		{"put to missing", "PUT", o + "/object/nope", j, `{"id":"nope"}`, 404},
		{"post invalid JSON", "POST", o + "/object/new", j, `{not json`, 400},
		{"post empty body", "POST", o + "/object/new", j, ``, 400},
		{"put invalid JSON", "PUT", o + "/object/k1", j, `{not json`, 400},
		{"post text/plain", "POST", o + "/object/new", "text/plain", `hello`, 415},
		{"put text/plain", "PUT", o + "/object/k1", "text/plain", `hello`, 415},
		{"post without a content type", "POST", o + "/object/new", "", `{"id":"new"}`, 415},
		{"post over 1 MiB", "POST", o + "/object/big", j, over, 413},
		{"put over 1 MiB", "PUT", o + "/object/k1", j, over, 413},
		{"form post new", "POST", o + "/object/formnew", form, "id=formnew&v=1", 201},
		{"query ok", "GET", o + "/query/all", "", "", 200},
		{"query missing parameter", "GET", o + "/query/by_key/key", "", "", 400},
		{"query failing on user input", "GET", o + "/query/by_path/path?path=$[", "", "", 400},
		{"query name not defined", "GET", o + "/query/nope", "", "", 404},
		{"query without a name", "GET", o + "/query", "", "", 400},
		{"query post invalid JSON", "POST", o + "/query/by_key/key", j, `{not json`, 400},
		{"query post text/plain", "POST", o + "/query/all", "text/plain", "x", 415},
		{"query post over 1 MiB", "POST", o + "/query/by_key/key", j, over, 413},
		{"unknown collection", "GET", "/api/nope.ds/keys", "", "", 404},
		{"unknown route", "GET", o + "/nope", "", "", 404},
		{"disabled method", "POST", "/api/" + harnessReadOnlyCollection + "/object/new", j, `{"id":"new"}`, 405},
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

	t.Run("a 409 writes nothing", func(t *testing.T) {
		h.reset()
		h.do("POST", o+"/object/k1", j, `{"id":"k1","v":"changed"}`)
		obj := map[string]interface{}{}
		if err := h.api.CMap[harnessOpenCollection].Read("k1", obj); err != nil {
			t.Fatal(err)
		}
		if obj["v"] != "original" {
			t.Errorf("the record must be untouched, got %+v", obj)
		}
	})
	t.Run("a 404 on PUT creates nothing", func(t *testing.T) {
		h.reset()
		h.do("PUT", o+"/object/nope", j, `{"id":"nope"}`)
		if h.api.CMap[harnessOpenCollection].HasKey("nope") {
			t.Errorf("nope must not exist")
		}
	})
	t.Run("a form post to an existing key still updates it", func(t *testing.T) {
		h.reset()
		got, _, body := h.do("POST", o+"/object/k1", form, "id=k1&v=viaform")
		if got != 200 {
			t.Errorf("want 200, got %d %q", got, body)
		}
		obj := map[string]interface{}{}
		h.api.CMap[harnessOpenCollection].Read("k1", obj)
		if obj["v"] != "viaform" {
			t.Errorf("a form post is an upsert, got %+v", obj)
		}
	})
	t.Run("405 carries Allow", func(t *testing.T) {
		_, hdr, _ := h.do("POST", "/api/"+harnessReadOnlyCollection+"/object/new", j, `{}`)
		if hdr.Get("Allow") != "GET" {
			t.Errorf("Allow: got %q", hdr.Get("Allow"))
		}
	})
	t.Run("a successful read is JSON", func(t *testing.T) {
		_, hdr, _ := h.do("GET", o+"/object/k1", "", "")
		if !strings.HasPrefix(hdr.Get("Content-Type"), "application/json") {
			t.Errorf("Content-Type: got %q", hdr.Get("Content-Type"))
		}
	})
}

// TestStrictErrorBody checks the DR-0020 error body and Cache-Control.
func TestStrictErrorBody(t *testing.T) {
	h := newStatusHarness(t, strictOn)
	o := "/api/" + harnessOpenCollection
	cases := []struct {
		name        string
		method      string
		path        string
		contentType string
		body        string
		status      int
		outcome     string
	}{
		{"not found", "GET", o + "/object/nope", "", "", 404, "not_found"},
		{"conflict", "POST", o + "/object/k1", "application/json", `{"id":"k1"}`, 409, "conflict"},
		{"invalid", "POST", o + "/object/new", "application/json", `{no`, 400, "invalid"},
		{"unsupported media type", "POST", o + "/object/new", "text/plain", `x`, 415, "invalid"},
		{"too large", "POST", o + "/object/big", "application/json", `{"v":"` + strings.Repeat("x", jsonSizeLimit+1) + `"}`, 413, "too_large"},
		{"unknown route", "GET", "/api/nope.ds/keys", "", "", 404, "not_found"},
		{"method not allowed", "POST", "/api/" + harnessReadOnlyCollection + "/object/new", "application/json", `{}`, 405, "invalid"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h.reset()
			status, hdr, body := h.do(tc.method, tc.path, tc.contentType, tc.body)
			if status != tc.status {
				t.Fatalf("status: want %d, got %d, body %q", tc.status, status, body)
			}
			if !strings.HasPrefix(hdr.Get("Content-Type"), "application/json") {
				t.Errorf("Content-Type: got %q", hdr.Get("Content-Type"))
			}
			if hdr.Get("Cache-Control") != "no-store" {
				t.Errorf("Cache-Control: got %q", hdr.Get("Cache-Control"))
			}
			var e struct {
				Error struct {
					Status  int    `json:"status"`
					Outcome string `json:"outcome"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal([]byte(body), &e); err != nil {
				t.Fatalf("body is not JSON, %s: %q", err, body)
			}
			if e.Error.Status != tc.status || e.Error.Outcome != tc.outcome || e.Error.Message == "" {
				t.Errorf("error object: want status %d outcome %q and a message, got %+v", tc.status, tc.outcome, e.Error)
			}
		})
	}
}

// TestStrictQueryReason: the message for a 404 on a query says which one it is.
func TestStrictQueryReason(t *testing.T) {
	h := newStatusHarness(t, strictOn)
	_, _, body := h.do("GET", "/api/"+harnessOpenCollection+"/query/nope", "", "")
	if !strings.Contains(body, `is not defined`) {
		t.Errorf("the message should say why, got %q", body)
	}
}

// TestStrictFormPost: a form post that succeeds with create_success answers 303.
func TestStrictFormPost(t *testing.T) {
	h := newStatusHarness(t, func(dir string, s *Settings) {
		s.StrictStatus = true
		s.Collections[0].CreateSuccess = "/thanks"
	})
	o := "/api/" + harnessOpenCollection
	const form = "application/x-www-form-urlencoded"
	for _, tc := range []struct{ name, key string }{{"new record", "formnew"}, {"existing record", "k1"}} {
		t.Run(tc.name, func(t *testing.T) {
			got, hdr, body := h.do("POST", o+"/object/"+tc.key, form, "id="+tc.key+"&v=1")
			if got != http.StatusSeeOther || !strings.HasPrefix(hdr.Get("Location"), "/thanks") {
				t.Errorf("want 303 to /thanks, got %d Location %q body %q", got, hdr.Get("Location"), body)
			}
		})
	}
}

// TestStrictValidation: a record that fails validation is 422, with the
// X-Validation-Errors header kept; legacy keeps 400.
func TestStrictValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		strict bool
		want   int
	}{{"strict is 422", true, 422}, {"legacy is 400", false, 400}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newStatusHarness(t, func(dir string, s *Settings) {
				m := new(models.Model)
				if err := YAMLUnmarshal([]byte(`
id: needs_title
elements:
  - id: title
    type: text
    attributes:
      name: title
      required: "true"
    label: Title
`), m); err != nil {
					t.Fatal(err)
				}
				s.Collections[0].Model = m
				s.Collections[0].Validate = true
				s.StrictStatus = tc.strict
			})
			got, hdr, body := h.do("POST", "/api/"+harnessOpenCollection+"/object/v1", "application/json", `{"other":"x"}`)
			if got != tc.want {
				t.Fatalf("want %d, got %d, body %q", tc.want, got, body)
			}
			if hdr.Get("X-Validation-Errors") == "" {
				t.Errorf("X-Validation-Errors should be kept")
			}
			// A valid record is still 201.
			got, _, body = h.do("POST", "/api/"+harnessOpenCollection+"/object/v2", "application/json", `{"title":"ok"}`)
			if got != 201 {
				t.Errorf("a valid record: want 201, got %d %q", got, body)
			}
		})
	}
}

// TestStrictNoStore: every strict answer reflects the state of a collection,
// so successes carry Cache-Control: no-store too; legacy answers do not.
func TestStrictNoStore(t *testing.T) {
	o := "/api/" + harnessOpenCollection
	for _, strict := range []bool{true, false} {
		name := "legacy"
		if strict {
			name = "strict"
		}
		t.Run(name, func(t *testing.T) {
			h := newStatusHarness(t, func(dir string, s *Settings) { s.StrictStatus = strict })
			for _, path := range []string{o + "/object/k1", o + "/keys", o + "/query/all"} {
				_, hdr, _ := h.do("GET", path, "", "")
				got := hdr.Get("Cache-Control") == "no-store"
				if got != strict {
					t.Errorf("GET %s: no-store is %t, want %t", path, got, strict)
				}
			}
		})
	}
}

// TestStrictInternal: a store that fails is the server's fault, 500 in strict
// mode (legacy answered 400). The collection is closed under the service.
func TestStrictInternal(t *testing.T) {
	for _, tc := range []struct {
		name   string
		strict bool
		want   int
	}{{"strict is 500", true, 500}, {"legacy is 400", false, 400}} {
		t.Run(tc.name, func(t *testing.T) {
			h := newStatusHarness(t, func(dir string, s *Settings) { s.StrictStatus = tc.strict })
			h.api.CMap[harnessOpenCollection].Close()
			got, _, body := h.do("POST", "/api/"+harnessOpenCollection+"/object/new", "application/json", `{"id":"new"}`)
			if got != tc.want {
				t.Errorf("want %d, got %d, body %q", tc.want, got, body)
			}
			if tc.strict && strings.Contains(body, "closed") {
				t.Errorf("the body must not carry the store's error text: %q", body)
			}
		})
	}
}
