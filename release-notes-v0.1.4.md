
This release has breaking changes with release v0.1.3 and early.
Many functions in libdataset have been renamed to prevent collisions
in the environments using libdataset C-shared library. Most function
names now have two parts separated by a underscore (e.g. status
has become collection_status, repair has become collection_repair).

Google Sheet integration has been dropped. It was just more trouble
then it was worth to maintain.

Tests from py_dataset now have been ported to the test library for
libdataset.

Redundant functions have been removed (we had accumulated multiple 
definitions for the same thing in libdataset). Where possible
code has been simplified.

Most libdataset functions will cause an "open" on a dataset collection
automatically. Some additional functions around collections have been
added primarily to make testing eaisier (e.g. open_collection(), is_open(),
close_collection(), close_all()).

Functions that were overloaded via optional parameters have been simplified.
E.g. keys() now returns all keys in collection, use key_filter() and key_sort() accordingly.

