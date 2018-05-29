
# Google Spreadsheet Integration

_dataset_ provides for importing from and export to a Google spreadsheet (i.e. GSheet). 
There is required setup for this to work.  _dataset_ needs to beable to access
the Google Sheets API for reading and writing. You can find documentation on setting
up access in "step 1" at https://developers.google.com/sheets/api/quickstart/go.

You'll need a "client_secret.json" file and OAuth authorization for access to be permitted.
If credentials for the OAuth part are usually stored in your `$HOME/.credentials` directory
as sheets.googleapis.com-dataset.json.  If this file doesn't exist then the first time you
run the _dataset_ command with a GSheet option it'll prompt you to use your web browser
to authorize _dataset_ to access your Google spreadsheet.


## import-gsheet

This is place holder for documentation on import Google Sheets into a _dataset_ collction.

### Syntax

```
    dataset COLLECTION_NAME import-gsheet SHEET_ID SHEET_NAME CELL_RANGE COL_NO_FOR_ID
```

+ COLLECTION_NAME is the collection we are going to import into
+ SHEET_ID is the hash id Google assignes, it looks like a long string with numbers and letters in 
  the URL when you edit your sheet
+ SHEET_NAME is a string name of the sheet. The default name is usually "Sheet1" it is seen at the 
  lower part of the spreadsheet page in Google Sheets edit view
+ CELL_RANGE is a range of cells to import, typically this is "A1:Z" but maybe adjusted (e.g. if you
  want to skip the first row then you might use "A2:Z")
+ COL_NO_FOR_ID is the column number to use for the unique ID name of the JSON document. It should
  be an integer starting with "1".


### Description

_dataset_ supports importing data from a single sheet at a time from a Google Sheets document. To
do this you need to beable to authenticate with the Google Sheets v4 API and an account with the
permissions allowing it to read the Google Sheets document.  Google Sheets like Excel workbooks
include multiple talbes in a single document. This is usually called a _sheet_. When importing
a Google Sheet into a _dataset_ collection the collection needs to exist and you need to identity
the source of the key. If none is provided the key will be created as the row number of each 
JSON document constructed from the column header and cell value. This is problematic if someome
sorts the sheet differently and then re-imports the data into the collection.  So usually you
want to explicitly set the column that will be used as as the record key in the collection. That
way you can re-import the sheet's data into your collection and replacing the stale data.

### Example

In this example we're using the example Google Sheet from the Golang Google Sheets API v4 
Quickstart. You'll first need to have created a *client_secret.json* file as described in
the Step 1 of the [Google Cloud SDK docs](https://developers.google.com/sheets/api/quickstart/go)
and placed it in *etc/client_secret.json*.  Our collection name is "DemoStudentList.ds".

```shell
    export GOOGLE_CLIENT_SECRET_JSON="etc/client_secret.json"
    dataset DemoStudentList.ds init
    dataset DemoStudentList.ds import-gsheet "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" "Class Data" "A1:Z" 1
    dataset DemoStudentList.ds keys | while read KY; do dataset DemoStudentList.ds read "${KY}"; done
```

In this example we've used the row number as the ID for the JSON document created. This isn't
ideal in production as someone may re-sort the spreadsheet thus changing the number relationship
between the row number and the document in your _dataset_ collection.

In this version we've not used the first row as field names in the JSON record. How does 
it look different? What does "-use-header-row=false" mean? Why is the range different?

```shell
    dataset -use-header-row=false DemoStudentList.ds import-gsheet "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" "Class Data" "A2:Z" 1
    dataset DemoStudentList.ds keys | while read KY; do dataset DemoStudentList.ds read "${KY}"; done
```


## export-gsheet

### Syntax

```
    dataset COLLECTION_NAME export-gsheet SHEET_ID SHEET_NAME CELL_RANGE FILTER_EXPR FIELDS_TO_EXPORT [COLUMN_NAMES]
```

### Description

export-gsheet will write the exported records and exported fields to a Google Sheets sheet in the given cell range.

SHEET_ID is the google sheet id, usually a very long alpha numeric string. If your URL looks like

```
    https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
```

The Sheet ID would be

```
    1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms
```

CELL_RANGE is expressed in a letter colum row fashion E.g. A1:B10 would refer to 
the grid starting in the first column (i.e. A) and row (i.a. 1) moving to the lower left with the grid
ending in column 2 (i.e. B) row 10. This forms the grid that will be written out. Typically this would be something
like "A1:Z" which would translate start in the upper right cell of the spreadsheet and replace all cells
to column Z going down.
 
FILTER_EXPR is an expression that evaluates to _true_ or _false_ based on Golang template expressions
(see `dataset -help filter` for more explanation).

FIELDS_TO_EXPORT is a comma separated list of dotpaths (e.g. .id,.title,.pubDate) in the JSON documents
in the collection (see `dataset -help dotpath` for more explanation)

### Usage

In the following examples we will "filter" for all records in a collection so we use the string "true". 
The following fields are being exported - .name and .contact with the following headings --
Name, Contact. Collection name is "people.ds".

```shell
	dataset people.ds export-gsheet "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" Sheet1 "A1:Z" true '.name,.contact' 'Name,Contact'
```


Related topics: [dotpath](../docs/dotpath.html), [../docs/export-csv](export-csv.html), [import-csv](../docs/import-csv.html), [import-gsheet](../docs/import-gsheet.html) and [export-gsheet](../docs/export-gsheet.html)

