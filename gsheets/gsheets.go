//
// gsheets.go is a part of the dataset package written to allow import/export of records
// to/from dataset collections.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
package gsheets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	// Google Sheets packages
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-dataset.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// ColNoToColLetters converts a zero based column index to a spreadsheet
// style letter sequence (e.g. 0 -> A, 26 -> AB, 52 -> BA, ..)
// If colNo is legative then an empty string is returned.
func ColNoToColLetters(colNo int) string {
	alpha := []string{
		"A", "B", "C", "D", "E",
		"F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y",
		"Z",
	}
	c := len(alpha)
	out := []string{}
	i := colNo
	m := 0
	for i >= 0 {
		if i < c {
			out = append([]string{alpha[i]}, out...)
			break
		}
		m = i % c
		i = (i - c) / c
		out = append([]string{alpha[m]}, out...)
	}
	return strings.Join(out, "")
}

func ReadSheet(clientSecretJSON, spreadSheetId, sheetName, cellRange string) ([][]interface{}, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecretJSON)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %s", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-dataset.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %s", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets Client %s", err)
	}

	// Prints the columns from sheet described by spreadSheetId
	readRange := fmt.Sprintf("%s!%s", sheetName, cellRange)
	//resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, readRange).ValueRenderOption("FORMULA").Do()
	resp, err := srv.Spreadsheets.Values.Get(spreadSheetId, readRange).ValueRenderOption("UNFORMATTED_VALUE").Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet. %s", err)
	}

	if len(resp.Values) > 0 {
		return resp.Values, nil
	}
	return nil, fmt.Errorf("No data found")
}

func WriteSheet(clientSecretJSON, spreadSheetId, sheetName, cellRange string, table [][]interface{}) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecretJSON)
	if err != nil {
		return fmt.Errorf("Unable to read client secret file: %s", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-dataset.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return fmt.Errorf("Unable to parse client secret file to config: %s", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets Client %s", err)
	}

	// Prints the columns from sheet described by spreadSheetId
	var (
		vr sheets.ValueRange
	)
	for _, row := range table {
		vr.Values = append(vr.Values, row)
	}
	writeRange := fmt.Sprintf("%s!%s", sheetName, cellRange)
	if _, err := srv.Spreadsheets.Values.Update(spreadSheetId, writeRange, &vr).ValueInputOption("USER_ENTERED").Do(); err != nil {
		return fmt.Errorf("Unable to write sheet %s. %s", writeRange, err)
	}
	return nil
}
