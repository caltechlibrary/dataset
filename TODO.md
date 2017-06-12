
# Action Items

## Bugs

+ [ ] keys.json and collection.json's keymap are empty in some cases
    + [x] check dataset
    + [ ] check cait usage
    + [ ] check epgo usage
+ [ ] repair and check will fail on S3 without warning or indication why

## Next

+ [ ] Merge results.tmpl changes into defaults from dr2
+ [ ] Add RSS2 and BibTeX format support via templates
    + currently required templates are page.tmpl (for HTML pages), include.tmpl (for HTML includable output)
    + if fmt value has a related template then format is "supported" by dsws
+ [ ] Add specific index search under /api (e.g. /api gives you everything, /api/authors would limit search to the authors.bleve index)
+ [ ] Add option for batch indexing in dsindexer
+ [ ] Reconfigure Makefile to build individual releases for each supported platform
+ [ ] add repair and check support for S3 storage

## Someday, Maybe

+ [ ] Should the keymap in collection.json be a separate file(s)?
+ [ ] optional strageties for including arrays in a single column of CSV output
    + provide a hint for eaching express such as quoted comma delimited list, semi-column delimited list, pipe delimited list, etc.
+ [ ] Bug? Need to include optional stimmers (e.g. search for Adventure should also spot Adventures)
+ [ ] prototype what a web service might look like for a dataset collection (including search)
    + [ ] template HTML results and search forms
    + [ ] support static pages in site
    + [ ] evaluate including SparQL support
+ [ ] Improve internal stringToGeoPoint support a few more string notations of coordinates
    + [ ] N35.0000,W118.0000 or S35.000,E118.000
    + [ ] slice notation (GeoJSON) with longitude as cell 0, latitude as cell 1
+ [ ] Bleve search support for dataset
    + [ ] integrate batch indexing to speed things up
    + [ ] generate a select list from search results
    + [ ] add facet support
+ [ ] implementing select lists as CSV files using Go's encoding/csv package 
+ [ ] take KeyMap out of collection.json so collection.json is smaller
    + support for segmented key maps (to limit memory consuption for very large collections)
+ sparql cli interface for searching collection
+ cli to convert collection into JSON-LD
+ dsselect would generate select lists based on query results in the manner of dsfind
+ dstoscv would take a select list and a list of "column name/dot path" pairs or a list of dot paths writing the results into a CSV file
    + header line would be optional 
    + dot paths that point at array, objects would be joined with a multi-value delimiter based on type 
    + mult-value delimiters would be configurable indepentantly
        + a object k/v might be delimited by colon which each pair delimited by newline
        + an array might be delimited by a pipe or semi-colon
+ dataset "versioning" support via something like libgit2
+ dsserver would allow HTTPS REST access do a collection server, it would support multi-user access and with group acls
    + authentication would be through an external system (e.g. Shibboleth, PAM, or OAuth2)
    + groups would contain a list of users
    + permissions (CRUD) would be based on group and collection (permissions would be collection wide, not record specific)
+ dsbagit would generate a "BagIt" bag for preservation of collection objects
+ collection.json should hold a list of available indexes and their definitions to automate repair
+ OAI-PMH importer to prototype iiif service based on Islandora content driven by a dataset collection
+ merge dsindexer and dsfind into dataset cli and depreciate individual programs
+ RSS importer (example RSS as JOSN: http://scripting.com/rss.json)


## Completed

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
