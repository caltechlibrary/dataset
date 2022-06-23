File system layout
==================

dataset provides a way to manage your JSON documents. This
can be done on your local disk using a technique called a
[pairtree](https://tools.ietf.org/html/draft-kunze-pairtree-01). 
Optionally can also be done by storing the JSON Objects in a 
SQL database table.

The pairtree layout is described by "collection.json" and "keymaps.json"
documents. These are located in the root folder of the collection. 

A JSON document may also have one or more "attachments". These are
stored in their own pairtree under the "attachments" sub-directory of the
collection. Attachments are supported for both pairtree storage
and SQL storage of JSON objects in a collection.

Pairtrees ensures sequence of characters in a object's key
will not collide with and is legal on common file systems. 
E.g. storing the document "hello-world.json" with the attachment
"smiles.png" in a collection named "C.ds" would result in paths like 

    `C.ds/pairtree/he/ll/o-/wo/rl/d/hello-world.json` 

and 

    `C.ds/attachments/he/ll/o-/wo/rl/d/smiles.png".

Attachments are experimental and how they are handled
may change in the future. 

In a collection's root directory you will find two or three
JSON documents.

- codemeta.json hold general metadata about the collection
- collection.json describes the operational metadata for a collection
- keymap.json holds the key to pairtree path map for pairtree
  collections (it will be missing for for SQL storage collections)

If you're using SQLite3 as your JSON document storage engine you
can choose to include your SQLite3 database file in the directory
too, if so "collection.db" is a good name. That way if you zip
up your collection and share it with a friend the JSON documents
will travel appropriately.


Pairtree
--------

The directory layout looks like:

- collection (directory on the file system)
    - [namaste](https://confluence.ucop.edu/display/Curation/Namaste) 
      records identifying the collection
        - these will get used to generate things like index.md and codemeta.json files 
    - a file, collection.json, holding metadata about the collection
    - a directory named "_frames" holding frame definitions for the 
      collection
    - a directory named "pairtree" holding the pairtree where the 
      JSON document and attachments are stored.


