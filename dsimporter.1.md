%dsimporter(1) dataset user manual | version 2.1.23 00278aa
% R. S. Doiel and Tom Morrell
% 2025-04-03

# NAME

dsimporter

# SYNOPSIS

dsimporter [OPTIONS] C_NAME CSV_FILENAME KEY_COLUMN

# DESCRIPTION

__dsimporter__ is a tool to import CSV content into a dataset collection
where the column headings become the attribute names and the row values
become the attribute values.

# PARAMETERS

C_NAME
: If harvesting the dataset collection name to harvest the records to.

CSV_FILENAME
: The name of the CSV file to import

KEY_COLUMN
: The column name to use the they object key. If none is provided then
the first column is used as the object key. Keys values must be unique.


# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-comma
: Set column delimiter

-comment
: Set row comment delimiter

-overwrite
: Overwrite objects on key collision

# EXAMPLES

Import a file with three columns

- item_code
- title
- location

The "item_code" is unique for each row. The data is stored
in a file called "books.csv". We are importing the CSV file
into a collections called. "shelves.ds"

~~~
dsimporter shelves.ds books.csv "item_code"
~~~


