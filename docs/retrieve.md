retrieve
========

Syntax
------

```shell
    dataset retrieve COLLECTION_NAME KEY ATTACHMENT_NAME
```

Description
-----------

__retrieve__ writes out (to local disc) the items that have been 
attached to a JSON record in the collection with the matching KEY

Usage
-----

Write out all the attached files for k1 in collection named 
"publications.ds"

```shell
    dataset retrieve publications.ds k1
```

Write out only the *stats.xlsx* file attached to k1

```shell
    dataset retrieve publications.ds k1 stats.xlsx
```

