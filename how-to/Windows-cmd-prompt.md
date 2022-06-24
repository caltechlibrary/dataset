Windows 11 Notes
================

Windows 11 command prompt presents some challenges for working
with JSON on the command line. This is partularly true if you're not
using the Linux subsystem shell. A command line

```shell
    dataset create T1.ds one '{"one":1}'
```

which would work in a POSIX shell fails. The command prompt makes
the JSON expression look like `{one:1}` which is NOT JSON and also not a filename.  As a result working with dataset at the Windows command prompt requires conforming to the command prompt's expectation on quoting. This will work.

```shell
    dataset create T1.ds one "{"""one""":1}
```
