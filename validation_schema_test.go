// validation_schema_test.go — integration tests for schema validation and
// generator behaviour in the datasetd HTTP API.
//
// Run with:
//   go test -v -run TestSchema ./...
//   go test -v -run TestGenerator ./...
//
// Known bugs exercised by these tests:
//   - TestSchemaValidation_PUTJsonSkipsValidation: JSON PUT does not call
//     ValidateRecord; invalid data is accepted (should return 400, returns 200).
//   - TestGenerator_NowIsUnrecognized: generator:"now" is not handled in
//     api_routes.go; the field is left unpopulated (should behave like "timestamp").
package dataset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

const (
	validationTestPort = "8587"
	validationTestHost = "localhost:8587"
)

// validationTestSettingsYAML is a minimal datasetd config with the
// stewardship schema (matching content_dashboard.yaml) and validate: true.
const validationTestSettingsYAML = `
host: "localhost:8587"
htdocs: "testout/htdocs"
schemas:
  stewardship_model:
    id: stewardship_model
    description: "Test schema for stewardship validation"
    elements:
      - id: pageId
        type: text
        attributes:
          name: pageId
          required: "true"
        pattern: "^[0-9]+$"
        is_primary_id: true
        label: "Page ID"
      - id: expert
        type: text
        attributes:
          name: expert
        label: "Expert"
      - id: lastUpdated
        type: datetime-local
        attributes:
          name: lastUpdated
          required: "true"
        generator: "timestamp"
        label: "Last Updated"
      - id: updatedBy
        type: text
        attributes:
          name: updatedBy
        label: "Updated By"
  stewardship_now_model:
    id: stewardship_now_model
    description: "Schema with generator:now to test unrecognized generator"
    elements:
      - id: pageId
        type: text
        attributes:
          name: pageId
          required: "true"
        pattern: "^[0-9]+$"
        is_primary_id: true
        label: "Page ID"
      - id: lastUpdated
        type: datetime-local
        attributes:
          name: lastUpdated
          required: "true"
        generator: "now"
        label: "Last Updated"
collections:
  - dataset: "testout/stewardship_validate.ds"
    schema: stewardship_model
    validate: true
    keys: true
    read: true
    create: true
    update: true
    delete: true
  - dataset: "testout/stewardship_now.ds"
    schema: stewardship_now_model
    validate: true
    keys: true
    read: true
    create: true
    update: true
    delete: true
`

// setupValidationTest initialises the test collections and writes the settings
// YAML. Returns the settings file path and a cleanup function.
func setupValidationTest(t *testing.T) (string, func()) {
	t.Helper()
	dir := "testout"
	if err := os.MkdirAll(path.Join(dir, "htdocs"), 0775); err != nil {
		t.Fatalf("MkdirAll htdocs: %s", err)
	}

	// Remove and re-initialise test collections so each run starts clean.
	for _, cName := range []string{
		path.Join(dir, "stewardship_validate.ds"),
		path.Join(dir, "stewardship_now.ds"),
	} {
		os.RemoveAll(cName)
		c, err := Init(cName, "sqlite://collection.db")
		if err != nil {
			t.Fatalf("Init(%q): %s", cName, err)
		}
		c.Close()
	}

	fName := path.Join(dir, "validation_settings.yaml")
	if err := os.WriteFile(fName, []byte(validationTestSettingsYAML), 0664); err != nil {
		t.Fatalf("WriteFile(%q): %s", fName, err)
	}

	return fName, func() { os.Remove(fName) }
}

// startValidationServer launches datasetd in a goroutine and waits for it to
// be ready. Callers must ensure the server is not already running on the port.
func startValidationServer(t *testing.T, fName string) {
	t.Helper()
	go func() {
		if err := RunAPI(os.Args[0], fName, false); err != nil {
			// RunAPI blocks; errors here are fatal for the goroutine only.
			fmt.Fprintf(os.Stderr, "validation test server error: %s\n", err)
		}
	}()
	// Wait for the server to be accepting connections.
	deadline := time.Now().Add(8 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := http.Get(fmt.Sprintf("http://%s/api/version", validationTestHost))
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("validation test server did not start in time")
}

func doJSON(t *testing.T, method, url string, body map[string]interface{}) *http.Response {
	t.Helper()
	src, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(src))
	if err != nil {
		t.Fatalf("NewRequest %s %s: %s", method, url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Do %s %s: %s", method, url, err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll: %s", err)
	}
	return strings.TrimSpace(string(b))
}

