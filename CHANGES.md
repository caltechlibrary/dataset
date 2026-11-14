Unreleased

- libdataset (the WASM module, `libdataset.wasm`, and its Python and TypeScript wrappers) is deprecated. Releases
  after v2.5.2 no longer attach `dataset-vVERSION-libdataset.zip`: `make release` does not build it and
  `release-check` no longer requires it. No C shared library is provided either (it was dropped in the 2.0.0
  series). `libdataset.md` now says so and explains how to build `libdataset.wasm` yourself
  (`make libdataset.wasm`, or `GOOS=wasip1 GOARCH=wasm go build ./libdataset/`); `make libdataset-dist` still lays
  out the old zip from your own build. The py_dataset repository is archived. A different WASM module may replace
  libdataset in the future (dataset DR-0004).
- `libdataset.md` said Go 1.24 was needed; it is the version in `go.mod`.
- The README's link to ts_dataset pointed at `caltechlibraray`; fixed.

Release 2.5.2

- Added `strict_status` to the datasetd YAML (global, default false). When true the object, query and keys
  routes answer with the precise HTTP status codes of workspace DR-0020 and a JSON error body, with
  `Cache-Control: no-store`: 404 for a missing record on GET, DELETE and PUT (instead of 200, and instead of
  `200 {}` on GET), 409 for a JSON POST of a key that exists (nothing is written; before it silently
  overwrote), 400 for invalid JSON, an empty body or a bad query parameter (instead of 406), 415 for an
  unsupported content type, 413 for a body over 1 MiB (instead of 400), 422 for a record that fails schema
  validation (instead of 400), 404 for an unknown collection or route (instead of 403) and for an undefined
  query name, 500 for a store that fails, and 303 (instead of 302) after a successful form post. Form posts to
  an existing key still update it. Left out, every answer is as in 2.5.1. See `datasetd_api.5.md`.
- datasetd: a query that cannot be routed (missing name, unknown collection, no queries defined, name not
  defined) still answers 400 but now says which in the body.
- datasetd: a 405 now carries an `Allow` header.
- datasetd: a POST to a query with an unsupported content type, or a body that cannot be read, answered
  `303 See Other` to a `Location` made of the error text. It now answers `400` with the reason. This is the
  one status code that changes with `strict_status` off; it was a broken redirect.
- Fixed a `dataset` collection configured with an absolute `dataset:` path: named queries answered `400`
  and the collection's `create_success` and `create_error` redirects were silently ignored, because the
  configuration was looked up by the whole path while the collection is routed by its base name. The lookup
  now matches the base name. A deployment with an absolute path and a `create_success` setting will now
  redirect, which is what its settings asked for.
- datasetd now refuses to start when two collections have the same base name (for example `/a/etd.ds` and
  `/b/etd.ds`); the second used to silently replace the first.
- `dataset load` accepts both `key`/`object` (what `dump` writes) and `__key__`/`__object__` (what the help
  said); the help now says both.
- Fixed `Load` panicking on every call when built with Go 1.27 (`json.NewDecoder(nil)`).
- Fixed datasetd panicking on a form PUT to a collection with no model (a dropped connection for the client).
- Dropped the 32-bit Raspberry Pi build (`Linux-armv7l`) from the release. 64-bit Raspberry Pi OS is covered by
  `Linux-aarch64`. Releases already published are unchanged.
- Build and release brought to the workspace's current practice (workspace DR-0017, DR-0018): a plain `make`
  and `make build` no longer run `cmt`; `make generate` regenerates `version.go`, `CITATION.cff` and the
  installers; `make release` regenerates, tests, builds the six platform zips and the libdataset zip and runs
  `release-check`, which stops if anything does not carry the version and the commit hash or a zip lacks a
  program. `release.bash` and `release.ps1` are the shared versions (`--dry-run`, `--yes`, refuses with
  uncommitted or unpushed work, writes checksums, creates a draft release) and are tested by
  `release_test.bash` (`make test-release`).
- The manual page sources (`dataset.1.md`, `datasetd.1.md`, `dsquery.1.md` and the three `datasetd_*.5.md`) moved
  from the repository root to `docs/`, beside the rest of the documentation site. They are generated from the
  programs' help text, which is where to edit them.
