
Keys (end point)
================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

This end point lists keys available in a collection.

~~~
    http://localhost:8485/<COLLECTION_ID>/keys
~~~

Requires a "GET" method.

The keys are turned as a JSON array or http error if not found.

Example
-------

In this example `<COLLECTION_ID>` is "t1".

~~~shell
    curl http://localhost:8485/t1/keys
~~~

The document return looks some like

~~~json
    [
        "one",
        "two",
        "three"
    ]
~~~

For a "t1" containing the keys of "one", "two" and "three".