// ── Schema validation tests ────────────────────────────────────────────────

// TestSchemaValidation_Suite runs all schema and generator sub-tests against a
// single server instance to avoid port conflicts.
func TestSchemaValidation_Suite(t *testing.T) {
	fName, cleanup := setupValidationTest(t)
	defer cleanup()
	startValidationServer(t, fName)

	t.Run("POST_valid_data_returns_201", testPOSTValidReturns201)
	t.Run("POST_bad_pattern_returns_400", testPOSTBadPatternReturns400)
	t.Run("POST_missing_required_returns_400", testPOSTMissingRequiredReturns400)
	t.Run("PUT_valid_json_returns_200", testPUTValidReturns200)
	t.Run("PUT_invalid_json_returns_400", testPUTInvalidJSONReturns400KnownBug)
	t.Run("generator_timestamp_set_on_create", testGeneratorTimestampOnCreate)
	t.Run("generator_timestamp_updated_on_put", testGeneratorTimestampUpdatedOnPUT)
	t.Run("generator_now_set_on_create", testGeneratorNowUnrecognized)
}

const validateBase = "http://" + validationTestHost + "/api/stewardship_validate.ds"
const nowBase = "http://" + validationTestHost + "/api/stewardship_now.ds"

// POST with a valid record should return 201 Created.
func testPOSTValidReturns201(t *testing.T) {
	record := map[string]interface{}{
		"pageId":    "9001",
		"expert":    "Jane Smith",
		"updatedBy": "test",
	}
	resp := doJSON(t, http.MethodPost, validateBase+"/object/9001", record)
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("POST valid record: expected 201, got %d\nbody: %s", resp.StatusCode, body)
	}
}

// POST with a pageId that does not match the pattern ^[0-9]+$ should return 400.
func testPOSTBadPatternReturns400(t *testing.T) {
	record := map[string]interface{}{
		"pageId":    "not-a-number",
		"expert":    "Jane Smith",
		"updatedBy": "test",
	}
	resp := doJSON(t, http.MethodPost, validateBase+"/object/not-a-number", record)
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("POST bad pattern: expected 400, got %d\nbody: %s\nX-Validation-Errors: %s",
			resp.StatusCode, body, resp.Header.Get("X-Validation-Errors"))
	}
}

// POST with required field pageId missing should return 400.
func testPOSTMissingRequiredReturns400(t *testing.T) {
	record := map[string]interface{}{
		"expert":    "Jane Smith",
		"updatedBy": "test",
	}
	// Key in URL is unrelated to the model's primary id validation.
	resp := doJSON(t, http.MethodPost, validateBase+"/object/9002", record)
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("POST missing required: expected 400, got %d\nbody: %s\nX-Validation-Errors: %s",
			resp.StatusCode, body, resp.Header.Get("X-Validation-Errors"))
	}
}

// PUT with a valid JSON body should return 200 OK.
func testPUTValidReturns200(t *testing.T) {
	// Ensure the record exists first.
	doJSON(t, http.MethodPost, validateBase+"/object/9010", map[string]interface{}{
		"pageId": "9010", "expert": "Alice", "updatedBy": "setup",
	})

	record := map[string]interface{}{
		"pageId":    "9010",
		"expert":    "Alice Updated",
		"updatedBy": "test",
	}
	resp := doJSON(t, http.MethodPut, validateBase+"/object/9010", record)
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("PUT valid: expected 200, got %d\nbody: %s", resp.StatusCode, body)
	}
}

