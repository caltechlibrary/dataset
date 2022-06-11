package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
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

func clientTestVersion(t *testing.T, cfg *Config) {
	fmt.Printf("starting client test verisons\n")
	// Run through a set of the
	u := fmt.Sprintf("http://%s/api/version", cfg.Host)
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

func clientTestKeys(t *testing.T, cfg *Config) {
	fmt.Printf("starting client test verisons\n")
	// Run through a set of the
	u := fmt.Sprintf("http://%s/api/keys", cfg.Host)
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
}

func TestRunAPI(t *testing.T) {
	// Setup up a test collection
	cName := path.Join("testout", "apitest.ds")
	dsnURI := "sqlite://testout/apitest.ds/collection.db"
	if err := SetupTestCollection(cName, dsnURI, records); err != nil {
		t.Errorf("failed to setup %q %q, %s", cName, dsnURI, err)
		t.FailNow()
	}

	// Setup a test configuration
	fName := path.Join("testout", "settings.json")
	cfg := new(Config)
	cfg.Host = "localhost:8001"
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
	if err := cfg.SaveConfig(fName); err != nil {
		t.Errorf("failded to save config %q, %s", fName, err)
		t.FailNow()
	}

	setupWait := "3s"
	wait, _ := time.ParseDuration(setupWait)
	fmt.Printf(`
Launching API Tests at http://%q

Press Ctr-C to quit on hang

Client testings starts in %s (s = seconds, m = minutes)
`, cfg.Host, setupWait)
	appName := os.Args[0]
	// Send our web service into the back ground so we can run
	// a client test.
	go func() {
		if err := RunAPI(appName, fName); err != nil {
			t.Errorf("Expected API to setup and run, %s", err)
		}
	}()
	time.Sleep(wait)
	clientTestVersion(t, cfg)
	clientTestKeys(t, cfg)
}
