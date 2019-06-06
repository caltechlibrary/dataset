
# frame-labels

Set the labels associated with the frame. The number of labels 
provide must match the number of *dot paths* defined in the frame. 
In this example the collection name is `example.ds`, frame name is 
"f1", the labels are 'Column A', 'Column B', and 'Column C' 
coresponding to the three dotpaths defined in `examples.ds`.

```shell
    dataset frame-labels example.ds f1 'Column A' 'Column B' 'Column C'
```

In python

```python
    err = dataset.frame_labels('example.ds', 'f1', ["Column A", "Column B", "Column C"])
```

Note that setting the labels will cause the object list in the frame to be regenerated from the current state of the collection. This is because the
labels are used as the attribute names for the objects in a frame.

Related topics: [frame](frame.html), [frames](frames.html), [frame-types](frame-types.html), [reframe](reframe.html), [delete-frame](delete-frame.html)

