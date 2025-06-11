reframe (depreciated)
=====================

This command replaces the current keys/objects in a frame based
on the new keys provided. 

In the following example the frame name is \"f1\", the collection is
\"examples.ds\". The first example is reframing an existing frame using
existing keys coming from standard input, the second example performs
the same thing but is taking a filename to retrieve the list of keys.

```shell
    cat f1-updated.keys | dataset reframe example.ds f1
    dataset reframe example.ds f1 f1-updated.keys
```

