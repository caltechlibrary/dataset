/*
Package dataset includes the operations needed for processing collections of JSON documents and their attachments.

Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>

Copyright (c) 2024, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Dataset Project, v3
===================

The Dataset Project provides tools for working with collections of
JSON Object documents stored on the local file system or via a dataset
web service.  Two tools are provided, a command line interface (__dataset3__)
and a web service (__dataset3d__).

dataset3 command line tool
--------------------------

__dataset3__ is a command line tool for working with collections of JSON
objects.  Collections are stored in a SQL storage engine that has JSON
column supoprt (SQLite3 is the default storage engine). The __dataset3__
collections can be accessed via __dataset3d__ web service.

The __dataset3__ command line tool supports common data management operations
such as initialization of collections; document creation, reading,
updating and deleting; listing keys of JSON objects in the collection;
and associating non-JSON documents (attachments) with specific JSON
documents in the collection.

### enhanced features include

- aggregate objects using SQL queries

dataset3d, dataset3 as a web service
------------------------------------

__dataset3d__ is a web service implementation of the __dataset3__ command line
program. It features a sub-set of capability found in the command line
tool. This allows dataset collections to be integrated safely into web
applications or used concurrently by multiple processes. It achieves this
by storing the dataset collection in a SQL database using JSON columns.

Design choices
--------------

__dataset3__ and __dataset3d__ are intended to be simple tools for managing
collections JSON object documents in a predictable structured way.

__dataset3__ is guided by the idea that you should be able to work with
JSON documents as easily as you can any plain text document on the Unix
command line. __dataset3__ is intended to be simple to use with minimal
setup (e.g.  `dataset init mycollection.ds` creates a new collection
called 'mycollection.ds').

  - __dataset3__ and __dataset3d__ store JSON object documents in collections.
    The storage of the JSON documents differs.
  - dataset collections are defined in a directory containing a
    collection.json file
  - collection.json metadata file describing the collection,
    e.g. storage type, name, description
  - collection objects are accessed by their key which are stored lower case
  - collection names lowered case and usually have a `.ds` extension
    for easy identification the directory must be lower case folder
    contain

__datatset3__ stores JSON object documents in a JSON column of a SQL storage engine
  - versioned JSON documents are in a "history" table using the same structure as
  the primary table holding the collection's objects. History is enabled by default.
  
__dataset3d__ stores JSON object documents in a table named for the collection
  - objects are versioned into a collection history table using the version column (integer)
    and key as a the "primary key".
  - attachments are not supported

The choice of plain UTF-8 is intended to help future proof reading dataset
collections.  Care has been taken to keep __dataset3__ simple enough and light
weight enough that it will run on a machine as small as a Raspberry Pi Zero
while being equally comfortable on a more resource rich server or desktop
environment. __dataset3__ can be re-implement in any programming language
supporting file input and output, common string operations and along with
JSON encoding and decoding functions. The current implementation is in the
Go language.

Features
--------

__dataset3__ supports
- Initialize a new dataset collection
  - Define metadata about the collection using a codemeta.json file
  - Define a keys file holding a list of allocated keys in the collection
  - Creates a pairtree for object storage

- Listing __keys__ in a collection
- __query__ allows returning lists of objects via SQL queries
- Object level actions
  - create
  - read
  - update
  - delete

__datasetd3__ supports

- List collections available from the web service
- List or update a collection's metadata
- List a collection's keys
- Object level actions
  - create
  - read
  - update
  - delete
  - query

Both __dataset3__  and __dataset3d__ maybe useful for general data science
applications needing JSON object management or in implementing repository
systems in research libraries and archives.

Limitations of __dataset3__ and __dataset3d__
---------------------------------------------

__dataset3__ has many limitations, some are listed below

  - it is not a general purpose database system
  - it stores all keys in lower case
  - it stores collection names as lower case to deal with file systems that
    are not case sensitive
  - it does not have a built-in query language, search or sorting but SQL is available
    via the SQL storage engine
  - it should NOT be used for sensitive or secret information

__dataset3d__ is a simple web service intended to run on localhost only (e.g. "localhost:8485")

  - it is a semi-RESTful service
  - it does not include support for authentication
  - it does not support access control by users or roles
  - it does not provide auto key generation
  - it limits the size of JSON documents stored to the size supported by
    with host SQL JSON columns
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
