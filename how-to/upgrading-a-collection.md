Upgrading your dataset collection(s)
====================================

The __dataset__ Go package is still rapidly evolving though it is
now commonly used at Caltech Library. As a result we've developed
an easy method for migration collections from an old version to
a new version. In version 1 dataset the "usual" to upgrading your
was to use use the "check" and "repair" features of the __dataset__
command line tool.

```shell
    dataset check mycollection.ds
    # you'll get a verbose report to the console
    dataset repair mycollection.ds
    # dataset will not attempt to "repair" including upgrade, 
    # your collection
```

In the version 2 of dataset a "migrate" verb has been created.
There are allot of structural changes inside a collection between
version one and two not to mention different storage engine options.
The new recipe for upgrading is create an empty target collection
with the storage engine you prefer (e.g. Pairtree or SQL).
Run "check" as before to make sure we can read the old collection
and then use "migrate" to populate the new collection with the old
collection's contrent.

```shell
    # Create a new empty (v2) new collection,
    # in this case we stuck with pairtree storage.
    dataset init new_collection.ds 

    # Double check the old one to make sure we can read it
    dataset check old_collection.ds

    # Migrate the content verion old to new.
    dataset migrate old_collection.ds new_collection.ds
```

Once migrated you can run your usual check and repair on
any pairtree storage collection. For SQL storage collections
you will need to use your SQL database tools to check and
fix broken tables.



