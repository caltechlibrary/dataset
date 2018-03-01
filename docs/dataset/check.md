
# check

## syntax

```shell
    dataset check
```

## Description

Check reviews a collection and reports if any problems are identified based on the 
`collection.json` file found in the folder holding the collection's buckets. Check
only works on local disc based collections. If you are storing your collection in
the cloud (e.g. S3 or Google Cloud Storage) then download a copy before running
check.

## Usage

If multiple instances of dataset write (e.g. create or update) to a collection then
it is possible that the JSON file `collection.json` will become inaccurate.

```shell
    dataset -c MyBrokenCollection.ds check
```


