
# import-csv

## Syntax

```
    dataset COLLECTION_NAME import-csv CSV_FILENAME COLUMN_NO_AS_KEY
```

## Description

_import-csv_ adds JSON documents to a collection from a CSV table. 

## Usage

In the following examples the CSV filename is _data.csv_.
The first column (column 1) is used as the value for KEY if
specified.  Our collection is named "data.ds".

```shell
    dataset data.ds import-csv data.csv 1
```

By default the header row of the table (the first row of the table) 
is used as the attribute names of the JSON document you create on 
import.  If you don't want that behavior you can use 
the "-use-header-row=false" option and the fields will be in the
form of "column_IDNO" where IDNO is replaced with a left zero 
padded column number (e.g. column_001, column_002, column_003).


Related topics: [export-csv](export-csv.html), [import-gsheets](import-gsheets.html), [export-gsheets](export-gsheet.html)

