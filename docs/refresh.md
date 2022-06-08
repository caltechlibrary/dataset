refresh
=======

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

In python

```python
    err = dataset.frame_refresh('example.ds', 'f1')
```

Releted topics: [frame](frame.html), [reframe](reframe.html),
[frame-objects](frame-objects.html), [frame-grid](frame-grid.html),
[frames](frames.html), [delete-frame](delete-frame.html)
