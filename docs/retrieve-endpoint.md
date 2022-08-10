
Retrieve (end point)
====================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Retrieves an s attached document from a JSON record using `<KEY>` and `<FILENAME>`.

    `http://localhost:8485/<COLLECTION_ID>/attach/<KEY>/<FILENAME>`

Requires a POST method and expects a multi-part web form providing the filename. The document will be written the JSON document directory by `<KEY>` in attachments sub directory using a pairtree path.

Example
-------

In this example we`re retieving the `<FILENAME>` of "a1.png" into `<COLLECTION_ID>` of "t1" and `<KEY>`
of "one" using curl.

```shell
    curl http://localhost:8485/t1/retrieve/one/a1.png
```

This should trigger a download of the "a1.png" image file in the
collection for document "one".

