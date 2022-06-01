Release 2.x:

This release is a rewrite of version 1 focusing on removing features, adding useful concepts and abstracting the storage engines cleanly. The latter was done to allow the web implementation of dataset to achieve an appropriate performance and be able to scale to a larger number of collections and size of collections.


The dataset collections structures have changed

- a dataset is a directory containing a collection.json and codemeta.json file
- the collection.json no longer contains general metadata or maps to the keys and pairtree, it focuses on operational settings (e.g. storage type and access information)
    - when running dataset as a web service or on a shared user machine you can setup the database connection through the environment. I.e. set DATASET_DSN_URI value (DSN URI is formed with a protocol named for the SQL driver, a "://" and the DSN for that driver, e.g. "mysql://DB_USER:DB_PASSWD@/DB_NAME")
- a codemeta.json file is now used for holding collection level metadata as this has been adopted by some of the data science community
- additional JSON configuration files may be used to manage the collection dependent on storage engine


Golang package changes:

- The dataset v2 package has been substantially reorganized and simplified, most things have changed
    - Collection.Init() now takes two parameters, collection name, an an optional DSN URI, if a DSN URI is provided it'll define the storage engine, e.g. a dataset.SQLSTORE)
    - Collection.DocPath() removed, doesn't make sense anymore since JSON may be stored in a SQL table
    - Collection.Read() only takes two parameters, not three
    - Collection.Keys() returns a list of keys and an error value
    - Collection.KeyExists() was renamed Collection.HasKey() to be more idiomatic in Go
    - Collection.Length() returns an int64 rather than an int
    - Collection.MetadataJSON() renamed Collection.Metadata() returns the codemeta JSON for the collection


libdataset:

- This has be dropped for now and is no longer supported due to the challenges in cross compiling for all supported platforms

CLI changes:

- options have be restructured so that some come after the verb
- help has been restructured to better support focusing the help text on the task needed
- the command line version is single user, single process for pairtree storage but can also use a SQL database with JSON column support for storage of JSON objects


Web Service changes:

- the web service does not support the pairtree implementation except export, it still uses a pairtree for storing attachments
- the web service is now RESTful so the end points no longer map directly to the client syntax
- the web service requires running a SQL database with JSON column support for storing JSON documents, currently this means SQLite3 and MySQL
- Access to the SQL storage engine is through either the environment or a URI expressing a storage type as protocol and a data source name to making the connection


Release 1.1.0:

Added attachment support for __datasetd__.

Updated the metadata fields to include richer PersonOrOrg data structures
for author, contributor, funder as well as added an annotation map field for custom metadata.

Added "MetadataJSON()" function for Collection to quickly copy out
the metadata values from a collection.

```
    c, err := dataset.Open("MyData.ds")
    ...
    defer c.Close()
    fmt.Printf("%s", c.MetadataJSON())
```

Added "MetadataUpdate()" function to update a collection's metadata.

```
    c, err := dataset.Open("MyData.ds")
    ...
    defer c.Close()
    meta := new(Collection)
    meta.Description = "A test dataset"
    meta.Version = "1.0.0"
    meta.Author = []*PersonOrOrg{
        &Person{
            Type: "Person",
            GivenName: "Jane",
            FamilyName: "Doe",
        },
    }
    err = c.MetadataUpdate(meta)
    ...
```

Depreciated dependency on namaste package and Namaste support in command line tools. Removed "collections.go and collections_test.go" from repository (redundent code). Updated libdataset/libdataset.go to hold functions that were needed for the C-Shared library from collections.go. The Namaste fields in the collection's metadata are now depreciated.

The dataset.Init() now places a lock file in the collection directory and leaves the collection in an "Open" state, it should be explicitly closed after Init is called.

E.g. 

```
   c, err := dataset.Init("MyData.ds")
   ...
   defer c.Close()
```

Removed "set_*" for collection metadata fields from libdataset.go. These should be set using the dataset command line tool only.

The dataset.Analzyer() and dataset.Repair() commands expect the dataset collections to be closed before being called. E.g..

```
    c, err := dataset.Open("MyData.ds")
    ...
    c.Close()
    err := dataset.Analyzer("MyData.ds", true)
    if err == nil {
        c, err = dataset.Open("MyData.ds")
        ...
    }
```

Release 1.0.2:

Added support for __datasetd__, a localhost web service for
dataset collections. The web service supports a subset of
the command line tool.

Both __datasetd__ and __dataset__ command line program now
include a "lock.pid" file in the collection root. This is to
prevent multiple processes from clashing when maintaining the
"collections.json" file in the collection root.

