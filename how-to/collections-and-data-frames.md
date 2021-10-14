COLLECTIONS, GRIDS AND FRAMES
=============================

__dataset__ stores JSON objects and can store also data frames similar to
that used in Python, R and Julia. This document outlines the ideas
behind __dataset__\'s implementation of data frames.

COLLECTIONS
-----------

Collections are at the core of the __dataset__ tool. A collection is a
pairtree directory structure storing JSON objects in plain text with
optional attachments. The root folder for the collection contains a
*collection.json* file with the metadata associating a name to the
pairtree path where the json object is stored. One of the guiding ideas
behind dataset was to keep everything in plain text (i.e. UTF-8)
whenever reasonable. The dataset project provides Go package for working
with dataset collections, a python package (based on a C-shared library
included in the Go package) and a command line tool.

Dataset collections are typically stored on your local disc but may be
stored easily in Amazon\'s S3 (or compatible platform) or Google\'s
cloud storage using operating systems integration (e.g. [fuse file
system tools](https://en.wikipedia.org/wiki/Filesystem_in_Userspace)).
Dataset can also import and export CSV files.

Dataset isn\'t a database (there are plenty of JSON oriented databases
out there, e.g. CouchDB, MongoDB and No SQL storage systems for MySQL
and Postgresql). __dataset__\'s focus is on providing a mechanism to
manage JSON objects, group them and to provide alternative data shapes
for the viewing the collection (e.g. data frames and grids).

DATA FRAMES
-----------

Working with subsets of data in a collection is useful, particularly
ordered subsets. Implementing this started me thinking about the
similarity to data frames in Python, Julia and Octave. A *frame* is an
ordered list of objects. It\'s like a grid except that rather than have
columns and row you have a list of objects and attribute names mapped to
values. Frames can be retrieved as a list of objects or a *grid* (2D
array). Frames contain a additional metadata to help them persist.
Frames include enough metadata to efficiently refresh objects in the
list or even replace all objects in the list. If you want to get back a
\"Grid\" of a frame you can optionally include a header row as part of
the 2D array returned.

__dataset__ stores frames with the collection so they are is available for
later processing. The objects in a frame reflect the objects as they
existed when the frame was generated.

Frames become handy when moving data from JSON documents (tree like) to
other formats like spreadsheets (table like). Date frames provide a one
to one map between a 2D representation and a list of objects containing
key/value pairs. Frames will become the way we define synchronization
relationships as well as potentially the way we define indexing should
dataset re-aquire a search ability.

The map to frame names is stored in our collection\'s collection.json
Each frame itself is stored in a sub directory of our collection. If you
copy/clone a collection the frames can travel with it.

FRAME OPERATIONS
----------------

-   frame-create (define a frame)
-   frame (read a frame back)
-   frames (return a list of frame names)
-   frame-reframe (replace all frame objects given a list of keys)
-   frame-refresh (update objects in a frame pruning objects no longer
    in the collection)
-   frame-exists (check to see if a frame exists in the collection)
-   frame-delete


### Create a frame

Example creating a frame named \"dois-and-titles\"

```shell
    dataset keys Pubs.ds >pubs.keys
    dataset frame-create -i pubs.keys Pubs.ds dois-and-titles \
        ".doi=DOI" \
        ".title=Title"
```

Or in python

```python
    keys = dataset.keys('Pubs.ds')
    frame = dataset.frame_crate('Pubs.ds', 'dois-and-titles', keys, {
        '.doi': 'DOI', 
        '.title': 'Title'
        })
```

### Retrieve an existing frame

Example of getting the contents of an existing frame with all the
metadata.

```shell
    dataset frame Pubs.ds dois-and-titles
```

An example of getting the frame\'s object list only.

```shell
    dataset frame-objects Pubs.ds dois-and-titles
```

Or in python getting the full frame with metadata

```python
    (frame, err) = dataset.frame('Pubs.ds', 'dois-and-titles')
    if err != '':
        print(f'Something went wront {err}')
```

Or only the object list (note: we\'re going to check for the frame\'s
existence first).

```python
    if dataset.frame_exists('Pub.ds', 'dois-and-titles'):
        object_list = dataset.frame_objects('Pubs.ds', 'dois-and-titles')
```

### Regenerating a frame

Regenerating \"dois-and-titles\".

```shell
    dataset reframe Pubs.ds dois-and-titles
```

Or in python

```python
    keys = dataset.keys('Pubs.ds')
    keys.sort()
    frame = dataset.frame_reframe('Pubs.ds', 'dois-and-titles', keys)
```

### Updating keys associated with the frame

```shell
    dataset Pubs.ds keys >updated.keys
    dataset frame-refresh -i updated.keys Pubs.ds reframe titles-and-dios
```

In python

```python
    frame = dataset.frame-refresh('Pubs.ds', 'dois-and-titles', updated_keys)
```

### Updating labels in a frame

Labels are represented as a JSON array, when we set the labels
explicitly we\'re replacing the entire array at once. In this example
the frame\'s grid has two columns in addition the required `_Key` label.
The `_Key` column is implied and with be automatically inserted into the
label list. Additionally using `frame-labels` will cause the object list
stored in the frame to be updated.

```shell
    dataset frame-labels Pubs.ds dois-and-titles '["Column 1", "Column 2"]'
```

In python

```python
    err = dataset.frame_labels('Pubs.ds', 'dois-and-titles', ["Column 1", "Column 2"])
```

### Removing a frame

```shell
    dataset frame-delete Pubs.ds titles-and-dios
```

Or in python

```python
    err = dataset.frame_delete('Pubs.ds', 'dois-and-titles')
```

Listing available frames
------------------------

```shell
    dataset frames Pubs.ds
```

Or in python

```python
    frame_names = dataset.frames('Pubs.ds')
```
