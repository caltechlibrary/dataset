
Delete (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Delete a JSON document in the collection. Requires the document key and collection name.

    `http://localhost:8485/<COLLECTION_ID>/delete/<KEY>`

Requires a `GET` HTTP method.

Deletes a JSON document for the `<KEY>` in collection `<COLLECTION_ID>`. On success it returns HTTP 200 OK. Otherwise an HTTP error if creation fails.

Example
-------

The `<COLLECTION_ID>` is "t1", the `<KEY>` is "one" The content posted is

Posting using CURL is done like

```shell
    curl -X GET -H 'Content-Type: application.json' \
      http://locahost:8485/t1/delete/one
```

