
Update (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Replace a JSON document in the collection. Requires the key of an existing document in the URL and the content must be JSON less than 1 MiB in size.

    `http://localhost:8485/<COLLECTION_ID>/object/<KEY>`

Update requires using the "PUT" HTTP method.

Replaces the JSON document for the `<KEY>` in collection `<COLLECTION_ID>`. On success it returns a JSON status object and HTTP 200 OK. Otherwise an HTTP error if the update fails.

The "PUT" needs to be JSON encoded and using a Content-Type of "application/json" in the request header.

If `<KEY>` does not exist the default answer (as in datasetd 2.5.1) is HTTP 200 and nothing is created. When the service is run with `strict_status: true` in its YAML settings it is HTTP 404 with a JSON error object. Malformed JSON is HTTP 400, an unsupported content type is HTTP 415, a body over 1 MiB is HTTP 413 and a record that fails schema validation is HTTP 422 in strict mode. See [datasetd_api(5)](datasetd_api.5.md), "HTTP status codes".

Example
-------

The `<COLLECTION_ID>` is "t1", the `<KEY>` is "one". The content sent is

```json
    {
       "one": 1,
       "two": 2
    }
```

Sending it using curl is done like

```shell
    curl -X PUT -H 'Content-Type: application/json' \
      http://localhost:8485/t1/object/one \
      -d '{"one": 1, "two": 2}'
```

