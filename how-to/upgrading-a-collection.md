Upgrading your dataset collection(s)
====================================

The __dataset__ provides a "dump" and "load" action for a collection since version v2.2. This replaces the need for a promblematic functions such as clone, clone-sample, and dsimporter. The dump operations renders a __dataset__ collection as a [JSON lines](https://jsonlines.org) stream. This can easily be redirected from standard out to a file.  The load operation is the reverse, it reads a dump JSON lines stream and creates records in the collection. Starting in the upcoming v3 of dataset the command line program will include the major version number. That means when it is time to migrate your v2 collection to a v3 collection it can be done as a very simple operation using a standard pipe. This is true if you're using a Unix system or Windows and PowerShell.

~~~shell
dataset dump data_v2.ds | dataset3 load data_v3.ds
~~~

The JSON lines dump file are far faster than cloning or repairing a repository which was the old method of migrasting content. It also opens up opportunties to use `jq` or other JSON filters to process the content. This could simplify metadata migrations.
