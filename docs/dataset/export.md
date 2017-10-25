
# export

## Syntax

```
    dataset export OUTPUT_NAME FILTER_EXPR FIELDS_TO_EXPORT COLUMN_HEADINGS
```

## Description

export will render the contents of a collection as a tabular file. Supported tabular formats include
CSV and xlsx. Format is determined by file suffix (e.g. .csv for CSV format, .xlsx for
workbook format).

FILTER_EXPR is an expression that evaluates to _true_ or _false_ based on Golang template expressions
(see `dataset -help filter` for more explanation).

FIELDS_TO_EXPORT is a comma separated list of dotpaths (e.g. .id,.title,.pubDate) in the JSON documents
in the collection (see `dataset -help dotpath` for more explanation)

Note for workbooks you can set the sheet name with the option "-sheet".

## Usage

In the following examples we will "filter" for all records in a collection so we use the string "true". 
The following fields are being exported - ._id,.title, and .pubDate with the following headings --
id, title and publication date. 

The example blow creates a CSV file then creates a Workbook with a sheet named "Title List".


```shell
	dataset export titles.csv true '._id,.title,.pubDate' 'id,title,publication date'
	dataset -sheet 'Title List' export titles.xlsx true '._id,.title,.pubDate' 'id,title,publication date'
```

Related topics: extract, import, import-gsheets, export-gsheets

