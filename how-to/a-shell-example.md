A Shell Example
===============

dataset
-------

Below is a simple example of shell based interaction with dataset a collection using the command line dataset tool.

```shell
    # Create a collection "friends.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    dataset init friends.ds
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    dataset create friends.ds frieda \
        '{"name":"frieda","email":"frieda@inverness.example.org"}'
    # If successful then you should see an OK otherwise an error message

    # Read a JSON document
    dataset read friends.ds frieda
    
    # Update a JSON document
    dataset update friends.ds frieda \
        '{"name":"frieda","email":"frieda@zbs.example.org", "count": 2}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys friends.ds

    # Delete a JSON document
    dataset delete friends.ds frieda

    # To remove the collection just use the Unix shell command
    rm -fR friends.ds
```

datasetd
--------

We need to have two shell sessions running for this example.

For this example we're going to use the "friends.ds" collection created in the previous example.  We need to create a "settings.json" file in the same directory where you have "friends.ds".  We will use it to start _datasetd_. That file should contain

```json
    {
        "host": "localhost:8485",
        "collections": {
            "friends": {
                "dataset": "friends.ds",
                "keys": true,
                "create": true,
                "read": true,
                "update": true,
                "delete": true
            }
        }
    }
```

We start up _datasetd_ with the following command.

```shell
    datasetd settings.json
```

In this first session you will see log output to the console. We can use that to see how the service handles the requests.

In a second shell session we're going to use the [curl](https://curl.se/) command to interact with our collections.

```shell
    # Create a JSON document 
    curl -X POST -H 'application/json' \
    'http://localhost:8485/friends/object/frieda' \
    -d '{"name":"frieda","email":"frieda@inverness.example.org"}'

    # Read a JSON document
    curl 'http://localhost:8485/friends/object/frieda'
    
    # Update a JSON document
    curl -X PUT -H 'application/json' \
    'http://localhost:8485/friends/object/frieda' \
    -d '{"name":"frieda","email":"frieda@zbs.example.org", "count": 2}'

    # List the keys in the collection
    curl 'http://localhost:8485/friends/keys'

    # Delete a JSON document
    curl -X DELETE 'http://localhost:8485/friends/object/frieda'
```

