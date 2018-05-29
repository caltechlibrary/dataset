
## Filters and sorts

_dataset keys_ can be used to to filter and sort a key list. 
Here's is a simple case for match records in a collection named **characters.ds** 
where the given name is equal to "Mojo". We will save the result in a file called _mojo.keys_.

```shell
   dataset characters.ds keys '(eq .given "Mojo")' > mojo.keys
```

You can also use an existing key list (e.g. _mojo.keys_)
to sub select based on a new filter and/or sort expression. 
We will filter for a family name of "Sam" and sort by the age field.

```shell
   dataset -c characters.ds -key-file=mojo.keys keys '(eq .family "Sam")' '+.age'
```

