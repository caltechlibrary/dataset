
# NAME

libdataset

# SYNOPSIS

Use via C.

~~~
include "libdataset.h"
~~~

Use via Python.

~~~
from py_dataset import dataset
~~~

# DESCRIPTION

libdataset is a C shared library based on the Go package called
dataset from Caltech Library.  The dataset package provides a unified
way of working with JSON documents as collections. libdataset was
create better integrate working with dataset collection from Python
via the [py_dataset](https://pypi.org/project/py-dataset/) Python package.


# Requirements

To compile libdataset you need both a the Go 1.20 toolchain and
a POSIX C toolchain including GNU Make. You need to understand your
C compiler options to include the dataset shared libraries generated
via the Go compiler. 

# Compiling libdataset

The following includes some general notes about compiling libdataset.

## Linux/Darwin

To compile you need to have Go 1.20 and 
GNU Make. Running `make` in this directory will generate the compiled 
shared library and create header file (e.g. libdataset.so, libdataset.dll, 
or libdataset.dylib and libdataset.h).  You can then copy the shared 
library and header file to an appropriate on your system.

## Windows 11

Install Go 1.20 or better from the Golang website using the provided 
Windows binaries. Install Miniconda (from Anaconda). Using Miniconda 
install git, gcc (i.e. m2w64-gcc) and zip (m2-zip). Run "make.bat" to 
compile DLL. Modify and run release.bat to generate a release version.


## Issues

+ Windows 11: Need to install gcc and git via Miniconda after installing 
the Go binaries for Windows from the golang.org website
+ Mac OS X (Darwin): only one Go c-shared library seems possible at in a 
python session, the Go code doesn't seem to be movable in memory, this is 
related to a long standing issue in Mac OS X only supporting xcode's 
linker

# LICENSE

This software is licensed under a varation of the BSD license. See LICENSE
in the source repository [LICENSE](https://github.com/caltechlibrary/dataset/master/LICENSE) for details.



