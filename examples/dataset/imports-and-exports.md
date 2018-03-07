
# imports and exports

## importing data into a collection

We can import data from a CSV file and store each row as a JSON document in dataset. In
this example we're generating a UUID for the key name of each row. Our dataset is called
*mydata.ds*. Because we have used the `-uuid` option, no columm number is needed to use
for the key.

```shell
    dataset init mydata.ds
    dataset -uuid mydata.ds import my-data.csv
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.


## exporting data from a collection

```shell
   dataset mydata.ds export titles.csv true '.id,.title,.pubDate' 'id,title,publication date'
```

If you wanted to restrict to a subset (e.g. publication in year 2016)

```shell
   dataset mydata.ds export titles2016.csv '(eq 2016 (year .pubDate))' \
           '.id,.title,.pubDate' 'id,title,publication date'
```

