Python Basics
-------------

This is an example of creating a dataset called *fiends.ds*, saving a
record called \"littlefreda.json\" and reading it back.

```python
    import sys
    import json
    from py_dataset import dataset

    c_name = 'friends.ds'
    err = dataset.init(c_name)
    if err != '':
        print(f"init error, {err}")
        sys.exit(1)
    key = 'littlefreda'
    record = {"name":"Freda","email":"little.freda@inverness.example.org"}
    err = dataset.create(c_name, key, record)
    if err != '':
        print(f"create error, {err}")
        sys.exit(1)
    keys = dataset.keys(c_name)
    for key in keys:
        p = dataset.path(c_name, key)
        print(p)
        record, err := dataset.read(c_name, key)
        if err != '':
            print(f"read error, {err}")
            sys.exit(1)
        print(f"Doc: {record}")
```

Notice that the command `dataset.init(c_name)` and
`dataset.create(c_name, key)`. Many of the dataset command will require
the collection name as the first parameter. Likewise many also return a
tuple where the first value is the object you are fetching and the
second part of the tuple is any error messages.

Now check to see if the key, littlefreda, is in the collection

```python
   dataset.haskey(c_name, 'littlefreda')
```

You can also read your JSON formatted data from a file but you need to
convert it first to a Python dict. In theses examples we are creating
for Mojo Sam and Capt. Jack then reading back all the keys and
displaying their paths and the JSON document created.

```python
    with open("mojosam.json") as f:
        src = f.read().encoding('utf-8')
        dataset.create(c_name, "mojosam", json.loads(src))

   with open("capt-jack.json") as f:
      src = f.read()
      dataset.create("capt-jack", json.loads(src))

   for key in dataset.keys(c_name):
        print(f"Path: {dataset.path(c_name, key)}")
        print(f"Doc: {dataset.read(c_name, key)}")
        print("")
```

It is also possible to filter and sort keys from python by providing
extra parameters to the keys method. First we\'ll display a list of keys
filtered by email ending in \"example.org\" then sorted by email.

```python
    print(f"Filtered only")
    keys = dataset.keys(c_name, '(has_suffix .email "example.org")')
    for key in keys:
        print(f"Path: {dataset.path(c_name, key)}")
        print(f"Doc: {dataset.read(c_name, key)}")
        print("")
    print(f"Filtered and sorted") 
    keys = dataset.keys(c_nane, '(has_suffix .email "example.org")', '.email')
    for key in keys:
        print(f"Path: {dataset.path(c_name, key)}")
        print(f"Doc: {dataset.read(c_name, key)}")
        print("")
```

Filter and sorting a large collection can take time due to the number of
disc reads. It can also use allot of memory. It is more effecient to
first filter your keys then sort the filtered keys.

```python
    print(f"Filtered, sort by stages")
    all_keys = dataset.keys(c_name)
    keys = dataset.key_filter(c_name, keys, '(has_suffix .email "example.org")')
    keys = dataset.key_sort(c_name, keys, ".email")
    for key in keys:
        print(f"Path: {dataset.path(c_name, key)}")
        print(f"Doc: {dataset.read(c_name, key)}")
        print("")
```
