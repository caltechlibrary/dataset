
## File system layout

dataset provides a two ways to organize your JSON Objects. The original
was a "buckets" oriented layout. The newer versions of file layout uses
a [pairtree](https://tools.ietf.org/html/draft-kunze-pairtree-01). 
Both are managed/described by the collection.json document
at in the root folder the the collection. Both file layouts currently
support "attachments" as a tar ball of with the same basename as the JSON
object document (e.g. hello-world.json would have attachments stored as
hello-world.tar). Attachments are experimental and how they are handled
will likely change in the future. If so the repair/analyzer abilities
of dataset should ease the migration process.

## Pairtree

The directory layout looks like:

- collection (directory on the file system)
    - [namaste](https://confluence.ucop.edu/display/Curation/Namaste) records identifying the collection
    - a file, collection.json, holding metadata about the collection
    - a directory named "_frames" holding frame definitions for the collection
    - a directory named "pairtree" holding the pairtree where the JSON document and attachmetns are stored.

## Buckets 

The directory layout looks like:

- collection (directory on file system)
    - [namaste](https://confluence.ucop.edu/display/Curation/Namaste) records identifying the collection
    - a file, collection.json, holding metadata about collection
        - includes a map of filenames and buckets
        - frame names mapped to frame metadata files
    - a dircectory named "_frames" - any defined frames for the collection.
    - BUCKETS - a sequence of alphabet names (AA to ZZ) for buckets holding JSON documents and their attachments
        - Buckets let supporting common commands like ls, tree, etc. when the doc count is high

BUCKETS are names without meaning normally using Alphabetic characters. A dataset defined with four buckets
might looks like aa, ab, ba, bb. These directories will contains JSON documents and a tar file if the document
has attachments.

