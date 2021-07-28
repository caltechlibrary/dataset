frame-grid
==========

USAGE
-----

```
    frame-grid COLLECTION FRAME_NAME
```

Returns the object list as a 2D array.

OPTIONS
-------

-p, -pretty
: pretty print JSON output

-use-header-row
: Include labels as a header row

Example
-------

If I have a collection named "photos.ds" and a previously
defined frame name "captions-dates-locations" I can get that
as a 2D JSON array with the following---

```
    dataset frame-grid -p photos.ds captions-dates-locations
```

The `-p` is the pretty print option for JSON output.

