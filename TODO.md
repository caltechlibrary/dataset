
Action Items
============

X
: completed

D
: declined, decided not to implement

Bugs
----

Next (prep for v2.1.21)
-----------------------

- [ ] Update datasetd to support urlencoded data submissions in additional to application/json
  - this would allow a simple data entry system to be build directly from HTML without the need for JavaScript in the browser
  - the urlencoded data should support embedded YAML in text areas for extrapolating more complex data structures
  - [ ] Handle POST with 'application/x-www-form-urlencoded'
    - [X] Handle the submissions mapping create/update to POST
    - [ ] Handle the success or failure of the create/update of POST
  - [D] Handle PUT with 'application/x-www-form-urlencoded'
    - Browsers only honor GET, POST and DIALOG in 2024. Frustrating.
    - Modified POST to work for both Create and Update actions, delete will get handled like GET but I need to decide on symantics (e.g. `?delete=true`)
  - [ ] Integrate models package into dataset cli and datasetd
   - [X] Modify Create and Update in datasetd to use the models package
   - [X] Handle generated elements on Create and Update
   - [ ] For forms submited with URL encoding reply I currently reply with JSON to indicate success or failure, should return HTML
     - [ ] Success response should provide record view of submitted content
     - [ ] Failure should redirect back to the form that was submitted
     - [ ] It'd be nice to style/theme the HTML for better integration with website
      - Is this a configuration chioice (e.g. success, fail pages in model attributes?)
   - [ ] Can I can skip the handlebars templates and just support HTML?
     - Skipping the templates has several advantages
       - one less thing to document
       - less dependency for datasetd
     - If I only support HTML generation then I need to include JavaScript in the generated code
       - If I do the then PUT and DELETE would work 
       - Downside is it requires JavaScript to update records and submit them
   - [ ] Decide if it is exceptible to render HTML with JavaScript to adjust form behavior between create and update
- [D] Update datasetd to handle YAML submission for create and update: This didn't work in my experiments, not sure why.
  - Form handling in browser restricts the mime types submitted, I'd have to use text/plain to submit YAML then check server side to make sure I had YAML

Someday, Maybe
--------------

- [ ] Should the project be renamed "collections"?
- [ ] Update datasetd to allow multipart form subission treating file(s) upload as an attachment request
- [ ] Rewrite py_dataset, drop support for libdataset
  - [ ] Figure out correct approach
    - [ ] Generate WASM module for libdataset
    - [ ] Use ts_dataset approach and required datasetd for Python support
    - [ ] Rewrite dataset, datasetd and libdataset in Rust and continue shared library support without built in GC
- [ ] create a cli named `ds` that wraps all the cli except datasetd similar to how the Go command or Git works
- [ ] My current approach to versioning is too confusing, causing issues in implementing py_dataset, versioning needs to be automatic with a minimum set of methods explicitly supporting it otherwise versioning should just happen in the back ground and only be supported at the package and libdataset levels.
  - [ ] create, read, update, list operations should always reflect the "current" version (objects or attachments), delete should delete all versions of objects as should prune for attachments, this is because versioning suggests things never really get deleted, just replaced.
