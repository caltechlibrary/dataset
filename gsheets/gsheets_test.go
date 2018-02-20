package gsheets

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestReadSheet(t *testing.T) {
	// The following sheet id is taken from https://developers.google.com/sheets/api/quickstart/go
	clientSecretJSON := os.Getenv("GOOGLE_CLIENT_SECRET_JSON") //"../etc/client_secret.json"
	if len(clientSecretJSON) == 0 {
		fmt.Fprintf(os.Stderr, "Skipping TestReadSheet, GOOGLE_CLIENT_SECRET_JOSN environment variable not set\n")
		return
	}
	clientSecretJSON = path.Join("..", clientSecretJSON)
	spreadSheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	sheetName := "Class Data"
	cellRange := "A2:E"
	table, err := ReadSheet(clientSecretJSON, spreadSheetId, sheetName, cellRange)
	if err != nil {
		t.Errorf("Expect success, got %s", err)
		t.FailNow()
	}
	src, _ := json.Marshal(table)
	e := 30
	r := len(table)
	if e != r {
		t.Errorf("expected %d rows, got %d rows\n\t\t%s", e, r, src)
	}
	for i, row := range table {
		e = 4 // or five in some cases
		r = len(row)
		src, _ = json.Marshal(row)
		if e != r && (e+1) != r {
			t.Errorf("expected %d cols, got %d cols in row %d\n\t\t%s", e, r, i, src)
		}
	}
}