Migrated cli package into dataset repository sub-package "github.com/caltechlibrary/dataset/cli". Eventually this package will be replaced by "datasetCli.go" in the root folder.

In the dataset command line program the verb "detach" has been
renamed "retrieve" better describe the action. "detach" is depreciated
and will be removed in upcoming releases.

Release 1.0.1:

- Keys are stored lowercase
- Removed filtering and sorting options from dataset and libdataset
- Use pairtree 1.0.2 configurable separator
- Added check and repair for migrating to case insensitive keys and path
- Updated required packages to latest releases
- Added notes about Windows cmd prompt issues when providing JSON objects on command line
- Added M1 support for libdataset

Release 1.0.0:

- Initial Stable Release

Release 0.1.11:

- Requires go1.16 compilation
- Add macOS M1 compiled binaries

Release 0.1.10:

- Improved memory handling when handling for large attachments

Release 0.1.8:

This release focuses on minor bug fixes in libdataset.

- Removing duplicate functions:
    - `delete_frame()` has been superseded by `frame_delete()`
- Renamed functions:
    - `make_objects()` has been renamed `create_objects()` to be more consistent with naming scheme.
- Build Notes:
    - Golang v1.14
        - Caltech library go packages
            - storage v0.1.0
            - namaste v0.0.5
            - pairtree v0.0.4
    - OS used to compiled and test
         - macOS Catalina
         - Windows 10
         - Ubuntu 18.04 LTS
    - Python 3.8 (from Miniconda 3)
        - zip has replaced tar in the releases of libdataset
- Some tests fail on Windows 10 for libdataset. These will be addressed in future releases.

Release 0.1.6:

This release focuses on minor bug fixes in libdataset.
All functions which returned an error string only now return
True for success and False otherwise.  The error string
can be retrieved with `dataset.error_message()`.

- Build Notes:
    - Golang v1.14
    - Caltech library go packages
    - storage v0.1.0
    - namaste v0.0.5
    - pairtree v0.0.4
- OS used to compiled and test
    - macOS Catalina
    - Windows 10
    - Ubuntu 18.04 LTS
- Python 3.8 (from Miniconda 3)
- zip has replaced tar in the releases of libdataset
- Some tests fail on Windows 10 for libdataset. These will be addressed in future releases.

Release 0.1.5:

This release focuses on refine function names, simplification
and easy of testing for Windows 10 deployments.

- Build Notes:
    - Golang v1.14
    - Caltech library go packages
        - storage v0.1.0
        - namaste v0.0.5
        - pairtree v0.0.4
    - OS used to compiled and test
        - macOS Catalina
        - Windows 10
        - Ubuntu 18.04 LTS
    - Python 3.8 (from Miniconda 3)
    - zip has replaced tar in the releases of libdataset
- Renamed functions:
    - collection_status() is now collection_exists()
- Depreciated functions and features:
    - S3, Google Cloud Storage support dropped.
    - grid(), if you need this create a frame first and use frame_grid().
- Some tests fail on Windows 10 for libdataset. These will be addressed in future releases.

Release 0.1.4:

This release has breaking changes with release v0.1.3 and early.
Many functions in libdataset have been renamed to prevent collisions
in the environments using libdataset C-shared library. Most function
names now have two parts separated by a underscore (e.g. status
has become collection_status, repair has become collection_repair).

Google Sheet integration has been dropped. It was just more trouble
then it was worth to maintain.

Tests from py_dataset now have been ported to the test library for
libdataset.

Redundant functions have been removed (we had accumulated multiple 
definitions for the same thing in libdataset). Where possible
code has been simplified.

Most libdataset functions will cause an "open" on a dataset collection
automatically. Some additional functions around collections have been
added primarily to make testing easier (e.g. open_collection(), is_open(),
close_collection(), close_all()).

Functions that were overloaded via optional parameters have been simplified.
E.g. keys() now returns all keys in collection, use key_filter() and key_sort() accordingly.

- Dropped support for GSheet integration
- Only support pairtree layout of collection
- cleaned up libdataset API focusing on removing overloaded functions

Release 0.1.3:

- Bug fixes

Release 0.1.2:

- Persisting _Attachments metadata when updating with clean objects using the same technique as _Key

Release 0.1.1:

- Fixed problem where keys_exist called before an open command.

Release 0.1.0:

- Updated libdataset API, simplified func names and normalized many of the calls (breaking change)
- libdataset now manages opening dataset collections, inspired by Oberon System file riders (breaking change)
- Added Python test code for libdataset to make sure libdataset works
- Added support for check and repair when working on S3 deployed collections
- Refactored and simplified frame behavior (breaking change)


