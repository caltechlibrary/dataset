hasframe
========

Check to see if a frame name exists in a collection.

``` {.shell}
    dataset hasframe pubs.ds f1
```

In python

``` {.python}
    if dataset.has_frame('pubs.ds', 'f1') == true:
        print('We have frame f1 in pubs.ds')
```

Related topics: [frame](frame.html), [frame-grid](frame-grid.html),
[frame-objects](frame-objects.html), [frames](frames.html),
[reframe](reframe.html), [delete-frame](delete-frame.html)
