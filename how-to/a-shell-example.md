
A Shell Example using dataset and datasetd
==========================================

dataset
--------

Below is a simple example of shell based interaction with dataset a collection using the command line dataset tool.

~~~shell
    # Create a collection "friends.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    dataset init friends.ds
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    dataset create friends.ds lady.jowls \
        '{"name":"lady.jowls","email":"lady.jowls@inverness.example.org"}'
    # If successful then you should see an OK otherwise an error message

    # Read a JSON document
    dataset read friends.ds lady.jowls
    
    # Update a JSON document
    dataset update friends.ds lady.jowls \
        '{"name":"lady.jowls","email":"lady.jowls@zbs.example.org", "current_residence": "Inverness"}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys friends.ds

    # Delete a JSON document
    dataset delete friends.ds lady.jowls

    # To remove the collection just use the Unix shell command
    rm -fR friends.ds
~~~

dataset
-------

We need to have two shell sessions running for this example.

Before we begin let's create an populate our "friends.ds" collection.

1. Create our "friends.ds" collection
2. Load the "friends.ds" from [zbs_cast_list.jsonl](zbs_cast_list.jsonl)

~~~shell
dataset3 init friends.ds
dataset3 load friends.ds <zbs_cast_list.jsonl
~~~

For this example we're going to use the "friends.ds" collection created in the previous example.  We need to create a "friends_api.yaml" file in the same directory where you have "friends.ds".  We will use it to start __datasetd__. That file should contain

~~~yaml
host: "localhost:8485",
collections:
  - dataset: "friends.ds"
    query:
      cast_list: |
        select src
        from friends
        order by src->>'family_name'
    keys: true
    create: true
    read: true
    update: true
    delete: true
~~~

We start up __dataset3d__ with the following command.

~~~shell
    datasetd start friends_api.yaml
~~~

In this first session you will see log output to the console. We can use that to see how the service handles the requests.

In a second shell session we're going to use the [curl](https://curl.se/) command to interact with our collections.

~~~shell
    # Create a JSON document 
    curl -X POST -H 'application/json' \
    'http://localhost:8485/friends/object/lord.jowls' \
    -d '{"name":"Lord Jowls","email":"lord.jowls@inverness.example.org"}'

    # Read a JSON document
    curl 'http://localhost:8485/friends/object/lord.jowls'
    
    # Update a JSON document
    curl -X PUT -H 'application/json' \
    'http://localhost:8485/friends/object/lord.jowls' \
    -d '{"name":"Lord Jowls","email":"lord.jowls@zbs.example.org", "current_residency": "astroplanes"}'

    # List the keys in the collection
    curl 'http://localhost:8485/friends/keys'

    # Delete a JSON document
    curl -X DELETE 'http://localhost:8485/friends/object/lord.jowls'
~~~
