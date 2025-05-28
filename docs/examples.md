
# EXAMPLES

~~~
   dataset3 help init

   dataset3 init my_objects.ds 

   dataset3 help create

   dataset3 create my_objects.ds "123" '{"one": 1}'

   cat <<EOT | dataset3 create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   dataset3 update my_objects.ds "123" '{"one": 1, "two": 2}'

   dataset3 delete my_objects.ds "345"

   dataset3 keys my_objects.ds

   dataset3 hasKey my_objects.ds "345"

   dataset3 dump my_objects.ds >objects.jsonl

   dataset3 load my_objects.ds <objects.jsonl

   cat <<SQL | dataset3 query my_objects.ds 
   select json_object('key', _Key, 'version', version) as obj
   from my_objects_history
   where _Key "345"
   order by version desc
   SQL
~~~

