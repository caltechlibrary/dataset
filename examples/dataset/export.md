
## Export

```shell
   dataset export titles.csv true '.id,.title,.pubDate' 'id,title,publication date'
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```shell
   dataset export titles2016.csv '(eq 2016 (year .pubDate))' \
           '.id,.title,.pubDate' 'id,title,publication date'
```

