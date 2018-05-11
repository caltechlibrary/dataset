
# Data Grids

Often when processing data it is useful to pull date into a grid format.
_dataset_ provides a verb "grid" for doing just that. Below we're going
to create a small dataset collection called `grid_test.ds`, populate
it with some simple asymetric data (i.e. each record doesn't have the
same fields) then turn this into a 2D JSON array suitable for further
processing in a language like Python, R or Julia. The *grid* verb
is available in the Python module for dataset so we'll show that too.

In both examples the JSON representing our raw data seen in the file
[data-grids.json](data-grids.json)

## command line example 

### Generate a key list

From an existing collection, `grid_test.ds`, create a list of keys.

```shell
    dataset grid_test.ds keys > grid_test.keys
```

### Check a few records to see which will go into our grid.

We have the following keys  in our collection "gutenberg:21489",
"gutenberg:2488", "gutenberg:21839", "gutenberg:3186", "hathi:uc1321060001561131". Let's pick the first one and see what fields we might want
in our grid (notice we're using the `-p` option to pretty print
the JSON record).

```shell
    dataset -p grid_test.ds read "gutenberg:21489"
```

The fields that we're interested in are "._Key", ".title", ".authors",

### Create our grid from our collection

Now that we have a list of keys we're interested and and know the 
dot paths to the fields we're interested in we can create our grid.

```shell
    dataset -p grid_test.ds grid grid_test.keys "._Key" ".title" ".authors" 
```

The results are a 2D array wich rows for each key and cells matching the
contents of the dot paths. Note that a cell may have a complex structure
like that shown with ".authors"

## python 3 example

In this example we're use the _dataset_ python module to read in our
raw JSON data (e.g. [data-grids.json](data-grids.json)) and convert
it into a *dataset* collection called "grid_test.ds". Next
we'll generate our set of keys and finally generate our grid as
a python list of lists.

```python3
    import sys
    import json
    import dataset

    # Read in our test data and convert from JSON into an array of dicts
    f_name = 'data-grids.json'
    with open(f_name, mode = 'r', encoding = 'utf-8') as f:
        src = f.read()
    data = json.loads(src)

    # create our collection
    c_name = 'grid_test.ds'
    err = dataset.init(c_name)
    if err != '':
        print(err)
        sys.exit(1)
    
    # load our test data
    for key in data:
        rec = data[key]
        err = dataset.create(c_name, key, rec)
        if err != '':
            print(err)
            sys.exit(1)
    
    # Create a list of keys and list of dot paths
    keys = dataset.keys(c_name)
    dot_paths = ["._Key", ".title", ".authors"]
    # now we can create our grid
    (g, err) = dataset.grid(c_name, keys, dot_paths)
    if err != '':
        print(err)
        sys.exit(1)

    # Now pretty print our grid
    print(json.dumps(g, indent = 4))
```
