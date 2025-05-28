Storage Engines
===============

With the introduction of v3 of dataset we only use SQL databases with JSON column support for storage engines. SQLite3 is the currently supported. PostgreSQL and MySQL 8 are possible but not web tested yet.

- SQL Storage (SQLite3), experimental but more performant

With the introduction of SQL Storage dataset can be used in a multi-process/multi-user mode via a RESTful API.  The SQL storage is experimental and as it gets you more various considerations are coming to the surface

- SQLite3 works fine for single process multi user (via the web API) 
