
# detach

## Syntax

```
    dataset detach COLLECTION_NAME KEY [SEMVER]
    dataset detach COLLECTION_NAME KEY [SEMVER] ATTACHMENT_NAME
```

## Description

_detach_ writes out (to local disc) the items that have been 
attached to a JSON record in the collection with the matching KEY

## Usage

Write out all the attached files for k1 in collection named 
"publications.ds"

```shell
    dataset detach publications.ds k1
```

Write out only the *stats.xlsx* file attached to k1

```shell
    dataset detach publications.ds k1 stats.xlsx
```

Write out only the v0.0.1 *stats.xlsx* file attached to k1

```shell
    dataset detach publications.ds k1 v0.0.1 stats.xlsx
```

Related topics: [attach](attach.html), [attachments](attachments.html), and [prune](prune.html)

