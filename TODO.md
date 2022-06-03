
Action Items
============

Bugs
----

Next (prep for v2.x)
--------------------

- [ ] Common dataset verbs (dataset/datasetd)
    - [cli] keys
        - list the keys in a collection
    - [cli] has-key
        - return "true "(w/OS exit 0 in CLI) if key is in collection,
          "false" otherwise (w/OS exit 1 in CLI)
    - [cli] sample
        - return a sample of keys from a collection
    - [cli] create
        - add an new object to the collection if key does not exist,
          return false if object already exists or unable to create
          the new object
    - [cli] read
        - return the object with nil error in the collection with the
          provided key, nil object and error value if not found
    - [cli] read-version
        - return the object with nil error in the collection with the
          provided key and version, nil object and error value if not found
    - [cli] update
        - replace the object in the collection for given key, return false
          is object does not to replace or replacement fails
    - [cli] delete
        - delete the object in the collection for given key, return true
          if deletion was successful, false if the object was not deleted
          (e.g. key not found or the collection is read only)
    - [cli] versioning
        - set the versioning on a collection, the following strings enable
          versioning "major", "minor", "patch". Any other value disables
          versioning on the collection
    - [ ] frames
        - list the names of the frames currently defined in the collection
    - [ ] frame
        - define a new frame in the collection, if frame exists replace it
    - [ ] frame-meta
        - return the frame definition and metadata about the frame (e.g.
          how many objects and attributes)
    - [ ] frame-objects
        - return the frame's list of objects
    - [ ] refresh
        - update all the objects in the frame based on current state of
          the collection
    - [ ] reframe
        - replace the frame definition but using the existing frame's keys
          refresh the frame with the new object describe
    - [ ] delete-frame
    - [ ] has-frame
- [ ] Evaluate the following end points for datasetd, how to we manage
      concurent update of attachments from multiple client requests?
    - [ ] attach
    - [ ] attachments
    - [ ] retrieve (aka detach)
    - [ ] prune
- [ ] Verbs support by cli only
    - [cli] sample
    - [cli] clone
    - [cli] clone-sample
    - [ ] check
    - [ ] repair
- [ ] Document example Python 3 http client support for web API
- [x] Document example Shell access to datasetd via cURL
- [x] take KeyMap out of collection.json so collection.json is smaller
    - support for segmented key maps (to limit memory consumption for very
      large collections)
- [x] Auto-version attachments by patch, minor or major release per
      settings in collection.json using keywords of patch, minor, major

Someday, Maybe
--------------

- [ ] Missing tests for AttachStream()
- [ ] Implement a wrapping logger that takes a verboseness level for
      output (e.g. 0 - quiet, 1 progress messages, 2 warnings, errors
      should always show)
- [ ] Memory consumption is high for attaching, figure out how to improve
      memory usage
  - Currently the attachment process generates the tar ball in memory
    rather than a tmp file on disc
  - for each attached filename process as stream instead of
    ioutil.ReadFile() and ioutil.ReadAll()
  - for size info, call Stats first to get the filesize to include in
    tarball header
- [ ] Add support for https:// based datasets (in addition to local disc
      and s3://)
- [ ] dsbagit would generate a "BagIt" bag for preservation of collection
      objects
- [ ] OAI-PMH importer to prototype iiif service based on Islandora
      content driven by a dataset collection


