
# import

## Syntax

```
    dataset import TABLE_FILENAME COLUMN_NO_AS_KEY
```

## Description

import adds JSON documents to a collection from a table. Tables can be either
CSV files, or Workbook sheets.

## Usage

In the following examples the CSV filename is _data.csv_, the Workbook filename is 
_data.xlsx_ The first column (column 1) is used as the value for KEY if
specified. For the Workbook example the option "-sheet" specifies the name of the
sheet to be imported. In our example the sheet name is "Title List".

Note if no ID column is specified then row number becomes the ID.  Import will replace 
any records with the same ID, if the "-update-only" option is used then
it'll only add records and skip importing rows that have an existing KEY
in the collection.

In the following examples the first one imports all the contents of _data.csv_ using the
row number as KEY. In the second one all the rows are import from _data.csv_ 
using the first column as the KEY (overwriting records with the same id). 
In the third version we're importing all the rows of _data.xlsx_ using column 1 as 
the KEY. In the final example we're only adding new records from _data.xlsx_
workbook sheet named "Title List" where the KEY is taken from column 1.

```shell
    dataset import data.csv
    dataset import data.csv 1
    dataset -sheet "Title List" import data.xlsx 1
    dataset -sheet "Title List" -update-only import data.xlsx 1
```

By default the header row of the table (the first row of the table) is used as the
attribute names of the JSON document you create on import.  If you don't want that
behavior you can use the "-use-header-row=false" option and the fields will be in the
form of "column_IDNO" where IDNO is replaced with a left zero padded column number
(e.g. column_001, column_002, column_003).

Related topics: extract and export

