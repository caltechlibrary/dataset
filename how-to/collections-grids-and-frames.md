
# COLLECTIONS, GRIDS AND FRAMES

*dataset* stores JSON objects and can store a form of data frame based 
on grids of data. This document outlines the ideas behings *grid* and 
*frame* support in *dataset*.

## COLLECTIONS

Collections are at the core of the *dataset* tool. A collection is a 
pairtree directory structure storing a JSON objects in plaintext with 
optional attachments. The root folder for the collection contains a 
_collection.json_ file with the metadata associating a name to the 
pairtree path where the json object is stored. One of the guiding 
ideas behind dataset was to keep everything as plain text whenever 
reasonable.  The dataset project provides Go package for working with 
dataset collections, a python package (based on a shared library with 
the Go package) and command line tool.

Dataset collections are typically stored on your local disc but may be 
stored easily in Amazon's S3 (or compatible platform) or Google's cloud 
storage. Dataset and also import and export to/from a Google sheet or CSV
file.

Dataset isn't a database (there are plenty of JSON oriented databases out 
there, e.g. CouchDB, MongoDB and No SQL storage systems for MySQL and 
Postgresql). Rather the focus is on providing a mechanism to manage 
JSON objects, group them and to provide alternative 
data shapes for the viewing the collection (e.g. frames, grids).


## GRIDS

A *grid* is a 2D JSON array based on combining a set of keys (rows) and a 
list of dot paths (columns).  It is similar to the data shape you'd use in 
spreadsheets. It is a convenient data shape for build indexes, filtering 
and sorting.  *grid* support is also available in *dataset*'s Python 3.7 
package  

Here's an example of a possible grid for titles and authors.

```json
    [
        ["title", "authors"],
        ["20,000 Leagues under the Sea", "Verne, Jules"],
        ["All Around the Moon", "Verne, Jules"],
        ["The Short Reign of Pippin IV", "Steinbeck, John"]
    ]
```

If a column is missing a value then you should see a "null" for that cell. Here is an expanded example where
we've added a link to Project Gutenberg as a third column.

```json
    [
        ["title", "authors", "gutenberg_href"],
        ["20,000 Leagues under the Sea", "Verne, Jules", "http://www.gutenberg.org/ebooks/6538"],
        ["All Around the Moon", "Verne, Jules", "http://www.gutenberg.org/ebooks/16457"],
        ["The Short Reign of Pippin IV", "Steinbeck, John", null]
    ]
```


### A SIMPLE GRID EXAMPLE

This example creates a two column grid with *DOI* and *titles* from a 
dataset collection called *Pubs.ds* using the *dataset* command. Step one, 
generate a list of keys piping them into dataset using the grid verb.
If you didn't want to use a pipe you could also use an option to read
the keys from a file or to use all keys. The dataset keys command sends
the keys to standard out one key per line, the dataset grid command reads
the keys from standard input (one per line) and then creates a 
corresponding grid based on the dotpaths provided. In this example
we're using the paths ".doi" and ".title" from our "Pub.ds" collection.
If either “.doi” or “.title” is missing in a JSON object then a “null” 
value will be used. This way the grid rows retain the same number of 
cells.


```shell
    dataset keys Pubs.ds |\
        dataset grid Pubs.ds .doi .title
```

The 2D JSON array is easy to process in programming languages like Python. 
Below is an example of using a *grid* for sorting across an entire 
collection leveraging Python's standard sort method for lists.

```python
    import sys
    from py_dataset import dataset
    from operator import itemgetter
    keys = dataset.keys("Pubs.ds")
    (g, err) = dataset.grid("Pubs.ds", [".doi", ".title"])
    # g holds the 2D arrary
    if err != '':
        print(f'{err}')
        sys.exit(1)
    # sort by title
    g.sort(key=itemgetter (1))
    for row in g:
        (doi, title) = row
        print(f'{doi} {title}')
```


## THINKING ABOUT FRAMES

Implementing the grid verb started me thinking about the similarity to 
data frames in Python, Julia and Octave. A *frame* is defined as like
a grid in that it is an ordered list of objects. Frames can be retrieved
as a *grid* (2D array) or as a list of Objects. Frames also contain
a additional metadata to help them persist. Frames include enough
metadata (e.g. the attribute names, labels, used in the list of objects
as well as the dot paths used to provide a value to those labels).
If you want to get back a "Grid" of a frame you can optionally include
a header row as part of the 2D array returned.
*dataset* stores frames with the collection so unlike a *grid* it
is available for later processing.

Frames become handy when moving data from JSON documents (tree like)
to other formats like spreadsheets (table like). Date frames provide
a one to one map between a 2D representation and a list of objects
containing key/value pairs. Frames will become the way we define 
syncronization relationships as well as potentionally the way we 
define indexing should dataset re-aquire a search ability.

The map to frame names is stored in our collection's collection.json
Each frame itself is stored in a subdirectory of our collection. If you
copy/clone a collection the frames can travel with it.

## FRAME OPERATIONS

+ frame (read, define a frame)
+ frames (return a list of frame names)
+ reframe (refresh the content of a frame’s grid optionally replacing 
  the keys associated in the frame)
