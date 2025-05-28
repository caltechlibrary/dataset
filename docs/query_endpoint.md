Query endpoint
==============

This provides a means to create custom lists of objects based on the contents of a dataset collection. The query is written in SQL and may return one column of data. The data returned is an array.

Query endpoints are defined in your YAML configuration used with __dataset3d__. Here's an example of the YAML used to configure a **recipes.ds** collection display. This file could be called "recipe_api.yaml".

~~~yaml
host: localhost:8483
collections:
  - dataset: recipes.ds
  - query:
    show_reciept_by_name: |
      select src
      from recipes
      order by src->>'name'
keys: true
create: true
read: true
update: true
delete: false
~~~

Running the __dataset3d__ service with "recipe_api.yaml".

~~~
dataset3d recipe_api.yaml
~~~
