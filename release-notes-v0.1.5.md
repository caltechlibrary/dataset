
This release focuses on refine function names, simplification
and easy of testing for Windows 10 deployments.

Build Notes:

+ Golang v1.14
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

Renamed functions:

+ collection_status() is now collection_exists()

Depreciated functions and features:

+ S3, Google Cloud Storage support dropped.
+ grid(), if you need this create a frame first and use frame_grid().

Some tests fail on Windows 10 for libdataset. These will be addressed in future releases.

