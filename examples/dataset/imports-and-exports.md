
# imports and exports

## importing data into a collection

We can import data from a CSV file and store each row as a JSON document in dataset. You
need to pick a column with unique values to be the key for each record in the collection.
In this example we assume column one has the key value.

```shell
    dataset init mydata.ds
    dataset mydata.ds import-csv my-data.csv 1
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.


## exporting data from a collection

```shell
   dataset mydata.ds export-csv titles.csv true '.id,.title,.pubDate' 'id,title,publication date'
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```shell
   dataset mydata.ds export-csv titles2016.csv '(eq 2016 (year .pubDate))' \
           '.id,.title,.pubDate' 'id,title,publication date'
```

