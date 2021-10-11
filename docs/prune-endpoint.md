Prune (end point)
=================

Removes an attached document from a JSON Document using `<KEY>`, `<SEMVER>` and `<FILENAME>`.

    `http://localhost:8485/<COLLECTION_ID>/attach/<KEY>/<SEMVER>/<FILENAME>`

Requires a GET method. Returns an HTTP 200 OK on success or an HTTP error code if not.

See https://semver.org/ for more information on semantic version numbers.

Example
-------

In this example `<COLLECTION_ID>` is "t1", `<KEY>` is "one", `<SEMVER>` is "0.0.1" and `<FILENAME>` is "a1.png". Once again our example uses curl.

```
    curl http://localhost:8485/t1/prune/one/0.0.1/a1.png
```

This will cause the attached file to be removed from the record
and collection.

