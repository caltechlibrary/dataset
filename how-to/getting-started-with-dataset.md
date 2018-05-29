
# Getting started with dataset

**dataset** is designed to easily manage collections of JSON objects. Objects are associated with a unique key you provide. The objects themselves are stores on disc in a bucket folder inside the collection folder. The collection folder contains a JSON object document called *collection.json*. This file stores metadata about the collection including the association of keys with paths to their objects.  **dataset** comes in several flavors — a command line program called **dataset**, a Go langauge package also called dataset, a shared library called libdataset and a Python 3.6 package called **dataset**. This tutorial talks about the command line program and the Python package.


## Create a collection with init

To create a collection you use the init verb. In the following examples you will see how to do this with both the command line tool called dataset as well as the Python module of the same name.

Let's create a collection called **friends.ds**. At the command line type the following.


```bash
    dataset init friends.ds
```

Notice that when you typed this in you see an "OK" response. If there had been an error then you would have seen an error message instead. 

Working in Python is similar to the command line we do need to import some modules and for these exercises we'll be importing the following modules **sys**, **os**, **json** and of course **dataset**.


```python
    import sys
    import os
    import json
    import dataset
    
    # stop is a convenience function
    def stop(msg):
        print(msg)
        sys.exit(1)
        
    err = dataset.init("friends.ds")
    if err != "":
        stop(err)
```

In Python the error message is an empty string if everything is ok, otherwise we call stop which prints the message and exits. You will see this pattern followed in a number of upcoming Python examples.


### removing friends.ds

There is no dataset  verb to remove a collection. A collection is just a folder with some files  in it. You can delete the collection by throwing the folder in the trash (Mac OS X and Windows) or using a recursive remove in the Unix shell.

## create, read, update and delete

As with many systems that store information dataset provides for basic operations of creating, updating and deleting. In the following section we will work with the **friends.ds** collection and **favorites.ds** collection we created previously.

