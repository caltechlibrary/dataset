
Read (end point)
================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

Retrieve a JSON document from a collection.

    `http://localhost:8485/<COLLECTION_ID>/object/<KEY>`

Requires a "GET" HTTP method.

Returns the JSON document for given `<KEY>` found in `<COLLECTION_ID>` or a HTTP error if not found.

Example
-------

Curl accessing "t1" with a key of "one"

```shell
    curl http://localhost:8485/t1/object/one
```

An example JSON document (this example happens to have an attachment) returned.

```json
{
   "_Attachments": [
      {
         "checksums": {
            "0.0.1": "bb327f7bcca0f88649f1c6acfdc0920f"
         },
         "created": "2021-10-11T11:09:51-07:00",
         "href": "T1.ds/pairtree/on/e/0.0.1/a1.png",
         "modified": "2021-10-11T11:09:51-07:00",
         "name": "a1.png",
         "size": 32511,
         "sizes": {
            "0.0.1": 32511
         },
         "version": "0.0.1",
         "version_hrefs": {
            "0.0.1": "T1.ds/pairtree/on/e/0.0.1/a1.png"
         }
      }
   ],
   "_Key": "one",
   "four": "four",
   "one": 1,
   "three": 3,
   "two": 2
}
```

