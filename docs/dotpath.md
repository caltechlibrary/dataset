
# dotpath 

## Description

### Basic

A dotpath is a notation to reference elements of a JSON object. It includes
the ability to traverse nested arrays.  The root of a is a leading dot. A dot is
typed as a period (i.e. "."). Given the following the dot path to the "name" field
would be ".name".

```json
    {"name":"Jane Doe"}
```

The dotpath `.name` would return the value "Jane Doe".

### Arrays

Arrays are designated with square brackets (E.g. [0] would reference the first
element of an array, arrays are number from zero).

```json
    ["one", "two", "three"]
```

The dotpath of `[0]` would correspond to the value "one", `[1]` would
correspond to the value "two" and `[2]` would refer to the value "three". If
you wish to include all the elements of an array you would use `[:]`. This
would return the full array.  Likewise if you need the second until end of the 
array you would get the values with `[2:]`.  Finally if you only wanted the
first and second element you could refere to it with the dotpath `[0:1]`.

### Putting it all together.

Often you have more complex objects including some level of nesting.  Element(s)
can be reference by combine the dotpaths into more complex expressions.

```json
    {
        "title": "Introducing dataset",
        "authors":[
            {"given_name":"Tom", "family_name":"Morrel"},
            {"given_name":"Robert","family_name": "Doiel"}
        ]
    }
```

You would reference the title with `.title`, the first author's family name with
`.authors[0].family_name` or get an array of authors family names with 
`.authors[:].fmaily_name`.

Related topics: [export-csv](export-csv.html), [frame](frame.html), [grid](grid.html), [import-csv](import-csv.html), [indexer](indexer.html)
