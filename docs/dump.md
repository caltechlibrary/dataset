
dump
====

This will export an intire collection as a JSON lines document. Useful for migrating content between dataaset versions or collections. It also can be used with tools like `jq` to filter. `dump` is the counter part to `load`.

`dump` and `load` replace collection cloning as they are storage agnostic and faster than the cold clone methods. Dumps can be easily processed with most data science tools.

Example
-------

Dump "data.ds" to "data.jsonl". Next pipe the result through `jq`.  

~~~shell
dataset dump data.ds >data.jsonl
cat data.jsonl | jq .
~~~



