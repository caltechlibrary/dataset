
# attachments

## Syntax

```
    dataset COLLECTION_NAME attachments KEY
```

## Description

List the files attached to the JSON record matching the KEY
in the collection.

## Usage

List all the attachments for _k1_ in collection "stats.ds".

```shell
    dataset stats.ds attachments k1
```

Related topics: [attach](attach.html), [detach](detach.html) and [prune](prune.html)