+ hasframe (check to see if a frame exists in the collection)
+ delete-frame


### Create a frame

Example creating a frame named "titles-and-dois"


```shell
    dataset keys Pubs.ds >pubs.keys
    dataset frame -i pubs.keys Pubs.ds titles-and-dois DOI=.doi Title=.title
```

Or in python


```python
    keys = dataset.keys('Pubs.ds')
    frame = dataset.frame('Pubs.ds', 'titles-and-dois', keys, ['.doi', '.title'], ["DOI", "Title"])
```


### Retrieve an existing frame

Example of getting the contents of an existing frame.

```shell
    dataset frame Pubs.ds titles-and-dois
```

Or in python

```python
    frame = dataset.frame('Pubs.ds', 'titles-and-dois')
```


### Regenerating a frame

Regenerating "titles-and-dois".

```shell
    dataset reframe Pubs.ds titles-and-dois
```

Or in python

```python
    frame = dataset.reframe('Pubs.ds', 'titles-and-dois')
```


### Updating keys associated with the frame

```shell
    dataset Pubs.ds keys >updated.keys
    dataset reframe -i updated.keys Pubs.ds reframe titles-and-dios
```

In python

```python
    frame = dataset.reframe('Pubs.ds', 'titles-and-dois', updated_keys)
```


### Updating labels in a frame

Labels are represented as a JSON array, when we set the labels explicitly we’re replacing the entire array at once. In this example the frame’s grid has two columns in addition the required `_Key` label. The `_Key` column is implied and with be automatically inserted into the label list. Additionally using `frame-labels` will cause the object list stored in the frame to be updated.

```shell
    dataset frame-labels Pubs.ds titles-and-dois '["Column 1", "Column 2"]'
```

In python

```python
    err = dataset.frame_labels('Pubs.ds', 'titles-and-dois', ["Column 1", "Column 2"])
```


### Removing a frame

```shell
    dataset delete-frame Pubs.ds titles-and-dios
```

Or in python

```python
    err = dataset.delete_frame('Pubs.ds', 'titles-and-dois')
```

## Listing available frames

```shell
    dataset frames Pubs.ds
```

Or in python

```python
    frame_names = dataset.frames('Pubs.ds')
```


# Data Grids

Often when processing data it is useful to pull date into a grid format.
_dataset_ provides a verb "grid" for doing just that. Below we're going
to create a small dataset collection called `grid_test.ds`, populate
it with some simple asymetric data (i.e. each record doesn't have the
same fields) then turn this into a 2D JSON array suitable for further
processing in a language like Python, R or Julia. The *grid* verb
is available in the Python module for dataset so we'll show that too.

In both examples the JSON representing our raw data seen in the file
[data-grids.json](data-grids.json)

## command line example 

### Generate a key list

From an existing collection, `grid_test.ds`, create a list of keys.

```shell
    dataset keys grid_test.ds >grid_test.keys
```

### Check a few records to see which will go into our grid.

We have the following keys  in our collection "gutenberg:21489",
"gutenberg:2488", "gutenberg:21839", "gutenberg:3186", 
"hathi:uc1321060001561131". Let's pick the first one and see what fields 
we might want in our grid (notice we're using the `-p` option to pretty 
print the JSON record).

```shell
    dataset read -p grid_test.ds "gutenberg:21489"
```

The fields that we're interested in are "._Key", ".title", ".authors",

### Create our grid from our collection

Now that we have a list of keys we're interested and and know the 
dot paths to the fields we're interested in we can create our grid.

```shell
    dataset grid -p -i=grid_test.keys grid_test.ds "._Key" ".title" ".authors" 
```

The results are a 2D array wich rows for each key and cells matching the
contents of the dot paths. Note that a cell may have a complex structure
like that shown with ".authors"

## python 3 example

In this example we're use the _dataset_ python module to read in our
raw JSON data (e.g. [data-grids.json](data-grids.json)) and convert
it into a *dataset* collection called "grid_test.ds". Next
we'll generate our set of keys and finally generate our grid as
a python list of lists.

```python3
    import sys
    import json
    from py_dataset import dataset

    # Read in our test data and convert from JSON into an array of dicts
    f_name = 'data-grids.json'
    with open(f_name, mode = 'r', encoding = 'utf-8') as f:
        src = f.read()
    data = json.loads(src)

    # create our collection
    c_name = 'grid_test.ds'
    err = dataset.init(c_name)
    if err != '':
        print(err)
        sys.exit(1)
    
    # load our test data
    for key in data:
        rec = data[key]
        err = dataset.create(c_name, key, rec)
        if err != '':
            print(err)
            sys.exit(1)
    
    # Create a list of keys and list of dot paths
    keys = dataset.keys(c_name)
    dot_paths = ["._Key", ".title", ".authors"]
    # now we can create our grid
    (g, err) = dataset.grid(c_name, keys, dot_paths)
    if err != '':
        print(err)
        sys.exit(1)

    # Now pretty print our grid
    print(json.dumps(g, indent = 4))
```
