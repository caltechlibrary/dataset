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
go_keys.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_keys.restype = ctypes.c_char_p

go_key_filter = lib.key_filter
go_key_filter.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_key_filter.restype = ctypes.c_char_p

go_key_sort = lib.key_sort
go_key_sort.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_key_sort.restype = ctypes.c_char_p

go_count = lib.count
go_count.argtypes = [ctypes.c_char_p]
go_count.restype = ctypes.c_int

go_extract = lib.extract
go_extract.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_extract.restype = ctypes.c_char_p

go_indexer = lib.indexer
go_indexer.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
go_indexer.restype = ctypes.c_int


go_deindexer = lib.deindexer
go_deindexer.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
go_deindexer.restype = ctypes.c_int

go_find = lib.find
go_find.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_find.restype = ctypes.c_char_p

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
    ok = go_create_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if ok == 1:
        return True
    return False
    
# Read a JSON record from a Dataset collection
def read_record(name, key):
    '''read a JSON record from a collection with the given name and record key, returns a dict'''
    value = go_read_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    rval = value.decode()
    if rval == "":
        return {}
    return json.loads(rval)
    

# Update a JSON record from a Dataset collection
def update_record(name, key, value):
    '''update a JSON record from a collection with the given name, record key, JSON string returning True/False'''
    ok = go_update_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if ok == 1:
        return True
    return False
    

# Delete a JSON record from a Dataset collection
def delete_record(name, key):
    '''delete a JSON record (and any attachments) from a collection with the collectin name and record key, returning True/False'''
    ok = go_delete_record(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if ok == 1:
        return True
    return False
    

# Keys returns a list of keys from a collection optionally applying a filter or sort expression
def keys(name, filter_expr = "", sort_expr = ""):
    '''keys returns a list of keys, optionally apply a filter and sort expression'''
    value = go_keys(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(filter_expr.encode('utf8')), ctypes.c_char_p(sort_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)
    
# Key filter takes a list of keys and filter expression and returns a filtered list of keys
def key_filter(name, keys, filter_expr):
    '''key_filter takes a list of keys (if empty or * then it uses all keys in collection) and a filter expression returning a filtered list'''
    key_list = json.dumps(keys)
    value = go_key_filter(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key_list.encode('utf8')), ctypes.c_char_p(filter_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)
    
# Key sort takes sort expression and an optional list of keys and returns a sorted list of keys
def key_sort(name, keys, sort_expr):
    '''key_filter takes a list of keys (if empty or * then it uses all keys in collection) and a filter expression returning a filtered list'''
    key_list = json.dumps(keys)
    value = go_key_sort(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(key_list.encode('utf8')), ctypes.c_char_p(filter_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)

# Count returns an integer of the number of keys in a collection
def count(name, filter = ''):
    '''count returns an integer of the number of keys in a collection'''
    value = go_count(ctypes.c_char_p(name.encode('utf8')))
    return value

# Extract unique values from the JSON records in a collection given a filter expression and dot path
def extract(name, filter_expr, dot_expr):
    '''extract unique values from the JSON records in a collection given a filter expression and dot path'''
    value = go_extract(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(filter_expr.encode('utf8')), ctypes.c_char_p(dot_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)
    

# Indexer takes a collection name, an index name, an index map file name, and an optional keylist 
# and creates/updates a Bleve index on disc.
def indexer(name, index_name, index_map_name, key_list = [], batch_size = 0):
    '''indexes a collection given a collection name, bleve index name, index map filename, and optional key list'''
    key_list_src = json.dumps(key_list)
    ok = go_indexer(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(index_name.encode('utf8')), ctypes.c_char_p(index_map_name.encode('utf8')), ctypes.c_char_p(key_list_src.encode('utf8')), ctypes.c_int(batch_size))
    if ok == 1:
        return True
    return False

# Deindexer takes a collection name, an index name, key list and optional batch size deleting the provided keys from 
# the index.
def deindexer(name, index_name, key_list, batch_size = 0):
    '''indexes a collection given a collection name, bleve index name, index map filename, and optional key list'''
    key_list_src = json.dumps(key_list).encode('utf8')
    ok = go_deindexer(ctypes.c_char_p(name.encode('utf8')), ctypes.c_char_p(index_name.encode('utf8')), ctypes.c_char_p(key_list_src), ctypes.c_int(batch_size))
    if ok == 1:
        return True
    return False

# Find takes an index name, query string an optional options dict and returns a search result
def find(index_names, query_string, options = {}):
    '''Find takes an index name, query string an optional options dict and returns a search result'''
    option_src = json.dumps(options)
    value = go_find(ctypes.c_char_p(index_names.encode('utf8')), ctypes.c_char_p(query_string.encode('utf8')), ctypes.c_char_p(option_src.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return {}
    return json.loads(rval)