I have some friends who are characters in [ZBS](https://zbs.org) radio plays. I am going to create
save some of their info in our collection called **friends.ds**. I am going to store their name and email address so I can contact them. Their names are Little Frieda, Mojo Sam and Jack Flanders.


```bash
    dataset friends.ds create frieda '{"name":"Little Frieda","email":"frieda@inverness.example.org"}'
```

Notice the "OK". Just like **init** the **create** verb returns a status. "OK" means everything is good, otherwise an error is shown. Doing the same thing in Python would look like.


```python
    err = dataset.create("friends.ds", "frieda", {"name":"Little Frieda","email":"frieda@inverness.example.org"})
    if err != "":
        stop(msg)
```

With create we need to provide a collection name, a key (e.g. "frieda") and Python
dict (which becomes our JSON object). Now let's add records for Mojo Sam and Jack Flanders.

command line -- 


```bash
    dataset friends.ds create "mojo" '{"name": "Mojo Sam, the Yudoo Man", "email": "mojosam@cosmic-cafe.example.org"}'
    dataset friends.ds create "jack" '{"name": "Jack Flanders", "email": "capt-jack@cosmic-voyager.example.org"}'
```

in python -- 


```python
    err = dataset.create("friends.ds", "mojo", {"name": "Mojo Sam, the Yudoo Man", "email": "mojosam@cosmic-cafe.example.org"})
    if err != "": 
        stop(err)
    err = dataset.create("friends.ds", "jack", {"name": "Jack Flanders", "email": "capt-jack@cosmic-voyager.example.org"})
    if err != "": 
        stop(err)
```


### read

We have three records in our **friends.ds** collection — "frieda", "mojo", and "jack".  Let's see what they look like with the **read** verb.

command line -- 


```bash
    dataset friends.ds read frieda
```

This command emitts a JSON object. The JSON  is somewhat hard to read. To get a pretty version of the JSON object used the "-p"  option.


```bash
    dataset -p friends.ds read frieda
```

On the command line you can easily pipe the results to a file for latter modification. Let's do this for each of the records we have created so far.


```bash
    dataset -p friends.ds read frieda > frieda-profile.json
    dataset -p friends.ds read mojo > mojo-profile.json
    dataset -p friends.ds read jack > jack-profile.json
```

Working in python is similar but rather than write out our JSON structures to a file we're going to 
keep them in memory as Python dict.

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

In Python, just like with **init** and **create** the **read** verb returns a tuple of the value and err. Notice a pattern?


### update

Next we can modify the profiles (the *.json files for the command line version). We're going to add a key/value pair for "catch_phrase" associated with each JSON object in **friends.ds**.  For 
Little Frieda edit freida-profile.json to look like -- 


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

On the command line we can read in the updated JSON objects and save the results in the collection with the **update** verb. Like with **init** and **create**  the **update** verb will return an “OK” or error message. Let's update each of our JSON objects.


```bash
    dataset friends.ds update freida frieda-profile.json
    dataset friends.ds update mojo mojo-profile.json
    dataset friends.ds update jack jack-profile.json
```

**TIP**: By providing a filename ending in “.json” the dataset command knows to read the JSON object from disc. If the object had stated with a "{" and ended with a "}" it would assume you were using an explicit JSON expression.

In Python we can work with each of the dictionaries variables we save from our previous **read** example.  We add our “catch_phrase” attribute then **update** each record.


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

As an exercise how would you read back the updated version on the command line or in Python?


### delete

Eventually you might want to remove a JSON object from the collection. Let's remove Jack Flander's record for now.

command line -- 


```bash
    dataset friends.ds delete jack
```

Notice the “OK” in this case it means we've successfully delete the JSON object from the collection.

An perhaps as you've already guessed working in Python looks like -- 


```python
    err = dataset.delete("friends.ds", "jack")
    if err != "":
       stop(err)
```


## keys and count

Eventually you have lots of objects in your collection. You are not going to be able to remember all the keys. dataset provides a **keys** function for getting a list of keys as well as a **count** to give you a 
total number of keys.

Now that we've deleted a few things let's see how many keys are in **friends.ds**. We can do that with the **count** verb.

Command line -- 

 
```bash
    dataset friends.ds count
```

In Python -- 


```python
    cnt = dataset.count("friends.ds")
    print(f"Total Records Now: {cnt}")
```

Likewise we can get a list of the keys with the **keys** verb. 


```bash
    dataset friends.ds keys
```

If you are following along in Python then you can just save the keys to a variable called keys.


```python
    keys = dataset.keys("friends.ds")
    print("\n".join(keys))
```


## grids and frames

One of the challenges in working on JSON objects is their tree like structure. When tabulating or
comparing values it is often easier to work in a spreadsheet like grid.  **grid** is dataset's verb for taking a list of keys, a list of dot paths into the JSON objects and returning a 2D grid of the results. This is handy when generating reports. A **grid** unlike **frame** which we will see shortly doesn't enforce any specifics on the columns and rows. It only contains the values you specify.


### grid

Let's create a **grid** from our *friends.ds* collection.


```bash
    dataset friends.ds keys > fiends.keys
    dataset friends.ds grid friends.keys .name .email .catch_phrase
```

As with **read** the **grid** verb can take the “-p” option to make the JSON grid a little easier to read.


```bash
    dataset -p friends.ds grid friends.keys .name .email .catch_phrase
```

Notice we make a list of keys first and save those to a file. Then we use that list of keys and create our grid.  The grid output is in JSON notation. In Python making a grid follows a similar patter, generate a list of keys, use those keys and a list of dot paths to define the grid.


```python
    keys = dataset.keys("friends.ds")
    (g, err) = dataset.grid("friends.ds", keys, [".name", ".email", "catch_phrase"])
    if err != "":
        stop(err)
    print(json.dumps(g, indent = 4))
```

In python **grid** like **create** and **update** returns a tuple that has your result and an error status. Finally we print our result using the JSON module's **dumps**.


### frame

dataset also comes with a **frame** verb.  A **frame** is like a grid plus additional matadata. It enforces a structure such on its grid. Column 1 of the **frame**'s internal grid element always has the keys associated with the collection. A **frame** will also derive heading labels from the dot paths used to define the frame and will include metadata about the collection, keys used to define the frame and default types of data in the columns. The extra information in a **frame** stays with the collection. Frames are persistent and can be easily recalculated based on collection updates. Finally frames as used by more complex verbs such as **export-csv**, **export-gsheet**, and **indexer** we'll be covering later. 

To define a frame we only need one additional piece of information besides what we used for a grid. We need a name for the frame. 

Working from our previous **grid** example, let's call this frame "name-and-email".


```bash
    dataset friends.ds frame "name-and-email" fiends.keys .name .email .catch_phrase
```

In python it would look like


```python
    keys = dataset.keys("friends.ds")
    err = dataset.frame("friends.ds", "name-and-email",  keys, [ ".name", ".email", ".catch_phrase"])
    if err != "":
        stop(err)
```

To see the contents of a frame we only need to support the collection name and frame name.


```bash
    dataset friends.ds frame "name-and-email"
```

In Python it'd look like


```python
    (f, err) = dataset.frame("friends.ds", "name-and-email")
    if err != "":
        stop(err)
    print(json.dumps(f, indent = 4))
```

Looking at the resulting JSON object you see many other attribute beyond the grid of values. These are what simplify some of dataset more complex interactions.



Let's add back the Jack record we deleted a few sections again and “reframe” our “name-and-email” frame.


```bash
    # Adding back Jack
    dataset friends.ds create jack jack-profile.json
    # Save all the keys in the collection
    dataset friends.ds keys > friends.keys
    # Now reframe "name-and-email" with the updated friends.keys
    dataset friends.ds reframe "name-and-email" friends.keys
    # Now let's take a look at the frame
    dataset -p friends.ds frame  "name-and-email"
```

Like with **grid** and **read** before it the “-p” option will cause the JSON representation of the frame to be pretty printed.

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

We can list the frames in the collection using the **frames** verb.


```bash
    dataset friends.ds frames
```

In Python


```python
    frame_names = dataset.frames("friends.ds")
    print("\n".join(frame_names))
```

In our frame we have previously defined three columns, looking at the JSON representation of the frame we also see three labels and three “types”.  Labels are used when exporting frames to spreadsheets. They are also used as the field names when we get to indexes and search. The types are used when defining indexes for searching. The values in types should correspond to either JSON types or the types supported by the search system (e.g. keyword, datetime, geolocation).have three fields in our frame. We will work with both labels and types when we are using the other commands to export and indexes.

Finally the last thing we need to be able to do is delete a frame. Delete frames work very similar to deleting a JSON record.


```bash
    dataset friends.ds delete-frame "name-and-email"
```

Or in Python


```python
    err = dataset.delete_frame("friends.ds", "name-and-email")
    if err != "":
          stop(err)
```

**TIP**: Frames like collections have a number of operations. Here's the list

1. **frame** will set you define a frame
2. **frame** will also let you read back a frame
3. **frames** will list the frames defined in the collection
4. **frame-labels** will let you replace the labels values for all columns in a frame
5. **frame-types** will let you replace the type values for all columns in a frame
6. **delete-frame** will remove the frame from the collection



Continue exploring dataset with

- [Indexing and Search](indexing-and-search.html)
- [Working with CSV](working-with-csv.html)
- [Working with GSheets](working-with-gsheets.html)
- [Working with Cloud Storage](working-with-cloud-storage.html)


