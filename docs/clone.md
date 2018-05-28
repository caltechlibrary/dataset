
# CLONE

Clone a collection from a list of keys into a new collection.

In this example we create a list of keys using the `-sample` option
and then clone those keys into a new collection called *sample.ds*.

```shell
    dataset -sample=3 mycollection.ds keys > sample.keys
    dataset mycollection.ds clone sample.keys sample.ds
```

Related topics: [clone-sample](clone-sample.html)

