
## Import

Import can take a CSV file and store each row as a JSON document in dataset. In
this example we're generating a UUID for the key name of each row

```shell
   dataset -uuid import my-data.csv
```

You can create a CSV export by providing the dot paths for each column and
then givening columns a name.