- Documentation: the read, create, update and delete end point pages describe the status codes in both modes;
  `update-endpoint.md` was a copy of the create page and is rewritten for PUT; `load.md` names both spellings;
  `datasetd.md` no longer shows a made-up URL and a 404 the default mode does not send.

Release 2.5.1

- Split SQLite driver by build target: normal builds (`!wasip1`) use glebarez/go-sqlite
  (pure-Go, no WASM runtime overhead, eliminates spurious stderr warnings on startup);
  WASM builds (`wasip1`) continue to use ncruces/go-sqlite3 with the SQLite binary
  embedded via the `/embed` sub-package
- Fixed nil pointer dereference in `dsquery` when querying a pairtree collection with
  `--pt-index`: `ds.SQLStore` was never initialised for pairtree collections, causing a
  panic at runtime

Release 2.5.0

- Added libdataset WASM module (GOOS=wasip1 GOARCH=wasm) replacing the old cgo-based C shared library
  - Supports the full dataset CLI operation set (26 ops: lifecycle, CRUD, query, versioning, attachments, bulk, maintenance)
  - Named SQL queries pre-configured via settings.yaml-style config, matching datasetd's query approach
  - Python wrapper using wasmtime-py (`wrappers/python/libdataset/`)
  - TypeScript/Deno wrapper (`wrappers/typescript/libdataset.ts`)
  - New release artifact: `dataset-vVERSION-libdataset.zip`
  - Build with: `make libdataset.wasm`
- Switched SQLite driver from glebarez/go-sqlite3 (SQLite 3.41.2) to ncruces/go-sqlite3 (SQLite 3.53.1)
  - Enables WASM/WASIP1 compilation
  - FTS5 and full SQLite feature set retained
  - Existing collection files and DSN URIs (`sqlite://...`) unchanged
- Fixed datasetd request body size limits not being enforced (jsonSizeLimit 1 MiB, attachmentSizeLimit 250 MiB now applied via http.MaxBytesReader)
- Fixed datasetd error redirect using HTTP 304 (Not Modified) instead of HTTP 303 (See Other)
- Replaced deprecated io/ioutil with io.ReadAll and os.ReadFile/os.WriteFile throughout

Release 2.4.0

- Removed MySQL support
- Added schema-based validation for dataset collections using the models package
- Added support for nested object and list structures in schema definitions
- Added comprehensive identifier types for metadata validation:
  - ISBN, ISSN for publication identifiers
  - DOI for digital object identifiers
  - ORCID, ROR, ISNI for person/organization identifiers
  - PMID, PMCID for biomedical literature
  - FundRef, LCNAF, VIAF, SNAC for authority control
  - ArXiv for preprint identifiers
  - EAN for product barcodes
- Added `schemas` configuration to settings.yaml for defining reusable validation schemas
- Added `schema_name` and `validate` fields to collection configuration for enabling per-collection validation
- API now returns `X-Validation-Errors` header with detailed validation error information on create/update failures
- Backward compatible: existing collections continue to work without schema configuration

Release 2.3.4

- Deprecated MySQL support in documentation, will depreciate MySQL support in future release
- Improved loading large JSON objects from jsonl files
- Fixed issue #164 where the queries in COLD would work in v2.2.0 but fail in v2.3.x.
- Removed support for SQL parameters in dsquery due to encoding issues and lack of practical use cases
- Removed duplicated code from dsquery.go and api_routes.go in favor of collection.go's implementation of query functionality.
- Added tailing semi-colon removal for SQL queries due to changes in behavior of SQLite3 driver

Release 2.3.3

- Issue #161 fix for handling GET with query were data is passed via URL parameters.
- Removed support for frame, clone, sample, sync and join support removed. 
- The dsimporter cli removed (use jsonl dump and load instead).
- Improved loading large JSON objects from jsonl files

Release 2.3.2

Frame, clone, sample, sync and join support removed. The dsimporter cli removed (use jsonl dump and load instead). Fixed a bug in using the JSON API with GET where command line parameters were not getting processed. Also added support for standard encoding of web forms in queries.

