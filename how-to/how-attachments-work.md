How Attachments Work
====================

The primary use case of the __dataset__ tool is managing JSON documents.
There exist a common secondary use case of including support for
\"attached\" non-JSON documents. Example 1, when we harvest content from
a system who does not support JSON natively it is handy to keep a
version of the harvested content for audit purposes. The EPrints system
has a REST API that returns XML. Storing the original EPrint XML
document gives the developer an ability to verify that their JSON
rendering matches the EPrint XML should their JSON needs change in the
future.

This raises questions of how to keep things simple while supporting an
arbitrary number of attachments for JSON object document? How do you
handle versioning when some types of collections need it for attachments
and others don\'t?

The __dataset__ command line tool and related Go package store the
attachments un-versioned by default in the pairtree.  The un-versioned
attached document is stored in a pairtree in the "attachments" folder
of the collection. This is true regardless of the storage engine used
(e.g. pairtree storage, sql storage). The attached document is
stored by its basename. The basename must be unique among the documents
attached otherwise it will be overwritten when attaching another
document using the same basename.

If you need versioning you create your collection with versioning
support. Attaching the documents will automatically version based
on the basename of the attachment. When retrieving a specific version
you need to support a [semver](https://semver.org) using the appropraite
versioned verb.  By default reads will be the current version of the
document, meaning the version with the "largest" semver value.
In a versioned collection two files with the same basename will result
in different "versions" of the document with the highest semver 
reflecting the most recent addition.

It is easier to see with this example. We have a dataset collection
called \"Sea-Mamals.ds\". We have a JSON object stored called
\"walrus\". We want to attach \"notes-on-walrus.docx\" which is on our
local drive under \"/Users/fred/Documents/notes-on-walrus.docx\".

Using the __dataset__ cli you issue the follow command \--

```shell
    dataset create Sea-Mamals.ds walrus '{"description": "may have tusks", "size": "impressive"}'
    dataset attach Sea-Mamals.ds walrus \
       /Users/fred/Documents/notes-on-walrus.docx
```

The results in a simple directory stricture for the JSON object and
attachment.

        Sea-Mamanls/pairtree/wa/lr/us/walrus.json
        Sea-Mamanls/attachments/wa/lr/us/notes-on-walrus.docx

The directory structured for versioned attachments and JSON document
is more complex. In the case of the JSON document the semver gets
embedded in the JSON document name while the attachments are
stored in subfolders by version. The assignment of the semver
is automatic based on the collection's original setup.


How Attachments look in the JSON Object
---------------------------------------

In version 2 of dataset the JSON document remain unmodified. You will
nolonger see added attributes like `_Key` or `_Attachments` in the object.
Likewise attachments will remain unaltered beyond remaining the
file path to the basename when the attachment is made. Versions are
in their own version folder.

