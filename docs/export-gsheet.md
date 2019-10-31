
# export

## Syntax

```
    dataset export COLLECTION_NAME FRAME_NAME GSHEET_ID GSHEET_NAME [CELL_RANGE]
```

## Description

_export_ will render the contents of a collection as a CSV file
or export to a Google Sheet based on a frame defined in the 
collection. 

## Usage

In the following examples we will be using a newly defined
"frame" named "my-report".  The frame will have the following fields are 
being exported - ._Key,.title, and .pubDate with the following 
labels for those fields -- id, title and publication date. 

```shell
    dataset frame publications.ds my-report \
        "._Key=id" ".title=title" \
        ".pubDate=publication date"
```

The example blow creates a CSV file named 'output.csv'. The collection 
is "publications.ds".

```shell
	dataset export publications.ds my-report > output.csv
```

Likewise we can export to a Google Sheet.  SHEET_ID is the google 
sheet id, usually a very long alpha numeric string. If your URL 
looks like

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
 
```shell
	dataset export publications.ds my-report "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" Sheet1 "A1:Z" 
```

Related topics: [frame](frame.html), [import-csv](import-csv.html), [export-csv](export-csv.html), [import-gsheet](import-gsheet.html)

