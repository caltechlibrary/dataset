
# import-gsheet

This is place holder for documentation on import Google Sheets into a _dataset_ collction.

## Syntax

```
    dataset import-gsheet SHEET_ID SHEET_NAME CELL_RANGE [COL_NO_FOR_ID]
```

+ SHEET_ID is the hash id Google assignes, it looks like a long string with numbers and letters in 
  the URL when you edit your sheet
+ SHEET_NAME is a string name of the sheet. The default name is usually "Sheet1" it is seen at the 
  lower part of the spreadsheet page in Google Sheets edit view
+ CELL_RANGE is a range of cells to import, typically this is "A1:Z" but maybe adjusted (e.g. if you
  want to skip the first row then you might use "A2:Z")
+ COL_NO_FOR_ID is the column number to use for the unique ID name of the JSON document. It should
  be an integer starting with "1".


## Description

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

## Example

In this example we're using the example Google Sheet from the Golang Google Sheets API v4 
Quickstart. We've assumed you've already create a *client_secret.json* file as described in
the Google Cloud SDK docs and placed in it *etc/client_secret.json*. 

```shell
    export GOOGLE_CLIENT_SECRET_JSON="etc/client_secret.json"
    dataset init DemoStudentList
    export DATASET="DemoStudentList"
    dataset import-gsheet "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" "Class Data" "A1:Z"
    dataset keys | while read KY; do dataset read "${KY}"; done
```

In this example we've used the row number as the ID for the JSON document created. This isn't
ideal in production as someone may re-sort the spreadsheet thus changing the number relationship
between the row number and the document in your _dataset_ collection.

In this version we've not used the first row as field names in the JSON record. How does 
it look different? What does "-use-header-row=false" mean? Why is the range different?

```shell
    dataset -use-header-row=false import-gsheet "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" "Class Data" "A2:Z"
    dataset keys | while read KY; do dataset read "${KY}"; done
```

