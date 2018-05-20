
# frames-types

Set the types associated with the dotpath in the frame. If 
the frame "f1" has the following dotpaths - .title, .authors, and .year
then set the types as string, string and year. Not the order
the dotpaths are defined is the order you're apply the types.

```shell
    dataset pubs.ds frame-types f1 string string year
```

In python

```python
    err = dataset.frame_types('pubs.ds', 'f1', ['string', 'string', 'year'])
```


