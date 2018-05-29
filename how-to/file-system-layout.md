
## File system layout

dataset provides a standardize way of organizing your JSON Objects, their attachments and 
data frames.

The directory layout looks like:

- collection (directory on file system)
    - a file, collection.json, holding metadata about collection
        - includes a map of filenames and buckets
        - frame names mapped to frame metadata files
    - a dircectory named "_frames" - any defined frames for the collection.
    - BUCKETS - a sequence of alphabet names (AA to ZZ) for buckets holding JSON documents and their attachments
        - Buckets let supporting common commands like ls, tree, etc. when the doc count is high

BUCKETS are names without meaning normally using Alphabetic characters. A dataset defined with four buckets
might looks like aa, ab, ba, bb. These directories will contains JSON documents and a tar file if the document
has attachments.

