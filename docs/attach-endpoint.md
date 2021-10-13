
Attach (end point)
==================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Attaches a document to a JSON Document using `<KEY>`, `<SEMVER>` and `<FILENAME>`.

    `http://localhost:8485/<COLLECTION_ID>/attach/<KEY>/<SEMVER>/<FILENAME>`

Requires a "POST" method. The "POST" is expected to be a multi-part web form providing the source filename in the field "filename".  The document will be written the JSON document directory by `<KEY>` in sub directory indicated by `<SEMVER>`.

See https://semver.org/ for more information on semantic version numbers.

Example
=======

In this example the `<COLLECTION_ID>` is "t1", the `<KEY>` is "one" and
the content upload is "a1.png" in the home directory "/home/jane.doe".
The `<SEMVER>` is "0.0.1".

```shell
    curl -X POST -H 'Content-Type: multipart/form-data' \
       -F 'filename=@/home/jane.doe/a1.png' \
       http://localhost:8485/t1/attach/one/0.0.1/a1.png
```

NOTE: The URL contains the filename used in the saved attachment. If
I did not want to call it "a1.png" I could have provided a different
name in the URL path.

