//
// gsheets.go is a part of the dataset package written to allow import/export of records
// to/from dataset collections.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
package gsheets

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	clientSecretFName string
	spreadsheetID     string
)

func TestReadSheet(t *testing.T) {
	// The following sheet id is taken from https://developers.google.com/sheets/api/quickstart/go
	if len(clientSecretFName) == 0 {
		fmt.Fprintf(os.Stderr, "Skipping TestReadSheet, CLENT_SECRET_JSON filename or SPREADSHEET_ID not provided\n")
		fmt.Fprintln(os.Stderr, "USAGE: go test -client-secret CLIENT_SECRET_JSON -spreadsheet-id SPREADSHEET_ID")
		return
	}
	sheetName := "Staff Data"
	cellRange := "A2:B"
	table, err := ReadSheet(clientSecretFName, spreadsheetID, sheetName, cellRange)
	if err != nil {
		t.Errorf("Expect success, got %s", err)
		t.FailNow()
	}
	src, _ := json.Marshal(table)
	e := 4 // expected four rows of data
	r := len(table)
	if e != r {
		t.Errorf("expected %d rows, got %d rows\n\t\t%s", e, r, src)
	}
	for i, row := range table {
		e = 2 // two cells per row
		r = len(row)
		src, _ = json.Marshal(row)
		if e != r && (e+1) != r {
			t.Errorf("expected %d cols, got %d cols in row %d\n\t\t%s", e, r, i, src)
		}
	}
}

func TestMain(m *testing.M) {
	flag.StringVar(&clientSecretFName, "client-secret", "", "Set path/filename for client_secret.json")
	flag.StringVar(&spreadsheetID, "spreadsheet-id", "", "Set spreadsheet id to use for testing")
	flag.Parse()
	os.Exit(m.Run())
}
