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
attachments un-versioned by default in the pairtree. The metadata about
the attached document is stored in a sub folder `_docs`. The un-versioned
attached document is stored in `v0.0.0` folder. The attached document is
stored by its basename. The basename must be unique among the documents
attached otherwise it will be overwritten when attaching another
document using the same basename.

If you need versioning you MUST supply a valid
[semver](https://semver.org) when attaching the document. The metadata
for the attached document will be in `_docs` as before but the document
will be stored in a sub directory indicated by the semver. The basename
must be unique to the semver provided otherwise the document with the
same basename using that semver will be overwritten.

It is easier to see with this example. We have a dataset collection
called \"Sea-Mamals.ds\". We have a JSON object stored called
\"walrus\". We want to attach \"notes-on-walrus.docx\" which is on our
local drive under \"/Users/fred/Documents/notes-on-walrus.docx\".

Using the __dataset__ cli you issue the follow command \--

``` {.shell}
    dataset create Sea-Mamals.ds walrus '{"description": "may have tusks", "size": "impressive"}'
    dataset attach Sea-Mamals.ds walrus \
       /Users/fred/Documents/notes-on-walrus.docx
```

The results in a simple directory stricture for the JSON object and
attachment.

        Sea-Mamanls/pairtree/wa/lr/us/walrus.json
        Sea-Mamanls/pairtree/wa/lr/us/v0.0.0/notes-on-walrus.docx

In this example the metadata for the attachment is updated in the
walrus.json file. Since no versioning was specified for
\"notes-on-walrus.docx\" it is stored as version v0.0.0.

If we had added our attachment including a semver the directory
structure will be slightly more complex.

``` {.shell}
    dataset attach Sea-Mamals.ds walrus v0.0.1 /Users/fred/Documents/notes-on-walrus.docx
```

This will cause additional sub directories to exist (if they haven\'t be
created before). Our \"un-versioned\" version still exists as v0.0.0 but
now we have v0.0.1. Our attachment metadata file in our JSON object file
will now include an href pointing to v0.0.1 and a map to all versions
including v0.0.0.

        Sea-Mamanls/pairtree/wa/lr/us/walrus.json
        Sea-Mamanls/pairtree/wa/lr/us/v0.0.0/notes-on-walrus.docx
        Sea-Mamanls/pairtree/wa/lr/us/v0.0.1/notes-on-walrus.docx

If we later add a v0.0.2 of \"notes-on-walrus.docx\" it\'d looks like

        Sea-Mamanls/pairtree/wa/lr/us/walrus.json
        Sea-Mamanls/pairtree/wa/lr/us/v0.0.0/notes-on-walrus.docx
        Sea-Mamanls/pairtree/wa/lr/us/v0.0.1/notes-on-walrus.docx
        Sea-Mamanls/pairtree/wa/lr/us/v0.0.2/notes-on-walrus.docx

All the metadata about the files attached are stored in the primary JSON
document under the attribute `_Attachments`. In the metadata we include
an \"href\" string and \"version_hrefs\" map. The version_href will
point to all known versions keyed by the semver. The href string will
point to the last version added, in this case v0.0.2.

IMPORTANT: If you provide the same semver and attach a file with the
same basename the previously stored version will be overwritten. Example
if we issue our original un-versioned command the v0.0.0 copy of
\"notes-on-walrus.docx\" will be overwritten!

__dataset__ attachment versioning is user driven. The only implicit
version is v0.0.0 if no semver is provided. __dataset__ is not a
substitute for a version control system like [Subversion]() or [Git]()
and is not substitute for a versioned file systems like [ZFS](). If your
program needs to avoid overwriting an existing version or to \"auto
increment\" the semver you need to check the existing versions and
decide what the new version will be before attaching the new version of
the document.

The semver versioned directories may contain more than one attached
document. The documents attached can be of various versions though if
you attach more than one document at a time they will carry the same
semver. This is because their is an implied semver is v0.0.0 when using
the command line without semver __dataset__ otherwise the first valid
semver is used for all files being attached in that command execution.

NOTE: The href in the attachments metadata always points at the last
attached version.

How Attachments look in the JSON Object
---------------------------------------

When you retrieve a JSON object __dataset__ will add some internal
fields. The first is a `_Key` and if you have any attachments a
`_Attachments` array will be added. The later holds the metadata we
create during the attachment process.

Let\'s look at our first example again in detail.

``` {.shell}
    dataset create Sea-Mamals.ds walrus '{"description": "may have tusks", "size": "impressive"}'
    dataset attach Sea-Mamals.ds walrus \
       /Users/fred/Documents/notes-on-walrus.docx
```

The JSON object created by the two command looks like

``` {.json}
    {
        "_Key": "walrus",
        "description": "may have tusks",
        "size": "impressive",
        "_Attachments": [
            {
                "name": "notes-on-walrus.docx",
                "href": "v0.0.0/notes-on-walrus.docx",
                "version_hrefs": {
                    "v0.0.0": "v0.0.0/notes-on-walrus.docx"
                },
                ...
            }
        ]
    }
```

When we added v0.0.1 the object would change shape and be something like

``` {.json}
    {
        "_Key": "walrus",
        "description": "may have tusks",
        "size": "impressive",
        "_Attachments": [
            {
                "name": "notes-on-walrus.docx",
                "href": "v0.0.1/notes-on-walrus.docx",
                "version_hrefs": {
                    "v0.0.0": "v0.0.0/notes-on-walrus.docx",
                    "v0.0.1": "v0.0.1/notes-on-walrus.docx"
                },
                ...
            }
        ]
    }
```

If you have a program that moves old versions off to Glacier you\'ll
want to update the value in the version_hrefs effected. In this example
we\'ve moved v0.0.0 off to

    "s3://sea-mamals/walrus/v0.0.0/notes-on-walrus.docx"

The JSON should look something like\--

``` {.json}
    {
        "_Key": "walrus",
        "description": "may have tusks",
        "size": "impressive",
        "_Attachments": [
            {
                "name": "notes-on-walrus.docx",
                "size": "1041",
                "href": "v0.0.1/notes-on-walrus.docx",
                "version_hrefs": {
                    "v0.0.0": "s3://sea-mamals/v0.0.0/notes-on-walrus.docx",
                    "v0.0.1": "v0.0.1/notes-on-walrus.docx"
                },
                ...
            }
        ]
    }
```
