keys
====

Syntax
------

```shell
    dataset keys COLLECTION_NAME
```

Description
-----------

List the JSON_DOCUMENT_ID available in a collection. Key order is not
guaranteed. Keys are forced to lower case when the record is created
in the dataset (as of version 1.0.2). Note combining "keys" with
a pipe and POSIX commands like "sort" can given a rich pallet of
ways to work with your dataset collection's keys.

Examples
--------

Here are three examples usage. Notice the sorting is handled by
the POSIX sort command which lets you sort ascending or descending
including sorting number strings.

```shell
    dataset keys COLLECTION_NAME
    dataset keys COLLECTION_NAME | sort
    dataset keys COLLECTION_NAME | sort -n
```

Getting a "sample" of keys
--------------------------

The __dataset__ command respects an option named `-sample N` where N 
is the size (number) of the keys to include in the sample. The sample 
is taken after any filters are applied but may be less than requested 
size if the the filtered results are few than the sample size.  The 
basic process is to get a set of keys, randomly sort the keys, then 
return the top N number of those keys.


Related topics: [count](count.md), [clone](clone), [clone-sample](clone-sample.md), [frame](frame.md), [frame-grid](frame-grid.md), [frame-objects](frame-objects.md)


