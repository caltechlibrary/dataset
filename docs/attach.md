
# attach

## Syntax 

```
    dataset attach COLLECTION_NAME KEY [SEMVER] FILENAME(S)
```

## Description

Attach a file to a JSON record. Attachments are stored in a tar ball
related to the JSON record key.

## Usage

Attaching a file named *start.xlsx* to the JSON record with id _t1_ in 
collection "stats.ds"

```shell
    dataset attach stats.ds t1 start.xlsx
```

Attaching the file as version v0.0.1

```shell
    dataset attach stats.ds t1 v0.0.1 start.xlsx
```

Related topics: [attachments](attachments.html), [detach](detach.html) and [prune](prune.html)

