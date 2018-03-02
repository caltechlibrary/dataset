/*
Package dataset provides a common approach for storing JSON object documents
on local disc or on S3 and Google Cloud Storage. It is intended as a
single user system for intermediate processing of JSON content for analysis
or batch processing.  It is not a database management system (if you need
a JSON database system I would suggest looking at Couchdb, Mongo and Redis
as a starting point).

The approach dataset takes to storing buckets is to maintain a JSON document
with keys (document names) and bucket assignments. JSON documents (and
possibly their attachments) are then stored based on that assignment.
Conversely the collection.json document is used to find and retrieve
documents from the collection. The layout of the metadata is as follows

+ Collection
	+ Collection/collection.json - metadata for retrieval
	+ Collection/[Buckets] - usually an "aa" to "zz" list of buckets
	+ Collection/[Bucket]/[Document]

A key feature of dataset is to be Posix shell friendly. This has
lead to storing the JSON documents in a directory structure that
standard Posix tooling can traverse. It has also mean that the JSON
documents themselves remain on "disc" as plain text. This has facilitated
integration with many other applications, programming langauages and
systems.

Attachments are non-JSON documents explicitly "attached" that share the
same basename but are placed in a tar ball (e.g. document Jane.Doe.json
attachements would be stored in Jane.Doe.tar).

Additional operations beside storing and reading JSON documents are also
supported. These include creating lists (arrays) of JSON documents from
a list of keys, listing keys in the collection, counting documents in the
collection, indexing and searching by indexes.

The primary use case driving the development of dataset is harvesting
API content for library systems (e.g. EPrints, Invenio, ArchivesSpace,
ORCID, CrossRef, OCLC). The harvesting needed to be done in such a
way as to leverage existing Posix tooling (e.g. grep, sed, etc) for
processing and analysis.

Initial use case:

Caltech Library has many repository, catelog and record management systems
(e.g. EPrints, Invenion, ArchivesSpace, Islandora, Invenio). It is common
practice to harvest data from these systems for analysis or processing.
Harvested records typically come in XML or JSON format. JSON has proven a
flexibly way for working with the data and in our more modern tools the
common format we use to move data around. We needed a way to standardize
how we stored these JSON records for intermediate processing to allow us
to use the growing ecosystem of JSON related tooling available under
Posix/Unix compatible systems.

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
directories will contains JSON documents and a tar file if the document
has attachments.


Operations

+ Collection level
    + InitCollection (collection) - creates or opens collection structure on disc, creates collection.json and keys.json if new
    + Open (collection) - opens an existing collections and reads collection.json into memory
    + Close (collection) - writes changes to collection.json to disc if dirty
    + Keys (collection) - list of keys in the collection
+ JSON document level
    + Create (JSON document) - saves a new JSON blob or overwrites and existing one on  disc with given blob name, updates keys.json if needed
    + Read (JSON document)) - finds the JSON document in the buckets and returns the JSON document contents
    + Update (JSON document) - updates an existing blob on disc (record must already exist)
    + Delete (JSON document) - removes a JSON blob from its disc
    + Path (JSON document) - returns the path to the JSON document
+ Select list level
    + Count (select list) - returns the number of keys in a select list

Example

Common operations using the *dataset* command line tool

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

Example Bash script usage

    # Create a collection "mystuff.ds" inside the directory called demo
    dataset init mystuff.ds
    # if successful an expression to export the collection name is show
    export DATASET="mystuff.ds"

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
    # /bin/rm -fR mystuff.ds

Common operations shown in Golang

+ create collection
+ create a JSON document to collection
+ read a JSON document
+ update a JSON document
+ delete a JSON document

Example Go code

    // Create a collection "mystuff" inside the directory called demo
    collection, err := dataset.InitCollection("mystuff.ds")
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

*/
package dataset
