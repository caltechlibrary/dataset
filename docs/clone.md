
# CLONE

Clone a collection from a list of keys into a new collection.

In this example we create a list of keys using the `-sample` option
and then clone those keys into a new collection called *sample.ds*.

```shell
    dataset keys -sample=3 mycollection.ds > sample.keys
    dataset clone -i sample.keys mycollection.ds sample.ds
```

Related topics: [clone-sample](clone-sample.html)

