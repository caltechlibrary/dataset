
# Action Items

## Bugs

## Next (prep for v1.0.x)

+ [ ] Force keys to lowercase
+ [ ] Fix pathing issues running in cmd prompt under Windows 10
+ [ ] Remove dependency on tmplfn


## Someday, Maybe

+ [ ] Changing metadata for Namaste and Codemeta should re-render both
+ [ ] Add a service object to support building services with libdataset
+ [ ] Missing tests for AttachStream()
+ [ ] Auto-version attachments by patch, minor or major release per settings in collection.json
+ [ ] datasetd - a deamon for an http/https service for accessing dataset collections with support for multi-user public or restricted collections
+ [ ] Add Experimental Julia _dataset_ module for script collection management in Julia 
+ [ ] Add Experimental R _dataset_ module for scripting collection management in R
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
+ [ ] Memory consumption is high for attaching, figure out how to improve memory usage
    - Currently the attachment process generates the tar ball in memory rather than a tmp file on disc
    - for each attached filename process as stream instead of ioutil.ReadFile() and ioutil.ReadAll()
    - for size info, call Stats first to get the filesize to include in tarball header
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
- dataset "versioning" support via something like libgit2

