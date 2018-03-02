
# detach

## Syntax

```
    dataset detach KEY
    dataset detach KEY ATTACHMENT_NAME
```

## Description

detach writes out (to local disc) the items that have been attached to a JSON record in the collection with
the matching KEY

## Usage

Write out all the attached files for k1

```shell
    dataset detach k1
```

Write out only the *stats.xlsx* file attached to k1

```shell
    dataset detach k1 stats.xlsx
```

Related topics: attach, attachments, and prune

