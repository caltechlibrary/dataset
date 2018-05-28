
# export-csv

## Syntax

```
    dataset COLLECTION_NAME export-csv OUTPUT_NAME FILTER_EXPR FIELDS_TO_EXPORT COLUMN_HEADINGS
```

## Description

_export-csv_ will render the contents of a collection as a CSV file. 

FILTER_EXPR is an expression that evaluates to _true_ or _false_ based on Golang template expressions
(see `dataset -help filter` for more explanation).

FIELDS_TO_EXPORT is a comma separated list of dotpaths (e.g. .id,.title,.pubDate) in the collection's JSON documents.
(see `dataset -help dotpath` for more explanation of dotpaths)

## Usage

In the following examples we will "filter" for all records in a collection so we use the string "true". 
The following fields are being exported - ._id,.title, and .pubDate with the following headings --
id, title and publication date. 

The example blow creates a CSV file named 'output.csv'. The
collection is "publications.ds".

```shell
	dataset publications.ds export-csv titles.csv true '._id,.title,.pubDate' 'id,title,publication date' > output.csv
```

Related topics: [import-csv](import-csv.html), [import-gsheets](import-gsheet.html), [export-gsheets](export-gsheet.html)

