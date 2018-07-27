
# export-gsheet

## Syntax

```
    dataset COLLECTION_NAME export-gsheet SHEET_ID SHEET_NAME CELL_RANGE FILTER_EXPR FIELDS_TO_EXPORT [COLUMN_NAMES]
```

## Description

export-gsheet will write the exported records and exported fields to 
a Google Sheets sheet in the given cell range.

SHEET_ID is the google sheet id, usually a very long alpha numeric 
string. If your URL looks like

```
    https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
```

The Sheet ID would be

```
    1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms
```

CELL_RANGE is expressed in a letter colum row fashion E.g. A1:B10 
would refer to the grid starting in the first column (i.e. A) and row 
(i.e. 1) moving to the lower left with the grid ending in column 2 
(i.e. B) row 10. This forms the grid that will be written out. 
Typically this would be something like "A1:Z" which would translate 
start in the upper right cell of the spreadsheet and replace all cells
to column Z going down.
 
FILTER_EXPR is an expression that evaluates to _true_ or _false_ 
based on Golang template expressions (see `dataset -help filter` 
for more explanation).

FIELDS_TO_EXPORT is a comma separated list of dotpaths (e.g. 
.id,.title,.pubDate) in the JSON documents in the collection 
(see `dataset -help dotpath` for more explanation)

## Usage

In the following examples we will "filter" for all records in a 
collection so we use the string "true".  The following fields are 
being exported - .name and .contact with the following 
headings -- Name, Contact. Collection name is "people.ds".

```shell
	dataset people.ds export-gsheet "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" Sheet1 "A1:Z" true '.name,.contact' 'Name,Contact'
```

Related topics: [import-csv](import-csv.html), [export-csv](export-csv.html), [import-gsheet](import-gsheet.html)

