Storage Engines
===============

With the introduction of v2 of dataset you now have a choice of storage
engines.  Currently offered in the 2.0 release is

- pairtree (the engine implementation of v1, stable)
- SQL Storage (via MySQL 8, Postgres 14 and SQLite3 ), experimental but more performant

The pairtree storage engine is stable. The primary limitations are
the file system (where it stores the JSON documents) case limitations and lack of record/field locking. Pairtree work fine for batch operations, single user/single process operations with less than 100k documents.

With the introduction of SQL Storage dataset can be used in a multi-process/multi-user mode via a RESTful API.  The SQL storage is experimental and as it gets you more various considerations are coming to the surface

- SQLite3 works fine for single process multi user (via the web API) 
- MySQL 8 works well when you need multi-user and multi-process acecss
- Postgres 14 is relatively untested so likely has quirks not yet addressed in the 2.0 release of dataset


