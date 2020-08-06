reframe
=======

This command replace the current keys and object list in a frame based
on the keys provided.

In the following example the frame name is \"f1\", the collection is
\"examples.ds\". The first example is reframing an existing frame using
existing keys coming from standard input, the second example performs
the same thing but is taking a filename to retrieve the list of keys.

``` {.shell}
    cat f1-updated.keys | dataset reframe example.ds f1
    dataset reframe example.ds f1 f1-updated.keys
```

In python

``` {.python}
    f1_updated_keys = generate_updates_keys()
    err = dataset.frame_reframe('example.ds', 'f1', f1_updated_keys)
```

Releted topics: [frame](frame.html), [refresh](refresh.html),
[frame-objects](frame-objects.html), [frame-grid](frame-grid.html),
[frames](frames.html), [delete-frame](delete-frame.html)
