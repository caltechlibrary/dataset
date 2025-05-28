
delete
======

Syntax
------

~~~shell
    dataset3 delete COLLECTION_NAME KEY
~~~

Description
-----------

- delete - removes a JSON document from primary table holding collection but records the removal in the history.
  - requires JSON document name

Usage
-----

This usage example will delete the JSON document withe the key _r1_ in 
the collection named "publications.ds".

~~~shell
    dataset3 delete publications.ds r1
~~~

