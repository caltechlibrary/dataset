---
title: "libdataset (3) user manual"
pubDate: 2023-02-08
author: "R. S. Doiel"
---

# NAME

libdataset

# SYNOPSIS

Use via C.

~~~
include "libdataset.h"
~~~

Use via Python.

~~~
from py_dataset import dataset
~~~

# DESCRIPTION

libdataset is a C shared library based on the Go package called
dataset from Caltech Library.  The dataset package provides a unified
way of working with JSON documents as collections. libdataset was
create better integrate working with dataset collection from Python
via the [py_dataset](https://pypi.org/project/py-dataset/) Python package.


# METHODS

The following are the exported C methods available in the C-shared
library generated from `libdataset.go`.


## error_clear

error_clear will set the global error state to nil.

## error_message

error_message returns an error message previously recorded or an empty string if no errors recorded

## use_strict_dotpath

use_strict_dotpath sets the library option value for enforcing strict dotpaths. 1 is true, any other value is false.

## is_verbose

is_verbose returns the library options' verbose value.

## verbose_on

verbose_on set library verbose to true

## verbose_off

verbose_off set library verbose to false

## dataset_version

dataset_version returns the version of libdataset.

## init_collection

init_collection intializes a collection and records as much metadata as it can from the execution environment (e.g. username, datetime created). NOTE: New parameter required, storageType. This can be either "pairtree" or "sqlstore".


## is_collection_open

is_collection_open returns true (i.e. one) if a collection has been opened by libdataset, false (i.e. zero) otherwise


## open_collection

open_collection returns 0 on successfully opening a collection 1 otherwise. Sets error messages if needed.


## collections

collections returns a JSON list of collection names that are open otherwise an empty list.


## close_collection

close_collection closes a collection previously opened.


## close_all_collections

close_all_collections closes all collections previously opened


## collection_exists

collection_exits checks to see if a collection exists or not.


## check_collection

check_collection runs the analyzer over a collection and looks for problem records.


## repair_collection

repair_collection runs the analyzer over a collection and repairs JSON objects and attachment discovered having a problem. Also is useful for upgrading a collection between dataset releases.

## clone_collection

clone_collection takes a collection name, a JSON array of keys and creates a new collection with a new name based on the origin's collections' objects. NOTE: If you are using pairtree dsn can be an empty string otherwise it needs to be a dsn to connect to the SQL store.


## clone_sample

clone_sample is like clone both generates a sample or test and training set of sampled of the cloned collection. NOTE: The training name and testing name are followed by their own dsn values.  If the dsn is an empty string then a pairtree store is assumed.


#  import_csv

import_csv - import a CSV file into a collection

Syntax: COLLECTION CSV_FILENAME ID_COL

Options that should support sensible defaults:

- cUseHeaderRow
- cOverwrite

## export_csv

export_csv - export collection objects to a CSV file

Syntax: COLLECTION FRAME CSV_FILENAME

## sync_send_csv

sync_send_csv - synchronize a frame sending data to a CSV file
returns 1 (True) on success, 0 (False) otherwise.

## sync_recieve_csv
 
sync_recieve_csv - synchronize a frame recieving data from a CSV file
returns 1 (True) on success, 0 (False) otherwise.

## has_key

has_key returns 1 if the key exists in collection or 0 if not.

## keys

keys returns JSON source of an array of keys from the collection

## create_object

create_object takes JSON source and adds it to the collection with
the provided key.

## read_object

read_object takes a key and returns JSON source of the record

## update_object

update_object takes a key and JSON source and replaces the record
in the collection.

## delete_object

delete_object takes a key and removes a record from the collection


## join_objects

join_objects takes a collection name, a key, and merges JSON source with an
existing JSON record. If overwrite is 1 it overwrites and replaces
common values, if not 1 it only adds missing attributes.

## count_objects

count_objects returns the number of objects (records) in a collection.
if an error is encounter a -1 is returned.

## object_path

object_path returns the path on disc to an JSON object document
in the collection.

## create_objects

create_objects - is a function to creates empty a objects in batch.
It requires a JSON list of keys to create. For each key present
an attempt is made to create a new empty object based on the JSON
provided (e.g. `{}`, `{"is_empty": true}`). The reason to do this
is that it means the collection.json file is updated once for the
whole call and that the keys are now reserved to be updated separately.
Returns 1 on success, 0 if errors encountered.

## update_objects

update_objects - is a function to update objects in batch.
It requires a JSON array of keys and a JSON array of
matching objects. The list of keys and objects are processed
together with calls to update individual records. Returns 1 on
success, 0 on error.

## list_objects

list_objects returns JSON array of objects in a collections based on a
JSON array of keys.

## attach

attach will attach a file to a JSON object in a collection. It takes
a semver string (e.g. v0.0.1) and associates that with where it stores
the file.  If semver is v0.0.0 it is considered unversioned, if v0.0.1
or larger it is considered versioned.

## attachments

attachments returns a list of attachments and their size in
associated with a JSON obejct in the collection.

## detach

detach exports the file associated with the semver from the JSON
object in the collection. The file remains "attached".

## prune

prune removes an attachment by semver from a JSON object in the
collection. This is destructive, the file is removed from disc.

## frame

frame retrieves a frame including its metadata. NOTE:
if you just want the object list, use frame_objects().

## has_frame

has_frame returns 1 (true) if frame name exists in collection, 0 (false) otherwise

## frame_keys

frame_keys takes a collection name and frame name and returns a list of keys from the frame or an empty list.  The list is expressed as a JSON source.

## frame_create

frame_create defines a new frame an populates it.

## frame_objects

frame_objects retrieves a JSON source list of objects from a frame.

## frame_refresh

frame_refresh refresh the contents of the frame using the
existing keys associated with the frame and the current state
of the collection.  NOTE: If a key is missing
in the collection then the key and object is removed.

## frame_reframe

frame_reframe will change the key and object list in a frame based on
the key list provided and the current state of the collection.

## frame_clear

frame_clear will clear the object list and keys associated with a frame.

## frame_delete

frame_delete will removes a frame from a collection

## frame_names

frame_names returns a JSON array of frames names in the collection.

## frame_grid

frame_grid takes a frames object list and returns a grid
(2D JSON array) representation of the object list.
If the "header row" value is 1 a header row of labels is
included, otherwise it is only the values of returned in the grid.

## get_versioning

get_version will get the dataset "versioning" setting.

## set_versioning

This will setting the versioning on a collection. The settings can be 
"", "none", "patch", "minor", "major".

