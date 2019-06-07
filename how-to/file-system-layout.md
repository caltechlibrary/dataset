
## File system layout

dataset provides a way to organize your JSON Objects on disc. The original
was a "buckets" oriented layout. The newer and current layout is a 
[pairtree](https://tools.ietf.org/html/draft-kunze-pairtree-01). 
The layout managed/described by the collection.json document
located in the root folder of the collection. The file pairtree 
supports "attachments" by creating a sub directory next the the JSON
document. The sub directory name is `_docs`. The sequence of characters
will not collide with pairtree semantics and is legal on common
file systems.  E.g. storing the document "hello-world.json" with 
the attachment "smiles.png" in a collection named "C" would result 
in paths like 

    `C/pairtree/he/ll/o-/wo/rl/d/hello-world.json` 

and 

    `C/pairtree/he/ll/o-/wo/rl/d/_docs/smiles.png".

Attachments are experimental and how they are handled
will likely change in the future. 


## Pairtree

The directory layout looks like:

+ collection (directory on the file system)
    + [namaste](https://confluence.ucop.edu/display/Curation/Namaste) 
      records identifying the collection
        + these will get used to generate things like index.md and codemeta.json files 
    + a file, collection.json, holding metadata about the collection
    + a directory named "_frames" holding frame definitions for the 
      collection
    + a directory named "pairtree" holding the pairtree where the 
      JSON document and attachmetns are stored.


