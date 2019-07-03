
# frame

This command will define a frame or return the contents of a frame.
To define a new frame you need to provide a frame name 
followed by a list of label/dotpaths paris. The labels are used as
object attribute names and the dot paths as the source of data.
You also need a list of keys.  By default the keys are read from 
standard input. With options you can include a specific file or 
even indicate to use all the keys in a collection.  In this example 
we are creating a frame called "title-authors-year" based on the 
titles, authors and publication year from a dataset collection 
called `pubs.ds`. Not the labels of "Title", "Authors", "PubYear"
are on the left side the an equal sign and the dot paths to the 
right (defining the data source). 

```shell
    dataset keys pubs.ds |\
        dataset frame pubs.ds "title-authors-year" \
                Title=.title \
                Authors=.authors \
                PubYear=.publication_year
```

On the object list of the frame in this example the objects 
look like

```json
    {
        "Title": ...,
        "Authors": ...,
        "PubYear": ...,
    }
```

This allows you to create convient names for otherwise deep dot paths.

In python

```python
    keys = dataset.keys('pubs.ds')
    (frame, err) = dataset.frame('pubs.ds', 'title-authors-year', 
         keys, ['.title', '.authors', '.publication_year'], 
         ['Title', 'Authors', 'PubYear'])
```

Related topics: [frames](frames.html), [frame-types](frame-types.html), [reframe](reframe.html), [delete-frame](delete-frame.html)

