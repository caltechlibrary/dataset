%dataset(1) user manual | version 2.2.7
% R. S. Doiel and Tom Morrell
% 2025-06-11

# Compatibity

AS of 2.2.7 the features around frames is no longer going to be supported. I practice the SQL query support replaced it. Also cloning and samples have been replace by SQL query support and the new "dump" and "load" feature.

As of 2.2.0 you can "dump" and "load" are availabe to create a portable
export and import using JSON lines.

As of 2.2.0 the default Dataset collection uses SQLite3 databases for 
the JSON document store. As of 2.2.1 all tests are passing again with
this change.

As of 2.2.1 dataset cli now supports dump and load verbs. This allows
a fast way to export/import an entire dataset collection as a JSONL stream.
This will likely replace cloning in the future which is considerably slower.

As of 2.2.2 libdataset is no longer available. The pain of compiling
native DLL was too high. The Dataset Project provides [datasetd](datasetd.1.md)
which exposes Dataset Collections as a JSON API and web service. This
is available to any language that support http access. Of course you
can also shell out to the cli. You have lots of option.

[ts_dataset](https://github.com/caltechlibrary/ts_dataset) is a go
example of writing a wrapper around the datasetd API.


