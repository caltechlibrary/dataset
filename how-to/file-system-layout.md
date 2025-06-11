File system layout
==================

dataset provides a way to manage your JSON documents. In version we usually use a SQL engine for storage. SQLite3 is the default in v2.2. The means you can have your JSON objects on disk and easily preservable without running a database management system. In version one the default store was a [pairtree](https://tools.ietf.org/html/draft-kunze-pairtree-01).  While  those are still supported in v2.2 they are depreciated and will not be support in v3.

The layout of a dataset collection on disk is as follows

~~~
- COLLECTION_NAME
  - collection.json (holds operational metadata about the collection)
  - codemeta.json (holds citation metadata about the collection)
  - collection.db (for SQLite3, the SQL data)
  - attachments (optional, holds a pairtree of attached files)
~~~

If the legacy pairtree implementation you would have a directory named "pairtree" holding the JSON objects.

This simple data structure is easy to [Bag](https://en.wikipedia.org/wiki/BagIt). The collection as a whole would be placed the the "data" sub directory of the bag and along with the required metadata files supporting the Bag format.

If you are using MySQL or PostgreSQL to store your data then you have two options when bagging. You could dump the collection contents along side the collections directory or you could created an SQLite3 version of the collection (trivial using the dump and load feature of v2.2) and Bag that.
