//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2022, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package dataset

const (
	CLIDescription = `
USAGE

   {app_name} [OPTIONS] VERB COLLECTION_NAME [PRAMETER ...]

SYNOPSIS

{app_name} command line interface supports creating JSON object
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
- frame, will add a data frame to a collection if a definition is provided or return an existing frame if just the frame name is provided
- reframe, will recreate a frame based on the current state of objects in the collection, if keys are provide with the reframe request then the objects in the frame will be replaces by objects associated with the new keys provided
- delete-frame, will remove a frame from the collection
- has-frame, will return true (exit 0) if frame exists, false (exit 1) if not
- attachments, will list any attachments for a JSON document
- attach, will add an attachment to a JSON document
- retrieve, will copy out the attachment to a JSON document 
  into the current directory 
- prune, will remove an attachment from the JSON document
- versions, will list the versions known for a JSON document if versioning is enabled for collection
- read-version, will return a specific version of a JSON document if versioning is enabled for collection
- versioning,  will set the versioning of a collection, can be "none", "major", "minor", or "patch"

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
    cat JSON_DOCNAME | {app_name} create COLLECTION_NAME KEY
    {app_name} create -i JSON_DOCNAME COLLECTION_NAME KEY
    {app_name} create COLLECTION_NAME KEY JSON_VALUE
    {app_name} create COLLECTION_NAME KEY JSON_FILENAME
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
    cat jane-doe.json | {app_name} create people.ds r1
    {app_name} create -i blob.json people.ds r1
    {app_name} create people.ds r1 '{"name":"Jane Doe"}'
    {app_name} create people.ds r1 jane-doe.json
` + "```" + `

Related topics: [update](update.html), [read](read.html), and [delete](delete.html)

`

	cliRead = `
read
====

Syntax
------

` + "```" + `shell
    {app_name} read COLLECTION_NAME KEY
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
    {app_name} read data.ds r1
` + "```" + `

Options
-------

Normally {app_name} adds two values when it stores an object, ` + "`._Key`" + `
and possibly ` + "`._Attachments`" + `. You can get the object without these
added attributes by using the ` + "`-c` or `-clean`" + ` option.


` + "```" + `shell
    {app_name} read -clean data.ds r1
` + "```" + `

`

	cliUpdate = `
update
======

Syntax
------

` + "```" + `shell
    {app_name} update COLLECTION_NAME KEY
` + "```" + `

Description
-----------

_update_ will replace a JSON document in a {app_name} collection for 
a given KEY.  By default the JSON document is read from standard 
input but you can specific a specific file with the "-input" 
option. The JSON document should already exist in the collection
when you use update.


Usage
------

In this example we assume there is a JSON document on local disc 
named _jane-doe.json_. It contains ` + "`{\"name\":\"Jane Doe\"}`" + ` and the 
KEY is "jane.doe". In the first one we specify the full JSON document 
via the command line after the KEY.  In the second example we read the 
data from _jane-doe.json_. Finally in the last we read the JSON 
document from standard input and save the update to "jane.doe".
The collection name is "people.ds".

` + "```" + `shell
    {app_name} update people.ds jane.doe '{"name":"Jane Doiel"}'
    {app_name} update people.ds jane.doe jane-doe.json
    {app_name} update -i jane-doe.json people.ds jane.doe
    cat jane-doe.json | {app_name} update people.ds jane.doe
` + "```" + `

`

	cliDelete = `
delete
======

Syntax
------

` + "```" + `shell
    {app_name} delete COLLECTION_NAME KEY
` + "```" + `

Description
-----------

- delete - removes a JSON document from collection
  - requires JSON document name

Usage
-----

This usage example will delete the JSON document withe the key _r1_ in 
the collection named "publications.ds".

` + "```" + `shell
    {app_name} delete publications.ds r1
` + "```" + `

`

	cliKeys = `
keys
====

Syntax
------

` + "```" + `shell
    {app_name} keys COLLECTION_NAME
` + "```" + `

Description
-----------

List the JSON_DOCUMENT_ID available in a collection. Key order is not
guaranteed. Keys are forced to lower case when the record is created
in the {app_name} (as of version 1.0.2). Note combining "keys" with
a pipe and POSIX commands like "sort" can given a rich pallet of
ways to work with your {app_name} collection's keys.

Examples
--------

Here are three examples usage. Notice the sorting is handled by
the POSIX sort command which lets you sort ascending or descending
including sorting number strings.

` + "```" + `shell
    {app_name} keys COLLECTION_NAME
    {app_name} keys COLLECTION_NAME | sort
    {app_name} keys COLLECTION_NAME | sort -n
` + "```" + `

Getting a "sample" of keys
--------------------------

The __{app_name}__ command respects an option named ` + "`-sample N`" + ` where N 
is the size (number) of the keys to include in the sample. The sample 
is taken after any filters are applied but may be less than requested 
size if the the filtered results are few than the sample size.  The 
basic process is to get a set of keys, randomly sort the keys, then 
return the top N number of those keys.

`

	cliHasKey = `
haskey
======

Syntax
------

` + "```" + `shell
    {app_name} [OPTIONS] haskey COLLECTION_NAME KEY_TO_CHECK_FOR
` + "```" + `

Description
-----------

Checks if a given key is in the a collection. Returns "true" if 
found, "false" otherwise. The collection name is "people.ds"

Usage
-----

` + "```" + `shell
    {app_name} haskey people.ds '0000-0003-0900-6903'
    {app_name} haskey people.ds r1
` + "```" + `

In python

` + "```" + `shell
    {app_name}.has_key('people.ds', '0000-0003-0900-6903')
    {app_name}.has_key('people.ds', 'r1')
` + "```" + `

`

	cliCount = `
count
=====

Syntax
------

` + "```" + `shell
    {app_name} count COLLECTION_NAME [FILTER EXPRESSION]
` + "```" + `

Description
-----------

This returns a count of the keys in the collection. It is reasonable 
quick as only the collection metadata is read in. *count* also can 
accept a filter expression. This is slower as it iterates over all 
the records and counts those which evaluate to true based on the
filter expression provided.

Usage
-----

Count all records in collection "publications.ds"

` + "```" + `shell
    {app_name} count "publications.ds"
` + "```" + `

Count records where the ".published" field is true.

` + "```" + `shell
    {app_name} count "publications.ds" '(eq .published true)'
` + "```" + `

`

	cliInit = `
init
====

Syntax
------

` + "```" + `shell
    {app_name} init COLLECTION_NAME
` + "```" + `

Description
-----------

_init_ creates a collection. Collections are created on local 
disc.

Usage
-----

The following example command create a dataset collection 
named "data.ds".

` + "```" + `shell
    {app_name} init data.ds
` + "```" + `

NOTE: After each evocation of ` + "`{app_name} init`" + ` if all went well 
you will be shown an ` + "`OK`" + ` if everything went OK, otherwise
an error message. 

`

	cliVersioning = `
versioning
==========

Collections can support a simplistic form of versioning for JSON documents
and their attachments.  It is a collection wide setting and if enabled
JSON documents and attachments will associated with a semver (symantic
version number). The implementation details are based on the storage engine.

The versioning can be set to increment on the patch, minor or major 
semver values creating or updating a JSON document or attachment.  The 
value before creation is assumed to be "0.0.0". 

To return a list of JSON documents versions stored use the "versions"
verb. This will return a list of version numbers available for a given
key. To read a specific version of a JSON document use the "read-version"
verb which takes a key and a semver version string.
create a copy of the JSON document or attachment using a semver 
then and

` + "```" + `shell
  # List versions in the data.ds collection for key "123"
  {app_name} versions data.ds 123
  # To read a specific vesion, e.g. 0.0.3
  {app_name} version data.ds 123 0.0.3
` + "```" + `

`
)
