
# migrate

## Syntax

```
    dataset migrate COLLECTION_NAME LAYOUT_NAME
```

## Description

_migrate_ trys to convert a collection's file layout. Two layouts
are currently support. The older, default layout is "buckets" and
the new layout is "pairtree".  It is expected that "pairtree" will
become the default layout in the future as the _dataset_ command
evolves.

## Usage

Our collection name is "MyCollectiond.ds". The type
of layout we migrating to in this case is "pairtree".
"pairtree" and "buckets" are the two currently supported
file layout scheme.

```
   dataset migrate MyCollection.ds pairtree
```

Related topic: [check](check.html), [repair](repair.html)

