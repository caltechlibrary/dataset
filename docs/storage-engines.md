Storage Engines
===============

With the introduction of v2 of dataset you now have a choice of storage engines. In v2.2 SQLite3 became the default storage engine.  Current supported storage engines are.

- SQL Storage (via MySQL 8, Postgres >= 12 and SQLite3 >= 3.4)
  - SQLite3 is the default storage engine since it requires no configuration
- pairtree (the original storage engine of v1, depricated)

With the introduction of SQL Storage dataset can be used in a multi-process/multi-user mode via a RESTful API.  The SQL storage is experimental and as it gets you more various considerations are coming to the surface

- SQLite3 works fine across single and multie process scenarios. It has the advantage of requiring no configuration and avoiding the need to run a database manage system such as MySQL or Postgres. 
- Postgres 14.5 is preferred and most test implementation when a database management system is needed.
- MySQL 8 works well when you need multi-user and multi-process access. (depreciated)

Cautions
--------

The pairtree storage engine is stable. The primary limitations are the file system (where it stores the JSON documents) case limitations and lack of record/field locking. Pairtree work fine for batch operations, single user/single process operations with less than 100k documents. It does not appropriate for concurrent access involving writes. I remains in v2.2 as a historical artifact. Will be removed in v3.

MySQL is used less and less in the software I work with. MySQL support maybe dropped in version 3 of dataset.
