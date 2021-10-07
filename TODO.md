
Action Items
============

Bugs
----

Next (prep for v2.x)
--------------------

- [ ] Provide a dataset service (datasetd)
    - [ ] keys
    - [ ] create
    - [ ] read
    - [ ] update
    - [ ] delete
- [ ] Evaluate the following end points for datasetd
    - [ ] attach
    - [ ] attechments
    - [ ] retrieve
    - [ ] remove
- [ ] Add pid lock support for processes accessing dataset collections
- [ ] Document example Shell access to datasetd via cURL
- [ ] Document dataesetd access from Python 3
- [x] Migrate cli package into dataset package

Someday, Maybe
--------------

- [ ] Drop Namaste from dataset
- [ ] Missing tests for AttachStream()
- [ ] Auto-version attachments by patch, minor or major release per settings in collection.json
- [ ] Add some additional metadata fields
    - [ ] version control on/off for attachments (we could version via Subversion or git depending...)
    - [ ] Date/time repair was done
    - [ ] Date/time clone was executed as well as basename name of cloned
        - [ ] clone should include info about where it was cloned from
- [ ] Documentation updates
    - Write up spec for storage indicating where it relates to other approaches (e.g. datacrate, bagit, Oxford Common File Layout, dflat, redd, pairtree)
- [ ] Implement a wrapping logger that takes a verboseness level for output (e.g. 0 - quiet, 1 progress messages, 2 warnings, errors should always show)
+ [ ] Integrate Lunrjs and an index.html file into the root folder of a collection, this could be used to provide a web browser read interface to the collection without installing dataset itself.
+ [ ] Memory consumption is high for attaching, figure out how to improve memory usage
    - Currently the attachment process generates the tar ball in memory rather than a tmp file on disc
    - for each attached filename process as stream instead of ioutil.ReadFile() and ioutil.ReadAll()
    - for size info, call Stats first to get the filesize to include in tarball header
- [ ] Add support for https:// based datasets (in addition to local disc and s3://)
- [ ] VCARD and VCAL importer
- [ ] Should the keymap in collection.json be a separate file(s)?
- [ ] optional strategies for including arrays in a single column of CSV output
    - provide a hint for   express such as quoted comma delimited list, semi-column delimited list, pipe delimited list, etc.
- [ ] Bug? Need to include optional stimmers (e.g. search for Adventure should also spot Adventures)
- [ ] Improve internal stringToGeoPoint support a few more string notations of coordinates
    - [ ] N35.0000,W118.0000 or S35.000,E118.000
    - [ ] slice notation (GeoJSON) with longitude as cell 0, latitude as cell 1
- [ ] take KeyMap out of collection.json so collection.json is smaller
    - support for segmented key maps (to limit memory consumption for very large collections)
- [ ] dsbagit would generate a "BagIt" bag for preservation of collection objects
- [ ] OAI-PMH importer to prototype iiif service based on Islandora content driven by a dataset collection
- dataset "versioning" support via something like libgit2

