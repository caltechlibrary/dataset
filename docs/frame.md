frame
=====

This command will define a data frame or return the contents and
metadata of a defined frame. To define a new frame you need to provide a
collection name, a frame name followed by a list of dotpath/label pairs.
The labels are used as object attribute names and the dot paths as the
source of data. You also need a list of keys.\
By default the keys are read from standard input. With options you can
include a specific file or even indicate to use all the keys in a
collection. In this example we are creating a frame called
\"title-authors-year\" based on the titles, authors and publication year
from a dataset collection called `pubs.ds`. Note the labels of
\"Title\", \"Authors\", \"PubYear\" are on the right side the an equal
sign and the dot paths to the left.

``` {.shell}
    dataset keys pubs.ds |\
        dataset frame pubs.ds "title-authors-year" \
                ".title=Title" \
                ".authors=Authors" \
                ".publication_year=PubYear"
```

The objects in the frame\'s object list will look like

``` {.json}
    {
        "Title": ...,
        "Authors": ...,
        "PubYear": ...,
    }
```

This allows you to create convenient names for otherwise deep dot paths.

In python we use a Dict to map the dotpaths to labels rather than an
embedded equal sign. Doing the same task as before would look like this
in Python.

``` {.python}
    keys = dataset.keys("pubs.ds")
    (frame, err) = dataset.frame("pubs.ds", "title-authors-year", 
         keys, { 
             ".title": "Title", 
             ".authors": "Authors",
              ".publication_year": "PubYear"
         })
```

Related topics: [frames](frames.html),
[frame-objects](frame-objects.html), [frame-grid](frame-grid.html),
[frame-types](frame-types.html), [reframe](reframe.html),
[delete-frame](delete-frame.html)
