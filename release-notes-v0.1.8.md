
This release focuses on minor bug fixes in libdataset.

Removing duplicate functions:

+ `delete_frame()` has been superceded by `frame_delete()`

Renamed functions:

+ `make_objects()` has been renamed `create_objects()` to be more consistant with naming scheme.

Build Notes:

+ golang v1.14
+ Caltech library go packages
    + storage v0.1.0
    + namaste v0.0.5
    + pairtree v0.0.4
+ OS used to compiled and test
    + macOS Catalina
    + Windows 10
    + Ubuntu 18.04 LTS
+ Python 3.8 (from Miniconda 3)
+ zip has replaced tar in the releases of libdataset

Some tests fail on Windows 10 for libdataset. These will be addressed in future releases.

