
# Upgrading your dataset collection(s)

The _dataset_ Go package is still rapidly evolving though it is
now commonly used at Caltech Library. As a result we've developed
an easy method for migration collections from an old version to
a new version. There "usual" to upgrading your is to use 
use the "check" and "repair" features of the _dataset_ command 
line tool.

```
    dataset check mycollection.ds
    # you'll get a verbose report to the console
    dataset repair mycollection.ds
    # dataset will not attempt to "repair" including upgrade, your collection
```
