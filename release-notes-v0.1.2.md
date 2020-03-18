
From v0.1.4 release

+ Dropped support for GSheet integration
+ Only support pairtree layout of collection
+ cleaned up libdataset API focusing on removing overloaded functions

From v0.1.0 release

+ Updated libdataset API, simplified func names and normalized many of the calls (breaking change)
+ libdataset now manages opening dataset collections, inspired by Oberon System file riders (breaking change)
+ Added Python test code for libdataset to make sure libdataset works
+ Added support for check and repair when working on S3 deployed collections
+ Refactored and simplified frame behavior (breaking change)

From v0.1.0 to v0.1.1

+ Fixed problem where keys_exist called before an open command.

From v0.1.1 to v0.1.2

+ Persisting _Attachments metadata when updating with clean objects using the same technique as _Key


