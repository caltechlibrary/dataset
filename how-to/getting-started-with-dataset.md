
# Getting started with dataset

_dataset_ is a set of tools for managing JSON (object) documents as a collection of key/value pairs stored on either your
local file system, AWS S3 or Google Cloud Storage. These documents can be interated over or retrieved individually.
There is also a full text indexer for supporting fielded or full text searches based on the index definitions.
One final feature of _dataset_ is the ability to add attachments to your JSON objects. These attachments are stored
in a simple archive format called [tar](https://en.wikipedia.org/wiki/Tar_(computing)). Basic metadata can be retrieved
, and the attachments can be retreive as a group or individually. Attachments can be removed.


## Getting dataset onto your computer

The command line _dataset_ is available for installation from https://github.com/caltechlibrary/dataset/releases/latest.
Find the zip file associated with your computer type and operating system then download it. Once downloaded you can unzip the zip
file and copy the programs into a local directory called "bin" on your comptuer. For full instructions on installation see
[INSTALL.md](../install.html). In addition to the command line tool a Python 3.6 package is also provide and can
be installed with the usual `python3 setup.py install --user --record files.txt`.


## Basic workflow with dataset

_dataset_'s focus is in storing JSON (object) documents in collections. The documents are stored in a bucketed directory structure and
named for the "key" provided. The documents remain plain text JSON on disc. When you first start working with a dataset you
will need to initialize the collection. This creates the bucket directories and associated metadata so you can easily
retrieve your documents. If you were to initialize a dataset collection called "FavoriteThings.ds" it would look like --

```shell
    dataset init FavoriteThings.ds
```

or in Python

```python
    import dataset

    dataset.init('FavoriteThings.ds')
```

Next you'll want to add some records to the collection of "FavoriteThings.ds".  The records we're going to add need
to be expressed as JSON objects. You need to decide on a key (the thing you'll used to retrieve the record later)
of the document to store.  For this example I'm going to use the key, "beverage" and a document that looks like
`{"thing": "coffee"}`.  If you've set the DATASET environment variable you can run the following command --

```shell
    dataset FavoriteThings.ds create beverage '{"thing":"coffee"}'
```

If all goes well you'll get a response of "OK".  If you forgot to set the environment variable you can 
explicitly include the collection name
```shell
    dataset FavoriteThings.ds create beverage '{"thing":"coffee"}'
```


In Python

```python
    # continued from the previous python example
    err = dataset.create('FavoriteThings.ds', 'beverage', {"thing": "coffee"})
    if err != '':
        print(f"create error, {err}")
    else:
        print("OK")
```


Later if your have forgotten what your favorite beverage was you can read it back with

```shell
    dataset FavoriteThings.ds read beverage
```

Or in Python

```python
    (record, err) = dataset.reac('FavoriteThings', 'beverage')
    if err != '':
        print(f"read error, {err}")
    else:
        print(record)
```

To list all your favorite things' keys try

```shell
    dataset FavoriteThings.ds keys
```

In Python

```python
    keys = dataset.keys('FavoriteThings.ds')
```

## Adding an existing JSON document to a collection

One of my favorite things is music. I happen to have a JSON document that I started currating a list of 
Jazz related songs and musicians.  The document is called `jazz-notes.json`. I can add this to my collection too.

Here's the JSON document,

```json
    {
       "songs": ["Blue Rondo al la Turk", "Bernie's Tune", "Perdido"],
       "pianist": [ "Dave Brubeck" ],
       "trumpet": [ "Dirk Fischer", "Dizzy Gillespie" ]
    }
```

Add this to my collection of *FavoriteThings.ds* this way using the key "jazz-notes". 

```shell
    dataset FavoriteThings.ds create "jazz-notes" jazz-notes.json
```

or in Python

```python
    import json

    with open('jazz-notes.json', mode = 'r', encoding = 'utf-8') as f:
        src = f.read()
    jazz_notes = json.loads(src)
    err = dataset.create('FavoriteThings.ds', 'jazz-notes', jazz_notes)
    if err != '':
        print(f"create error, {err}")
```

Notice that the organization of the JSON documents do not impose a common structure (though that is
often useful). We can list the documents using our key command.

```shell
    dataset FavoriteThings.ds keys
```


Would return something like

```
    beverage
    jazz-notes
```
or in Python like this 

```python
    keys = dataset.keys('FavoriteThings.ds')
    print(keys)
```

The should list out "beverage" and "jazz-notes". 

I can create a JSON list of the objects stored using the "list" command.

```shell
    dataset FavoriteThings.ds list beverage jazz-notes
```

Would return something like

```json
    [
        {
            "_Key": "beverage",
            "thing": "coffee"
        },
        {
            "_Key": "jazz-notes",
            "pianist": [
                "Dave Brubeck"
            ],
            "songs": [
                "Blue Rondo al la Turk",
                "Bernie's Tune",
                "Perdido"
            ],
            "trumpet": [
                "Dirk Fischer",
                "Dizzy Gillespie"
            ]
        }
    ]
```

Similarly in Python 

```python
    (l, err) = dataset.list('FavoriteThings.ds')
    if err != '':
        print(f"list error, {err}")
    else:
        print(json.dumps(l, indent = 4)
```


## A workflow in Bash

This is an example of creating a dataset called *fiends.ds*, saving
a record called "littlefreda.json" and reading it back. We'll be adding some
records, print things out to the screen as well as checking if they keys can 
be found in a collection.

```shell
   dataset init friends.ds
   dataset friends.ds create littlefreda '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Now check to see if the key, littlefreda, is in the collection

```shell
   dataset friends.ds haskey littlefreda
```

You can also read your JSON formatted data from a file or standard input.
In this example we are creating a mojosam record and reading back the contents
of fiends.ds

```shell
   dataset -i mojosam.json friends.ds create mojosam
   for KY in $(dataset friends.ds keys); do
      echo "Path: $(dataset friends.ds path $KY) 
      echo "Doc: $(dataset friends.ds read $KY)
   done
```

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

```shell
   cat capt-jack.json | dataset friends.ds create capt-jack
   for KY in $(dataset friends.ds keys); do
      echo "Path: $(dataset friends.ds path $KY) 
      echo "Doc: $(dataset friends.ds read $KY)
   done
```

Adding high-capt-jack.txt as an attachment to "capt-jack"

```shell
   echo "Hi Capt. Jack, Hello World!" > high-capt-jack.txt
   dataset friends.ds attach capt-jack high-capt-jack.txt
```

List attachments for "capt-jack"

```shell
   dataset friends.ds attachments capt-jack
```

Get the attachments for "capt-jack" (this will untar in your current directory)

```shell
   dataset friends.ds attached capt-jack
```

Writing out the attachment named *high-capt-jack.txt* from "capt-jack"

```shell
    dataset friends.ds detach capt-jack high-capt-jack.txt
```

Remove all (prune) attachments from "capt-jack"

```shell
   dataset friends.ds prune capt-jack
```

### Continuing a Bash workflow

"import-csv" can take a CSV file and store each row as a JSON document in dataset. 
There does need to be a column of unique values to use as a key (each row becomes and
object in the collection).  In this example the first column is going to hold a id number.
The file contains a list of cast member, the title of the story and year of production.
We're going to create a new empty collection called _characters.ds_ and populated it from
a CSV file.

```shell
    dataset init characters.ds
    dataset friends.ds import characters.csv 1
```

You can check the number of records in _characters.ds_ with *count*.

```shell
    dataset characters.ds count
```

Here's an example of looping through all the keys and displaying titles and years.
We're using a command line tool called `jsoncols` from the [datatools](https://caltechlibrary.github.io/datatools)
project. It lets us read in a JSON object and display selected fields as a column

```shell
    dataset characters.ds keys | while read KEY; do
        echo -n "Title and year: "
        dataset -new-line=true characters.ds read "${KEY}" | jsoncols -i - .title .year
    done
```

Keys can be used to filter and sort keys. Here's is a simple case for match 
records where name is equal to "Mojo Sam".

```shell
   dataset characters.ds keys '(eq .name "Mojo Sam")'
```

You can take one list of keys and then do futher filtering using
the `-key-file` option with the *keys* verb.

```shell
   dataset characters.ds keys '(eq .name "Mojo Sam") > mojo.keys
   dataset -key-file mojo.keys characters.ds keys '(contains .title "Morroco")'
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.

```shell
   dataset characters.ds export titles.csv true '.id,.title,.year' 'id,title,publication year'
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```shell
   dataset characters.ds export titles2016.csv '(eq 2016 (year .year))' \
           '.id,.title,.year' 'id,title,publication year'
```

Let's return back to our friends collection.  You can augement JSON key/value 
pairs for a JSON document in your collection using the join operation. This works similar to the datatools cli called jsonjoin.

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

```shell
    dataset join update jane.doe profile.json
```

The result would look like

```json
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42, "bio": "renowned geophysist"}
```

If you wanted to overwrite the common fields you would use 'join overwrite'

```shell
    dataset join overwrite jane.doe profile.json
```

Which would result in a record like

```json
    {"name":"Doe, Jane", "email": "jane.doe@example.edu", "age": 42, "bio": "renowned geophysist"}
```


## A workflow in Python

Like in the Bash example we're creating a dataset collection called *fiends.ds*, saving
a record called "littlefreda.json" and reading it back. We going to
use more variables, add a logging class and reference a few extra Python modules to
make it more like scripts you'll write in practice.

```python
    # Standard Python packages
    import sys
    import os
    import json
    from datetime import tzinfo, timedelta, datetime

    # Caltech Library packages
    import dataset

    class Logger:
        def __init__(self, pid, time_format = '%Y/%m/%d %H:%M:%S'):
            self.pid = pid
            self.time_format = time_format

        def sprint(self, msg):
            dt = datetime.now().strftime(self.time_format)
            pid = self.pid
            return (f'{dt} (pid: {pid}) {msg}')

        def print(self, msg):
            dt = datetime.now().strftime(self.time_format)
            pid = self.pid
            print(f'{dt} (pid: {pid}) {msg}', flush = True)

        def fatal(self, msg):
            dt = datetime.now().strftime(self.time_format)
            pid = self.pid
            print(f'{dt} (pid: {pid}) {msg}', flush = True)
            sys.exit(1)


        def read_json(filename):
            with open(filename, mode = 'r', encoding = 'utf-8') as f:
                src = f.read()
            return json.loads(src)

            
    log = Logger(os.getpid())

    # We are saving our collection name in the variable c_name to save typing.
    c_name = 'friends.ds'
    err = dataset.init(c_name)
    if err != '':
        log.fatal(f"init error, {err}")
   
    key = 'littlefreda'
    err = dataset.create(c_name, key, '{"name":"Freda","email":"little.freda@inverness.example.org"}')
    if err != '':
        log.fatal(f"create error, {err}")
    log.print(f"Displaying path and JSON notation for keys in {c_name}") 
    keys = dataset.keys(c_name)
    for key in keys:
        p = dataset.path(key)
        log.print(f"Path for {key}: {p}")
        (record, err) = dataset.read(key)
        if err != '':
            log.fatal(f"read error, {err}")
        else:
            log.printf(f"JSON Object: {record}")
```

We can read the file "capt-jack.json" off disc an add it too.


```python
    capt_jack = read_json('capt-jack.json')
    err = dataset.create(c_name, 'capt-jack')
    if err != '':
        log.fatal(f"create error, {err}")

    for key in [ 'littlefreda', 'capt-jack' ]:
        log.print(f"Double check if {key} is in {c_name}"
        ok = dataset.haskey(key)
        if ok == True:
            log.print("OK")
        else:
            log.print("Missing {key} in {c_name}")
```


Let's read in capt-jack.json and mojosam.json and add them to our friends collection.


```python
    c_name = 'friends.ds'
    for filename in [ 'capt-jack.json', 'mojo-sam.json' ]:
        key = filename[0:-5]
        record = read_json(filename)
        err = dataset.create(c_name, key, record)
        if err != '':
            log.fatal(f"create error, {err}")
```

Adding high-capt-jack.txt as an attachment to "capt-jack"

```python
    with open('high-capt-jack.txt', mode = 'w', encoding = 'utf-8') as f:
        f.write("Hi Capt. Jack, Hello World!")
    err = dataset.attach(c_name, 'capt-jack',  'high-capt-jack.txt')
    if err != '':
        log.fatal(f"create error, {err}")
```

List attachments for "capt-jack"

```python
   l = dataset.attachments(c_nanme, 'capt-jack')
   log.print(l)
```

Get the attachments for "capt-jack" (this will untar in your current directory)

Writing out the attachment named *high-capt-jack.txt* from "capt-jack"

```python
    err = dataset.detach(c_name, 'capt-jack', 'high-capt-jack.txt')
    if err != '':
        log.fatal(f"detach error, {err}")
```

Remove all (prune) attachments from "capt-jack"

```python
   err = dataset.prune(c_name, 'capt-jack')
   if err != '':
       log.fatal(f"prune error, {err}")
```

Keys can be used to filter and sort keys.  Here's is a simple case for match records 
where name is equal to
"Mojo Sam".

```python
    c_name = 'characters.ds'
    keys = dataset.keys_filter(c_name, filter = '(eq .name "Mojo Sam")')
```

You can take one list of keys and then do futher filtering using
the `keys_filter()`.

```python
    keys = dataset.keys(c_name, filter = '(eq .name "Mojo Sam")')
    morroco_keys = dataset.keys_filter(c_name, keys, '(contains .title "Morroco")')
```

Import can take a CSV file and store each row as a JSON document in dataset. A column
needs to contain unique keys and that column is specified with the import command.

```python
    c_name = 'characters.ds'
    err = dataset.init(c_name)
    if err != '':
        log.fatal(f"init error, {err}")
    err = dataset.import_csv(c_name, 'characters.csv', 1)
    if err != '':
        log.fatal(f"import_csv error, {err}")
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.

```python
    err = dataset.export_csv(c_name, 'titles.csv', 'true', [ '.id', '.title', '.year'], ['id','title','publication year'])
    if err != '':
        log.fata(f"export_csv error, {err}")
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```python
    err = dataset.export_csv(c_name, 'titles2016.csv', '(eq 2016 (year .pubDate))',
           ['.id', '.title', '.year'], [ 'id', 'title', 'publication year']
    if err != '':
        log.fatal(f"export_csv error, {err}")
```


Returning to our _friends.ds_ collection. You can augement JSON key/value pairs for a 
JSON document in your collection using the join operation. This works similar to the 
datatools cli called jsonjoin.

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

```python
    profile = read_json('profile.json')
    err = dataset.join(c_name, 'jane.doe', 'update', profile)
    if err != '':
        log.fatal(f"join error, {err}")
```

The result would look like

```json
    {"name":"Doe, Jane", "email": "jd@example.org", "age": 42, "bio": "renowned geophysist"}
```

If you wanted to overwrite the common fields you would use 'join overwrite'

```python
    profile = read_json('profile.json')
    err = dataset.join(c_name, 'jane.doe', 'overwrite', profile)
    if err != '':
        log.fatal(f"join error, {err}")
```

Which would result in a record like

```json
    {"name":"Doe, Jane", "email": "jane.doe@example.edu", "age": 42, "bio": "renowned geophysist"}
```

