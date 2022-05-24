/*
Package dataset provides a common approach for storing JSON object documents
on local disc or via a web service. The command line implementation is intended
as a single user system for intermediate processing of JSON content for analysis
or batch processing. The web implementation is intended for concurrent use by
users or processes.  dataset is not a database management system (if you need
a JSON database system I would suggest looking at Couchdb, Mongo and Redis
as a starting point).

The command line dataset stores JSON documents in a pairtree structure
under the collection folder. The keys are the JSON document names stored
in the pairtree. Keys are lowercase to avoid problems where the file system
is case insensitive (e.g. macOS default FS). In the root collection folder is
a codemeta JSON file describing the collection. There is a collection.json file
holding the dataset version number, an array of keys stored in the collection
and operational metadata (e.g.`"version_control": true`).
The layout of the metadata is as follows

+ Collection - a directory
	+ Collection/codemeta.json - metadata for retrieval
	+ Collection/collection.json - an JSON array of object keys in the collection
	+ Collection/[Pairtree] - holds individual JSON docs and attachments
	+ Collection/[_frames] - holds the data frames and their metadata

A key feature of dataset is to be Posix shell friendly. This has
lead to storing the JSON documents in a directory structure that
standard Posix tooling can traverse. It has also mean that the JSON
documents themselves remain on disk as plain text. This has facilitated
integration with many other applications, programming langauages and
systems.

Attachments are non-JSON documents explicitly "attached" that share the
same pairtree path but are placed in a sub directory called "_". If the
document name is "jane.doe.json" and the attachment is photo.jpg
the JSON document is "pairtree/ja/ne/.d/oe/jane.doe.json" and the photo
is in "pairtree/ja/ne/.d/oe/_/photo.jpg".

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
(e.g. EPrints, Invenio, ArchivesSpace, Islandora, InvenioRDM). It is common
practice to harvest data from these systems for analysis or processing.
Harvested records typically come in XML or JSON format. JSON has proven a
flexibly way for working with the data and in our more modern tools the
common format we use to move data around. We needed a way to standardize
how we stored these JSON records for intermediate processing to allow us
to use the growing ecosystem of JSON related tooling available under
Posix/Unix compatible systems.

The web API version of dataset stores metadata in a table named for the collection
using a JSON column to store the object.  The collection specific metadata
(e.g. versioning is enabled) is store in a row of a table called "_collections".
Metadata for the collection is stored in the "_collections"'s row as a JSON column.

*/
package dataset
