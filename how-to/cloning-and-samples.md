
# Cloning

There are times where it is useful to make a copy of a collection, take 
a random sample of a collection or even split a collection into a 
training set and test set.  Cloning is the general term used in dataset
for these types of operations.

## Cloning an entire collection

First since dataset is a folder with collections.json, pairtree and frames 
we can "clone" an the complete collection by simply copying the folder 
and its contents.  This can be done with standard operating system
tools (e.g. File managers, command line) or in the case of cloud storage 
with the tools provided by the vendor for managing content in the cloud. 
If you're working with a whole collection this is the best approach.

## Taking a sample

More frequently I've found I need to take a sample of a collection 
(e.g. because I am developing a bash or Python script that does some 
batch processing and testing on the whole collection takes too long). 
This is where you want to use _dataset_'s **clone** verb.  The basic 
idea is to get a sample list of keys then use the **clone**. In the 
command line version we use the "-sample" option with the **keys** verb, 
in Python you need to supply your own function to get a sample list of 
keys.

In the following examples the origin collection is _friends.ds_ or new 
sample collection will be _friends-sample.ds_

On the command line --

```shell
    dataset keys -sample=5 friends.ds > sample.keys
    dataset clone -i sample.keys friends.ds friends-sample.ds
```

In Python I am assuming you have defined a function called "get_sample_keys()" your self.

```python
    keys = get_sample_keys('friends.ds', 5)
    err = dataset.clone('friends.ds', keys, 'friends-sample.ds')
```

As you can see this version of clone works off a set of supplied keys. 
What if you don't want to calculate the key list first? e.g I am working 
in Python and I don't want to have to write "get_sample_keys()"! That's 
what **clone-sample** (or in Python **clone_sample**) is for.


## Clone sample

The **clone-sample** verb is about generating sample collections without 
having to take the extra step of generating a list of sample keys. As an 
added benefit **clone-sample** knows which keys were not selected in the 
sample so it is convienent for creating "training" and "test" collection 
if you are applying machine learning techniques.

Let's take a shorten version of generating a sample of size 5 for our 
friends collection.

```shell
    dataset clone-sample -size=5 friends.ds friends-sample.ds
```

Likewise in python this becomes

```python
    err = dataset.clone_sample('friends.ds', 5, 'friends-sample.ds')
    if err != '':
        print(err)
```

### Training and test collections

By adding a second target collection name we can use **clone-sample** to 
create both the training and test collection. Here's an example with our 
_friends.ds_ collection creating _training.ds_ and _test.ds_.


```shell
    dataset clone-sample -size=5 friends.ds training.ds test.ds
```

```python
    err = dataset.clone_sample('friends.ds', 5, 'training.ds', 'test.ds')
    if err != '':
        print(err)
```

