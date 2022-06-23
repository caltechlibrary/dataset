
attach
======

Syntax
------

```shell
    dataset attach COLLECTION_NAME KEY FILENAME(S)
```

Description
-----------

Attach a file to a JSON record. Attachments are stored in a tar ball
related to the JSON record key.

Usage
-----

Attaching a file named *start.xlsx* to the JSON record with id _t1_ in 
collection "stats.ds"

```shell
    dataset attach stats.ds t1 start.xlsx
```


