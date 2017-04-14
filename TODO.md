
# Action Items

## Bugs


## Next

+ [ ] integrate support for storing dataset collections in AWS S3
    + [x] figure out how to handle attachments with AWS S3 (e.g. download tar to temp file then work with it?)
    + [ ] dataset init s3://.... is not showing the correct export value
    + [ ] confirm I can perform all CRUD options on JSON blobs and attachments
    + [ ] confirm I can get a list of attachments back
    + [ ] confirm I can update attachments
    + [ ] confirm I can delete individual attachments
    + [ ] confirm I can delete all attachments
    + [ ] update docs, examples and how to for using AWS S3

## Someday, Maybe

+ [ ] add a _import_ verb to dataset where a single file can be rendered as many dataset records (e.g. spreadsheet rows as JSON objects)
    + syntax like `dataset import CSV csv_filename [column number to use for key value]`
+ [ ] implementing select lists as CSV files using Go's encoding/csv package 
+ [ ] take KeyMap out of collection.json so collection.json is smaller
+ [ ] add Bleve search support to *dataset* cli
    + default search would return IDS
    + detailed search would return values based on a list of dotpaths
+ [ ] implement a repair collection command that would allow replacing/re-creating collection.json and keys.json based on what is discovered on disc
    + `dataset repair COLLECTION_NAME` would rescan the disc and write a new keys.json and collection.json



## Completed

+ [x] support "attaching" non-JSON files to JSON record
    + [x] `dataset attach KEY FILENAME_LIST` would tar up FILENAME_LIST and place it next to the JSON record
    + [x] `dataset attachments KEY` returns a list of the tarballs content
    + [x] `dataset detach KEY` would remove the tarball from JSON record
    + [x] `dataset detach KEY FILENAME_LIST` would remove the selected file from tarball
    + [x] `dataset get KEY` get returns all the files in tarbal
    + [x] `dataset get KEY FILENAME_LIST` would return specific files from tarball
