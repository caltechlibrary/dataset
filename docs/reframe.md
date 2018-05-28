
# reframe

This command will regenerate the contents of a frame based on the
current records associated with the keys in the collection. Optionally
you can supply a list of new keys which will then replace the existing
list of keys before regenerating the content.

In the following example the frame name is "f1", the collection is
"examples.ds". The first example is reframing an existing frame using
existing keys while the anod the second is of replace the keys of 
the frame before regenerating it.

```shell
    dataset example.ds reframe f1
    dataset example.ds reframe f1 subset.keys
```

In python

```python
    err = dataset.reframe('example.ds', 'f1')
    subset_keys = generate_subset(dataset.keys('examples.ds'))
    err = dataset.reframe('example.ds', 'f1', subset_keys)
```

Releted topics: [frame](frame.html), [frames](frames.html), [frame-labels](frame-labels.html), [frame-types](frame-types.html), [delete-frame](delete-frame.html)

