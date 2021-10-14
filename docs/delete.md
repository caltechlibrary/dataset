delete
======

Syntax
------

```shell
    dataset delete COLLECTION_NAME KEY
```

Description
-----------

- delete - removes a JSON document from collection
  - requires JSON document name

Usage
-----

This usage example will delete the JSON document withe the key _r1_ in 
the collection named "publications.ds".

```shell
    dataset delete publications.ds r1
```

Related topics: [create](create.html), [read](read.html), and [update](update.html)

