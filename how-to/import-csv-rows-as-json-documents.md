
# How to import CSV rows as JSON documents

You can import rows of a CSV document as JSON documents. This is useful when
you have a large number of simple structures.

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

Save this file as _characters.csv_. To import this let's create a collection
named _characters_.

```shell
    dataset init characters-v1
```

Now we can populate our characters collection by importing _characters.csv_.
Then look at the keys.

```shell
    dataset import characters.csv
    dataset list keys
```

Notice the assigned ids.

```
    characters.csv_7
    characters.csv_2
    characters.csv_3
    characters.csv_4
    characters.csv_5
    characters.csv_6
```

The default ids are formed by combining the filename with the row number. By default the first
row is treated as the property names in the JSON document. Now take a look one of the rows
that now are JSON documents.

```shell
    dataset list keys | while read ID; do dataset read "${ID}"; done
```

Now let's make a new version of our characters collection but this time we'll column two (the email column)
as the key.

```shell
    dataset init characters-v2
    export DATASET=characters-v2
    dataset import characters.csv 2
    dataset list keys
```

Now our keys look a little different.

```
    ralph.rolf@zbs.example.org
    mojo.sam@zbs.example.org
    little.frieda@zbs.example.org
    old.art@zbs.example.org
    captain.jack@zbs.example.org
    zowie@zbs.example.org
```

Reading the records back we see we have the JSON same document structure.

```shell
    dataset list keys | while read ID; do dataset read "${ID}"; done
```

Our records look like...

```
    {"email":"captain.jack@zbs.example.org","name":"Jack Flanders"}
    {"email":"zowie@zbs.example.org","name":"Zowie"}
    {"email":"ralph.rolf@zbs.example.org","name":"Ralph Rolf"}
    {"email":"mojo.sam@zbs.example.org","name":"Mojo Sam"}
    {"email":"little.frieda@zbs.example.org","name":"Little Frieda"}
```

Again the header row becomes the property names of the JSON document. But what if you don't
have a unique ID and don't like the filename/row number in our first example?  You can generate
a UUID for each record by using the "-uuid" option. Let's create a third version of characters
and step through the results as before.


```shell
    dataset init characters-v3
    export DATASET=characters-v3
    dataset -uuid import characters.csv
    dataset list keys
    dataset list keys | while read ID; do dataset read "${ID}"; done
```

Notice that the UUID is inserted into the result JSON documents. This lets you easily keep
records straight even if you rename the keys when moving between collections.

```
    {"email":"little.frieda@zbs.example.org","name":"Little Frieda","uuid":"27a5295f-4a80-4855-a2d1-e8a3a1a4623f"}
    {"email":"old.art@zbs.example.org","name":"Far Seeing Art","uuid":"872f68fe-f96b-4ce0-83bb-5c255d28cae7"}
    {"email":"captain.jack@zbs.example.org","name":"Jack Flanders","uuid":"fa382371-9a9e-4ade-a63c-7ebf88ef266e"}
    {"email":"zowie@zbs.example.org","name":"Zowie","uuid":"c05fceaa-b5de-460a-9497-f38fd9434cef"}
    {"email":"ralph.rolf@zbs.example.org","name":"Ralph Rolf","uuid":"fb48731d-9da7-4cc0-990d-9a5d1e0b33ac"}
    {"email":"mojo.sam@zbs.example.org","name":"Mojo Sam","uuid":"5aea6f22-390c-4727-8235-b9cab5ea1180"}
```


## What if the CSV file has no header row?

Let's create a new collection and try the "-skip-header-row=false" option.

```shell
    dataset init characters-v4
    export DATASET=characters-v4
    dataset -skip-header-row=false import characters.csv
    dataset list keys
    dataset list keys | while read ID; do dataset read "${ID}"; done
```

Our ids are like in our first example because we chose to use the default JSON document key.


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

```
    {"col1":"Zowie","col2":"zowie@zbs.example.org"}
    {"col1":"Ralph Rolf","col2":"ralph.rolf@zbs.example.org"}
    {"col1":"Mojo Sam","col2":"mojo.sam@zbs.example.org"}
    {"col1":"Little Frieda","col2":"little.frieda@zbs.example.org"}
    {"col1":"Far Seeing Art","col2":"old.art@zbs.example.org"}
    {"col1":"Jack Flanders","col2":"captain.jack@zbs.example.org"}
    {"col1":"name","col2":"email"}
```

Instead of a _name_ and _email_ property name we have _col1_ and _col2_.  Setting "-skip-header-row" to false
can be used with column numbers and or "-uuid" option.  Give it a try with this final collection.

```shell
    dataset -skip-header-row=false import characters.csv 2
    dataset -skip-header-row=false -uuid import characters.csv
    dataset list keys
    dataset list keys | while read ID; do dataset read "${ID}"; done
```

Explore what you see.

