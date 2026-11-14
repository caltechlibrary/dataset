
Create (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Create a JSON document in the collection. Requires a unique key in the URL and the content most be JSON less than 1 MiB in size.

    `http://localhost:8485/<COLLECTION_ID>/object/<KEY>`

Create requires using the "POST" HTTP method with.

Creates a JSON document for the `<KEY>` in collection `<COLLECTION_ID>`. On success it returns a JSON status object and HTTP 201 Created. Otherwise an HTTP error if creation fails.

If `<KEY>` already exists the default answer (as in datasetd 2.5.1) is HTTP 200 and the existing document is replaced. When the service is run with `strict_status: true` in its YAML settings a JSON POST to a key that exists is HTTP 409 and nothing is written, so a repeated create cannot destroy the record that was there; use PUT (see the update end point) to replace a document. Malformed JSON is HTTP 400, an unsupported content type is HTTP 415 and a body over 1 MiB is HTTP 413 in strict mode, and a record that fails schema validation is HTTP 422. See [datasetd_api(5)](datasetd_api.5.md), "HTTP status codes".

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
    curl -X POST -H 'Content-Type: application.json' \
      http://locahost:8485/t1/object/one \
      -d '{"one": 1}'
```

