
EXAMPLES

~~~
   dataset help init

   dataset init my_objects.ds 

   dataset help create

   dataset create my_objects.ds "123" '{"one": 1}'

   cat <<EOT | dataset create my_objects.ds "345"
   {
	   "four": 4,
	   "five": "six"
   }
   EOT

   dataset update my_objects.ds "123" '{"one": 1, "two": 2}'

   dataset delete my_objects.ds "345"

   dataset keys my_objects.ds

   dataset hasKey my_objects.ds "345"

   dataset dump my_objects.ds >objects.jsonl

   dataset load my_objects.ds <objects.jsonl

   cat <<SQL | dataset query my_objects.ds 
   select json_object('key', _Key, 'version', version) as obj
   from my_objects_history
   where _Key "345"
   order by version desc
   SQL
~~~

