Cloning
=======

There are times where it is useful to make a copy of a collection, take
a random sample of a collection or even split a collection into a
training set and test set. Cloning is the general term used in dataset
for these types of operations.

Cloning an entire collection
----------------------------

First since dataset is a folder with collections.json, pairtree and
frames we can \"clone\" an the complete collection by simply copying the
folder and its contents. This can be done with standard operating system
tools (e.g. File managers, command line).

Taking a sample
---------------

More frequently I\'ve found I need to take a sample of a collection
(e.g. because I am developing a bash or Python script that does some
batch processing and testing on the whole collection takes too long).
This is where you want to use *dataset*\'s **clone** verb. The basic
idea is to get a sample list of keys then use the **clone**. In the
command line version we use the \"-sample\" option with the **keys**
verb, in Python you need to supply your own function to get a sample
list of keys.

In the following examples the origin collection is *friends.ds* or new
sample collection will be *friends-sample.ds*

On the command line \--

```shell
    dataset keys -sample=5 friends.ds > sample.keys
    dataset clone -i sample.keys friends.ds friends-sample.ds
```

Clone sample
------------

The **clone-sample** verb is about generating sample collections without
having to take the extra step of generating a list of sample keys. As an
added benefit **clone-sample** knows which keys were not selected in the
sample so it is convenient for creating \"training\" and \"test\"
collection if you are applying machine learning techniques.

Let\'s take a shorten version of generating a sample of size 5 for our
friends collection.

```shell
    dataset clone-sample -size=5 friends.ds friends-sample.ds
```

### Training and test collections

By adding a second target collection name we can use **clone-sample** to
create both the training and test collection. Here\'s an example with
our *friends.ds* collection creating *training.ds* and *test.ds*.

```shell
    dataset clone-sample -size=5 friends.ds training.ds test.ds
```

