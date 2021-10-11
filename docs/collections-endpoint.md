Collections (end point)
=======================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

This provides a JSON list of collections available from the running __datasetd__ service.

Example
=======

The assumption is that we have __datasetd__ running on port "8485" of "localhost" and a set of collections, "t1" and "t2", defined in the "settings.json" used at launch.

```{.json}
    [
      "t1",
      "t2"
    ]
```


