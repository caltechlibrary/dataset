
This release focuses on minor bug fixes in libdataset.
All functions which returned an error string only now return
True for success and False otherwise.  The error string
can be retreived with `dataset.error_message()`.

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

