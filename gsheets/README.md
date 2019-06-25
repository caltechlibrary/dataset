
# Setup Access to run GSheets tests

Reference Materials: 

+ [Overview of sheets API](https://developers.google.com/sheets/api/quickstart/go)
+ [API Console](https://console.developers.google.com/apis/dashboard) -- where you go to manage the credentials
+ ../cmd/gsheetssetup/gsheetsetup.go is a Go program to remind me how to set this up for testing!


This is what I do to setup access and remember the details how things
get authorized.

```shell
    go run cmd/gsheetaccess/gsheetaccess.go
```

Follow the prompts and suggestions until I get back the contents
of my "Sheet1".

Environment variables for testing

SPREADSHEET_ID 
: The Google Spreadsheet ID taken from the URL (long alpha numeric string)

SHEET_NAME
: Name the the sheet in the spread to test, e.g. "Sheet1"

CLIENT_SECRET
: Location of the "credentials.json" file.

