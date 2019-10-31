
# Action Items

## Bugs

## Next (prep for v0.1.0)

+ [ ] Refactor cli to match the exported functions from collections.go, add necessary tests
    + [ ] Add frame-exists
    + [ ] Add frame-create as distinct from frame (retrieving an existing frame)
    + [ ] Add frame-refresh
    + [ ] Rename refresh to frame-refresh
    + [x] Add frame-objects and frame-grid verbs to dataset cli
    + [ ] check SetWho, GetWho
    + [ ] check SetWhat, GetWhat
    + [ ] check SetWhen, GetWhen
    + [ ] check SetWhere, GetWhere
    + [ ] check SetContact, GetContact
    + [ ] check SetVersion, GetVersion
+ [ ] Add a service object to support building services with libdataset
+ [ ] Changing metadata for Namaste and Codemeta should re-render both
+ [x] Remove labels funcs
+ [x] Drop ._Key -> _Key requirement in frames
+ [ ] Confirm all documentation is current before v0.1.0 release
+ [ ] Confirm consensus on the release version v0.1.0 


## Someday, Maybe

+ [ ] Missing tests for AttachStream()
+ [x] Need semver.IncPatch(), semver.IncMinor(), semver.IncMajor() functions so we can auto increment version numbers easily
+ [ ] Auto-version attachments by patch, minor or major release per settings in collection.json
+ [ ] Switch from Bleve to Lunr indexes and search
+ [ ] Switch go go-cloud over our storage.go module
    + [ ] Make check and repair work in cloud storage
        + [ ] storage package should use go-cloud instead of individual SDKs
        + [ ] need a working file walker for S3 and GS storage ...
        + [ ] confirm check for s3:// collections
        + [ ] confirm repair for s3:// collections
        + [ ] confirm migrate for s3:// collections
+ [ ] datasetd - a deamon for an http/https service for accessing dataset collections with support for multi-user public or restricted collections
+ [ ] Add Experimental Julia _dataset_ module for script collection management in Julia 
+ [ ] Add Experimental R _dataset_ module for scripting collection management in R
+ [ ] Add Experimental PHP _dataset_ module for script collection management in PHP 
    - support JSON-LD for cross collection integration
+ [ ] sparql cli interface for searching collection
+ [ ] add *publish* command to generate index.md, index.html
    + [ ] Generate codemeta.json based on collection and any Namaste in collection folder
        + https://codemeta.github.io/terms 
    + [ ] Generate Lunr indexes for each frame
    + [ ] Generate a index.md based on codemata.json, namaste, and collection.json
    + [ ] Generate a index.html based on index.md plus a Lunrjs search
        + needs to support aggregate as well as selectable indexes
+ [ ] add *archive* command would do a *publish* then archive the collection
    + support adding relevant Namaste for preservation
    + archive should be suitable for ingesting in preservation systems
        + e.g. create tar, bag or web archive formatted instance
+ [ ] Remove dependency on github.com/caltechlibrary/tmplfn
+ [ ] Evaludate JMESPath as replacement/augmentation to dot path 
+ [ ] Add some additional metadata fields
    + [ ] version control on/off for attachments (we could verison via Subversion or git depending...)
    + [ ] Date/time repair was done
    + [ ] Date/time clone was executed as well as basename name of cloned
        + [ ] clone should include info about where it was cloned from
+ [ ] Documentation updates
    - Write up spec for storage indicating where it relates to other approaches (e.g. datacrate, bagit, Oxford Common File Layout, dflat, redd, pairtree)
