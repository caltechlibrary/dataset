join
====

Syntax
------

        dataset join [OPTION] COLLECTION_NAME KEY JSON_EXPRESSION
        dataset join [OPTION] COLLECTION_NAME KEY JSON_FILENAME
        dataset join -overwrite COLLECTION_NAME KEY JSON_FILENAME
        dataset join -i JSON_DOCUMENT_NAME COLLECTION_NAME KEY
        cat JSON_DOCUMENT_NAME | dataset join -i - COLLECTION_NAME KEY

Description
-----------

*join* will allow you to merge by appending or merge by overwriting to
an existing JSON document stored in a collection identified by KEY. With
\"append\" only new fields will be added to the record. If you specify
\"overwrite\" new fields will be added and existing fields in common
will be overwritten.

*join* is helpful in building up an aggregated record where you have a
common KEY.

Usage
-----

Let\'s assume you have a record in your collection with a key
\'jane.doe\'. It has three fields - name, email, age.

``` {.json}
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42}
```

You also have an external JSON document called profile.json. It looks
like

``` {.json}
    {"name": "Doe, Jane", "email": "jane.doe@example.edu", "bio": "world renowned geophysist"}
```

You can merge the unique fields in profile.json with your existing
jane.doe record (where the existing record id is \"jane.doe\"). The
collection is \"people.ds\"

``` {.shell}
    dataset join people.ds jane.doe profile.json
```

The result would look like

``` {.json}
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42, "bio": "renowned geophysist"}
```

If you wanted to overwrite the common fields you would use \'join
overwrite\'

``` {.shell}
    dataset join -overwrite people.ds jane.doe profile.json
```

Which would result in a record like

``` {.json}
    {"name":"Doe, Jane", "email": "jane.doe@example.edu", "age": 42, "bio": "renowned geophysist"}
```
