
# frame

This command will define a frame or return the contents of a frame.
To define a new frame you need to provide a filename with a list of keys (one per line)
and followed by a list of dotpaths. In this example we are creating a 
frame called "title-authors-year" based on the titles, authors and publication year from a dataset collection
called `pubs.ds`.

```shell
    dataset pubs.ds keys > title-authors-year.keys
    dataset pubs.ds frame  "title-authors-year" title-authors-year.keys .title .authors .publication_year
```

In python

```python
    keys = dataset.keys('pubs.ds')
    (frame, err) = dataset.frame('pubs.ds', 'title-authors-year', keys, ['.title', '.authors', '.publication_year'])
```


