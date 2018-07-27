
# grid

## Syntax

```
    dataset COLLECTION_NAME grid KEY_LIST_FILENAME DOTPATH [DOTPATH ...]
```

## Description

Creates a JSON structure representing a grid. The rows correspond to the
records identified in the key list and the columns are defined by the
list of dot paths.  If a dotpath isn't found then a null is placed in that 
cell.

DOTPATH provided (for DOTPATH see `dataset -help dotpath` and FITLER see `dataset -help filter`).

## Usage

In this example we're turning a small subset of fields available in 
collection called "publications.ds" into JSON structure suitable for 
sorting in python. We are pull the pub date, title, and orcid fields 
into a grid structure. Note that in our example below the orcid 
itself is an array.

```shell
    dataset publications.ds grid .pub_date .title .creators[:].orcid
```

The result is a 2D array of rows and cells (e.g. colums)

Related topics: [dotpath](dotpath.html), [frame](frame.html)