+ [ ] Consider implementing Sword importer(s)/exporter(s) (v3? when spec is settled)
+ [ ] Consider implementing an EPrint 3.x importer/exporter
+ [ ] `dataset index-frame COLLECTION_NAME FRAME_NAME INDEX_NAME` - generate a Lunrjs or Bleve Index for search
+ [ ] `dataset ccreate COLLECTION_NAME FRAME_NAME CRATE_NAME` - generate a [datacreate](http://ptsefton.com/2017/10/19/datacrate.htm) from a collection for given keys
+ [ ] Implement a wrapping logger that takes a verboseness level for output (e.g. 0 - quiet, 1 progress messages, 2 warnings, errors should always show)
+ [ ] dataset explorer tool, possibly electron base for single user exploration of dataset collections
    - Browser based for UI, localhost restrict server for interacting with file system
    - Interactively build up of command strings, display results and saving off commands to runnable Bash scripts
    - Support importing datasets from s3:// and gs:// to local disc for interactive work
+ [ ] Integrate lunrjs and an index.html file into the root folder of a collection, this could be used to provide a web browser read interface to the collection without installing dataset itself.
+ [ ] Depreciate _dsindexer_ in favor of Bleve native cli
+ [ ] Memory consumption is high for attaching, figure out how to improve memory usage
    - Currently the attachment process generates the tar ball in memory rather than a tmp file on disc
    - for each attached filename process as stream instead of ioutil.ReadFile() and ioutil.ReadAll()
    - for size info, call Stats first to get the filesize to include in tarball header
+ [ ] Migrate export functions into an appropriate sub-packages (e.g. like how subpackages work in Bleve)
+ [ ] Move indexes and definitions into folder with collection.json
+ [ ] Add support for https:// based datasets (in addition to local disc and s3://)
+ [ ] Inaddition to UUID, add support for ULID (https://github.com/oklog/ulid) or provide an option for using ulid instead of uuid
+ [ ] VCARD and VCAL importer
+ [ ] Should the keymap in collection.json be a separate file(s)?
+ [ ] optional strageties for including arrays in a single column of CSV output
    - provide a hint for eaching express such as quoted comma delimited list, semi-column delimited list, pipe delimited list, etc.
+ [ ] Bug? Need to include optional stimmers (e.g. search for Adventure should also spot Adventures)
+ [ ] Improve internal stringToGeoPoint support a few more string notations of coordinates
    + [ ] N35.0000,W118.0000 or S35.000,E118.000
    + [ ] slice notation (GeoJSON) with longitude as cell 0, latitude as cell 1
+ [ ] take KeyMap out of collection.json so collection.json is smaller
    - support for segmented key maps (to limit memory consuption for very large collections)
+ [ ] dsbagit would generate a "BagIt" bag for preservation of collection objects
+ [ ] OAI-PMH importer to prototype iiif service based on Islandora content driven by a dataset collection
+ [ ] RSS importer (example RSS as JSON: http://scripting.com/rss.json)
+ [ ] OPML importer
- dataset "versioning" support via something like libgit2


## Completed

+ [x] Create a attachment migration Python script example
    + Read out the jsonObject
    + Find the tarball and untar the listed files into a temp directory
    + For each untared file use dataset attach to attach it back onto object.
+ [x] Document GSheet cell length limitations and how that impact sync-send
+ [x] Need How-to covering key sort/filter and frame sort/filter plus limitations
    + [x] Missing documentation for filtering/sorting keys for Python package
+ [x] Need cli option/wrapping func to filter out `_Key` and `_Attachments` of read result
+ [x] Remove dependency on Blevesearch
    + we will be moving forward with Lunr
+ [x] Refactor grids to be an array of objects in frames
+ [x] Refactor attachments from tarball use`a semver sub
      directory scheme
+ [x] Remove bucket layout code
+ [x] Improve collection.json with the following additional fields
    + [x] Collection name
    + [x] Collection description set via Namaste
    + [x] Date/time init was run creating collection
+ [x] ValueRenderOption("FORMULA") is wrong given the newer option of `UNFORMATED_VALUE`
+ [x] Sort out cross compiling libdataset shared library for Python module
+ [x] sync command needs to be implemented for csv based tables
+ [x] sync command needs to be implemented for GSheet based tables
+ [x] grids seem to only work on buckets not pairtree, why?
    + if pairtree grid failes with an error message about buckets or returns nothing
+ [x] fixed: Namaste isn't being written to on repair
+ [x] Evaluate adding namaste verb for collections
    - `dataset who COLLECT_NAME "Doiel, R. S."`
    - namaste feilds should be added in collection.json too
+ [x] We have a regression in support S3 buckets, the AWS Bucket name isn't being saved in the "store" object
    + [x] confirm init for s3:// collections
+ [x] Add the ability to create a grid (array or records) with selected fields (e.g. `dataset -key-list=my.keys my.ds grid '.pub_date' '.title' '.authors'`), each contains the specific dotpath listed, be helpful to be able to read in from Python and leverage its sorting abilities
+ [x] Consider changing from aa-zz round robin buckets to a [pairtree](https://confluence.ucop.edu/display/Curation/PairTree) as buckets per OCFL
+ [x] Sync: When you add columns to Spreadsheet in Frame need to make sure the order is consistent and update the Header row to reflect it.
+ [x] Repair/check should work on S3 and Google Cloud Storage
    + needs testing
+ [x] A raw pairtree is being generated in the repo directory by test!
+ [x] Implement Buckets to Pair migration option
+ [x] Evaluate switching from aa to zz buckets to pairtree ppath under data
    + [x] Repair/check should handle old and new file layout (e.g. moving buckets to pairtree on upgrade) 
    + [x] Evaluate moving JSON object from [ID].json to [ID].json
    + [x] Evaluate moving "attachments" into a [collection_name]/[pairtree]/[ID]/[relative path for objects] (i.e. drop making tar balls) 
+ [x] If collections.json is missing layouts and buckets exists then type get set to BUCKETS_LAYOUT
+ [x] Implement repair and anlaysis for PAIRTREEs
+ [x] Moving object tree out of "data", leave "data" empty to be compatible with other bagit tools
+ [x] Evaluate moving buckets into a "payload" (i.e. "data") folder for easier Bagging
+ [x] Added namaste type and when on dataset init
+ [x] Fix attachment handling so listing attachment names are fast (move out of tarball and save as a subdirectory using ID as name)
+ [x] Add clone verb to _dataset_ command, clone will copy a repository if the -sample option is used it will copy a sample of the source repository if two destination repositories are provided and sample is choosen then the first will contain the sample (training set) and second records not included in the first (the test set)
+ [x] change dataset join update to dataset join append
+ [x] Merge _dsfind_ and _dsindexer_ into _dataset_ command
+ [x] Normalize Create, Read, Update to have CreateJSON, ReadJSON, UpdateJSON counter parts for working with non-map[string]interface{} objects
+ [x] Create an experimental Python native module for dataset package
+ [x] In _dsindexer_ 'delete' remove one or more records from an index using record ids
    - An array of ids should work as a batch delete
+ [x] Document creating/managing indexes using the Bleve native cli
+ [x] Update dataset documentation to use Bleve's JSON definitions for indexes
+ [x] Update demos to use Bleve's JSON definitions for indexes
+ [x] Re-write docs for JSON index definitions
+ [x] Re-write demos for JSON index definitions
+ [x] Re-write examples for JSON index definitions
+ [x] Re-write how-to for JSON index definitions
+ [x] Evaluate adding automatic Lunrjs index support for collections
+ [x] In _dsindexer_ adopt JSON map compatible with  `bleve create INDEX_NAME -m INDEX_DEF`
+ [x] In _dsindexer_ 'add' to add/update one or more records in an existing index
    - An array of objects should work as a Batch update
+ [x] Remove automated metadata for `_Attachments` when removing attachments from a JSON document
+ [x] Attachment metaphor still needs better alignment with idiomatic go
    + [x] AttachFile should be implemented with an io.Writer interface
+ [x] If you _dataset delete KEY_ it fails to remove any attachments before deleting the JSON file
+ [x] if you _dataset detach KEY_ a stale _Attachments remain
+ [x] _dataset_ collection records only store "objects" (e.g. start and end with curly brackets) rather than allow Arrays
+ [x] Add automatic metadata fields for `_Key` when creating a new JSON document in a collection
+ [x] Add automatic metadata field for `_Attachments` when attaching a file to a JSON document
+ [x] Use automated metadata when asking for list of attached files, e.g. `_Attachments` for a JSON document
+ [x] In _dsfind_ Add `-sample N` option
+ [x] -nl line should be defaulted to true in dataset
+ [x] -nl line should be defaulted to true in dsfind
+ [x] -nl line should be defaulted to true in dsindexer
+ [x] Migrate the cli funciton in _dsindexer_ to package level
+ [x] Migrate cli functions in _dsfind_ to package level
+ [x] Migrate cli functions in _dataset_ to package level
+ [x] Attachment listings are slow
    - Add an `_Attachments` attribute to _dataset_ document with metadata about the attached file
+ [x] dataset -p read ... doesn't indent JSON output
+ [x] In _dataset keys_ Add `-sample N` option
+ [x] -help isn't showing help topics, -help sample isn't showing the sample help page.
+ [x] 'dataset keys FILTER' should emit keys as they are found to match rather then be processed as a group (unless we're sorting)
+ [x] `dataset list` should return a list (JSON array) of keys, missing keys should be ignore, if no keys then an empty list should be returned
+ [x] Behavior of -timeout, -wait seem wrong in practice, on some cli when you want to explicitly read from stdin you pass a hyphen to -input or -i.
+ [x] dataset attachements error:  Renaming can produce a cross device link error for the tarballs, the code uses a rename to "move" the file, need to implement it as copy and delete if we have this error
    - fixed error is storage package, line 77 fs.go was using a os.Rename() with out handling the error directly.
+ [x] "keys" should support a single level sort of a dotpath that resolves to a simple JSON type (e.g. int, float or string)
+ [x] "read" should accept a list of keys and produce an ordered list of JSON list of records
+ [x] "keys" could accept an existing list of keys to provide a sub-select like feature when combined with filter and order expressions
+ [x] "count" should accept a filter to support sub counts
+ [x] "keys" should be extended to accept a filter 
+ [x] "filter" should support RegEx matching, e.g. `(match "*.md$" .filenames[:])`
    - add this support via tmplfn package
+ [x] Add composite fields to indexes by leveraging text templates to modify JSON structure
+ [x] Add template defined format support 
    - currently required templates are page.tmpl (for HTML page), include.tmpl (for HTML includable output)
+ [x] Add filter aware CSV export
+ [x] Add filter aware value list extraction (e.g. all the unique orcids in a collection of authors data)
+ [x] Depreciate select commands in favor of filter, export and extract
+ [x] Add a filter function to support listing keys for records where the filter evaluates to true
    - Use the pipeline filters available in Go's text templates's if clause
+ [x] Add _haskey_ for a fast check if the key exists (look inside collections.json/keys.json only)
+ [x] Add option for batch indexing in dsindexer
+ [x] Reconfigure Makefile to build individual releases for each supported platform
+ [x] Merge results.tmpl changes into defaults from dr2
+ [x] when adding a fielded search in default templates the query string breaks the HTML of the query input form
    - double quotes make <input ... value="{{- . -}}" ...> break
    	- is it better to just have query field be a textarea, or use the urlencode/urldecode functions from tmplfn
+ [x] implement a repair collection command that would allow replacing/re-creating collection.json and keys.json based on what is discovered on disc
    - `dataset repair COLLECTION_NAME` would rescan the disc or s3 bucket and write a new keys.json and collection.json
    - Should also serve as a means to update a collection from one version of dataset to another
+ [x] idxFields work for single indexes but fail on multiple indexes in an Alias, find a workaround
+ [x] Add check to make sure page.tmpl and include.tmpl are available, if not use the ones from defaults
+ [x] Add support for indexing arrays values and objects in index definitions
    + [x] code 
    + [x] test
+ [x] add Bleve search support to dataset
    + [x] paging options (starting from/to, all records)
        + [x] add option to return all results
    + [x] default search would return IDS
    + [x] detailed indexing should be configurable including which fields on a list of dotpaths and options
    + [x] search results should be able to merge multiple indexes
    + [x] sortable result options (e.g. sort by ascending,descending fields)
    + [x] output should support returning only ids 
    + [x] alternate output formats (e.g. JSON arrays, select lists, CSV exports)
        + [x] JSON output
        + [x] CSV output
        + [x] id only list
    + [x] handle specific typed data like dates and geo cordinates in index definition
        + [x] look at using dataset JSONDencode rather than json.Unmashal so numbers aren't all treated as float64
        + [x] think about handling common date formatting for indexing and query
        + [x] test GeoCoding and Sort in Bleve
+ [x] add a _import_ verb to dataset where a single file can be rendered as many dataset records (e.g. spreadsheet rows as JSON objects)
    - syntax like `dataset import csv_filename [column number to use for key value]`
+ [x] integrate support for storing dataset collections in AWS S3
    + [x] figure out how to handle attachments with AWS S3 (e.g. download tar to temp file then work with it?)
    + [x] dataset init s3://.... is not showing the correct export value
    + [x] confirm I can perform all CRUD options on JSON blobs and attachments
    + [x] confirm I can get a list of attachments back
    + [x] confirm I can update attachments
    + [x] confirm I can delete individual attachments
    + [x] confirm I can delete all attachments
    + [x] update docs, examples and how to for using AWS S3
+ [x] support "attaching" non-JSON files to JSON record
    + [x] `dataset attach KEY FILENAME_LIST` would tar up FILENAME_LIST and place it next to the JSON record
    + [x] `dataset attachments KEY` returns a list of the tarballs content
    + [x] `dataset detach KEY` would remove the tarball from JSON record
    + [x] `dataset detach KEY FILENAME_LIST` would remove the selected file from tarball
    + [x] `dataset get KEY` get returns all the files in tarbal
    + [x] `dataset get KEY FILENAME_LIST` would return specific files from tarball
+ [x] add verbose option for importing CSV file, for large files it would be nice to see activity
+ [x] keys.json and collection.json's keymap are empty in some cases
    + [x] check dataset
    + [x] check cait usage
    + [x] check epgo usage
+ [x] dstocsv would take a select list and a list of "column name/dot path" pairs or a list of dot paths writing the results into a CSV file
    - header line would be optional 
    - dot paths that point at array, objects would be joined with a multi-value delimiter based on type 
    - mult-value delimiters would be configurable indepentantly
        - a object k/v might be delimited by colon which each pair delimited by newline
        - an array might be delimited by a pipe or semi-colon
    - optional filter for specific JSON documents to flatten
+ [x] Titles don't seem to sort in deployment, triage problem - is it index definition or faulty search implementation?
+ [x] Add support for gs:// Google cloud storage as an alternative to disc and s3://
+ [x] Add Google Sheet import based on existing CSV import code
+ [x] Add Google Sheet export based on existing CSV export code
+ [x] dataset count to return a count of records
+ [x] Bleve search support for dataset
    + [x] integrate batch indexing to speed things up
    + [x] generate a select list from search results
+ [x] prototype what a web service might look like for a dataset collection (including search)
    + [x] template HTML results and search forms
    + [x] support static pages in site
    + [x] evaluate including SparQL support
+ [x] Titles don't seem to sort in deployment, triage problem - is it index definition or faulty search implementation?
+ [x] Add Google Sheet import based on existing CSV import code
+ [x] Add Google Sheet export based on existing CSV export code
+ [x] dataset count to return a count of records
+ [x] collection.Create() will replace an existing record. What do I want to want to do a Join style update instead of a replace? 
+ [x] Add support for gs:// Google cloud storage as an alternative to disc and s3://
+ [x] convert extract, etc to work on streams so we can leverage pipelines more effeciently
+ [x] Repair and check will fail on S3/Google Cloud Storage without warning or reason why it is failing
    - Add an error message or implementing repair and check for s3:// and gs:// storage systems
+ [x] _dsfind_ Implement simple field filters using a prefix notation (e.g. (and (gt pubDate "2017-06-01") (eq (has .authors_family[:] "Doiel") true)))
    + [x] explore using templates as filters for select lists and the like
    + [x] implement select lists that save results as CSV files (sorting then could be off loaded
+ [x] implementing select lists as CSV files using Go's encoding/csv package 
+ [x] Add Python3 _dataset_ module for scripting collection management in Python3
+ [x] Drop uuid integration for import/export
+ [x] Mark _indexer_, _deindexer_ and _find_ experimental features
+ [x] Remove extract as it is easier to filter via Python and grids or frames
+ [x] source collection, isn't being respected when using the -c, -collection option for collections that aren't s3, google cloud or `*.ds` in v0.0.39
    + collection detection should check for the collection.json document
    + namaste type
    + Use storage.StoreType() to determine type info from path
+ [x] Sizes map is not bring respected when add/prune records, need to add test in TestAttatchments() then fix
+ [x] Missing test for getAttachmentList() which translates the map[string]interface{} to an array of Attachment
+ [x] Need an AttachFile that will support reading a buffer for eprinttools support

