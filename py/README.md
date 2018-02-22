
# Python3 Module for Dataset

This directory includes code to build a Go C-Shared library to use with the python3 wrapper
for accessing _dataset_ functionality from Python3.

## Compiling

You need to have Go v1.10 or better installed and Python3.5 with ctypes installed.
If those are installed running _make_ in this directory should build the modules.
You can test the compiled version with _make test_ and build a release zip file with
_make release_.

