
# export

## Syntax

```
    dataset export COLLECTION_NAME FRAME_NAME [CSV_FILENAME]
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
    dataset frame publications.ds my-report "._Key" ".title" ".pubDate"
    dataset frame-labels publications.ds my-report "id" "title" "publication date"
```

The example blow creates a CSV file named 'output.csv'. The collection 
is "publications.ds".

```shell
	dataset export publications.ds my-report > output.csv
```

Related topics: [import-csv](import-csv.html), [import-gsheet](import-gsheet.html), [export-gsheet](export-gsheet.html)

