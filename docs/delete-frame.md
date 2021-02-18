delete-frame
============

This is used to removed a frame from a collection.

``` {.shell}
    dataset delete-frame example.ds f1
```

delete frame f1 from collection called example.ds

In python

``` {.python}
    err = dataset.delete_frame('example.ds', 'f1')
```

Related topics: [frame](frame.html), [frames](frames.html),
[frame-types](frame-types.html)
