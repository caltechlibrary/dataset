
## Extract

If wanted to extract a unqie list of all ORCIDs from a collection 

```shell
   dataset extract true .authors[:].orcid
```

If you wanted to extract a list of ORCIDs from publications in 2016.

```shell
   dataset extract '(eq 2016 (year .pubDate))' .authors[:].orcid
```

