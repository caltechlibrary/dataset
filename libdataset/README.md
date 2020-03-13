
# libdataset

This directory holds a C shared library version of the _dataset_ Go 
package. It is used to support using _dataset_ from other languages 
such as Python 3 via the ctypes binding.

## Linux/Darwin

To compile you need to have Go 1.13 and 
GNU Make. Running `make` in this directory will generate the compiled 
shared library and create header file (e.g. libdataset.so, libdataset.dll, 
or libdataset.dylib and libdataset.h).  You can then copy the shared 
library and header file to an appropriate on your system.

## Windows 10

Install Go 1.13 or better from the Golang website using the provided Windows binaries. Install Miniconda (from Anaconda). Using Miniconda install git and gcc (i.e. m2w64-gcc). Run "make.bat" to compile DLL.


## Issues

+ Windows 10: Need to install gcc and git via Miniconda after installing the Go binaries for Windows from the golang.org website
+ Mac OS X (Darwin): only one Go c-shared library seems possible at in a python session, the Go code doesn't seem to be movable in memory, this is related to a long standing issue in Mac OS X only supporting xcode's linker
