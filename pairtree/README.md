
Pairtree
========

This is a library for translate a UTF-8 string to/from a pairtree
notation. This is typically used in storing things on disc (e.g. repository filesystems). This code is based on the specification found at [OCLC's Confluence website](https://confluence.ucop.edu/download/attachments/14254128/PairtreeSpec.pdf?version=2&modificationDate=1295552323000&api=v2 "this site is not always available") which is cited on the [OCFL](https://github.com/OCFL/spec/wiki) wiki. A draft IETF spec by John Kunze et el. can be found at https://datatracker.ietf.org/doc/html/draft-kunze-pairtree-01. 

NOTE: If you are looking for Python or Java implementedations they can befound at [PyPi](https://pypi.org/project/Pairtree/) and on [GitHub](https://github.com/LibraryOfCongress/pairtree)

This Go package is managed as a sub-module of the dataset project developed
at Caltech Library.

Features
--------

- `Set()` will let you set the path separator
    - `Separator` is a readonly value of the file separator used by `Encode()` and `Decode()`
- `Encode()` will encode the provided string as a pairtree path
- `Decode()` will decode a pairtree path returning the unencoded string

Example
-------

```
    import (
        "fmt"
        "os"

        "github.com/caltechlibrary/pairtree"
    )

    func main() {
        key := "12mIEERD11"
        fmt.Printf("Key: %q\n", key)
        pairPath := pairtree.Encode(key)
        fmt.Printf("Endoded key %q -> %q\n", key, pairPath)
        key = Decode(pairPath)
        fmt.Printf("Decoded path %q -> %q\n", pairPath, key)
    }
```


