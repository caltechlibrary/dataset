status
======

Syntax
------

```shell
    dataset status COLLECTION_NAME [COLLECTION_NAME ...]
```

Description
-----------

Checks to see if a `collection.json` file is associated with 
the COLLECTION_NAME. Can work on multiple collection names. 
Returns "OK" if it is.

Usage
-----

Collection names are "MyRecordCollection.ds" and "MyBookCollection.ds".

```shell
    dataset status MyRecordCollection.ds
    dataset status MyRecordCollection.ds MyBookCollection.ds
```

Related topic: [init](init.html)

