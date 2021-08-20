Datasetd
========

__datasetd__ is a minimal web service based on the command line tool called dataset. It exposes dataset collections over HTTP/HTTPS. It's only supports a extreme subset of the command line functionality. If you need more than that in a web based object storage I suggest looking at existing mature NoSQL style solutions like Couchbase, MongoDB, MySQL (which now supports JSON objects) or Postgres (which also support JSON objects).

__datasetd__ can make one or more dataset collections visible over HTTP/HTTPS hosted in a common directory. NOTE: it can server collections in read only mode (i.e. not supporting "init", "create", "update" and "delete"). In this way you could vend our a public dataset collection.

Normally __datasetd__ intended to run on localhost behind another web service which would control access and activities. If you need something more than this look at other solutions such as Couchbase, MongoDB, MySQL (with JSON support) or Postgres (with JSON support).

__datasetd__ provide a limitted subset of actions support by the standard datset command line tool. It only supports the following verbs

1. init (create a new collection)
2. keys (return a list of all keys in the collection)
3. create
4. read
5. update
6. delete

__datasetd__ is intended to be combined with other services like Solr 8.9.
__datasetd__ only implements the simplest of object storage.
