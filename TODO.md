
# Action Items

## Bugs

+ [ ] Memory consumption is high, figure out how to improve memory usage
+ [ ] Attachment listings are slow
    + idea: assume all collection documents are an object, attach a `._attachments` to each document with attachment metadata, this would allow retrieval at same spead as document
+ [ ] `dataset list` should return a list (JSON array) of keys, missing keys should be ignore, if no keys then an empty list should be returned
+ [ ] Migrate cli functions in cmds to package level and convert from exported to private functions used only to support cmds cli
+ [ ] 'dataset keys FILTER' should emit keys as they are found to match rather then be processed as a group (unless we're sorting)


## Next (v0.1.x)

+ [ ] Add specific index search, e.g. path is  /api/INDEX_NAME/q? ...
+ [ ] Add /api/COLLECTION_NAME/records end point to get ALL keys in collection
+ [ ] Add /api/COLLECTION_NAME/records/RECORD_ID end point for fetch an individual collection record
+ [ ] Provide a mechanism to synchronize (only update matching rows, appending new rows) a Google Sheet with dataset collection.

## Roadmap (v0.2.x)

+ [ ] dataset explorer tool, possibly electron base for single user exploration of dataset collections
    + Browser based for UI, localhost restrict server for interacting with file system
    + Interactively build up of command strings, display results and saving off commands to runnable Bash scripts
    + Support importing datasets from s3:// and gs:// to local disc for interactive work
+ [ ] sparql cli interface for searching collection
    + support JSON-LD for cross collection integration
+ [ ] Add faceted support to search (dsfind, dsws)
+ [ ] Add Fast CGI support in _dsws_ to allow custom development in Python, PHP or R
+ [ ] Python3 native dataset module for scripting collection management in Python3
+ [ ] R native dataset module for scripting collection management in R
+ [ ] PHP native dataset module for script collection management in PHP 
+ [ ] A zero or negative length for result size will be treated as a request for all results in _dsws_ and _dsfind_

## Someday, Maybe

+ [ ] Move indexes and definitions into folder with collection.json
+ [ ] Fix attachment handling so listing attachment names are fast (move out of tarball and save as a subdirectory using ID as name)
+ [ ] Add support for https:// based datasets (in addition to local disc and s3://)
+ [ ] Inaddition to UUID, add support for ULID (https://github.com/oklog/ulid) or provide an option for using ulid instead of uuid
+ [ ] VCARD and VCAL importer
+ [ ] _dsfind_ Implement simple field filters using a prefix notation (e.g. (and (gt pubDate "2017-06-01") (eq (has .authors_family[:] "Doiel") true)))
    + [ ] explore using templates as filters for select lists and the like
    + [ ] implement select lists that save results as CSV files (sorting then could be off loaded
+ [ ] Should the keymap in collection.json be a separate file(s)?
+ [ ] optional strageties for including arrays in a single column of CSV output
    + provide a hint for eaching express such as quoted comma delimited list, semi-column delimited list, pipe delimited list, etc.
+ [ ] Bug? Need to include optional stimmers (e.g. search for Adventure should also spot Adventures)
+ [ ] Improve internal stringToGeoPoint support a few more string notations of coordinates
    + [ ] N35.0000,W118.0000 or S35.000,E118.000
    + [ ] slice notation (GeoJSON) with longitude as cell 0, latitude as cell 1
+ [ ] implementing select lists as CSV files using Go's encoding/csv package 
+ [ ] take KeyMap out of collection.json so collection.json is smaller
    + support for segmented key maps (to limit memory consuption for very large collections)
+ [ ] dsbagit would generate a "BagIt" bag for preservation of collection objects
+ [ ] OAI-PMH importer to prototype iiif service based on Islandora content driven by a dataset collection
+ [ ] RSS importer (example RSS as JSON: http://scripting.com/rss.json)
+ [ ] OPML importer
+ dsselect would generate select lists based on query results in the manner of dsfind
+ dataset "versioning" support via something like libgit2


## Completed

+ [x] Behavior of -timeout, -wait seem wrong in practice, on some cli when you want to explicitly read from stdin you pass a hyphen to -input or -i.
+ [x] dataset attachements error:  Renaming can produce a cross device link error for the tarballs, the code uses a rename to "move" the file, need to implement it as copy and delete if we have this error
    + fixed error is storage package, line 77 fs.go was using a os.Rename() with out handling the error directly.
+ [x] "keys" should support a single level sort of a dotpath that resolves to a simple JSON type (e.g. int, float or string)
+ [x] "read" should accept a list of keys and produce an ordered list of JSON list of records
+ [x] "keys" could accept an existing list of keys to provide a sub-select like feature when combined with filter and order expressions
+ [x] "count" should accept a filter to support sub counts
+ [x] "keys" should be extended to accept a filter 
+ [x] "filter" should support RegEx matching, e.g. `(match "*.md$" .filenames[:])`
    + add this support via tmplfn package
+ [x] Add composite fields to indexes by leveraging text templates to modify JSON structure
+ [x] Add template defined format support 
    + currently required templates are page.tmpl (for HTML page), include.tmpl (for HTML includable output)
    + if format parameters' value matches a known template name then it should treated as a "supported" format by dsws instance
+ [x] Add filter aware CSV export
+ [x] Add filter aware value list extraction (e.g. all the unique orcids in a collection of authors data)
+ [x] Depreciate select commands in favor of filter, export and extract
+ [x] Add a filter function to support listing keys for records where the filter evaluates to true
    + Use the pipeline filters available in Go's text templates's if clause
+ [x] Add _haskey_ for a fast check if the key exists (look inside collections.json/keys.json only)
+ [x] Add option for batch indexing in dsindexer
+ [x] Reconfigure Makefile to build individual releases for each supported platform
+ [x] Merge results.tmpl changes into defaults from dr2
+ [x] CSV and JSON output not sending correct Content-Type header in _dsws_
+ [x] when adding a fielded search in default templates the query string breaks the HTML of the query input form
    + double quotes make <input ... value="{{- . -}}" ...> break
    	+ is it better to just have query field be a textarea, or use the urlencode/urldecode functions from tmplfn
+ [x] implement a repair collection command that would allow replacing/re-creating collection.json and keys.json based on what is discovered on disc
    + `dataset repair COLLECTION_NAME` would rescan the disc or s3 bucket and write a new keys.json and collection.json
    + Should also serve as a means to update a collection from one version of dataset to another
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
    + syntax like `dataset import csv_filename [column number to use for key value]`
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
    + header line would be optional 
    + dot paths that point at array, objects would be joined with a multi-value delimiter based on type 
    + mult-value delimiters would be configurable indepentantly
        + a object k/v might be delimited by colon which each pair delimited by newline
        + an array might be delimited by a pipe or semi-colon
    + optional filter for specific JSON documents to flatten
+ [x] Titles don't seem to sort in deployment, triage problem - is it index definition or faulty search implementation?
+ [x] Fix CORS setting in _dsws_ (Let's Encrypt support implemented, not needed)
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
    + Add an error message or implementing repair and check for s3:// and gs:// storage systems

