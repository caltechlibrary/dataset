
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

```python
    #!/usr/bin/env python3
    import sys
    
    print("Starting import test")
    import dataset
    
    if len(sys.argv) > 1:
        collection_name = sys.argv[1]
        print("Initializing", collection_name)
        print(dataset.init_collection(collection_name))
    else:
        print("To run tests provide a collection name for testing,", sys.argv[0], '"test_collection2.ds"')
        exit(1)
    dataset.verbose_on()
    key = "2488"
    value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}
    
    print("Populating collection")
    if dataset.has_key(collection_name, key) == True:
        print("updating record", key)
        err = dataset.update_record(collection_name, key, value)
        if err == False:
            print("Could not update record",key)
    else:
       print("creating record", key)
       err = dataset.create_record(collection_name, key, value)
       if err == False:
           print("Could not create record",key)
    
    keyCount = dataset.count(collection_name)
    print("Record Count", keyCount)
    keyList = dataset.keys(collection_name)
    print("Keys", keyList)
    rec = dataset.read_record(collection_name, key)
    print("original record", rec)
    for k, v in value.items():
       if not isinstance(v, list):
            if k in rec and rec[k] == v:
                print("found", k, " -> ", v)
       else:
            print("Error, expected", k, "with v",v)
    value["verified"] = True
    err = dataset.update_record(collection_name, key, value)
    if err == False:
       print("Count not update record", key, value)
    rec = dataset.read_record(collection_name, key)
    print("updated record", rec)
    for k, v in value.items():
       if not isinstance(v, list):
           if k in rec and rec[k] == v:
               print("found", k, " -> ", v)
       else:
           print("Error, expected", k, "with v",v)
    err = dataset.delete_record(collection_name, key)
    if err == False:
        print("could not delete record", key)
    cnt = dataset.count(collection_name)
    if cnt != 0:
        print("expected zero records, got", cnt)
    print("All Done!")
```
