
Create (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Create a JSON document in the collection. Requires a unique key in the URL and the content most be JSON less than 1 MiB in size.

    `http://localhost:8485/<COLLECTION_ID>/object/<KEY>`

Requires a "POST" HTTP method with.

Creates a JSON document for the `<KEY>` in collection `<COLLECTION_ID>`. On success it returns HTTP 201 OK. Otherwise an HTTP error if creation fails.

The "POST" needs to be JSON encoded and using a Content-Type of "application/json" in the request header.

Example
-------

The `<COLLECTION_ID>` is "t1", the `<KEY>` is "one" The content posted is

```json
    {
       "one": 1
    }
```

Posting using CURL is done like

```shell
    curl -X POST -H `Content-Type: application.json` \
      -d `{"one": 1}` \
      http://locahost:8485/t1/object/one
```

