import
======

Syntax
------

```shell
    dataset import COLLECTION_NAME CSV_FILENAME ID_COL_NUMER
```

Description
-----------

_import_ adds JSON documents to a collection from a CSV table. 

Usage
-----

In the following examples the CSV filename is _data.csv_.
The first column (column 1) is used as the value for KEY if
specified.  Our collection is named "data.ds".

```shell
    dataset import data.ds data.csv 1
```

By default the header row of the table (the first row of the table) 
is used as the attribute names of the JSON document you create on 
import.  If you don't want that behavior you can use 
the "-use-header-row=false" option and the fields will be in the
form of "column_NO" where "NO" is replaced with a left zero 
padded column number (e.g. column_001, column_002, column_003).


Related topics: [export-csv](export-csv.html)

