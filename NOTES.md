
# Attachment ideas

S3/Google Cloud Storage brings additional overhead for attachments it works 
like a k/v store where the operations are on files not within them.

## Naive implementation steps for S3/Google Cloud Storage

### Attach (no other attachments)

1. calc basename of file to be attached as well as the 
   pairtree path including a `_docs` element before the basename
2. copy the file into place on attachment directory using the basename

### List attachments

1. scan for filenames using pairtree path plus `_docs` suffix 

### Delete specific attached file

1. calc pairtree path
2. delete item with path

### Delete all attachments

1. remove objects attachments from the pairtree path for containing `_docs` 

## Reference Google API integration

+ [Google Sheets API v4](https://developers.google.com/sheets/)
    + [REST methods](https://developers.google.com/sheets/api/reference/rest/)
    + [Golang Quickstart docs](https://developers.google.com/sheets/api/quickstart/go)
        + where to go to setup credentials and project specifics, we're using Go for our project

## Extending dataset's reach with shared libraries

### Python

Use [py_dataset](https://github.com/caltechlibrary/py_dataset).


### Julia

### R

+ [Writing R Extensions](https://cran.r-project.org/doc/manuals/R-exts.html)

## Metadata for collections

+_ ANVL/ERC are related to Namaste, these could be included in a collections-info.txt file that intern would then be expressed as codemeta.json, CATALOG.json and index.html
    + ERC: is human editable in a simple text editor, fields could be supplied collectively or individually, simplifying further the curration of the metadata, ERC is similar to the expression of Namaste focusing on who, what, whem, where and can be extended in a like manner
+ THUMP would be an interesting query option to support in addition to a simple REST API for listing keys, returning lists of objects or full objects


## Namaste support

+ initial implementation is to replace metadata, but if we called out to an editor we could implement editable metadata (e.g. write data to tmp file, read in with a restricted editor like nano, red, rvi, then recieve update)

