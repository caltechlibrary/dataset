
load
====

This will read a JSON lines document holding an object made of a "key" attribute and an "object" attribute and populate the collection using those objects. By default objects are not overwritten but there is an option for allowing that.

Example
-------

Load the objects from "data.jsonl" into the collection called "data.ds". Then update the collection using "new-data.jsonl" overwriting existing objects in "data.ds"

~~~shell
cat data.jsonl | dataset load data.ds
cat new-data.jsonl | dataset load --overwrite data.ds
~~~



