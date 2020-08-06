A Shell Example
---------------

Below is a simple example of shell based interaction with dataset a
collection using the command line dataset tool.

``` {.shell}
    # Create a collection "friends.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    dataset init friends.ds
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    dataset create friends.ds frieda \
        '{"name":"frieda","email":"frieda@inverness.example.org"}'
    # If successful then you should see an OK otherwise an error message

    # Read a JSON document
    dataset read friends.ds frieda
    
    # Path to JSON document
    dataset path friends.ds frieda

    # Update a JSON document
    dataset update friends.ds frieda \
        '{"name":"frieda","email":"frieda@zbs.example.org", "count": 2}'
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    dataset keys friends.ds

    # Get keys filtered for the name "frieda"
    dataset keys friends.ds '(eq .name "frieda")'

    # Join frieda-profile.json with "frieda" adding unique key/value pairs
    dataset join friends.ds frieda frieda-profile.json

    # Join frieda-profile.json overwriting in commont key/values adding
    # unique key/value pairs from frieda-profile.json
    dataset join -overwrite  friends.ds frieda frieda-profile.json

    # Delete a JSON document
    dataset delete friends.ds frieda

    # Import data from a CSV file using column 1 as key
    dataset import friends.ds my-data.csv 1

    # To remove the collection just use the Unix shell command
    rm -fR friends.ds
```
