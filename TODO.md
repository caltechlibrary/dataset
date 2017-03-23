
# Action Items

## Next

+ [ ] implement a repair collection command that would allow replacing/re-creating collection.json and keys.json based on what is discovered on disc
    + `dataset repair COLLECTION_NAME` would rescan the disc and write a new keys.json and collection.json
+ [ ] integrate support for storing dataset collections in AWS S3

## Someday, Maybe

+ [ ] implementing select lists as CSV files using Go's encoding/csv package 
+ [ ] take KeyMap out of collection.json so collection.json is smaller
+ [ ] add Bleve search support to *dataset* cli
    + default search would return IDS
    + detailed search would return values based on a list of dotpaths


## Completed

+ [x] support "attaching" non-JSON files to JSON record
    + [x] `dataset attach KEY FILENAME_LIST` would tar up FILENAME_LIST and place it next to the JSON record
    + [x] `dataset attachments KEY` returns a list of the tarballs content
    + [x] `dataset detach KEY` would remove the tarball from JSON record
    + [x] `dataset detach KEY FILENAME_LIST` would remove the selected file from tarball
    + [x] `dataset get KEY` get returns all the files in tarbal
    + [x] `dataset get KEY FILENAME_LIST` would return specific files from tarball
