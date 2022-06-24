
Prune (end point)
=================

Removes an attached document from a JSON record using `<KEY>` and `<FILENAME>`.

    `http://localhost:8485/<COLLECTION_ID>/attachment/<KEY>/<SEMVER>/<FILENAME>`

Requires a DELETE method. Returns an HTTP 200 OK on success or an HTTP error code if not.

Example
-------

In this example `<COLLECTION_ID>` is "t1", `<KEY>` is "one", and `<FILENAME>` is "a1.png". Once again our example uses curl.

```shell
    curl -X DELETE http://localhost:8485/t1/attachment/one/a1.png
```

This will cause the attached file to be removed from the record
and collection.


