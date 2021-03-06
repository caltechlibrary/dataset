
## Filters and sorts

_dataset keys_ can be used to to filter and sort a key list. 
Here's is a simple case for match records in a collection named 
**characters.ds** where the given name is equal to "Mojo". We will 
save the result in a file called _mojo.keys_.

```shell
   dataset keys characters.ds '(eq .given "Mojo")' > mojo.keys
```

You can also use an existing key list (e.g. _mojo.keys_)
to sub select based on a new filter and/or sort expression. 
We will filter for a family name of "Sam" and sort by the age field.

```shell
   dataset keys -key-file=mojo.keys characters.ds \
                '(eq .family "Sam")' '+.age'
```

You can improve the performance of filtering/sorting by
breaking it down to steps for large collections. First filter
the keys you want. Then sort the filtered list.

