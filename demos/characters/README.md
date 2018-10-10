
# The Character Demo

This is a demonstration of a "character.ds" dataset in different forms. 
Each version of the dataset is derived from the _character.csv_ file with 
different options for importing the content.

Additionally a set of example index definition files are included for
exploring _dataset indexer_ and _dataset find_.

## Try things out

Try _demo.bash_

```shell
    #!/bin/bash
    if [[ -d characters.ds ]]; then
        rm -fR characters.ds
    fi
    dataset init characters.ds
    dataset import characters.ds characters.csv 1
    dataset indexer characters.ds characters.json
    dataset indexer characters.ds names.json
    dataset indexer characters.ds emails.json
    dataset find characters.bleve "Mojo Sam"
    dataset find emails.bleve "Mojo Sam"
    dataset find emails.bleve "mojo.sam"
    dataset find emails.bleve "zbs.example.org"
    dataset find names.bleve "Mojo Sam"
    dataset find names.bleve:emails.bleve "Mojo Sam" 
    dataset find -sort='-name'  characters.bleve:names.bleve:emails.bleve "email:zbs.example.org"
```

## the files

+ characters.csv a small CSV file for generating a character collection
+ characters.json an index definition file for use with _dataset indexer_ and _dataset find_
+ names.json an index definition file for use with _dataset indexer_ and _dataset find_
+ emails.json an index definition file for use with _dataset indexer_ and _dataset find_

