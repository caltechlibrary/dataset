
## File system layout

dataset (directory on file system)

- collection (directory on file system)
    - collection.json - metadata about collection
        - maps the filename of the JSON blob stored to a bucket in the collection
        - e.g. file "mydocs.jons" stored in bucket "aa" would have a map of {"mydocs.json": "aa"}
    - keys.json - a list of keys in the collection (it is the default select list)
    - BUCKETS - a sequence of alphabet names for buckets holding JSON documents and their attachments
        - Buckets let supporting common commands like ls, tree, etc. when the doc count is high
    - SELECT_LIST.json - a JSON document holding an array of keys
        - the default select list is "keys", it is not mutable by Push, Pop, Shift and Unshift
        - select lists cannot be named "keys" or "collection"

BUCKETS are names without meaning normally using Alphabetic characters. A dataset defined with four buckets
might looks like aa, ab, ba, bb. These directories will contains JSON documents and a tar file if the document
has attachments.

