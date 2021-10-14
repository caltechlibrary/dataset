Working with CSV
================

[CSV](https://en.wikipedia.org/wiki/Comma-separated_values) files are 
commonly used to share data. Most spreadsheets and many database systems 
can export and import from CSV files.  _datatset_ can use a spreadsheet 
in CSV format to populate JSON objects in a collection. The header row 
of the CSV file will become the object attribute names and the rows will 
become their values. _dataset_ requires a column of unique values to 
become the keys for the JSON objects stored in the collection. 

You can export CSV directly from a collection too. The paths to the 
elements in the objects become the header row and the values exported 
from the objects become the subsequent rows.

    NOTE: In an upcoming release of data the specific command line 
    parameters and Python method definitions may change.  Now that 
    _dataset_ supports the concept of frames import and export to 
    a tabular structure (or even synchronizing with a tabular 
    structure) can be simplified. This document describes the current
    release method for working with CSV files.


Import objects from a CSV file
------------------------------

You can import rows of a CSV document as JSON documents. This is 
useful when you have a large number of simple structures.

In this example we have a simple CSV file containing the following

```csv
    name,email
    Zowie,zowie@zbs.example.org
    Ralph Rolf,ralph.rolf@zbs.example.org
    Mojo Sam,mojo.sam@zbs.example.org
    Little Frieda,little.frieda@zbs.example.org
    Far Seeing Art,old.art@zbs.example.org
    Jack Flanders,captain.jack@zbs.example.org
```

Save this file as _characters.csv_. To import this let's create a 
collection named _characters_.

```shell
    dataset init characters-v1.ds
```

Now we can populate our characters collection by importing 
_characters.csv_.  Then look at the keys.

```shell
    dataset import characters-v1.ds characters.csv 2
    dataset keys characeter-v1.ds 
```

Notice the assigned ids. We used the second column, the one with th 
email heading to be our keys.

```shell
    ralph.rolf@zbs.example.org
    zowie@zbs.example.org
    captain.jack@zbs.example.org
    little.frieda@zbs.example.org
    mojo.sam@zbs.example.org
    old.art@zbs.example.org
```

```shell
    dataset keys characters-v1.ds |\
       while read ID; do 
           dataset read characters-v1.ds "${ID}"; 
       done
```


Now let's make a new version of our characters collection but this time 
we'll column one (the name column) as the key.

```shell
    dataset init characters-v2.ds
    dataset import characters-v2.ds characters.csv 1
    dataset keys
```

Now our keys look a little different.

```shell
```

Reading the records back we see we have the JSON same document structure.

```shell
    dataset keys characters-v2.ds | \
        while read ID; do 
            dataset read characters-v2.ds "${ID}"; 
        done
```

Our records look like...

```shell
    {"email":"captain.jack@zbs.example.org","name":"Jack Flanders"}
    {"email":"zowie@zbs.example.org","name":"Zowie"}
    {"email":"ralph.rolf@zbs.example.org","name":"Ralph Rolf"}
    {"email":"mojo.sam@zbs.example.org","name":"Mojo Sam"}
    {"email":"little.frieda@zbs.example.org","name":"Little Frieda"}
```


What if the CSV file has no header row?
---------------------------------------

Let's create a new collection and try the "-use-header-row=false" option.

```shell
    dataset init characters-v4.ds
    dataset import -skip-header-row=false characters-v4.ds characters.csv
    dataset keys characters-v4.ds
    dataset keys characters-v4.ds | \
        while read ID; do 
            dataset read characters-v4.ds "${ID}"; 
        done
```

Our ids are like in our first example because we chose to use the default 
JSON document key.


```
    characters.csv_2
    characters.csv_3
    characters.csv_4
    characters.csv_5
    characters.csv_6
    characters.csv_7
    characters.csv_1
```

Now take a look at the records output

```json
    {"col1":"Zowie","col2":"zowie@zbs.example.org"}
    {"col1":"Ralph Rolf","col2":"ralph.rolf@zbs.example.org"}
    {"col1":"Mojo Sam","col2":"mojo.sam@zbs.example.org"}
    {"col1":"Little Frieda","col2":"little.frieda@zbs.example.org"}
    {"col1":"Far Seeing Art","col2":"old.art@zbs.example.org"}
    {"col1":"Jack Flanders","col2":"captain.jack@zbs.example.org"}
    {"col1":"name","col2":"email"}
```

Instead of a _name_ and _email_ property name we have _col1_ and _col2_.  
Setting "-use-header-row" to false can be used with column numbers and or 
"-uuid" option. Give it a try with this final collection.

```shell
    dataset import -use-header-row=false characters-v5.ds characters.csv 2
    dataset keys characters-v5.ds
    dataset keys characters-v5.ds | \
        while read ID; do 
            dataset read characters-v5.ds "${ID}"; 
        done
```

Explore what you see.


imports and exports
===================

importing data into a collection
--------------------------------

We can import data from a CSV file and store each row as a JSON document 
in dataset. You need to pick a column with unique values to be the key 
for each record in the collection.  In this example we assume column one 
has the key value.

```shell
    dataset init mydata.ds
    dataset import mydata.ds my-data.csv 1
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.


exporting data from a collection
--------------------------------

```shell
   dataset frame-create -all mydata.ds export-frame \
       '.id=id' \
       '.title=title' 
       '.publication=publication' 
       '.pubDate=date'
   dataset export mydata.ds export-frame 
```

If you wanted to restrict to a subset (e.g. publication in year 2016)
You just need to create a frame with that restriction.

```shell
   dataset keys mydata.ds '(eq 2016 (year .pubDate))' | \
      dataset frame-create mydata.ds published-2016 \
           '.id=id' '.title=title' '.pubDate=date' 
   dataset export mydata.published-2016 ds 
```

