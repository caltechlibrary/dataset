
# frame

This command will define a frame or return the contents of a frame.
To define a new frame you need to provide a frame name 
followed by a list of dotpaths. You also need a list of keys. By default
the keys are read from standard input. With options you can include
a specific file or even indicate to use all the keys in a collection.
In this example we are creating a frame called "title-authors-year" based 
on the titles, authors and publication year from a dataset 
collection called `pubs.ds`.

```shell
    dataset keys pubs.ds |\
        dataset frame pubs.ds "title-authors-year" \
                .title .authors .publication_year
```

In python

```python
    keys = dataset.keys('pubs.ds')
    (frame, err) = dataset.frame('pubs.ds', 'title-authors-year', 
         keys, ['.title', '.authors', '.publication_year'])
```

Related topics: [frames](frames.html), [frame-labels](frame-labels.html), [frame-types](frame-types.html), [reframe](reframe.html), [delete-frame](delete-frame.html)

