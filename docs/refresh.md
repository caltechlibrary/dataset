refresh (depricated)
====================

This command will update an objects in a frame based on the current
state of the collection.  

NOTE: If any keys/objects have been deleted from the collection then
the object associated with those keys in the frame will also
be removed.

In the following example the frame name is \"f1\", the collection is
\"examples.ds\". The example is refreshing the object list.

```shell
    dataset refresh example.ds f1
```

