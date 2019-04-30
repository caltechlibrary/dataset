
## Python Basics

This is an example of creating a dataset called *fiends.ds*, saving
a record called "littlefreda.json" and reading it back.

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
the collection name as the first parameter.  Likewise many also return 
a tuple where the first value is the object you are fetching and the 
second part of the tuple is any error messages. 

Now check to see if the key, littlefreda, is in the collection

```python
   dataset.haskey(c_name, 'littlefreda')
```

You can also read your JSON formatted data from a file or standard 
input.  In this example we are creating a mojosam record and reading 
back the contents of fiends.ds

```python
   dataset -i mojosam.json create mojosam
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

```python
   cat capt-jack.json | dataset create capt-jack
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

