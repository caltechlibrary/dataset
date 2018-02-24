#!/usr/bin/env python3
# 
# py/dataset.go is a C shared library targetting support in Python for dataset
# 
# @author R. S. Doiel, <rsdoiel@library.caltech.edu>
#
# Copyright (c) 2017, Caltech
# All rights not granted herein are expressly reserved by Caltech.
# 
# Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
# 
# 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
# 
# 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
# 
# 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
# 
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
# 
import ctypes
import os
import json

# Figure out shared library extension
go_basename = 'libdataset'
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

go_create_record = lib.create_record
go_create_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_create_record.restype = ctypes.c_int

go_read_record = lib.read_record
go_read_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
go_read_record.restype = ctypes.c_char_p

go_update_record = lib.update_record
go_update_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_update_record.restype = ctypes.c_int

go_delete_record = lib.delete_record
go_delete_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
go_delete_record.restype = ctypes.c_int

go_has_key = lib.has_key
go_has_key.argtypes = [ctypes.c_char_p,ctypes.c_char_p]
go_has_key.restype = ctypes.c_int

go_keys = lib.keys
go_keys.argtypes = [ctypes.c_char_p]
go_keys.restype = ctypes.c_char_p

go_count = lib.count
go_count.argtypes = [ctypes.c_char_p]
go_count.restype = ctypes.c_int

verbose_on = lib.verbose_on
verbose_off = lib.verbose_off

#
# Now write our Python idiomatic function
#

# Initializes a Dataset Collection
def init_collection(name):
    '''initialize a dataset collection with the given name'''
    value = go_init_collection(ctypes.c_char_p(name.encode('utf8')))
    if value == 1:
        return True
    return False

# Has key, checks if a key is in the dataset collection
def has_key(name, key):
    value = go_has_key(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if value == 1:
        return True
    return False

# Create a JSON record in a Dataset Collectin
def create_record(name, key, value):
    '''create a new JSON record in the collection based on collection name, record key and JSON string, returns True/False'''
    err = go_create_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if err == 1:
        return True
    return False
    
# Read a JSON record from a Dataset collection
def read_record(name, key):
    '''read a JSON record from a collection with the given name and record key, returns a dict'''
    value = go_read_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return json.loads(value.decode())
    

# Update a JSON record from a Dataset collection
def update_record(name, key, value):
    '''update a JSON record from a collection with the given name, record key, JSON string returning True/False'''
    err = go_update_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if err == 1:
        return True
    return False
    

# Delete a JSON record from a Dataset collection
def delete_record(name, key):
    '''delete a JSON record (and any attachments) from a collection with the collectin name and record key, returning True/False'''
    err = go_delete_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if err == 1:
        return True
    return False
    

# Keys returns a list of keys from a collection
def keys(name):
    '''keys returns a list of keys'''
    value = go_keys(ctypes.c_char_p(name.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    return json.loads(value.decode())
    

# Count returns an integer of the number of keys in a collection
def count(name, filter = ''):
    '''count returns an integer of the number of keys in a collection'''
    value = go_count(ctypes.c_char_p(name.encode('utf8')))
    return value

if __name__ == '__main__':
    import sys
    if len(sys.argv) > 1:
        print(init_collection(sys.argv[1]))
        collection_name = sys.argv[1]
    else:
        print("To run tests provide a collection name for testing,", sys.argv[0], '"TestCollection"')
        exit(1)
    verbose_on()
    key = "2488"
    value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}

    if has_key(collection_name, key) == True:
        print("updating record", key)
        err = update_record(collection_name, key, value)
        if err == False:
            print("Could not update record",key)
    else:
        print("creating record", key)
        err = create_record(collection_name, key, value)
        if err == False:
            print("Could not create record",key)
    keyCount = count(collection_name)
    print("Record Count", keyCount)
    keyList = keys(collection_name)
    print("Keys", keyList)
    rec = read_record(collection_name, key)
    print("original record", rec)
    for k, v in value.items():
        if not isinstance(v, list):
            if k in rec and rec[k] == v:
                print("found", k, " -> ", v)
            else:
                printf("Error, expected", k, "with v",v)
    value["verified"] = True
    err = update_record(collection_name, key, value)
    if err == False:
        print("Count not update record", key, value)
    rec = read_record(collection_name, key)
    print("updated record", rec)
    for k, v in value.items():
        if not isinstance(v, list):
            if k in rec and rec[k] == v:
                print("found", k, " -> ", v)
            else:
                printf("Error, expected", k, "with v",v)
    err = delete_record(collection_name, key)
    if err == False:
        print("could not delete record", key)
    cnt = count(collection_name)
    if cnt != 0:
        print("expected zero records, got", cnt)

