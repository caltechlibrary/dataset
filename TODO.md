
# Action Items

## Next

## Someday, Maybe

+ support "attaching" non-JSON files to JSON record
    + e.g. `dataset attach KEY FILENAME_LIST` would tar up FILENAME_LIST and place it next to the JSON record
    + e.g. `dataset attachments KEY` returns a list of the tarballs content
    + e.g. `dataset detach KEY` would remove the tarball from JSON record
    + e.g. `dataset detach KEY FILENAME_LIST` would remove the selected file from tarball
    + e.g. `dataset get KEY` get returns all the files in tarbal
    + e.g. `dataset get KEY FILENAME_LIST` would return specific files from tarball
+ implementing select lists as CSV files using Go's encoding/csv package 
+ take KeyMap out of collection.json so collection.json is smaller
+ implement a repair collection command that would allow replacing/re-creating collection.json and keys.json based on what is discovered on disc
    + `dataset repair COLLECTION_NAME` would rescan the disc and write a new keys.json and collection.json
+ add a repl
+ add Bleve search support to *dataset* cli
    + default search would return IDS
    + detailed search would return values based on a list of dotpaths
+ add a JavaScript/Python/Shell integration for defining functions and custom sorts
