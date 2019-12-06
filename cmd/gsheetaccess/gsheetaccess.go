//
// This is based on Google demo code to access It's Google Sheets API.
// It is setup to help me remember how to authorize access to run the
// tests for dataset and GSheet integration. It is UGLY code!
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// Google APIs
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

//
// This is standard boiler place connection stuff from Google
//

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

//
// This is dataset specific stuff.
//

func setupSheetAccess(spreadsheetID string, sheetName string, sheetRange string) error {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints cells in the range
	if sheetRange == "" {
		sheetRange = "A1:E"
	}
	readRange := sheetName + "!" + sheetRange
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			for j, cell := range row {
				if j > 0 {
					fmt.Printf(",")
				}
				fmt.Printf("%q", cell)
			}
			fmt.Printf("\n")
		}
	}
	return nil
}

func main() {
	spreadsheetID := ""
	sheetName := "Sheet1"
	flag.StringVar(&spreadsheetID, "spreadsheet-id", "", "GSheet Spreadsheet ID")
	flag.StringVar(&sheetName, "sheet-name", "", "GSheet Name, e.g. Sheet1")

	flag.Parse()

	if spreadsheetID == "" {
		spreadsheetID = os.Getenv("SPREADSHEET_ID")
	}
	if sheetName == "" {
		sheetName = os.Getenv("SHEET_NAME")
		if sheetName == "" {
			sheetName = "Sheet1"
		}
	}

	if spreadsheetID == "" {
		fmt.Fprintf(os.Stderr, "Missing spreadsheet id\n")
		os.Exit(1)
	}
	if sheetName == "" {
		fmt.Fprintf(os.Stderr, "Missing Sheet name, e.g. \"Sheet1\"\n")
		os.Exit(1)
	}

	_, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, `

You need a credentials.json file to authenticate. See documentation
at 

   https://developers.google.com/identity/protocols/OAuth2?hl=en_US

Go to GoogleAPI Credentials page 

	https://console.developers.google.com/apis/dashboard

Click "Credentials" menu item on left.
Pick or create credentials. Follow the instructions, they change allot.

Download the appropriate credentials file you've previously created.
Place the file in your current working directory 
and name it credentials.json
`)
	}

	if err := setupSheetAccess(spreadsheetID, sheetName, ""); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "OK\n")
}
