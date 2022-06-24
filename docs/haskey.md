has-key
=======

Syntax
------

```shell
    dataset [OPTIONS] has-key COLLECTION_NAME KEY_TO_CHECK_FOR
```

Description
-----------

Checks if a given key is in the a collection. Returns "true" if 
found, "false" otherwise. The collection name is "people.ds"

Usage
-----

```shell
    dataset has-key people.ds '0000-0003-0900-6903'
    dataset has-key people.ds r1
```

