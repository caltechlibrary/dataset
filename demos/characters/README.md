
# The Character Demo

This is a demonstration of a "character" dataset in different forms. Each 
version of the dataset is derived from the _character.csv_ file with 
different options for importing the content.

Additionally a set of example index definition files are included for
exploring _dsindexer_, _dsfind_ and _dsws_ utilities.

## Try things out

Try _demo.bash_

```shell
    #!/bin/bash
    if [ -d characters ]; then
        rm -fR characters
    fi
    $(dataset init characters)
    dataset -uuid import characters.csv
    dsindexer characters.json
    dsindexer names.json
    dsindexer emails.json
    dsfind "Mojo Sam"
    dsfind -indexes=characters.bleve "Mojo Sam"
    dsfind -indexes=emails.bleve "Mojo Sam"
    dsfind -indexes=emails.bleve "mojo.sam"
    dsfind -indexes=emails.bleve "zbs.example.org"
    dsfind -indexes=names.bleve "Mojo Sam"
    dsfind -indexes=names.bleve:emails.bleve "Mojo Sam" 
    dsfind -sort='-name'  -indexes=characters.bleve:names.bleve:emails.bleve "email:zbs.example.org"
```

Or try the website version (you'll need a web browser)

```shell
    #!/bin/bash
    if [ -d characters ]; then
        rm -fR characters
    fi
    $(dataset init characters)
    dataset -uuid import characters.csv
    dsindexer characters.json
    dsindexer names.json
    dsindexer emails.json
    dsws . characters.bleve names.bleve emails.bleve
```

## the files

+ characters.csv a small CSV file for generating a character collection
+ characters.json an index definition file for use with _dsindexer_ and _dsfind_
+ names.json an index definition file for use with _dsindexer_ and _dsfind_
+ emails.json an index definition file for use with _dsindexer_ and _dsfind_

