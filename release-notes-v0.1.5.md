
This release focuses on refine function names, simplification
and easy of testing for Windows 10 reployments.

Build Notes:

+ golang v1.14
+ Python 3.8 (from Miniconda 3)
+ zip has replaced tar in the releases of libdataset

Renamed functions:

+ collection_status() is now collection_exists()

Depreciated functions and features:

+ S3, Google Cloud Storage support dropped.
+ grid(), if you need this create a frame first and use frame_grid().


