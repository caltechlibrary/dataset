
# Workflow

This is an example a workflow using the _dataset_ command to creat andmanage a collection called *fiends.ds*.
We start by saving a record called "littlefreda.json" and reading it back.

```shell
   dataset init fiends.ds
   export DATASET=fiends.ds
   dataset create littlefreda '{"name":"Freda","email":"little.freda@inverness.example.org"}'
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Now check to see if the key, littlefreda, is in the collection

```shell
   dataset haskey littlefreda
```

You can also read your JSON formatted data from a file or standard input.
In this example we are creating a mojosam record and reading back the contents
of fiends.ds

```shell
   dataset -i mojosam.json create mojosam
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Or similarly using a Unix pipe to create a "capt-jack" JSON record.

```shell
   cat capt-jack.json | dataset create capt-jack
   for KY in $(dataset keys); do
      echo "Path: $(dataset path $KY) 
      echo "Doc: $(dataset read $KY)
   done
```

Adding high-capt-jack.txt as an attachment to "capt-jack"

```shell
   echo "Hi Capt. Jack, Hello World!" > high-capt-jack.txt
   dataset attach capt-jack high-capt-jack.txt
```

List attachments for "capt-jack"

```shell
   dataset attachments capt-jack
```

Get the attachments for "capt-jack" (this will untar in your current directory)

```shell
   dataset attached capt-jack
```

Writing out the attachment named *high-capt-jack.txt* from "capt-jack"

```shell
    dataset detach capt-jack high-capt-jack.txt
```

Remove all (prune) attachments from "capt-jack"

```shell
   dataset prune capt-jack
```

Keys can be used to filter and sort keys.
Here's is a simple case for match records where name is equal to
"Mojo Sam".

```shell
   dataset keys '(eq .name "Mojo Sam")'
```

You can take one list of keys and then do futher filtering using
the `-key-file` option with the *keys* verb.

```shell
   dataset keys '(eq .name "Mojo Sam") > mojo.keys
   dataset -key-file mojo.keys keys '(contains .title "Morroco")'
```

Import can take a CSV file and store each row as a JSON document in dataset. You must
indicate which column to use as the key.  We're using column zero in this example.

```shell
   dataset import-csv my-data.csv 0
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.

```shell
   dataset export-csv titles.csv true '.id,.title,.pubDate' 'id,title,publication date'
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```shell
   dataset export-csv titles2016.csv '(eq 2016 (year .pubDate))' \
           '.id,.title,.pubDate' 'id,title,publication date'
```

You can augment JSON key/value pairs for a JSON document in your collection
using the join operation. This works similar to the datatools cli called jsonjoin.

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

