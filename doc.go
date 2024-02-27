/*
Package dataset includes the operations needed for processing collections of JSON documents and their attachments.

Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>

Copyright (c) 2022, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Dataset Project
===============

The Dataset Project provides tools for working with collections of
JSON Object documents stored on the local file system or via a dataset
web service.  Two tools are provided, a command line interface (dataset)
and a web service (datasetd).

dataset command line tool
-------------------------

_dataset_ is a command line tool for working with collections of JSON
objects.  Collections are stored on the file system in a pairtree
directory structure or can be accessed via dataset's web service.
For collections storing data in a pairtree JSON objects are stored in
collections as plain UTF-8 text files. This means the objects can be
accessed with common Unix text processing tools as well as most
programming languages.

The _dataset_ command line tool supports common data management operations
such as initialization of collections; document creation, reading,
updating and deleting; listing keys of JSON objects in the collection;
and associating non-JSON documents (attachments) with specific JSON
documents in the collection.

### enhanced features include

- aggregate objects into data frames
- generate sample sets of keys and objects

datasetd, dataset as a web service
----------------------------------

_datasetd_ is a web service implementation of the _dataset_ command line
program. It features a sub-set of capability found in the command line
tool. This allows dataset collections to be integrated safely into web
applications or used concurrently by multiple processes. It achieves this
by storing the dataset collection in a SQL database using JSON columns.

Design choices
--------------

_dataset_ and _datasetd_ are intended to be simple tools for managing
collections JSON object documents in a predictable structured way.

_dataset_ is guided by the idea that you should be able to work with
JSON documents as easily as you can any plain text document on the Unix
command line. _dataset_ is intended to be simple to use with minimal
setup (e.g.  `dataset init mycollection.ds` creates a new collection
called 'mycollection.ds').

  - _dataset_ and _datasetd_ store JSON object documents in collections.
    The storage of the JSON documents differs.
  - dataset collections are defined in a directory containing a
    collection.json file
  - collection.json metadata file describing the collection,
    e.g. storage type, name, description, if versioning is enabled
  - collection objects are accessed by their key which is case insensitive
  - collection names lowered case and usually have a `.ds` extension
    for easy identification the directory must be lower case folder
    contain

_datatset_ stores JSON object documents in a pairtree
  - the pairtree path is always lowercase
  - a pairtree of JSON object documents
  - non-JSON attachments can be associated with a JSON document and
    found in a directories organized by semver (semantic version number)
  - versioned JSON documents are created sub directory incorporating a
    semver

_datasetd_ stores JSON object documents in a table named for the collection
  - objects are versioned into a collection history table by semver and key
  - attachments are not supported
  - can be exported to a collection using pairtree storage (e.g. a zip
    file will be generated holding a pairtree representation of the
    collection)

The choice of plain UTF-8 is intended to help future proof reading dataset
collections.  Care has been taken to keep _dataset_ simple enough and light
weight enough that it will run on a machine as small as a Raspberry Pi Zero
while being equally comfortable on a more resource rich server or desktop
environment. _dataset_ can be re-implement in any programming language
supporting file input and output, common string operations and along with
JSON encoding and decoding functions. The current implementation is in the
Go language.

Features
--------

_dataset_ supports
- Initialize a new dataset collection
  - Define metadata about the collection using a codemeta.json file
  - Define a keys file holding a list of allocated keys in the collection
  - Creates a pairtree for object storage

- Listing _keys_ in a collection
- Object level actions
  - create
  - read
  - update
  - delete
  - Documents as attachments
  - attachments (list)
  - attach (create/update)
  - retrieve (read)
  - prune (delete)
  - The ability to create data frames from while collections or based on
    keys lists
  - frames are defined using a list of keys and a lost
    "dot paths" describing what is to be pulled out
    of a stored JSON objects and into the frame
  - frame level actions
  - frames, list the frame names in the collection
  - frame, define a frame, does not overwrite an existing frame with
    the same name
  - frame-def, show the frame definition (in case we need it for some
    reason)
  - frame-objects, return a list of objects in the frame
  - refresh, using the current frame definition reload all the objects
    in the frame
  - reframe, replace the frame definition then reload the objects in
    the frame using the old frame key list
  - has-frame, check to see if a frame exists
  - delete-frame remove the frame

_datasetd_ supports

- List collections available from the web service
- List or update a collection's metadata
- List a collection's keys
- Object level actions
  - create
  - read
  - update
  - delete
  - Documents as attachments
  - attachments (list)
  - attach (create/update)
  - retrieve (read)
  - prune (delete)
  - A means of importing to or exporting from pairtree based dataset
    collections
  - The ability to create data frames from while
    collections or based on keys lists
  - frames are defined using "dot paths" describing
    what is to be pulled out of a stored JSON objects

Both _dataset_  and _datasetd_ maybe useful for general data science
applications needing JSON object management or in implementing repository
systems in research libraries and archives.

Limitations of _dataset_ and _datasetd_
-------------------------------------------

_dataset_ has many limitations, some are listed below

  - the pairtree implementation it is not a multi-process, multi-user
    data store
  - it is not a general purpose database system
  - it stores all keys in lower case in order to deal with file systems
    that are not case sensitive, compatibility needed by pairtrees
  - it stores collection names as lower case to deal with file systems that
    are not case sensitive
  - it does not have a built-in query language, search or sorting
  - it should NOT be used for sensitive or secret information

_datasetd_ is a simple web service intended to run on "localhost:8485".

  - it is a RESTful service
  - it does not include support for authentication
  - it does not support a query language, search or sorting
  - it does not support access control by users or roles
  - it does not provide auto key generation
  - it limits the size of JSON documents stored to the size supported by
    with host SQL JSON columns
  - it limits the size of attached files to less than 250 MiB
  - it does not support partial JSON record updates or retrieval
  - it does not provide an interactive Web UI for working with dataset
    collections
  - it does not support HTTPS or "at rest" encryption
  - it should NOT be used for sensitive or secret information

Authors and history
-------------------

- R. S. Doiel
- Tommy Morrell
*/
package dataset