- [ ] Common dataset verbs (dataset/datasetd)
  - [X] keys
    - list the keys in a collection
    - at the package level keys returns a list of keys and an error value
  - [X] has_key
    - return "true "(w/OS exit 0 in CLI) if key is in collection,
      "false" otherwise (w/OS exit 1 in CLI)
  - [ ] sample
    - return a sample of keys from a collection
    - [ ] the newly create collections should have versioning disabled by default
  - [ ] create
    - add an new object to the collection if key does not exist,
      return false if object already exists or unable to create
      the new object
    - if versioning is enabled set the semver appropriately
  - [ ] read
    - return the object with nil error in the collection with the
      provided key, nil object and error value if not found
    - read always returns to the "current" object version
    - [ ] `read_versions()`, list the versions available for JSON object
    - [ ] `read_version()` list an JSON object for a specific version
      - return the object with nil error in the collection with the
        provided key and version, nil object and error value if not
        found
  - [ ] update
    - replace the object in the collection for given key, return false
      is object does not to replace or replacement fails
      - if collection has versioning turned on then version the object
    - [ ] `update()` update the current record respecting the version settings for collection
    - [X] delete
      - delete the object in the collection for given key, return true
        if deletion was successful, false if the object was not deleted
        (e.g. key not found or the collection is read only)
      - if collection has versioning turned on then delete **all objects**, if you want to revert you just update the object with the revised object values
      - [ ] `delete()` delete all versions of an object
      - If versioning is enabled the idea of "deleting" an object or attachment doesn't make sense, you only need to support Create, Read, Update and List, possibly with the ability to read versions available and retrieve the specific version, is this worth implementing in the CLI? Or is this just a lib dataset/package "feature"?
    - [ ] versioning, versioning is now set for the whole collection and effects JSON objects and their attachments (you're versioning both or neither), versioning will auto-increment for patch, minir and major semvere values if set
      - [ ] `set_versioning()`, set the versioning on a collection, the following strings enable
          versioning "major", "minor", "patch". Any other value disables
          versioning on the collection
      - [ ] `get_versioning()` on a colleciton (should return "major", "minor", "patch" or "")
    - [ ] Attachment support
      - [ ] `attach()` will add a basename file to the JSON object record, if versioning is enabled then it needs to handle the appropraite versioning setting
      - [ ] `attachments()` lists the attachments for a JSON object record
      - [ ] `attachment_versions()` list versions of a specific attachment
      - [ ] `detach()` retrieve "current" version of attachment
      - [ ] `detach_version()` retrieve a specific version of attachment
      - [ ] `prune()` remove all versions of attachments
    - [ ] Data Frame Support
      - [ ] frame_names
        - list the names of the frames currently defined in the collection
      - [ ] frame
        - define a new frame in the collection, if frame exists replace it
      - [ ] frame_meta
        - return the frame definition and metadata about the frame (e.g.  how many objects and attributes)
      - [ ] frame_objects
        - return the frame's list of objects
      - [ ] refresh
        - update all the objects in the frame based on current state of
          the collection
      - [ ] reframe
        - replace the frame definition but using the existing frame's keys
          refresh the frame with the new object describe
      - [ ] delete_frame
      - [ ] has_frame
- [ ] Verbs supported by cli only
  - [ ] set_versioning (accepts "", "patch", "minor", or "major" as values)
  - [ ] get_versioning (returns collection's version setting)
  - [ ] keys
  - [ ] create (if versioning is enable then handle versioning)
  - [ ] read (if versioning is enabled return the current version of an object)
  - [ ] update (if versioning is enable then handle versioning)
  - [ ] delete (if versioning is enable, delete all versions of object and attachments)
  - [ ] sample
  - [ ] clone
  - [ ] clone-sample
  - [ ] check
  - [ ] repair
  - [ ] frames (return a list of frames defined for collection)
  - [ ] frame, define a new frame in a collection
  - [ ] frame_objects, return the object list from a frame
  - [ ] refresh, refresh all objects in a frame based on the current state of the collection
  - [ ] reframe, replace the frame definition but using the existing frame's keys for the object listed in frame
  - [ ] delete_frame, remove a frame from the collection
  - [ ] has_frame return true if frame exists or false otherwise
  - [ ] attachments, list the attachments for a JSON object in the collection
  - [ ] attach, add an attachment to a JSON object in the collection, respect versioning if enabled
  - [ ] detach, retrieve an attachment from the JSON object in the collection
  - [ ] prune, remove attachments (including all versions) from an JSON object in the collection
- [ ] Add support for segmented key maps (to limit memory consumption for very large collections)
      settings in collection.json using keywords of patch, minor, major
- [ ] Auto-version attachments by patch, minor or major release per
- [ ] Need to add getting updated Man pages using the `dataset help ...` command
- [ ] Allow a WASM module to be used to validate objects in the collection. It needs to me integrate such that it "travels" will the dataset collection
  - this would let our JSON collections support explicit JSON structures as well as ad-hoc JSON objects
  - could use the YAML model approach in Newt to define the structures
- [ ] Document an example Python 3 http client support for web API implementing a drop in replacement for py_dataset using the web service or cli
- [ ] Implement a wrapping logger that takes a verboseness level for
      output (e.g. 0 - quiet, 1 progress messages, 2 warnings, errors
      should always show)
- [ ] Add support for https:// based datasets (in addition to local disc
      and s3://)
- [ ] dsbagit would generate a "BagIt" bag for preservation of collection
      objects
- [ ] dsgen would take a model described in YAML and generate HTML and browser side ES6 for quick prototyping with datasetd
- [ ] OAI-PMH importer to prototype iiif service based on Islandora
      content driven by a dataset collection
- [ ] Implement version support in the web service
- [ ] Implement an integrated a web UI for managing dataset collections and their data structures
  - [ ] Form pages could be expressed in Markdown+YAML for forms and embedded in the datasetd settings YAML file
    - See my notes on my text oreinted web experiment, yaml2webform.go
    - Forms could be render into the htdocs auto-magically saving development effort
    - The same forms could then be used server side for validation based on descriptors and JavaScript converted to WASM code
  - [ ] A standard JavaScript library could be used to knit the forms to the datasetd web service (sort of a mini-newt)
It would be nice if citesearch was defined by the citesearch.yaml file and some markdown documents taking a text oriented web approach to embedding forms in Markdown combined with some JS glue code to knit the two together
