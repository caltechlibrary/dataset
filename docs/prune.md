
# prune

## Syntax

```
    dataset COLLECTION_NAME prune KEY
    dataset COLLECTION_NAME prune KEY ATTACHMENT_NAME
```

## Description

prune removes all or specific attachments to a JSON document. If only
the key is supplied then all attachments are removed if an attachment
name is supplied then only the specific attachment is removed.

## Usage

In the following examples _r1_ is the KEY, *stats.xlsx* is the 
attached file. In the first example only *stats.xlsx* is removed in
the second all attachments are removed. Our collection name is "data.ds"


```shell
    dataset data.ds prune k1 stats.xlsx
    dataset data.ds prune k1
```

Related topics: [attach](attach.html), [detach](detach.html) and [attachments](attachments.html)

