
# libdataset

This directory holds a C shared library version of the _dataset_ Go 
package. It is used to support using _dataset_ from other languages 
such as Python 3 or from C. To compile you need to have Go 1.11 and 
GNU Make. Running `make` in this directory will generate the compiled 
shared library and create header file (e.g. libdataset.so, libdataset.dll, 
or libdataset.dylib and libdataset.h).  You can then copy the shared 
library and header file to an appropriate on your system.

