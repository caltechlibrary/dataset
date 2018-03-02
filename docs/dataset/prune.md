
# prune

## Syntax

```
    dataset prune KEY
    dataset prune KEY ATTACHMENT_NAME
```

## Description

prune removes all or specific attachments to a JSON document. If only
the key is supplied then all attachments are removed if an attachment
name is supplied then only the specific attachment is removed.

## Usage

In the following examples _r1_ is the KEY, *stats.xlsx* is the 
attached file. In the first example only *stats.xlsx* is removed in
the second all attachments are removed.


```shell
    dataset prune k1 stats.xlsx
    dataset prune k1
```

Related topics: attach, detach, and attachments

