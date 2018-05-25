
# The Character Demo

This is a demonstration of a "character" dataset in different forms. Each 
version of the dataset is derived from the _character.csv_ file with 
different options for importing the content.

Additionally a set of example index definition files are included for
exploring _dataset indexer_, _dataset find_ and _dsws_ utilities.

## Try things out

Try _demo.bash_

```shell
    #!/bin/bash
    if [ -d characters ]; then
        rm -fR characters
    fi
    $(dataset init characters)
    dataset -uuid import characters.csv
    dataset indexer characters.json
    dataset indexer names.json
    dataset indexer emails.json
    dataset find "Mojo Sam"
    dataset find -indexes=characters.bleve "Mojo Sam"
    dataset find -indexes=emails.bleve "Mojo Sam"
    dataset find -indexes=emails.bleve "mojo.sam"
    dataset find -indexes=emails.bleve "zbs.example.org"
    dataset find -indexes=names.bleve "Mojo Sam"
    dataset find -indexes=names.bleve:emails.bleve "Mojo Sam" 
    dataset find -sort='-name'  -indexes=characters.bleve:names.bleve:emails.bleve "email:zbs.example.org"
```

## the files

+ characters.csv a small CSV file for generating a character collection
+ characters.json an index definition file for use with _dataset indexer_ and _dataset find_
+ names.json an index definition file for use with _dataset indexer_ and _dataset find_
+ emails.json an index definition file for use with _dataset indexer_ and _dataset find_