// PUT with invalid JSON (bad pattern) should return 400.
func testPUTInvalidJSONReturns400KnownBug(t *testing.T) {
	// Ensure the record exists first.
	doJSON(t, http.MethodPost, validateBase+"/object/9020", map[string]interface{}{
		"pageId": "9020", "expert": "Bob", "updatedBy": "setup",
	})

	record := map[string]interface{}{
		"pageId":    "not-a-number", // violates pattern ^[0-9]+$
		"expert":    "Bob",
		"updatedBy": "test",
	}
	resp := doJSON(t, http.MethodPut, validateBase+"/object/9020", record)
	body := readBody(t, resp)
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("PUT invalid json (KNOWN BUG — json PUT skips ValidateRecord): expected 400, got %d\nbody: %s",
			resp.StatusCode, body)
	}
}

// ── Generator tests ────────────────────────────────────────────────────────

// generator:"timestamp" should auto-populate lastUpdated on POST.
func testGeneratorTimestampOnCreate(t *testing.T) {
	record := map[string]interface{}{
		"pageId":    "8001",
		"expert":    "Carol",
		"updatedBy": "test",
		// lastUpdated intentionally omitted — generator should supply it.
	}
	resp := doJSON(t, http.MethodPost, validateBase+"/object/8001", record)
	readBody(t, resp)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("generator POST: setup failed, got %d", resp.StatusCode)
	}

	// Read the record back and confirm lastUpdated was populated.
	getResp, err := http.Get(validateBase + "/object/8001")
	if err != nil {
		t.Fatalf("GET after create: %s", err)
	}
	body := readBody(t, getResp)
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal GET response: %s\nbody: %s", err, body)
	}
	val, ok := got["lastUpdated"]
	if !ok || val == "" || val == nil {
		t.Errorf("generator:timestamp on create: lastUpdated not populated; got record: %s", body)
	}
}

// generator:"timestamp" should update lastUpdated on every PUT.
func testGeneratorTimestampUpdatedOnPUT(t *testing.T) {
	// Create with a placeholder timestamp.
	doJSON(t, http.MethodPost, validateBase+"/object/8002", map[string]interface{}{
		"pageId": "8002", "expert": "Dave", "updatedBy": "setup",
	})

	// Read the initial value.
	getResp1, _ := http.Get(validateBase + "/object/8002")
	body1 := readBody(t, getResp1)
	var rec1 map[string]interface{}
	json.Unmarshal([]byte(body1), &rec1)
	firstTS, _ := rec1["lastUpdated"].(string)

	// Small sleep so the timestamps differ.
	time.Sleep(1100 * time.Millisecond)

	// PUT an update.
	doJSON(t, http.MethodPut, validateBase+"/object/8002", map[string]interface{}{
		"pageId": "8002", "expert": "Dave Updated", "updatedBy": "test",
	})

	// Read again.
	getResp2, _ := http.Get(validateBase + "/object/8002")
	body2 := readBody(t, getResp2)
	var rec2 map[string]interface{}
	json.Unmarshal([]byte(body2), &rec2)
	secondTS, _ := rec2["lastUpdated"].(string)

	if firstTS == "" {
		t.Errorf("generator:timestamp: lastUpdated not set on create; record: %s", body1)
	}
	if secondTS == "" {
		t.Errorf("generator:timestamp: lastUpdated not set on update; record: %s", body2)
	}
	if firstTS != "" && secondTS != "" && firstTS == secondTS {
		t.Errorf("generator:timestamp: lastUpdated not updated on PUT (stayed %q)", firstTS)
	}
}

// generator:"now" is an alias for "timestamp" — field must be auto-populated on create.
func testGeneratorNowUnrecognized(t *testing.T) {
	record := map[string]interface{}{
		"pageId": "7001",
		// lastUpdated omitted — generator:"now" should supply it.
	}
	resp := doJSON(t, http.MethodPost, nowBase+"/object/7001", record)
	readBody(t, resp)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("generator:now POST: setup failed, got %d", resp.StatusCode)
	}

	getResp, err := http.Get(nowBase + "/object/7001")
	if err != nil {
		t.Fatalf("GET: %s", err)
	}
	body := readBody(t, getResp)
	var got map[string]interface{}
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %s\nbody: %s", err, body)
	}
	val, ok := got["lastUpdated"]
	if !ok || val == "" || val == nil {
		t.Errorf("generator:now: lastUpdated not populated; record: %s", body)
	}
}
