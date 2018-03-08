
## Extract

If wanted to extract a unqie list of all ORCIDs from a collection. For
our example we assume that you have already created a collection called *publishers.ds".

```shell
   dataset publishers.ds extract true .authors[:].orcid
```

If you wanted to extract a list of ORCIDs from publications in 2016.

```shell
   dataset publishers.ds extract '(eq 2016 (year .pubDate))' .authors[:].orcid
```

