
Action Items
============

Bugs
----

Next (prep for v2.x)
--------------------

- [ ] Add backward compatibility for dataset v2 with v1, re-implement a libdatset for use with py_datset, bump version number to 2.1.0
- [ ] Review [Go-app](https://go-app.dev/) and see if this would be a way to create a local client UI for working with datasets and enabling LunrJS for search
- [ ] Help cleanup
    - [ ] remove help pages for depreciated features
- [ ] Common dataset verbs (dataset/datasetd)
    - [X] keys
        - list the keys in a collection
    - [X] has-key
        - return "true "(w/OS exit 0 in CLI) if key is in collection,
          "false" otherwise (w/OS exit 1 in CLI)
    - [pkg,cli] sample
        - return a sample of keys from a collection
    - [X] create
        - add an new object to the collection if key does not exist,
          return false if object already exists or unable to create
          the new object
    - [X] read
        - return the object with nil error in the collection with the
          provided key, nil object and error value if not found
    - [X] update
        - replace the object in the collection for given key, return false
          is object does not to replace or replacement fails
    - [X] delete
        - delete the object in the collection for given key, return true
          if deletion was successful, false if the object was not deleted
          (e.g. key not found or the collection is read only)
    - [pkg,cli] versioning
        - set the versioning on a collection, the following strings enable
          versioning "major", "minor", "patch". Any other value disables
          versioning on the collection
        - [ ] read-versions, list the versions available for JSON object
        - [pkg,cli] read-version
             - return the object with nil error in the collection with the
               provided key and version, nil object and error value if not
               found
        - [pkg,cli] update-version
        - [pkg,cli] delete-version
        - [pkg] attachment-versions list versions of an attachment
        - [ ] attach-version add/replace a specific version of attachment
        - [pkg] retrieve-version retrieve version of attachment
        - [pkg] prune-version remove version of attachment
    - [X] frames
        - list the names of the frames currently defined in the collection
    - [X] frame
        - define a new frame in the collection, if frame exists replace it
    - [X] frame-meta
        - return the frame definition and metadata about the frame (e.g.
          how many objects and attributes)
    - [X] frame-objects
        - return the frame's list of objects
    - [X] refresh
        - update all the objects in the frame based on current state of
          the collection
    - [X] reframe
        - replace the frame definition but using the existing frame's keys
          refresh the frame with the new object describe
    - [X] delete-frame
    - [X] has-frame
- [X] Attachment support
    - [X] attachments
    - [X] attach
    - [X] retrieve (aka detach)
    - [X] prune
- [X] Verbs support by cli only
    - [X] sample
    - [X] clone
    - [X] clone-sample
    - [X] check
    - [X] repair
- [X] Document example Shell access to datasetd via cURL
- [X] take KeyMap out of collection.json so collection.json is smaller
    - support for segmented key maps (to limit memory consumption for very
      large collections)
- [X] Auto-version attachments by patch, minor or major release per
      settings in collection.json using keywords of patch, minor, major

Someday, Maybe
--------------

- [ ] Document an example Python 3 http client support for web API implementing a drop in replacement for py_dataset using the web service or cli
- [X] Missing tests for AttachStream()
- [ ] Implement a wrapping logger that takes a verboseness level for
      output (e.g. 0 - quiet, 1 progress messages, 2 warnings, errors
      should always show)
- [X] Memory consumption is high for attaching, figure out how to improve
      memory usage, switched to using streams where possible
- [ ] Add support for https:// based datasets (in addition to local disc
      and s3://)
- [ ] dsbagit would generate a "BagIt" bag for preservation of collection
      objects
- [ ] OAI-PMH importer to prototype iiif service based on Islandora
      content driven by a dataset collection
- [ ] Implement version support in the web service
- [ ] Implement an integrated UI for datasetd
