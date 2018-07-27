
# detach

## Syntax

```
    dataset COLLECTION_NAME detach KEY
    dataset COLLECTION_NAME detach KEY ATTACHMENT_NAME
```

## Description

_detach_ writes out (to local disc) the items that have been 
attached to a JSON record in the collection with the matching KEY

## Usage

Write out all the attached files for k1 in collection named 
"publications.ds"

```shell
    dataset publications.ds detach k1
```

Write out only the *stats.xlsx* file attached to k1

```shell
    dataset publications.ds detach k1 stats.xlsx
```

Related topics: [attach](attach.html), [attachments](attachments.html), and [prune](prune.html)

