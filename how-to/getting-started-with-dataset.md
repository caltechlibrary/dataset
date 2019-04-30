
# Getting started with dataset

_dataset_ is designed to easily manage collections of JSON objects. 
Objects are associated with a unique key you provide. The objects 
themselves are stored on disc in a folder inside the collection folder. 
The collection folder contains a JSON object document called 
*collection.json*. This file stores metadata about the collection 
including the association of keys with paths to their objects.
_dataset_ comes in several flavors — a command line program called
_dataset_, a Go langauge package also called dataset, a shared library
called libdataset and a Python 3.7 package called 
[py_dataset](https://github.com/caltechlibrary/py_dataset). This
tutorial talks both the command line program and the Python package.
The command line is great for simple setup but Python is often more
convienent for more complex operations.


## Create a collection with init

To create a collection you use the init verb. In the following examples 
you will see how to do this with both the command line tool 
_dataset_ as well as the Python module of the same name.

Let's create a collection called _friends.ds_. At the command line type 
the following.


```bash
    dataset init friends.ds
```

Notice that after you typed this and press enter you see an 
"OK" response. If there had been an error then you would have 
seen an error message instead. 

Working in Python is similar to the command line. We import the 
modules needed then use them. For these exercises we'll be 
importing the following modules _sys_, _os_, _json_ and 
of course _dataset_ via `from py_dataset import dataset`.


```python
    import sys
    import os
    import json
    from py_dataset import dataset
    
    # stop is a convenience function
    def stop(msg):
        print(msg)
        sys.exit(1)
        
    err = dataset.init("friends.ds")
    if err != "":
        stop(err)
```

In Python the error message is an empty string if everything is ok, 
otherwise we call stop which prints the message and exits. You will see 
this pattern followed in a number of upcoming Python examples.


### removing friends.ds

There is no dataset  verb to remove a collection. A collection is just a 
folder with some files in it. You can delete the collection by throwing 
the folder in the trash (Mac OS X and Windows) or using a recursive 
remove in the Unix shell.

## create, read, update and delete

As with many systems that store information dataset provides for basic 
operations of creating, updating and deleting. In the following section we 
will work with the _friends.ds_ collection and _favorites.ds_ 
collection we created previously.

I have some friends who are characters in [ZBS](https://zbs.org) radio 
plays. I am going to create and save some of their info in our collection 
called _friends.ds_. I am going to store their name and email address 
so I can contact them. Their names are Little Frieda, Mojo Sam and Jack 
Flanders.


```bash
    dataset create friends.ds frieda \
      '{"name":"Little Frieda","email":"frieda@inverness.example.org"}'
```

Notice the "OK". Just like _init_ the _create_ verb returns a status. 
"OK" means everything is good, otherwise an error is shown. Doing the 
same thing in Python would look like.


```python
    err = dataset.create("friends.ds", "frieda", 
          {"name":"Little Frieda","email":"frieda@inverness.example.org"})
    if err != "":
        stop(msg)
```

With create we need to provide a collection name, a key (e.g. "frieda") 
and Python dict (which becomes our JSON object). Now let's add records 
for Mojo Sam and Jack Flanders.

command line -- 


```bash
    dataset create friends.ds "mojo" \
        '{"name": "Mojo Sam, the Yudoo Man", "email": "mojosam@cosmic-cafe.example.org"}'
    dataset create friends.ds "jack" \
        '{"name": "Jack Flanders", "email": "capt-jack@cosmic-voyager.example.org"}'
```

in python -- 


```python
    err = dataset.create("friends.ds", "mojo", 
          {"name": "Mojo Sam, the Yudoo Man", 
          "email": "mojosam@cosmic-cafe.example.org"})
    if err != "": 
        stop(err)
    err = dataset.create("friends.ds", "jack", 
          {"name": "Jack Flanders", 
          "email": "capt-jack@cosmic-voyager.example.org"})
    if err != "": 
        stop(err)
```


### read

We have three records in our _friends.ds_ collection — "frieda", 
"mojo", and "jack".  Let's see what they look like with the _read_ 
verb.

command line -- 


```bash
    dataset read friends.ds frieda
```

This command emitts a JSON object. The JSON  is somewhat hard to read. 
To get a pretty version of the JSON object used the "-p"  option.


```bash
    dataset read -p friends.ds frieda
```

On the command line you can easily pipe the results to a file for latter 
modification. Let's do this for each of the records we have created so far.


```bash
    dataset read -p friends.ds frieda >frieda-profile.json
    dataset read -p friends.ds mojo >mojo-profile.json
    dataset read -p friends.ds jack >jack-profile.json
```

Working in python is similar but rather than write out our JSON structures 
to a file we're going to keep them in memory as Python dict.

In Python -- 


```python
    (frieda_profile, err) = dataset.read("friends.ds", "frieda")
    if err != "":
        stop(err)
    (mojo_profile, err) = dataset.read("friends.ds", "mojo")
    if err != "":
        stop(err)
    (jack_profile, err) = dataset.read("friends.ds", "jack")
    if err != "":
        stop(err)
```

In Python, just like with _init_ and _create_ the _read_ verb 
returns a tuple of the value and err. Notice a pattern?


### update

Next we can modify the profiles (the \*.json files for the command line 
version). We're going to add a key/value pair for "catch_phrase" 
associated with each JSON object in _friends.ds_.  For Little Frieda 
edit freida-profile.json to look like -- 


```json
    {
        "_Key": "frieda",
        "email": "frieda@inverness.example.org",
        "name": "Little Frieda",
        "catch_phrase": "Woweee Zoweee"
    }
```

For Mojo's mojo-profile.json -- 


```json
    {
        "_Key": "mojo",
        "email": "mojosam@cosmic-cafe.example.org",
        "name": "Mojo Sam, the Yudoo Man",
        "catch_phrase": "Feet Don't Fail Me Now!"
    }
```

An Jack's jack-profile.json -- 


```json
    {
        "_Key": "jack",
        "email": "capt-jack@cosmic-voyager.example.org",
        "name": "Jack Flanders",
        "catch_phrase": "What is coming at you is coming from you"
    }

```

On the command line we can read in the updated JSON objects and save the 
results in the collection with the _update_ verb. Like with _init_ 
and _create_  the _update_ verb will return an “OK” or error message. 
Let's update each of our JSON objects.


```bash
    dataset update friends.ds freida frieda-profile.json
    dataset update friends.ds mojo mojo-profile.json
    dataset update friends.ds jack jack-profile.json
```

**TIP**: By providing a filename ending in “.json” the dataset command 
knows to read the JSON object from disc. If the object had stated with 
a "{" and ended with a "}" it would assume you were using an explicit 
JSON expression.

In Python we can work with each of the dictionaries variables we save 
from our previous _read_ example.  We add our “catch_phrase” 
attribute then _update_ each record.


```python
    frieda_profile["catch_phrase"] = "Wowee Zowee"
    mojo_profile["catch_phrase"] = "Feet Don't Fail Me Now!"
    jack_profile["catch_phrase"] = "What is coming at you is coming from you"
    
    err = dataset.update("friends.ds", "frieda", frieda_profile)
    if err != "":
        stop(err)
    err = dataset.update("friends.ds", "mojo", mojo_profile)
    if err != "":
        stop(err)
    err = dataset.update("friends.ds", "jack", jack_profile)
    if err != "":
        stop(err)
```

As an exercise how would you read back the updated version on the 
command line or in Python?


### delete

Eventually you might want to remove a JSON object from the collection. 
Let's remove Jack Flander's record for now.

command line -- 


```bash
    dataset delete friends.ds jack
```

Notice the “OK” in this case it means we've successfully delete the 
JSON object from the collection.

An perhaps as you've already guessed working in Python looks like -- 


```python
    err = dataset.delete("friends.ds", "jack")
    if err != "":
       stop(err)
```


## keys and count

Eventually you have lots of objects in your collection. You are not going 
to be able to remember all the keys. dataset provides a _keys_ function 
for getting a list of keys as well as a _count_ to give you a 
total number of keys.

Now that we've deleted a few things let's see how many keys are in 
_friends.ds_. We can do that with the _count_ verb.

Command line -- 

 
```bash
    dataset count friends.ds
```

In Python -- 


```python
    cnt = dataset.count("friends.ds")
    print(f"Total Records Now: {cnt}")
```

Likewise we can get a list of the keys with the _keys_ verb. 


```bash
    dataset keys friends.ds
```

If you are following along in Python then you can just save the keys to 
a variable called keys.


```python
    keys = dataset.keys("friends.ds")
    print("\n".join(keys))
```


## grids and frames

JSON objects are tree like. This structure can be inconvienent
for some types of analysis like tabulation, comparing values or
generating summarizing reports. A spreadsheet, table or 2D grid 
like structure is often a more familair format for these types
of tasks. _grid_ is dataset's verb for taking a list 
of keys, a list of dot paths to the JSON objects attributes and 
returning a 2D grid of the results. The 2D grid is easy to
iterate over.  A _grid_ doesn't enforce any specifics on the 
columns and rows. It only contains the values you specified in
the list of keys and dot paths.


### grid

Let's create a _grid_ from our *friends.ds* collection.


```bash
    dataset keys friends.ds keys >fiends.keys
    dataset grid -i=friends.keys friends.ds .name .email .catch_phrase
```

As with _read_ the _grid_ verb can take the “-p” option to make the 
JSON grid a little easier to read.


```bash
    dataset grid -p -i=friends.keys friends.ds .name .email .catch_phrase
```

Notice we make a list of keys first and save those to a file. Then we use 
that list of keys and create our grid.  The grid output is in JSON 
notation. In Python making a grid follows a similar patter, generate a 
list of keys, use those keys and a list of dot paths to define the grid.


```python
    keys = dataset.keys("friends.ds")
    (g, err) = dataset.grid("friends.ds", keys, 
               [".name", ".email", "catch_phrase"])
    if err != "":
        stop(err)
    print(json.dumps(g, indent = 4))
```

In python _grid_ like _create_ and _update_ returns a tuple that 
has your result and an error status. Finally we print our result using 
the JSON module's _dumps_.


### frame

dataset also comes with a _frame_ verb.  A _frame_ is like a grid plus 
additional matadata. It enforces a structure such on its grid. Column 1 
of the _frame_'s internal grid element always has the keys associated 
with the collection. A _frame_ will also derive heading labels from the 
dot paths used to define the frame and will include metadata about the 
collection, keys used to define the frame and default types of data in 
the columns. The extra information in a _frame_ stays with the 
collection. Frames are persistent and can be easily recalculated based on 
collection updates. Finally frames are used by more complex verbs such as 
_export_ and _indexer_ we'll be covering later. 

To define a frame we only need one additional piece of information besides 
what we used for a grid. We need a name for the frame. 

Working from our previous _grid_ example, let's call this frame 
"name-and-email".


```bash
    dataset frame -i=friends.keys friends.ds \
        "name-and-email" .name .email .catch_phrase
```

In python it would look like


```python
    keys = dataset.keys("friends.ds")
    err = dataset.frame("friends.ds", "name-and-email", 
          keys, [ ".name", ".email", ".catch_phrase"])
    if err != "":
        stop(err)
```

To see the contents of a frame we only need to supply the collection 
and frame names.


```bash
    dataset frame friends.ds "name-and-email"
```

In Python it'd look like


```python
    (f, err) = dataset.frame("friends.ds", "name-and-email")
    if err != "":
        stop(err)
    print(json.dumps(f, indent = 4))
```

Looking at the resulting JSON object you see many other attribute 
beyond the grid of values. These are created to simplify some of dataset 
more complex interactions.


Let's add back the Jack record we deleted a few sections ago and 
“reframe” our “name-and-email” frame.


```bash
    # Adding back Jack
    dataset create friends.ds jack jack-profile.json
    # Save all the keys in the collection
    dataset keys friends.ds >friends.keys
    # Now reframe "name-and-email" with the updated friends.keys
    dataset reframe -i=friends.keys friends.ds "name-and-email" 
    # Now let's take a look at the frame
    dataset frame -p friends.ds "name-and-email"
```

Like with _grid_ and _read_ before it the “-p” option will cause the 
JSON representation of the frame to be pretty printed.

Let's try the same thing in Python


```python
    err = dataset.create("friends.ds", "jack", jack_profile)
    if err != "":
        stop(err)
    keys = dataset.keys("friends.ds")
    err = dataset.reframe("friends.ds", "name-and-email", keys)
    if err != "":
        stop(err)
    (f, err) = dataset.frame("friends.ds", "name-and-email")
    if err != "":
        stop(err)
    print(json.dumps(f, indent = 4))
```

We can list the frames in the collection using the _frames_ verb.


```bash
    dataset frames friends.ds
```

In Python


```python
    frame_names = dataset.frames("friends.ds")
    print("\n".join(frame_names))
```

In our frame we have previously defined three columns, looking at the 
JSON representation of the frame we also see a "labels" attribute. 
Labels are used when exporting and synchronizing content between a 
CSV file, Google Sheet and a collection (labels become column names).

You can set the labels by providing replacement values for ALL
the labels in the column order defined in the frame. In our previous
example we provided the order of the columns for the frame
"name-and-email" as .name, .email, .catch_phrase. Also recall the
frames AWAYS has a first column named `._Key`. If we want
to have the labels "ID", "Name", "EMail", and "Catch Phrase" instead
we can set them with the `frame-labels` verb.

```bash
    dataset frame-labels friends.ds "name-and-email" \
        "ID" "Name" "EMail" "Catch Phrase"
```

In Python it look like

```python
    err = dataset.frame_labels("friends.ds", "name-and-email", 
          ["Name", "EMail", "Catch Phrase"])
    if err != "":
        stop(err)
```


Finally the last thing we need to be able to do is delete a frame. Delete 
frames work very similar to deleting a JSON record.


```bash
    dataset delete-frame friends.ds "name-and-email"
```

Or in Python


```python
    err = dataset.delete_frame("friends.ds", "name-and-email")
    if err != "":
          stop(err)
```

**TIP**: Frames like collections have a number of operations. Here's 
the list

1. _frame_ will set you define a frame
2. _frame_ will also let you read back a frame
3. _frames_ will list the frames defined in the collection
4. _frame-labels_ will let you replace the labels values for all 
   columns in a frame
6. _delete-frame_ will remove the frame from the collection



Continue exploring dataset with

- [Working with CSV](working-with-csv.html)
- [Working with GSheets](working-with-gsheets.html)
- [Working with Cloud Storage](working-with-cloud-storage.html)
- (__experimental__) [Indexing and Search](indexing-and-search.html "experimental Bleve search support") 


