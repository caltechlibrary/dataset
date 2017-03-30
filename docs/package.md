# dataset
--
    import "github.com/caltechlibrary/dataset"

Package dataset is a go package for managing JSON documents stored on disc

@author R. S. Doiel, <rsdoiel@caltech.edu>

Copyright (c) 2017, Caltech All rights not granted herein are expressly reserved
by Caltech.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation and/or
other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
### WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


Package dataset is a go package for managing JSON documents stored on disc

@author R. S. Doiel, <rsdoiel@caltech.edu>

Copyright (c) 2017, Caltech All rights not granted herein are expressly reserved
by Caltech.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation and/or
other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
### WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Package dataset provides a common approach of storing JSON documents and related
attachments systematically and predictably on the file systems. The driving
usecase behind dataset is creating a unified approach to harvesting metadata
from various hetrogenious systems using in Caltech Library (e.g. EPrints,
ArchivesSpace, Islandora and outside API like ORCID, CrossRef, OCLC). This
suggests that dataset as a go package and command line tool may have more
general applications where a database sytems might be more than you need and
ad-hoc collections on disc maybe inconvient to evolve as your data explortion
takes you in different directions in analysis. The dataset command line tool is
in intended to be easy to script both in Bash as well as more featureful
languages like Python.

Dataset is not a good choice if you need a fast key/value store or actual
database features. It doesn't support multiple users, record locking or a query
language interface. It is targetting the sweet spot between ad-hoc JSON document
storage on disc and needing a more complete system like Couch, Solr, or Fedora4
for managing JSON docs.


Use case

Caltech Library has many repository, catelog and record management systems (e.g.
EPrints, Invenion, ArchivesSpace, Islandora, Invenio). It is common practice to
harvest data from these systems for analysis or processing. Harvested records
typically come in XML or JSON format. JSON has proven a flexibly way for working
with the data and in our more modern tools the common format we use to move data
around. We needed a way to standardize how we stored these JSON records for
intermediate processing to allow us to use the growing ecosystem of JSON related
tooling available under Posix/Unix compatible systems.


Aproach to file system layout

+ /dataset (directory on file system)

    + collection (directory on file system)
        + collection.json - metadata about collection
            + maps the filename of the JSON blob stored to a bucket in the collection
            + e.g. file "mydocs.jons" stored in bucket "aa" would have a map of {"mydocs.json": "aa"}
        + keys.json - a list of keys in the collection (it is the default select list)
        + BUCKETS - a sequence of alphabet names for buckets holding JSON documents and their attachments
            + Buckets let supporting common commands like ls, tree, etc. when the doc count is high
        + SELECT_LIST.json - a JSON document holding an array of keys
            + the default select list is "keys", it is not mutable by Push, Pop, Shift and Unshift
            + select lists cannot be named "keys" or "collection"

BUCKETS are names without meaning normally using Alphabetic characters. A
dataset defined with four buckets might looks like aa, ab, ba, bb. These
directories will contains JSON documents and a tar file if the document has
attachments.


### Operations

+ Collection level

    + Create (collection) - creates or opens collection structure on disc, creates collection.json and keys.json if new
    + Open (collection) - opens an existing collections and reads collection.json into memory
    + Close (collection) - writes changes to collection.json to disc if dirty
    + Delete (collection) - removes a collection from disc
    + Keys (collection) - list of keys in the collection
    + Select (collection) - returns the request select list, will create the list if not exist, append keys if provided
    + Clear (collection) - Removes a select list from a collection and disc
    + Lists (collection) - returns the names of the available select lists

+ JSON document level

    + Create (JSON document) - saves a new JSON blob or overwrites and existing one on  disc with given blob name, updates keys.json if needed
    + Read (JSON document)) - finds the JSON document in the buckets and returns the JSON document contents
    + Update (JSON document) - updates an existing blob on disc (record must already exist)
    + Delete (JSON document) - removes a JSON blob from its disc
    + Path (JSON document) - returns the path to the JSON document

