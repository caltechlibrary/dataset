
# attached

## Syntax

```
    dataset attached JSON_RECORD_ID
    dataset attached JSON_RECORD_ID ATTACHMENT_NAME
```

## Description

Attached writes out (to local disc) the items that have been attached to a JSON record in the collection with
the matching JSON_RECORD_ID

## Usage

Write out all the attached files for k1

```shell
    dataset attached k1
```

Write out only the *stats.xlsx* file attached to k1

```shell
    dataset attached k1 stats.xlsx
```

Related topics: attach, attachments, and detach

