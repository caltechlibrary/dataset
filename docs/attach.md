
# attach

## Syntax 

```
    dataset attach COLLECTION_NAME KEY FILENAME(S)
```

## Description

Attach a file to a JSON record. Attachments are stored in a tar ball
related to the JSON record key.

## Usage

Attaching a file named *start.xlsx* to the JSON record with id _t1_ in 
collection "stats.ds"

```shell
    dataset stats.ds attach t1 start.xlsx
```

Related topics: [attachments](attachments.html), [detach](detach.html) and [prune](prune.html)