+ Select list level

    + First (select list) - returns the value of the first key in the select list (non-distructively)
    + Last (select list) - returns the value of the last key in the select list (non-distructively)
    + Rest (select list) - returns values of all keys in the select list except the first (non-destructively)
    + List (select list) - returns values of all keys in the select list (non-destructively)
    + Length (select list) - returns the number of keys in a select list
    + Push (select list) - appends one or more keys to an existing select list
    + Pop (select list) - returns the last key in select list and removes it
    + Unshift (select list) - inserts one or more new keys at the beginning of the select list
    + Shift (select list) - returns the first key in a select list and removes it
    + Sort (select list) - orders the select lists' keys in ascending or descending alphabetical order
    + Reverse (select list) - flips the order of the keys in the select list


### Example

Common operations using the *dataset* command line tool

+ create collection + create a JSON document to collection + read a JSON
document + update a JSON document + delete a JSON document

Example Bash script usage

    # Create a collection "mystuff" inside the directory called demo
    dataset init demo/mystuff
    # if successful an expression to export the collection name is show
    export DATASET_COLLECTION=demo/mystuff

    # Create a JSON document
    dataset create freda.json '{"name":"freda","email":"freda@inverness.example.org"}'
    # If successful then you should see an OK or an error message

    # Read a JSON document
    dataset read freda.json

    # Path to JSON document
    dataset path freda.json

    # Update a JSON document
    dataset update freda.json '{"name":"freda","email":"freda@zbs.example.org"}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys

    # Delete a JSON document
    dataset delete freda.json

    # To remove the collection just use the Unix shell command
    # /bin/rm -fR demo/mystuff


Common operations shown in Golang

+ create collection + create a JSON document to collection + read a JSON
document + update a JSON document + delete a JSON document

Example Go code

    // Create a collection "mystuff" inside the directory called demo
    collection, err := dataset.Create("demo/mystuff", dataset.GenerateBucketNames("ab", 2))
    if err != nil {
        log.Fatalf("%s", err)
    }
    defer collection.Close()
    // Create a JSON document
    docName := "freda.json"
    document := map[string]string{"name":"freda","email":"freda@inverness.example.org"}
    if err := collection.Create(docName, document); err != nil {
        log.Fatalf("%s", err)
    }
    // Attach an image file to freda.json in the collection
    if buf, err := ioutil.ReadAll("images/freda.png"); err != nil {
       collection.Attach("freda", "images/freda.png", buf)
    } else {
       log.Fatalf("%s", err)
    }
    // Read a JSON document
    if err := collection.Read(docName, document); err != nil {
        log.Fatalf("%s", err)
    }
    // Update a JSON document
    document["email"] = "freda@zbs.example.org"
    if err := collection.Update(docName, document); err != nil {
        log.Fatalf("%s", err)
    }
    // Delete a JSON document
    if err := collection.Delete(docName); err != nil {
        log.Fatalf("%s", err)
    }

Working with attachments in Go

        collection, err := dataset.Open("dataset/mystuff")
        if err != nil {
            log.Fatalf("%s", err)
        }
        defer collection.Close()

    	// Add a helloworld.txt file to freda.json record as an attachment.
        if err := collection.Attach("freda", "docs/helloworld.txt", []byte("Hello World!!!!")); err != nil {
            log.Fatalf("%s", err)
        }

    	// Attached files aditional files from the filesystem by their relative file path
    	if err := collection.AttachFiles("freda", "docs/presentation-article.pdf", "docs/charts-and-figures.zip", "docs/transcript.fdx") {
            log.Fatalf("%s", err)
    	}

    	// List the attached files for freda.json
    	if filenames, err := collection.Attachments("freda"); err != nil {
            log.Fatalf("%s", err)
    	} else {
    		fmt.Printf("%s\n", strings.Join(filenames, "\n"))
    	}

    	// Get an array of attachments (reads in content into memory as an array of Attachment Structs)
    	allAttachments, err := collection.GetAttached("freda")
    	if err != nil {
            log.Fatalf("%s", err)
    	}
    	fmt.Printf("all attachments: %+v\n", allAttachments)

    	// Get two attachments docs/transcript.fdx, docs/helloworld.txt
    	twoAttachments, _ := collection.GetAttached("fred", "docs/transcript.fdx", "docs/helloworld.txt")
    	fmt.Printf("two attachments: %+v\n", twoAttachments)

        // Get attached files writing them out to disc relative to your working directory
    	if err := collection.GetAttachedFiles("freda"); err != nil {
            log.Fatalf("%s", err)
    	}

    	// Get two selection attached files writing them out to disc relative to your working directory
    	if err := collection.GetAttached("fred", "docs/transcript.fdx", "docs/helloworld.txt"); err != nil {
            log.Fatalf("%s", err)
    	}

        // Remove docs/transcript.fdx and docs/helloworld.txt from freda.json attachments
    	if err := collection.Detach("fred", "docs/transcript.fdx", "docs/helloworld.txt"); err != nil {
            log.Fatalf("%s", err)
    	}

    	// Remove all attached files from freda.json
    	if err := collection.Detach("fred")
            log.Fatalf("%s", err)
    	}