Release 2.3.1:

Documentation and tests have been removed for frames, clone, and sample. Frame support was removed from datasetd.

Release 2.3.0:

This release fixed problems related to how a collection might incorrectly return a unknown type. v2.3.x cycle will depreciated features in preparation for v2.4.

Release 2.0.0-b06:

This release is a rewrite of version 1 focusing on removing features, clearifying useful concepts and abstracting the storage engines cleanly. The latter was done to allow the web implementation of dataset to achieve an appropriate performance and be able to scale to a larger number of collections and size of collections.


The dataset collection's structure has changed.

- a dataset is a directory containing a collection.json and codemeta.json file
- the collection.json no longer contains general metadata or maps to the keys and pairtree, it focuses on operational settings (e.g. storage type and access information)
    - when running dataset as a web service or on a shared user machine you can setup the database connection through the environment. I.e. set DATASET_DSN_URI value (DSN URI is formed with a protocol named for the SQL driver, a "://" and the DSN for that driver, e.g. "mysql://DB_USER:DB_PASSWD@/DB_NAME")
- a codemeta.json file is now used for holding general collection level metadata. [codemeta](https://codemeta.github.io/) has been adopted by the data science community for describing data and software
- additional JSON configuration files may be used to manage the collection dependent on storage engine
- Experimental features in dataset v2
    - Using SQL JSON columns to store our JSON documents in a relational database
        - SQLite 3
        - MySQL 8
        - Posgres 14.5 (very experimental, still learning Postgres)
    - Object level versioning 
    - Attachment versioning
    - dataset package to support Async, multi-process access via JSON stored in SQL databases


Golang package changes:

- Minimum Go version is now 1.19
- The dataset v2 package has been substantially reorganized and simplified, most things have changed
    - Collection.Init() now takes two parameters, collection name, an an optional DSN URI, if a DSN URI is provided it'll define the storage engine, e.g. a dataset.SQLSTORE)
    - Collection.DocPath() removed, doesn't make sense anymore since JSON may be stored in a SQL table
    - Collection.Read() only takes two parameters, not three
    - Collection.Keys() returns a list of keys and an error value
    - Collection.KeyExists() was renamed Collection.HasKey() to be more idiomatic in Go
    - Collection.FrameExists() was renamed Collection.HasFrame() to be more idiomatic in Go
    - Collection.Length() returns an int64 rather than an int
    - Collection.MetadataJSON() renamed Collection.Metadata() returns the codemeta JSON for the collection
    - Collection.UpdateMetatada() has been added, takes the name of a codemeta.json to replace the existing codemeta content for the collection.


libdataset:

- The C shared library implementation has been dropped for now do to the challenges of easily cross compiling releases

CLI changes:

- options have be restructured so that most come after the verb
- help has been restructured to better support focusing the help text on the task needed
- the command line version is single user, single process and be default assumes pairtree storage. 
- Experimental features
    - JSON stored in SQL database via JSON column types
        - SQLite3
        - MySQL 8
        - Postgres 14.5


Web Service changes:

- the web service is generally RESTful so the end points no longer map directly to the client syntax
- it is recommended to use SQL storage for your dataset collections explosed using the web service (e.q. SQLite3 or MySQL). This is because a pairtree collection doesn't provide docuemnt locking
- relies on SQL JSON columns for storage
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

Depreciated dependency on namaste package and Namaste support in command line tools. Removed "collections.go and collections_test.go" from repository (redundant code). Updated libdataset/libdataset.go to hold functions that were needed for the C-Shared library from collections.go. The Namaste fields in the collection's metadata are now depreciated.

The dataset.Init() now places a lock file in the collection directory and leaves the collection in an "Open" state, it should be explicitly closed after Init is called.

E.g. 

```
   c, err := dataset.Init("MyData.ds")
   ...
   defer c.Close()
```

Removed "set_*" for collection metadata fields from libdataset.go. These should be set using the dataset command line tool only.

The dataset.Analyzer() and dataset.Repair() commands expect the dataset collections to be closed before being called. E.g..

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


