
Read (end point)
================

Interacting with the __dataset3d__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Retrieve a JSON document from a collection.

~~~
    http://localhost:8485/<COLLECTION_ID>/object/<KEY>
~~~

Requires a "GET" HTTP method.

Returns the JSON document for given `<KEY>` found in `<COLLECTION_ID>` or a HTTP error if not found.

Example
-------

Curl accessing "t1" with a key of "one"

~~~shell
    curl http://localhost:8485/t1/object/one
~~~

An example JSON document.

~~~json
{
   "key": "one.two.three",
   "four": "four",
   "one": 1,
   "three": 3,
   "two": 2
}
~~~

