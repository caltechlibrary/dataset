
# import

## Syntax

```
    dataset import COLLECTION_NAME SHEET_ID SHEET_NAME ID_COL_NO [CELL_RANGE]
```

+ COLLECTION_NAME is the collection we are going to import into
+ SHEET_ID is the hash id Google assignes, it looks like a long string with numbers and letters in 
  the URL when you edit your sheet
+ SHEET_NAME is a string name of the sheet. The default name is usually "Sheet1" it is seen at the 
  lower part of the spreadsheet page in Google Sheets edit view
+ CELL_RANGE is a range of cells to import, typically this is "A1:Z" but maybe adjusted (e.g. if you
  want to skip the first row then you might use "A2:Z")
+ ID_COL_NO is the column number to use for the unique ID name of the JSON document. It should be an integer starting with "1".

## Options

+ -overwrite=true Allows dataset to overwrite existing values in a collection

## Description

_dataset_ supports importing data from a single sheet at a time 
from a Google Sheets document. To do this you need to beable to 
authenticate with the Google Sheets v4 API and an account with the
permissions allowing it to read the Google Sheets document.
Google Sheets like Excel workbooks include multiple talbes in a 
single document. This is usually called a _sheet_. When importing
a Google Sheet into a _dataset_ collection the collection needs to 
exist and you need to identity the source of the key. If none is 
provided the key will be created as the row number of each JSON 
document constructed from the column header and cell value. This 
is problematic if someome sorts the sheet differently and then 
re-imports the data into the collection.  So usually you want to 
explicitly set the column that will be used as as the record key in 
the collection. That way you can re-import the sheet's data into 
your collection and replacing the stale data.


## Example

In this example we're using the example Google Sheet from the 
Golang Google Sheets API v4 Quickstart. You'll first need to have 
created a *credentials.json* file as described in the Step 1 of the 
[Google Cloud SDK docs](https://developers.google.com/sheets/api/quickstart/go)
and placed it in *etc/credentials.json*.  Our collection name 
is "DemoStudentList.ds".

```shell
    export GOOGLE_CLIENT_SECRET_JSON="etc/credentials.json"
    dataset DemoStudentList.ds init
    dataset import DemoStudentList.ds "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" 1 "A1:Z" 
    dataset keys DemoStudentList.ds | while read KY; do dataset read DemoStudentList.ds "${KY}"; done
```

Related topics: [dotpath](dotpath.html), [export-csv](export-csv.html), [import-csv](import-csv.html), and [export-gsheet](export-gsheet.html)

