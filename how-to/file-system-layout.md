File system layout
==================

dataset provides a way to organize your JSON Objects. This
can be done on your local disk using a technique called a
[pairtree](https://tools.ietf.org/html/draft-kunze-pairtree-01). 
Optionally can also be done by storing the JSON Objects in a 
SQL database table.

The pairtree layout is described by the collection.json and 
keymaps.json documents. These are located in the root folder of
the collection. 

A JSON document may also have "attachments". These are stored in
their own pairtree under the "attachments" sub-directory of the
collection. 

The sequence of characters will not collide with pairtree semantics
and is legal on common file systems.  E.g. storing the document
"hello-world.json" with the attachment "smiles.png" in a collection
named "C.ds" would result in paths like 

    `C.ds/pairtree/he/ll/o-/wo/rl/d/hello-world.json` 

and 

    `C.ds/attachments/he/ll/o-/wo/rl/d/smiles.png".

Attachments are experimental and how they are handled
may change in the future. 

In a collection's directory you will find two or three
JSON documents.

- codemtea.json hold general metadata about the collection
- collection.json describes the operational metadata for a collection
- keymap.json hold a key to pairtree path map for pairtree collections (it will be missing for for SQL storage collections)

If you're using SQLite3 as your JSON document storage engine you
can choose to include your SQLite3 database file in the directory
too, if so "collection.db" is a good name.

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


