package dataset

const (
	CLIDescription = `
USAGE

   {app_name} [OPTIONS] VERB COLLECTION_NAME [PRAMETER ...]

SYNOPSIS

dataset command line interface supports creating JSON object
collections and managing the JSON object documents in a collection.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

SUPPORTED VERBS

- help will give this this documentation of help on a verb
- create, creates a new document in the collection
- read, retrieves a document from the collection writing it standard out
- update, updates a document in the collection
- delete, removes a document from the collection
- list, returns a list of keys in the collection
- codemeta, copies metadata a codemeta file and updates the collections metadata
- info, returns the metadata associated with collection
- import, imports another collecting into the current one
- export, exports the collection into another collection
- attach, attaches a document to a JSON object record
- attachments, lists the attachments associated with a JSON object record
- retrieve, creates a copy local of an attachement in a JSON record
- prune, removes and attachment from a JSON record
- frames, lists the frames defined in a collection
- frame will add a data frame to a collection if a definition is provided or return an existing frame if just the frame name is provided
- reframe will recreate a frame based on the current state of objects in the collection, if keys are provide with the reframe request then the objects in the frame will be replaces by objects associated with the new keys provided
- delete-frame will remove a frame from the collection
- has-frame will return true (exit 0) if frame exists, false (exit 1) if not
- attachments will list any attachments for a JSON document
- attach will add an attachment to a JSON document
- retrieve will copy out the attachment to a JSON document 
  into the current directory 
- prune will remove an attachment from the JSON document

You can get additional help 

`

	CLIExamples = `
EXAMPLES

   {app_name} help init

   {app_name} init my_objects.ds 

   {app_name} help create

   {app_name} create my_objects.ds "123" '{"one": 1}'

   {app_name} create my_objects.ds "234" mydata.json 
   
   cat <<EOT | {app_name} create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   {app_name} update my_objects.ds "123" '{"one": 1, "two": 2}'

   {app_name} delete my_objects.ds "345"

   {app_name} keys my_objects.ds

`

	//
	// cli specific help, not exported
	//

	// Taken from docs/create.md
	cliCreate = `
create
======

Syntax
------

` + "```" + `shell
    cat JSON_DOCNAME | dataset create COLLECTION_NAME KEY
    dataset create -i JSON_DOCNAME COLLECTION_NAME KEY
    dataset create COLLECTION_NAME KEY JSON_VALUE
    dataset create COLLECTION_NAME KEY JSON_FILENAME
` + "```" + `

Description
-----------

create adds or replaces a JSON document to a collection. The JSON 
document can be read from a standard in, a named file (with a 
".json" file extension) or expressed literally on the command line.

Usage
-----

In the following four examples *jane-doe.json* is a file on the 
local file system contains JSON data containing the JSON_VALUE 
of ` + "`" + `{"name":"Jane Doe"}` + "`" + `.  The KEY we will create is _r1_. 
Collection is "people.ds".  The following are equivalent in 
resulting record.

` + "```" + `shell
    cat jane-doe.json | dataset create people.ds r1
    dataset create -i blob.json people.ds r1
    dataset create people.ds r1 '{"name":"Jane Doe"}'
    dataset create people.ds r1 jane-doe.json
` + "```" + `

Related topics: [update](update.html), [read](read.html), and [delete](delete.html)

`

	cliRead = `
read
====

Syntax
------

` + "```" + `shell
    dataset read COLLECTION_NAME KEY
` + "```" + `

Description
-----------

The writes the JSON document to standard out (unless you've 
specific an alternative location with the "-output" option)
for the given KEY.

Usage
-----

An example we're assuming there is a JSON document with a KEY 
of "r1". Our collection name is "data.ds"

` + "```" + `shell
    dataset read data.ds r1
` + "```" + `

Options
-------

Normally dataset adds two values when it stores an object, ` + "`._Key`" + `
and possibly ` + "`._Attachments`" + `. You can get the object without these
added attributes by using the ` + "`-c` or `-clean`" + ` option.


` + "```" + `shell
    dataset read -clean data.ds r1
` + "```" + `

`
)