## Usage

```go
const (
	// Version of the dataset package
	Version = "v0.0.1-beta9"

	// License is a formatted from for dataset package based command line tools
	License = `
%s %s

Copyright (c) 2017, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	DefaultAlphabet = `abcdefghijklmnopqrstuvwxyz`

	ASC  = iota
	DESC = iota
)
```

#### func  Delete

```go
func Delete(name string) error
```
Delete an entire collection

#### func  GenerateBucketNames

```go
func GenerateBucketNames(alphabet string, length int) []string
```
GenerateBucketNames provides a list of permutations of requested length to use
as bucket names

#### type Attachment

```go
type Attachment struct {
	// Name is the filename and path to be used inside the generated tar file
	Name string
	// Body is a byte array for storing the content associated with Name
	Body []byte
}
```

Attachment is a structure for holding non-JSON content you wish to store
alongside a JSON document in a collection

#### type Collection

```go
type Collection struct {
	// Version of collection being stored
	Version string `json:"verison"`
	// Name of collection
	Name string `json:"name"`
	// Dataset is a directory name that holds collections
	Dataset string `json:"dataset"`
	// Buckets is a list of bucket names used by collection
	Buckets []string `json:"buckets"`
	// KeyMap holds the document name to bucket map for the collection
	KeyMap map[string]string `json:"keymap"`
	// SelectLists holds the names of available select lists
	SelectLists []string `json:"select_lists"`
}
```

Collection is the container holding buckets which in turn hold JSON docs

#### func  Create

```go
func Create(name string, bucketNames []string) (*Collection, error)
```
Create - create a new collection structure on disc name should be filesystem
friendly

#### func  Open

```go
func Open(name string) (*Collection, error)
```
Open reads in a collection's metadata and returns and new collection structure
and err

#### func (*Collection) Attach

```go
func (c *Collection) Attach(name string, attachments ...*Attachment) error
```
Attach a non-JSON document to a JSON document in the collection. Attachments are
stored in a tar file, if tar file exits then attachment(s) are appended to tar
file.

#### func (*Collection) AttachFiles

```go
func (c *Collection) AttachFiles(name string, fileNames ...string) error
```
AttachFiles a non-JSON documents to a JSON document in the collection.
Attachments are stored in a tar file, if tar file exits then attachment(s) are
appended to tar file.

#### func (*Collection) Attachments

```go
func (c *Collection) Attachments(name string) ([]string, error)
```
Attachments returns a list of files in the attached tarball for a given name in
the collection

#### func (*Collection) Clear

```go
func (c *Collection) Clear(name string) error
```
Clear removes a select list from disc and the collection

#### func (*Collection) Close

```go
func (c *Collection) Close() error
```
Close closes a collection, writing the updated keys to disc

#### func (*Collection) Create

```go
func (c *Collection) Create(name string, data interface{}) error
```
Create a JSON doc from an interface{} and adds it to a collection, if problem
returns an error name must be unique

#### func (*Collection) CreateAsJSON

```go
func (c *Collection) CreateAsJSON(name string, src []byte) error
```
CreateAsJSON adds or replaces a JSON doc to a collection, if problem returns an
error name must be unique (treated like a key in a key/value store)

#### func (*Collection) Delete

```go
func (c *Collection) Delete(name string) error
```
Delete removes a JSON doc from a collection

#### func (*Collection) Detach

```go
func (c *Collection) Detach(name string, filterNames ...string) error
```
Detach a non-JSON document from a JSON document in the collection. FIXME: Need
to add detaching specific filenames

#### func (*Collection) DocPath

```go
func (c *Collection) DocPath(name string) (string, error)
```
DocPath returns a full path to a key or an error if not found

#### func (*Collection) GetAttached

```go
func (c *Collection) GetAttached(name string, filterNames ...string) ([]Attachment, error)
```
GetAttached returns an Attachment array or error If no filterNames provided then
return all attachments or error

#### func (*Collection) GetAttachedFiles

```go
func (c *Collection) GetAttachedFiles(name string, filterNames ...string) error
```
GetAttachedFiles returns an error if encountered, side effect is to write file
to destination directory If no filterNames provided then return all attachments
or error

#### func (*Collection) Keys

```go
func (c *Collection) Keys() []string
```
Keys returns a list of keys in a collection

#### func (*Collection) Lists

```go
func (c *Collection) Lists() []string
```
Lists returns a list of available select lists, should always contain the
default keys list

#### func (*Collection) Read

```go
func (c *Collection) Read(name string, data interface{}) error
```
Read finds the record in a collection, updates the data interface provide and if
problem returns an error name must exist or an error is returned

#### func (*Collection) ReadAsJSON

```go
func (c *Collection) ReadAsJSON(name string) ([]byte, error)
```
ReadAsJSON finds a the record in the collection and returns the JSON source

#### func (*Collection) Select

```go
func (c *Collection) Select(params ...string) (*SelectList, error)
```
Select returns a select assocaited with a collection, it will be created if
neccessary and any keys included will be added before returning the updated list

#### func (*Collection) Update

```go
func (c *Collection) Update(name string, data interface{}) error
```
Update JSON doc in a collection from the provided data interface (note: JSON doc
must exist or returns an error )

#### func (*Collection) UpdateAsJSON

```go
func (c *Collection) UpdateAsJSON(name string, src []byte) error
```
UpdateAsJSON takes a JSON doc and writes it to a collection (note: Record must
exist or returns an error)

#### type SelectList

```go
type SelectList struct {
	FName        string   `json:"name"`
	Keys         []string `json:"keys"`
	CustomLessFn func([]string, int, int) bool
}
```

SelectList is an ordered set of keys

#### func (SelectList) First

```go
func (s SelectList) First() string
```
First select list returns the first item in the list (non-destructively)

#### func (*SelectList) Last

```go
func (s *SelectList) Last() string
```
Last select list returns the list item from the list (non-destructively)

#### func (*SelectList) Len

```go
func (s *SelectList) Len() int
```
Len returns the number of keys in the select list

#### func (*SelectList) Less

```go
func (s *SelectList) Less(i, j int) bool
```
Less compare two elements returning true if first is less than second, false
otherwise

#### func (*SelectList) List

```go
func (s *SelectList) List() []string
```
List returns all the keys in the select list (non-destructively)

#### func (*SelectList) Pop

```go
func (s *SelectList) Pop() string
```
Pop select list removes from the end of an array returning the element removed

#### func (*SelectList) Push

```go
func (s *SelectList) Push(val string)
```
Push select list appends an element to the end of an array

#### func (*SelectList) Reset

```go
func (s *SelectList) Reset()
```
Reset a select list to an empty state (file still exists on disc)

#### func (*SelectList) Rest

```go
func (s *SelectList) Rest() []string
```
Rest select list returns all but the first n items of the list
(non-destructively)

#### func (*SelectList) Reverse

```go
func (s *SelectList) Reverse()
```
Reverse flips the order of a select list

#### func (*SelectList) SaveList

```go
func (s *SelectList) SaveList() error
```
SaveList writes the .Keys to a JSON document named .FName

#### func (*SelectList) Shift

```go
func (s *SelectList) Shift() string
```
Shift select list removes from the beginning of and array returning the element
removed

#### func (*SelectList) Sort

```go
func (s *SelectList) Sort(direction int)
```
Sort sorts the keys in in ascending order alphabetically

#### func (*SelectList) String

```go
func (s *SelectList) String() string
```
String returns the Keys portion of the select list as a string delimited with
new lines

#### func (*SelectList) Swap

```go
func (s *SelectList) Swap(i, j int)
```
Swap updates the position of two compared keys

#### func (*SelectList) Unshift

```go
func (s *SelectList) Unshift(val string)
```
Unshift select list inserts an element at the start of an array
