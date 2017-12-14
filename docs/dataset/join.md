
# join

## Syntax

```
    dataset join JOIN_TYPE JSON_RECORD_ID [JSON_EXPRESSION]
    dataset -i JSON_DOCUMENT_NAME join JOIN_TYPE JSON_RECORD_ID
    cat JSON_DOCUMENT_NAME | dataset join JOIN_TYPE JSON_RECORD_ID
```

## Description

join will allow you to merge (updating or overwriting) an existing JSON document stored in
a collection. Two JOIN_TYPES are available -- "update" and "overwrite".  With "update" 
only new fields will be added to the record. If you specify "overwrite" new fields will be 
added and existing fields in common will be overwritten.

join us helpful in building up an aggregated record where you have a common JSON_RECORD_ID.

## Usage

Let's assume you have a record in your collection with a key 'jane.doe'. It has
three fields - name, email, age.

```json
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42}
```

You also have an external JSON document called profile.json. It looks like

```json
    {"name": "Doe, Jane", "email": "jane.doe@example.edu", "bio": "world renowned geophysist"}
```

You can merge the unique fields in profile.json with your existing jane.doe record
(where the existing record id is "jane.doe").

```shell
    dataset -i profile.json join update jane.doe
```

The result would look like

```json
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42, "bio": "renowned geophysist"}
```

If you wanted to overwrite the common fields you would use 'join overwrite'

```shell
    dataset -i profile.json join overwrite jane.doe
```

Which would result in a record like

```json
    {"name":"Doe, Jane", "email": "jane.doe@example.edu", "age": 42, "bio": "renowned geophysist"}
```

