
# Migrating your dataset collection file layout

The _dataset_ Go package is still rapidly evolving though it is
now commonly used in Caltech Library. As a result of this evolution
we are experimenting with two different file layouts currentely.
The older layout is called "buckets", the newer layout is a
"pairtree". Currently "buckets" is the default but you can migrate
your collection form one to the other.  Below is an example
of migrating to the "pairtree" file layout.

```
    # Migrating to a pairtree layout
    dataset migrate mycollection.ds pairtree
```

To migrate from a "pairtree" to a bucket follows the same process.

```
    # Migrating to a bucket layout
    dataset migrate mycollection.ds buckets
```

