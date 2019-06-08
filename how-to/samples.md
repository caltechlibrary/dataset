
# Samples

Getting a sample of a collection and be done by generating
a random sample of keys or by cloning a collection based on
a random sample of objects.

## Getting a random sample of keys

Create a random sample of keys. This uses the keys verb.

In this example 3 randomly selected keys would be returned 
from the collection named "mycollection.ds".

```shell
    dataset keys -sample 3 mycollection.ds
```

## Cloning a random sample of objects from a collection

When cloning a randomize sample of objects from a collection
you often want two collections. Typically this is a training 
collection and a test collection. The "clone-sample" verb
supports this when you provide two collection destinations.

In this example we create a training and testing collections 
based on a training sample size of 1000.

```shell
    dataset clone-sample -size=1000 mycollection.ds training.ds test.ds
```

If you only want a single sample collection skip the second collection
name.

```shell
    dataset clone-sample -size=1000 mycollection.ds small-sample.ds
```


