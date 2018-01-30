#!/usr/bin/env python3
import ctypes
import os

# Figure out shared library extension
go_basename = 'dataset'
uname = os.uname().sysname
ext = '.so'
if uname == 'Darwin':
    ext = '.dylib'
if uname == 'Windows':
    ext = '.dll'

# Find our shared library and load it
dir_path = os.path.dirname(os.path.realpath(__file__))
lib = ctypes.cdll.LoadLibrary(os.path.join(dir_path, go_basename+ext))

# Setup our Go functions to be nicely wrapped
go_init_collection = lib.init_collection
go_init_collection.argtypes = [ctypes.c_char_p]
go_init_collection.restype = ctypes.c_int

# Now write our Python idiomatic function
def init_collection(name):
    value = go_init_collection(ctypes.c_char_p(name.encode('utf8')))
    print("DEBUG value returned", value)
    if value == 1:
        return True
    return False


if __name__ == '__main__':
    import sys
    if len(sys.argv) > 1:
        print(init_collection(sys.argv[1]))
    else:
        print("To run tests provide a collection name for testing,", sys.argv[0], '"TestCollection"')

