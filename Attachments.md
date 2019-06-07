
# Attachments and how they work

The primary use case of the **dataset** tool is managing JSON documents.
There exist a common secondary use case of including support for "attached"
non-JSON documents. Example 1, when we harvest content from a system who
does not support JSON natively it is handy to keep a version of the 
harvested content for audit purposes. The EPrints system has a REST
API that returns XML.  Storing the original EPrint XML document gives
the developer an ability to verify  that their JSON rendering matches
the EPrint XML should their JSON needs change in the future. 

This raises questions of how to keep things simple while supporting
an arbitrary number of attatchments for JSON object document? How do
you handle versioning when some types of collections need it for attachments
and others don't? 

The **dataset** command line tool and related Go package store the 
attachments unversioned by default in the pairtree. The metadata
about the attached document is stored in a sub folder `_docs`.
The unversioned attached document is stored in `v0.0.0` folder.
The attached document is stored by its basename.  The basename must 
be unique among the documents attached otherwise it will be overwritten 
when attaching another document using the same basename.

If you need versioning you MUST supply a valid [semver](https://semver.org)
when attaching the document. The metadata for the attached document will 
be in `_docs` as before but the document will be stored in a sub directory 
indicated by the semver.  The basename must be unique to the semver 
provided otherwise the document with the same basename using that semver 
will be overwritten.

It is easier to see with this example. We have a dataset collection
called "Sea-Mamals.ds". We have a JSON object stored called "walrus".
We want to attach "notes-on-walrus.docx" which is on our local
drive under "/Users/fred/Documents/notes-on-walrus.docx".

Using the **dataset** cli you issue the follow commond (in this
example without a provided semver--

```shell
    dataset attach Sea-Mamals.ds walrus /Users/fred/Documents/notes-on-walrus.docx
```

The results in a simple directory stricture for the JSON object and attachment.

```
    Sea-Mamanls/pairtree/wa/lr/us/walrus.json
    Sea-Mamanls/pairtree/wa/lr/us/_docs/notes-on-walrus.json
    Sea-Mamanls/pairtree/wa/lr/us/v0.0.0/notes-on-walrus.docx
```

In this example "notes-on-walrus-json" is the metadata for attachment
"notes-on-walrus.docx". No versioning happening.

If we had added our attachment including a semver the directory structure
will be slightly more complex.

```shell
    dataset attach Sea-Mamals.ds walrus v0.0.1 /Users/fred/Documents/notes-on-walrus.docx
```

This will cause an additional sub directories to exist (if they haven't be created
before). Our "unversioned" version still exists as v0.0.0 but now we have v0.0.1.
Our metadata file will include an href pointing to v0.0.1 and a map to all versions
including v0.0.0.

```
    Sea-Mamanls/pairtree/wa/lr/us/walrus.json
    Sea-Mamanls/pairtree/wa/lr/us/_docs/notes-on-walrus.json
    Sea-Mamanls/pairtree/wa/lr/us/v0.0.0/notes-on-walrus.docx
    Sea-Mamanls/pairtree/wa/lr/us/v0.0.1/notes-on-walrus.docx
```

If we later add a v0.0.2 of "notes-on-walrus.docx" it'd looke like

```
    Sea-Mamanls/pairtree/wa/lr/us/walrus.json
    Sea-Mamanls/pairtree/wa/lr/us/_docs/notes-on-walrus.json
    Sea-Mamanls/pairtree/wa/lr/us/v0.0.0/notes-on-walrus.docx
    Sea-Mamanls/pairtree/wa/lr/us/v0.0.1/notes-on-walrus.docx
    Sea-Mamanls/pairtree/wa/lr/us/v0.0.2/notes-on-walrus.docx
```

The single metadata file in the `_docs` subdirectory holds the metadata.
In the metadata we include an "href" string and "version_href" map. The
version_href will point to all known versions keyed by the semver. The
href string will point at the last version added, in this case v0.0.2.

IMPORTANT: If you provide the same semver and attach a file with the same
basename that previously stored version will be overwritten. Example if we
issue our original unversioned command the v0.0.0 copy of "notes-on-walrus.docx"
will be overwritten!

**dataset** attachment versioning is user driven. The only implicit version
is v0.0.0 if no semver is provided. **dataset** is not a substitute
for a version control system like [Subversion]() or [Git]() and is not
substitute for a versioned file systems like [ZFS](). If your
program needs to avoid overwriting an existing version or to "auto increment"
the semver you need to check the existing versions and decide what the
new version will be before attaching the new version of the document.

The semver versioned dircetories may contain more than one attached document.
The documents attached can be of various versions though if you attach 
more than one document at a time they will carry the same semver. This is because
the a implied semver is v0.0.0 when using the command line **dataset** otherwise
the first valid semver is used for all files being attached.

NOTE: The href in the attachments metadata always points at the last attached 
version.
