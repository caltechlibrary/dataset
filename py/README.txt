
# Python3 Module for Dataset

This directory includes code to build a Go C-Shared library to use from
Python 3.  The makefile is used to compile the Go code to a shared library
format needed to take advantage of Python 3's `ctypes` package.

## Compiling

You need to have Go v1.10 or better installed and Python3.5 with `ctypes`i
installed.  If those are installed running _make_ in this directory should 
build the modules.  You can test the compiled version with _make test_ and 
build a release zip file with _make release_.

## Installation

The shared library (i.e. `libdataset.so`, `libdataset.dll` or `libdataset.dylib`) needs to be in the same directory as `dataset.py` which in turn
needs to be in your Python environment's search path.

## Usage

The file `dataset_test.py` shows you how to use the basic functions
in Python via importing the module `dataset`.

