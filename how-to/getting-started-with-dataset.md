Getting started with dataset
============================

**dataset** is designed to easily manage collections of JSON documents. A JSON object is associated with a unique key you provide. If you are using the default storage engine the objects themselves are stored on disc in a folder inside the collection folder. If you are using a SQL storage engine they are stored in a column of a table of the collection in your SQL database.

The collection folder contains a JSON object document called **collection.json**. This file stores operational metadata about the collection. If the collection is using a pairtree then a **keymap.json** file will include the association of keys with paths to their objects. 

When a collection is initialized a minimal codemeta.json file will created describing the collection. This can be update to a full codemeta.json file, follow the guideline and practice described at the [codemeta](https://codemeta.github.io) website.

**dataset** comes in several flavors --- a command line program called *dataset*, a web service called **datasetd** and the Go language package used to build for programs.

This tutorial talks both the command line program and the Go package. The command line is great for simple setup, the Go package allows you to build on other programs that use dataset collections for content persistence.

Create a collection with init
-----------------------------

To create a collection you use the init verb. In the following examples you will see how to do this with both the command line tool **dataset** as well as the Python module of the same name.

Let\'s create a collection called *friends.ds*. At the command line type the following.

```bash
    dataset init friends.ds
```

Notice that after you typed this and press enter you see an \"OK\" response. If there had been an error then you would have seen  an error message instead.

Working in Go is similar. We use the `dataset.Init()` func to create our new collection. We can import the "dataset" package using the import line `"github.com/caltechlibrary/dataset"`.  Here's a general code sketch.

```golang
   import (
      // import the packages your program needs ...
      "fmt"
      "os"

      // import dataset
      "github.com/caltechlibrary/dataset"
   )
        
   func main() {
       // The dataset collection is held in 'c'
       // This create the collection "friends.ds"
       collectionName := "frieds.ds"
       // "c" is a handle to the collection
       c, err := dataset.init(collectionName)
       if err != nil {
           fmt.Fprintf(os.Stderr, "Something went wrong, %s\n", err)
           os.Exit(1)
       }
       defer c.Close() // Remember to close your collection
       fmt.Printf("Created %q, ready to use\n", collectionName)
   }
```

In this Go example if the error is nil a statement is written to standard out saying the collection was created, if not an error is shown.

### removing friends.ds {#removing-friends.ds}

There is no dataset verb to remove a collection. A collection is just a folder with some files in it. You can delete the  collection by throwing the folder in the trash (Mac OS X and Windows) or using a recursive remove in the Unix shell.

```shell
    rm -fR friends.ds
```

Or using `os.RemoveAll()` in Go programs.

```
    if _, err := os.Stat(collectionName); err == nil {
        os.RemoveAll(collectionName)
    }
```


create, read, update and delete
-------------------------------

As with many systems that store information dataset provides for basic operations of creating, updating and deleting. In the following section we will work with the *friends.ds* collection and **favorites.ds** collection we created previously.

I have some friends who are characters in [ZBS](https://zbs.org) radio plays. I am going to create and save some of their info in our collection called *friends.ds*. I am going to store their name and email address so I can contact them. Their names are Little Frieda, Mojo Sam and Jack Flanders.

```bash
    dataset create friends.ds frieda \
      '{"name":"Little Frieda","email":"frieda@inverness.example.org"}'
```

Notice the \"OK\". Just like *init* the *create* verb returns a status. \"OK\" means everything is good, otherwise an error is shown. 

Doing the same thing in Go would look like. Note we have to explicitly `Open()` the collection to get a collection object then call `Create()` on the opened collection. `defer` make it easy for us to remember to close the collection when we're done.

```golang
    import (
        "fmt"
        "os"

        "github.com/caltechlibrary/dataset"
    )

    func main() {
        c, err := dataset.Open("fiends.ds")
        if err != nil {
            fmt.Fprintf(os.Stderr, "something went wrong, %s", err)
            os.Exit(1)
        }
        defer c.Close() // Don't forget to close the collection
        id := "frieda"
        m := map[string]interface{}{
            "id": id,
            "name":"Little Frieda",
            "email":"frieda@inverness.example.org",
        }
        // Create adds a map[string]interface{} to the collection.
        if err := dataset.Create(id, m); err != nil {
            fmt.Fprintf(os.Stderr, "%s",err)
            os.Exit(1)
        }
        fmt.Printf("OK")
        os.Exit(0)
    }
```

Go supports easy translation of struct types into JSON encoded byte slices. Can then use that store the JSON representations using the `CreateObject()` to create a JSON object from any Go type.

```golang
   import (
      "encoding/json"
      "fmt"
      "os"

      "github.com/caltechlibrary/dataset"
   )

   type Record struct {
       ID string `json:"id"`
       Name string `json:"name,omitempty"`
       EMail string `json:"email,omitempty"`
   }

   func main() {
       obj := &Record{
           ID: "frieda",
           Name: "Little Fieda",
           EMail: "frieda@inverness.example.org",
       }
       if err := dataset.CreateObject("friends.ds", obj.ID, obj); err != nil {
           fmt.Fprintf(os.Stderr, "%s", err)
           os.Exit(1)
       }
       fmt.Printf("OK")
       os.Exit(0)
   }
```

On the command line create requires us to provide a collection name, a key (e.g.  \"frieda\") and JSON markup to store the JSON  object. We can provide that either through the command line or by reading in a file or standard input.

command line \--

```bash
    cat <<EOT >mojo.json
    {
        "id": "mojo",
        "name": "Mojo Sam, the Yudoo Man", 
        "email": "mojosam@cosmic-cafe.example.org"
    }
    EOT

    cat mojo.json | dataset create friends.ds "mojo"

    cat <<EOT >jack.json
    {
        "id": "jack",
        "name": "Jack Flanders", 
        "email": "capt-jack@cosmic-voyager.example.org"
     
    EOT

    dataset create -i jack.json friends.ds "jack"
```

in Go  we can loop through records easily and add them \--

```golang
    // Open the collection
    c, err := dataset.Open("friends.ds")
    if err != nil {
        ...
    }
    defer c.Close()// Don't forget to close the collection

    // Create some new records
    newRecords := []Record{
        Record{
            ID: "mojo",
            Name: "Mojo Sam",
            EMail: "mojosam@cosmic-cafe.example.rog",
        },
        Record{
            ID: "jack",
            Name: "Jack Flanders",
            Email: "capt-jack@cosmic-voyager.example.org",
        },
    }
    // Save the new records into the collection
    for _, record := range newRecords {
        if err := dataset.CreateObject(record.ID, record); err != nil {
            fmt.Fprintf(os.Stderr, 
               "something went wrong add %q, %s\n", record.ID, key)
        }
    }
```

### read

We have three records in our *friends.ds* collection --- \"frieda\", \"mojo\", and \"jack\". Let\'s see what they look like with the **read** verb.

command line \--

```bash
    dataset read friends.ds frieda
```

On the command line you can easily pipe the results to a file for latter modification. Let\'s do this for each of the records we have created so far.

```bash
    dataset read -p friends.ds frieda >frieda-profile.json
    dataset read -p friends.ds mojo >mojo-profile.json
    dataset read -p friends.ds jack >jack-profile.json
```

Working in Go is similar but rather than write out our JSON structures to a file we\'re going to keep them in memory as an array of record structs before converting to JSON and writing it out.

In Go \--

```golang
    // Open our collection
    c, err := dataset.Open("friends.ds")
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
        os.Exit(1)
    }
    defer c.Close() // remember to close the collection

    // build our list of keys
    keys := []string{ "frieda", "mojo", "jack" }
    records := []*Record{}
    // loop through the list and write the JSON source to file.
    for _, key := range keys {
       obj := &Record{}
       if err := c.ReadObject(key, &obj); err != nil {
           fmt.Fprintf(os.Stderr, "%s\n", err)
           os.Exit(1)
       }  
       records = append(records, obj)
    }
    src, _ := json.MarshalIndent(records)
    fmt.Println("%s\n", src)
    os.Exit(0)
```

### update

Next we can modify the profiles (the \*.json files for the command line version). We\'re going to add a key/value pair for \"catch_phrase\" associated with each JSON object in *friends.ds*. For Little Frieda edit freida-profile.json to look like \--

```json
    {
        "_Key": "frieda",
        "email": "frieda@inverness.example.org",
        "name": "Little Frieda",
        "catch_phrase": "Woweee Zoweee"
    }
```

For Mojo\'s mojo-profile.json \--

```json
    {
        "_Key": "mojo",
        "email": "mojosam@cosmic-cafe.example.org",
        "name": "Mojo Sam, the Yudoo Man",
        "catch_phrase": "Feet Don't Fail Me Now!"
    }
```

An Jack\'s jack-profile.json \--

```json
    {
        "_Key": "jack",
        "email": "capt-jack@cosmic-voyager.example.org",
        "name": "Jack Flanders",
        "catch_phrase": "What is coming at you is coming from you"
    }
```

On the command line we can read in the updated JSON objects and save the results in the collection with the *update* verb. Like with **init** and **create** the **update** verb will return an "OK" or error message. Let\'s update each of our JSON objects.

```bash
    dataset update friends.ds freida frieda-profile.json
    dataset update friends.ds mojo mojo-profile.json
    dataset update friends.ds jack jack-profile.json
```

**TIP**: By providing a filename ending in ".json" the dataset command knows to read the JSON object from disc. If the object had stated with a \"{\" and ended with a \"}\" it would assume you were using an explicit JSON expression.

In Go we can work with each of the record as `map[string]interface{}` variables. We save from our previous *Read* example. We add our "catch_phrase" attribute then *Update* each record.

```golang
    c, err := dataset.Open("friends.ds")
    if err != nil { 
        // ... handle errors
    }
    defer c.Close()

    // Read our three profiles
    friedaProfile := map[string]interface{}{}
    if err := c.Read("frieda", fredaProfile); err != nil {
        // ... handle error
    }
    mojoProfile := map[string]interface{}{}
    if err :=  c.Read("mojo", mojoProfile); err != nil  {
        // ... handle error
    }
    jackProfile := map[string]interface{}{}
    if err := c.Read("jack", jackProfile); err != nil {
        // ... handle error
    }
    
    // Add our catch phrases
    friedaProfile["catch_phrase"] = "Wowee Zowee"
    mojoProfile["catch_phrase"] = "Feet Don't Fail Me Now!"
    jackProfile["catch_phrase"] = "What is coming at you is coming from you"
    
    // Update our records
    if err := c.Update("frieda", friedaProfile); err != "" {
        // ... handle error
    }
    if err := c.Update("mojo", mojoProfile); err != "" {
        // ... handle error
    }
    if err := c.Update("jack", jackProfile); err != nil {
        // ... handle error
    }
```

A better approach where we would be to use a Go struct to hold the profile records. This would ensure that they mapping of attribute names are consistently handled.

```golang
    import (
        "github.com/caltechlibrary/dataset"
    )

    type Profile struct {
        Name string `json:"name"`
        EMail string `json:"email,omitempty"`
        CatchPhrase string `json:"catech_phrase,omitempty"`
    }

    func main() {
        // Load our minimal records, i.e. name and email
        records := map[string]*Profile{}{
            "frieda": &Profile{ 
                Key: "frieda", 
                EMail: "frieda@inverness.example.org",
                Name: "Little Frieda", 
                },
            "mojo": &Profile{
                Key: "mojo",
                EMail: "mojosam@cosmic-cafe.example.org",
                Name: "Mojo Sam, the Yudoo Man",
            },
            "jack": &Profile{
                Key: "jack",
                EMail: "capt-jack@cosmic-voyager.example.org",
                Name: "Jack Flanders",
            },
        }

        // Create the collection and add our records
        c, err := dataset.Init("friends.ds", "")
        if err != nil {
            // ... handle errror
        }
        for key, record := range records {
            if err := c.CreateObject(key, recorrd); err != nil {
                // ... handle error
            }
        }

        // Add our catch phrases
    
        records["frieda"].CatchPhrase = "Wowee Zowee"
        records["mojo"].CatchPhrase = "Feet Don't Fail Me Now!"
        records["jack"].CatchPhrase = 
             "What is coming at you is coming from you"
    
        // Update our records
        for key, record := range records {
            if err := c.UpdateObject(key, record); err != "" {
                // ... handle error
            }
        }
    }
```


### delete

Eventually you might want to remove a JSON object from the collection. Let\'s remove Jack Flander\'s record for now.

command line \--

```bash
    dataset delete friends.ds jack
```

Notice the "OK" in this case it means we\'ve successfully delete the JSON object from the collection.

An perhaps as you\'ve already guessed working in Go looks like \--

```golang
   c, err := dataset.Open("friends.ds")
   if err != nil {
       // ... handle error
   }
   defer c.Close()

   if err := c.Delete("jack"); err != nil {
       fmt.Fprintf(os.Stderr, "%s\n", err)
       os.Exit(1)
   }
   fmt.Println("OK")
   os.Exit(0)
```

keys and count
--------------

Eventually you have lots of objects in your collection. You are not going to be able to remember all the keys. dataset provides a **keys** function for getting a list of keys as well as a **count** to give you a total number of keys.

Now that we\'ve deleted a few things let\'s see how many keys are in **friends.ds**. We can do that with the **count** verb.

Command line \--

```bash
    dataset count friends.ds
```

In Go \--

```golang
   c, err := dataset.Open("friends.ds")
   if err != nil {
       // ... handle error
   }
   defer c.Close()

   cnt = c.Length() // NOTE: this is an int64 value
   fmt.Printf("Total Records Now: %d\n", cnt)
```

Likewise we can get a list of the keys with the *keys* verb.

```bash
    dataset keys friends.ds
```

If you are following along in Go then you can just save the keys to a variable called keys.

```golang
   c, err := dataset.Open("friends.ds")
   if err != nil {
       // ... handle error
   }
   defer c.Close()

   keys, err = c.Keys()
   if err != nil {
       // ... handle error
   }
   fmt.Printf("%s\n", strings.Join(keys, "\n"))
```
