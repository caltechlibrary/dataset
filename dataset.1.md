---
title: "dataset (1) user manual"
pubDate: 2023-02-08
author: "R. S. Doiel"
---

# NAME

dataset 

# SYNOPSIS

dataset [GLOBAL_OPTIONS] VERB [OPTIONS] COLLECTION_NAME [PRAMETER ...]

# DESCRIPTION

dataset command line interface supports creating JSON object
collections and managing the JSON object documents in a collection.

When creating new documents in the collection or updating documents
in the collection the JSON source can be read from the command line,
a file or from standard input.

# SUPPORTED VERBS

- help will give documentation of help on a verb, e.g. "help create"
- create, creates a new document in the collection
- read, retrieves a document from the collection writing it standard out
- update, updates a document in the collection
- delete, removes a document from the collection
- list, returns a list of keys in the collection
- codemeta, copies metadata a codemeta file and updates the collections metadata
- info, returns the metadata associated with collection
- import, imports another collecting into the current one
- export, exports the collection into another collection
- attach, attaches a document to a JSON object record
- attachments, lists the attachments associated with a JSON object record
- retrieve, creates a copy local of an attachement in a JSON record
- prune, removes and attachment from a JSON record
- frames, lists the frames defined in a collection
- frame, will add a data frame to a collection if a definition is provided or return an existing frame if just the frame name is provided
- reframe, will recreate a frame using its existing definition but replacing objects based on a new set of keys provided
- refresh, will update all objects currently in the frame based on the current state of the collection. Any keys deleted in the collection will be delete from the frame.
- delete-frame, will remove a frame from the collection
- has-frame, will return true (exit 0) if frame exists, false (exit 1) if not
- attachments, will list any attachments for a JSON document
- attach, will add an attachment to a JSON document
- retrieve, will copy out the attachment to a JSON document 
  into the current directory 
- prune, will remove an attachment from the JSON document
- versions, will list the versions known for a JSON document if versioning is enabled for collection
- read-version, will return a specific version of a JSON document if versioning is enabled for collection
- versioning,  will set the versioning of a collection, can be "none", "major", "minor", or "patch"

A word about "keys". dataset uses the concept of key/values for
storing JSON documents where the key is a unique identifier and the
value is the object to be stored.  Keys must be lower case 
alpha numeric only.  Depending on storage engines there are issues
for keys with punctation or that rely on case sensitivity. E.g. 
The pairtree storage engine relies on the host file system. File
systems are notorious for being picky about non-alpha numeric
characters and some are not case sensistive.

A word about "GLOBAL_OPTIONS" in v2 of dataset.  Originally
all options came after the command name, now they tend to
come after the verb itself. This is because context counts
in trying to remember options (at least for the authors of
dataset).  The are three "GLOBAL_OPTIONS" that are exception
and they are `-version`, `-help`
and `-license`. All other options come
after the verb and apply to the specific action the verb
implements.


# OPTIONS

-help
: display help

-license
: display license

-version
: display version

# EXAMPLES

~~~
   dataset help init

   dataset init my_objects.ds 

   dataset help create

   dataset create my_objects.ds "123" '{"one": 1}'

   dataset create my_objects.ds "234" mydata.json 
   
   cat <<EOT | dataset create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   dataset update my_objects.ds "123" '{"one": 1, "two": 2}'

   dataset delete my_objects.ds "345"

   dataset keys my_objects.ds
~~~

dataset 2.1.0


