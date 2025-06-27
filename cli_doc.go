// This is part of the dataset package.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrell, <tmorrell@library.caltech.edu>
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
package dataset

const (
	// cliDescription describes how to use the cli
	cliDescription = `
USAGE

 {app_name} [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PARAMETER ...]

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
- keys, returns a list of keys in the collection
- has-key, returnss true if key if found in collection, false otherwise
- codemeta (deprecated), copies metadata a codemeta file and updates the collections metadata
- info, returns the metadata associated with collection
- import (deprecated), imports another collecting into the current one
- export (deprecated), exports the collection into another collection
- attach, attaches a document to a JSON object record
- attachments, lists the attachments associated with a JSON object record
- retrieve, creates a copy local of an attachment in a JSON record
- prune, removes and attachment from a JSON record
- frames (deprecated), lists the frames defined in a collection
- frame (deprecated), will add a data frame to a collection if a definition is provided or return an existing frame if just the frame name is provided
- reframe (deprecated), will recreate a frame using its existing definition but replacing objects based on a new set of keys provided
- refresh (deprecated), will update all objects currently in the frame based on the current state of the collection. Any keys deleted in the collection will be delete from the frame.
- delete-frame (deprecated), will remove a frame from the collection
- has-frame (deprecated), will return true (exit 0) if frame exists, false (exit 1) if not
- attachments, will list any attachments for a JSON document
- attach, will add an attachment to a JSON document
- retrieve, will copy out the attachment to a JSON document 
  into the current directory 
- prune, will remove an attachment from the JSON document
- versions, will list the versions known for a JSON document if versioning is enabled for collection
- read-version (deprecated), will return a specific version of a JSON document if versioning is enabled for collection
- versioning,  will set the versioning of a collection, can be "none", "major", "minor", or "patch"

A word about "keys". {app_name} uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys must be lower case 
alpha numeric only.  Depending on storage engines there are issues
for keys with punctuation or that rely on case sensitivity. E.g. 
The pairtree storage engine relies on the host file system. File
systems are notorious for being picky about non-alpha numeric
characters and some are not case sensitive.

A word about "GLOBAL_OPTIONS" in v2 of dataset.  Originally
all options came after the command name, now they tend to
come after the verb itself. This is because context counts
in trying to remember options (at least for the authors of
dataset).  The are three "GLOBAL_OPTIONS" that are exception
and they are ` + "`" + `-version` + "`, " + "`" + `-help` + "`" + `
and ` + "`" + `-license` + "`" + `. All other options come
after the verb and apply to the specific action the verb
implements.

`

	// cliExamples lists some examples of using the cli
	cliExamples = `
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

~~~shell
    cat JSON_DOCNAME | {app_name} create COLLECTION_NAME KEY
    {app_name} create -i JSON_DOCNAME COLLECTION_NAME KEY
    {app_name} create COLLECTION_NAME KEY JSON_VALUE
    {app_name} create COLLECTION_NAME KEY JSON_FILENAME
~~~

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

~~~shell
    cat jane-doe.json | {app_name} create people.ds r1
    {app_name} create -i blob.json people.ds r1
    {app_name} create people.ds r1 '{"name":"Jane Doe"}'
    {app_name} create people.ds r1 jane-doe.json
~~~

`

	cliRead = `
read
====

Syntax
------

~~~shell
    {app_name} read COLLECTION_NAME KEY
~~~

Description
-----------

The writes the JSON document to standard out (unless you've 
specific an alternative location with the "-output" option)
for the given KEY.

Usage
-----

An example we're assuming there is a JSON document with a KEY 
of "r1". Our collection name is "data.ds"

~~~shell
    {app_name} read data.ds r1
~~~

Options
-------

Normally {app_name} adds two values when it stores an object, ` + "`._Key`" + `
and possibly ` + "`._Attachments`" + `. You can get the object without these
added attributes by using the ` + "`-c` or `-clean`" + ` option.


~~~shell
    {app_name} read -clean data.ds r1
~~~

`

	cliUpdate = `
update
======

Syntax
------

~~~shell
    {app_name} update COLLECTION_NAME KEY
~~~

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

~~~shell
    {app_name} update people.ds jane.doe '{"name":"Jane Doiel"}'
    {app_name} update people.ds jane.doe jane-doe.json
    {app_name} update -i jane-doe.json people.ds jane.doe
    cat jane-doe.json | {app_name} update people.ds jane.doe
~~~

`

	cliDelete = `
delete
======

Syntax
------

~~~shell
    {app_name} delete COLLECTION_NAME KEY
~~~

Description
-----------

- delete - removes a JSON document from collection
  - requires JSON document name

Usage
-----

This usage example will delete the JSON document withe the key _r1_ in 
the collection named "publications.ds".

~~~shell
    {app_name} delete publications.ds r1
~~~

`

	cliKeys = `
keys
====

Syntax
------

~~~shell
    {app_name} keys COLLECTION_NAME
~~~

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

~~~shell
    {app_name} keys COLLECTION_NAME
    {app_name} keys COLLECTION_NAME | sort
    {app_name} keys COLLECTION_NAME | sort -n
~~~

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
has-key
======

Syntax
------

~~~shell
    {app_name} [OPTIONS] has-key COLLECTION_NAME KEY_TO_CHECK_FOR
~~~

Description
-----------

Checks if a given key is in the a collection. Returns "true" if 
found, "false" otherwise. The collection name is "people.ds"

Usage
-----

~~~shell
    {app_name} has-key people.ds '0000-0003-0900-6903'
    {app_name} has-key people.ds r1
~~~

`
	cliUpdatedKeys = `
updated-keys
============

(deprecated)

Syntax
------

~~~shell
    {app_name} update-keys COLLECTION_NAME START END
~~~

Description
-----------

List the JSON_DOCUMENT_ID available in a collection. Key order is not
guaranteed. Keys are forced to lower case when the record is created
in the {app_name} (as of version 1.0.2). Note combining "keys" with
a pipe and POSIX commands like "sort" can given a rich pallet of
ways to work with your {app_name} collection's keys.

Example
-------

Here is an example usage for select updates keys for record
created or update between a start and end time (inclusive).
The times are in the form of "YYYY-MM-DD HH:MM:SS" and are required.
The hours are in 24 hour notation.  The resulting keys are sorted
in ascending updated timestamp order.

~~~shell
    {app_name} updated-keys COLLECTION_NAME \
	           "2022-01-01 00:00:00"
	           "2022-12-31 23:23:59"
~~~

`

	cliCount = `
count
=====

Syntax
------

~~~shell
    {app_name} count COLLECTION_NAME [FILTER EXPRESSION]
~~~

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

~~~shell
    {app_name} count "publications.ds"
~~~

Count records where the ".published" field is true.

~~~shell
    {app_name} count "publications.ds" '(eq .published true)'
~~~

`

	cliInit = `
init
====

Syntax
------

~~~shell
    dataset init COLLECTION_NAME [DSN_URI|pairtree]
~~~

Description
-----------

_init_ creates a collection. Collections are created on local 
disc. By default it uses a SQLite3 database called "collection.db"
in the dataset directory for storing JSON Objects.

Usage
-----

The following example command create a dataset collection 
named "data.ds".

~~~shell
    dataset init data.ds
~~~

NOTE: After each evocation of ` + "`" + `dataset init` + "`" + ` if all went well 
you will be shown an ` + "`" + `OK` + "`" + ` if everything went OK, otherwise
an error message. 

You can still use a pairtree to store your dataset collection by using
the "pairtree" option when creating your collection. Example for creating
data.ds using a pairtree store.

~~~shell
    dataset init data.ds pairtree
~~~

By default dataset cli creates a collection and uses an SQLite3 database to
store the JSON objects.

~~~shell
    dataset init data.ds
~~~

You can optionally store your JSON documents in a
SQL in PostgreSQL or MySQL 8. This can improve performance for large collections
as well as support multi-user or multi-process concurrent use of a collection.

The SQL storage engine is designated using a "DSN_URI". The DSN_URI is formed
by setting the "protocol" of the URL to either "sqlite://", "mysql://",
or "postgres://" followed by a DSN or "data source name" The data source name
is usually a form of the database connection string. e.g. something like
"posgres://jane.doe:$SECRET@collections". The DSN is described in detail in the Go
database/sql package.

This examples shows using SQLite3 storage for the JSON documents in
a "collection.db" stored inside the "data.ds" collection.

~~~shell
    dataset init data.ds "sqlite://collection.db"
~~~

Here's a variation using MySQL 8 as the storage engine storing the
collection in the "collections" database.

~~~shell
    dataset init data.ds "mysql://$DB_USER:$DB_PASSWORD@/collections"
~~~

Here's an example of for using PostgreSQL running on localhost. 

~~~shell
    dataset init data.ds "postgres://$DB_USER:$DB_PASSWORD@localhost/collections"
~~~

If you need to pass the sslmode parameters to use PostgreSQL include them
in your DSN. E.g. disabling SSL when using localhost.

~~~shell
    dataset init data.ds "postgres://$DB_USER:$DB_PASSWORD@localhost/collections?sslmode=disable"
~~~

NOTE: If you are using a SQL based collection and you want to copy it
to a new collection use dataset's clone action with the '-all' option.


`

	cliVersioning = `
versioning
==========

Collections can support a simplistic form of versioning for JSON documents
and their attachments.  It is a collection wide setting and if enabled
JSON documents and attachments will associated with a semver (semantic
version number). The implementation details are based on the storage engine.

The versioning can be set to increment on the patch, minor or major 
semver values creating or updating a JSON document or attachment.  The 
value before creation is assumed to be "0.0.0". If versioning is enabled
it is automatically applied. 

Directly working with versioned documents or attachments requires writing
programs and using the Go dataset package or libdataset C-shared library.

Examples
--------

This example shows how to create a collection (versioning is turned
off by default). Setting patch level versioning, showing the versioning
setting and repeat for "minor", "major" and turning off versioning in the
collection.

~~~shell
   CNAME="mycollection.ds"
   {app_name} init $CNAME
   {app_name} set_versioning $CNAME patch
   {app_name} get_versioning $CNAME
   {app_name} set_versioning $CNAME minor
   {app_name} get_versioning $CNAME
   {app_name} set_versioning $CNAME major
   {app_name} get_versioning $CNAME
   {app_name} set_versioning $CNAME none
   {app_name} get_versioning $CNAME
~~~

`

	cliClone = `
clone
=====

(deprecated)

Clone a collection from a list of keys into a new collection. If the
` + "`" + `-all` + "`" + ` option is used then the whole collection will be
clone and a list of keys is not required.

In this example we create a list of keys using the ` + "`-sample`" + ` option
and then clone those keys into a new collection called *sample.ds*.

~~~shell
    {app_name} keys -sample=3 mycollection.ds > sample.keys
    {app_name} clone -i sample.keys mycollection.ds sample.ds
~~~

NOTE: If you are using SQL based collections clone is your tool
to copy the collection to a new collection rather than using your
operating system's version of copying.

`

	cliCloneSample = `
clone-sample
============

(deprecated)

Clone a collection into a sample size based training collection 
and test collection.

In this example we create a training and testing collections 
based on a training sample size of 1000.

~~~shell
    {app_name} clone-sample -size=1000 mycollection.ds training.ds test.ds
~~~

`

	cliFrames = `
frames
======

(deprecated)

Lists the frames available in a collection. In this example our
collection name is ` + "`pubs.ds`" + `.

~~~shell
   {app_name} frames pubs.ds
~~~

`

	cliFrame = `
frame
=====

(deprecated)

This command will define a data frame or return the contents and
metadata of a defined frame. To define a new frame you need to provide a
collection name, a frame name followed by a list of dotpath/label pairs.
The labels are used as object attribute names and the dot paths as the
source of data. You also need a list of keys.\
By default the keys are read from standard input. With options you can
include a specific file or even indicate to use all the keys in a
collection. In this example we are creating a frame called
\"title-authors-year\" based on the titles, authors and publication year
from a dataset collection called ` + "`pubs.ds`" + `. Note the labels of
\"Title\", \"Authors\", \"PubYear\" are on the right side the an equal
sign and the dot paths to the left.

~~~shell
    {app_name} keys pubs.ds |\
        {app_name} frame pubs.ds "title-authors-year" \
                ".title=Title" \
                ".authors=Authors" \
                ".publication_year=PubYear"
~~~

The objects in the frame\'s object list will look like

~~~json
    {
        "Title": ...,
        "Authors": ...,
        "PubYear": ...,
    }
~~~

This allows you to create convenient names for otherwise deep dot paths.

`

	cliFrameObjects = `
frame-objects
=============

(deprecated)

Usage
-----

~~~shell
    frame-objects COLLECTION FRAME_NAME
~~~

Returns the object list of a frame.

OPTIONS
-------

-p, -pretty
: pretty print JSON output

Example
-------

If I want to get a list of objects (JSON array of objects) 
for a frame named "captions-dates-locations" from my collection
called "photos.ds" I would do the following (will be using the
` + "`-p`" + ` option to pretty print the results)

~~~shell
    {app_name} frame-objects -p photos.ds captions-dates-locations
~~~

`

	cliReframe = `
reframe
=======

(deprecated)

This command replace the current keys/object list in a frame based
on the frame's current definition.

In the following example the frame name is \"f1\", the collection is
\"examples.ds\". The first example is reframing an existing frame using
existing keys coming from standard input, the second example performs
the same thing but is taking a filename to retrieve the list of keys.

~~~shell
    cat f1-updated.keys | {app_name} reframe example.ds f1
    {app_name} reframe example.ds f1 f1-updated.keys
~~~

`

	cliRefresh = `
refresh
=======

(deprecated)

Update the objects in a frame based on it's current set of keys and definition.  

NOTE: If any keys have been deleted from the collection then the object
associated with those keys in the frame will also be removed.

In the following example the frame name is \"f1\", the collection is
\"examples.ds\". The example is refreshing the object list.

~~~shell
    {app_name} refresh example.ds f1
~~~

`

	cliDeleteFrame = `
delete-frame
============

(deprecated)

This is used to removed a frame from a collection.

~~~shell
    {app_name} delete-frame example.ds f1
~~~

delete frame f1 from collection called example.ds

`

	cliAttachments = `
attachments
===========

Syntax
------

~~~shell
    {app_name} attachments COLLECTION_NAME KEY
~~~

Description
-----------

List the files attached to the JSON record matching the KEY
in the collection.

Usage
-----

List all the attachments for _k1_ in collection "stats.ds".

~~~shell
    {app_name} attachments stats.ds k1
~~~

`

	cliAttach = `
attach
======

Syntax
------

~~~shell
    {app_name} attach COLLECTION_NAME KEY [SEMVER] FILENAME(S)
~~~

Description
-----------

Attach a file to a JSON record. Attachments are stored in a tar ball
related to the JSON record key.

Usage
-----

Attaching a file named *start.xlsx* to the JSON record with id _t1_ in 
collection "stats.ds"

~~~shell
    {app_name} attach stats.ds t1 start.xlsx
~~~

Attaching the file as version v0.0.1

~~~shell
    {app_name} attach stats.ds t1 v0.0.1 start.xlsx
~~~

`

	cliRetrieve = `
retrieve
========

Syntax
------

~~~shell
    {app_name} retrieve COLLECTION_NAME KEY [SEMVER]
    {app_name} retrieve COLLECTION_NAME KEY [SEMVER] ATTACHMENT_NAME
~~~

Description
-----------

__retrieve__ writes out (to local disc) the items that have been 
attached to a JSON record in the collection with the matching KEY

Usage
-----

Write out all the attached files for k1 in collection named 
"publications.ds"

~~~shell
    {app_name} retrieve publications.ds k1
~~~

Write out only the *stats.xlsx* file attached to k1

~~~shell
    {app_name} retrieve publications.ds k1 stats.xlsx
~~~

Write out only the v0.0.1 *stats.xlsx* file attached to k1

~~~shell
    {app_name} retrieve publications.ds k1 v0.0.1 stats.xlsx
~~~

`

	cliPrune = `
prune
=====

Syntax
------

~~~shell
    {app_name} prune COLLECTION_NAME KEY [SEMVER]
    {app_name} prune COLLECTION_NAME KEY [SEMVER] ATTACHMENT_NAME
~~~

Description
-----------

prune removes all or specific attachments to a JSON document. If only
the key is supplied then all attachments are removed if an attachment
name is supplied then only the specific attachment is removed.

Usage
-----

In the following examples _r1_ is the KEY, *stats.xlsx* is the 
attached file. In the first example only *stats.xlsx* is removed in
the second all attachments are removed. Our collection name is "data.ds"


~~~shell
    {app_name} prune data.ds k1 v0.0.1 stats.xlsx
    {app_name} prune data.ds k1 stats.xlsx
    {app_name} prune data.ds k1
~~~

`
	cliSample = `
sample
======

(deprecated)

{app_name} supports the concept of generating a random sample
of keys from a collection. To do this you need to use the ` + "`sample`" + `
verb. ` + "`sample`" + ` expects the collection name followed by an
a positive integer value "N". It returns a randomly selected number of
keys.  If N is greater than the collection then all keys are returned
for the collection.

~~~shell
    {app_name} sample data.ds 100
~~~

`

	cliFrameDef = `
frame-def
=========

(deprecated)

{app_name} supports the concept of frames and the "frame-def" verb
lets you review the definition of an existing frame.

~~~shell
    {app_name} frame-def data.ds myframe
~~~shell

`

	cliHasFrame = `
has-frame
=========

(deprecated)

{app_name} supports the concept of frames and the "has-frame" verb
lets you check if a frame exists.

~~~shell
    {app_name} has-frame data.ds myframe
~~~shell

`

	cliCheck = `
check
=====

(deprecated)

syntax
------

~~~shell
    dataset check COLLECTION_NAME [COLLECTION_NAME ...]
~~~

Description
-----------

Check reviews one or more collections and reports if any problems 
are identified based on the ` + "`collection.json`" + ` file found in the 
folder holding the collection's pairtree. 

Usage
-----

If multiple instances of dataset write (e.g. create or update) to 
a collection then it is possible that the JSON file ` + "`collection.json`" + `
will become inaccurate.

~~~shell
    dataset check MyRecordCollection.ds
    dataset check MyBrokenCollection.ds MyRecordCollection.ds
~~~

`

	cliRepair = `
repair
======

(deprecated)

Syntax
------

~~~shell
    dataset repair COLLECTION_NAME
~~~

Description
-----------

_repair_ trys to repair a collection correcting as best it can 
the ` + "`collection.json`" + ` file defining where things are to be found.

Usage
-----

Our collection name is "MyCollectiond.ds".

~~~shell
   dataset repair MyCollection.ds
~~~

`

	cliCodemeta = `
codemeta
========

(deprecated)

The command imports a codemeta.json file into the collection replacing
it's existing metadata.

~~~shell
   {app_name} codemeta data.ds ./codemeta.json
~~~

Without the codemeta filename it returns the existing codemeta values.

`

	cliFrameKeys = `
frame-keys
==========

(deprecated)

This returns a list of keys associated with the frame. Keys are
returned one per line.


~~~shell
    {app_name} frame-keys data.ds
~~~

`

	cliMigrate = `
migrate
=======

(deprecated)

This will migrate content from a v1 dataset collection to a
v2 dataset collection.  Before migrating you need to create an
empty destination collection.

NOTE: attachments are not currently
migrated, just the JSON documents.

~~~shell
    {app_name} init new_collection.ds
    {app_name} migrate -verbose old_collection.ds new_collection.ds
~~~

`

	cliModel = `
# NAME

model

(experimental)

# SYNOPSIS

The model verb let's you describe a simple data model for validating objects an HTML web form
or the datasetd JSON API.  It currently doesn't validate JSON data are ready in a collection
or when the JSON data is added via the dataset cli.

The model is saved in the root the dataset collection directory. The file
is named "model.yaml" and is indeed a simple YAML file.

# DESCRIPTION

__model__ describes the JSON object shape for objects in the collection.  It is based on the
Go package "github.com/caltechlibrary/models". A model is expressed in YAML. It can be used to
validate the JSON data, render simple HTML as well as SQLite 3 SQL schema.

The __model__ is read when the collection is opened.  For datasetd this means the model is read
at startup.  The model isn't re-read until datasetd is restarted.

Models that can be described by a web form are shallow, i.e. web forms do not support lists,
arrays or dictionaries.

## A Model

id
: (required) The identifier for the model. Is the "id" given to the generated HTML web form.

description
: (optional) This is simple description of the model. It will be included as a comment in
generated HTML and the SQLite3 SQL. This is a text string or block.

elements
: This is a list of elements that describe the data attributes of your model. At least
one element should be defined otherwise your model does not make sense.

## Elements

The elements attribute holds a list of elements. You can think of these as HTML5 form
elements described in YAML. They will also be used to infer SQLite 3 column types.

Each element is made from the following properties.

type
: (required) This is a string and maps to the input element types available in HTML5[^1]. 

id
: (optional) This is the element's identifier. It should be unique with in the model. While
optional it is used to retrieve an element from a model. If is also required when rendering
column definitions in SQLite 3. A model that includes a submit or reset button would examples
of when to leave it blank.

attributes
: (optional) This is a list of key/value pairs that map to HTML5 input elements. Boolean HTML
element attributes like "required" and "checked" you are expressed as `+"`"+`required: true`+"`"+` and
`+"`"+`checked: true`+"`"+` in YAML. NOTE: attributes value's are resolved to quoted strings
when rendered as HTML.

pattern
: (optional) This is a regular expression pattern that is used to validate the input of
the element[^2].

options
: (optional) Are a list of key/value maps used to expression HTML5 select elements. They
can be be used in validation of a model's content as well as in render HTML selection
elements.

is\_primary\_id
: (optional) If set to true it indicates a given element holds the model's primary identifier.
If you are store model content in a SQLite 3 database or Dataset collection this would be
the unique identifier used to retrieve the modeled object.

label
: (optional) If set it is used as the text content of the label when rendering a web form.

[^1]: See <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input> for details.

[^2]: See <https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/pattern>
for details of how patterns are used in validation.

## Data type support in models package

The models package starts from the premise of supporting a YAML description of a web form
that then can be used to render HTML and SQL Schema. It also needs to be able to be a thin
layer in a Go API that can validate the elements of a model just like they are validated
browser side by the [HTML5 input types](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input).
The following are all implemented by the models package using a naive validation approach[^3].  

[^3]: "naive" in this case means overly simplistic validation, e.g. min max ranges don't
validate against step attributes. 


- [button](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/button)
- [checkbox](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/checkbox)
- [color](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/color)
- [date](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/date)
- [datetime-local](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/datetime-local)
- [email](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/email)
- [hidden](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/hidden)
- [image](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/image)
- [month](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/month)
- [number](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/number)
- [password](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/password)
- [radio](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/radio)
- [range](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/range)
- [reset](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/reset)
- [search](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/search)
- [submit](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/submit)
- [tel](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/tel)
- [text](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/text)
- [time](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/time)
- [url](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/url)
- [week](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/week)

Additional data types[^4] can be defined by using the `+"`"+`Model.Define`+"`"+` function provided in
this package. You need to provide a name for the new type as well as the func's name. The
"defined" data types are applied before the default types. This allows for improvements to
the defaults while retaining a fallback. Hopefully this mechanism can prove useful to expanding
the data types supported by models.

[^4]: The validation function is used server side only because it is written in Go. E.g. by
Dataset's JSON API.

NOTE: As the models package evolves the validation methods provided out of the box will
evolve too. Some may even be dropped if they prove problematic[^5].

[^5]: E.g. "week" input type is not widely used and is poorly supported by browsers in 2024.
"image" doesn't make a whole lot of sense.

# Example

This is an example model of a guest book entry used in a Dataset base guest book web application.

~~~yaml
id: test_model
elements:
  - id: id
    type: text
    attributes:
      name: id
      placeholder: Enter a unique string
      required: true
    is_primary_id: true
  - id: name
    type: text
    attributes:
      name: name
      placeholdertext: Enter your name
      required: true
  - id: msg
    type: textarea
    attributes:
      name: msg
      placehodertext: Enter a short message or comment
~~~

`

cliLoad = `load [OPTION]
============

This will read a JSONL stream (JSONL, see https://jsonlines.org) 
and store the objects in a dataset collection. The objects must have
two attributes, __key__ and __object__.  The key value is used as
the key assigned in the collection and object value is the JSON stored.
Load reads from standard input.

The object structure matches the schema used by dataset collections
that are using SQLite3 for their object store. The dataset collection
loading the objects can use any dataset collection storage format
supported by the version of dataset featuring load and dump verbs.

# OPTION

-o, -overwrite
: If an object exists in the collection with the same key replace it.

-m, -max-capacity INTEGER
: Objects can be large in JSONL so you have the option of setting the
maximum buffer size for a single object. The integer value should be
greater than zero. The unit value is measured in mega bytes. "1" is
one meta byte, "10" would be ten mega bytes.

# EXAMPLE

Load a JSONL file, duplicate objects will not be overwritten.

~~~shell
    {app_name} load mycollection.ds <mycollection.jsonl
~~~

Load a JSONL file, duplicate objects will be overwritten.

~~~shell
    {app_name} load -overwrite mycollection.ds <mycollection.jsonl
~~~

`

cliDump = `dump
============

This will dump all the JSON objects in a collection, one
object per line (see https://jsonlines.org) (aka JSONL). 

The objects are written to standard output. Dump is the complement of
load verb. The objects dumped reflect the structured using when storing
objects in an SQLite3 database regardless of the store of the specific
collection. Like clone it provides a means of easily moving your data out
of a dataset collection.

Example
-------

~~~shell
    {app_name} dump mycollection.ds >mycollection.jsonl
~~~

`

cliJoin = `join
============

(deprecated)

This will take a collection name, key and new JSON expression.
the object in the collection matching the key will be "joined"
with the new object.  If the overwrite option is set to true
then the JSON expressions attributes will overwrite the ones
in the collection. If false only the new attributes will be
added to the object in the collection.

~~~shell
    {app_name} join mycollection.ds my-key '{"one-o-one": 101}'
    {app_name} join -overwrite mycollection.ds my-key '{"one-o-one": 202}'
~~~

`

cliQuery = `query
============

This will run a SQL query against the SQLite3 JSON store. It returns
a list of JSON objects or an error. The SQL query must only return
one object per row.

Example
-------

Return a JSON array of all objects by descending created date.

~~~shell
    {app_name} query mycollection.ds \\
      "select src from mycollection order by created desc"
~~~

Read the SQL statement from a file called "report.sql".

~~~shell
    {app_name} query -sql report.sql mycollection.ds
~~~

If you are using PostgreSQL or MySQL as your SQL storage engine
then you'll need to adapt these examples to the SQL dialect for
used in those sytems.

`

)


