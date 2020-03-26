#!/usr/bin/env python3
# 
# libdataset is a wrapper around our C-Shared library of libdataset.go
# used for testing the C-Shared library functions.
# 
# @author R. S. Doiel, <rsdoiel@library.caltech.edu>
#
# Copyright (c) 2020, Caltech
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
import json
import sys
import os
from ctypes import c_char_p, c_int, c_bool
from .libdataset import libdataset

#
# These are our Python idiomatic functions
# calling the C type wrapped functions in libdataset.py
#

# error_clear() clears the errors previously set.
def error_clear():
    libdataset.error_clear()

# error_message() returns the current error message(s)
# accumulated.
def error_message():
    value = libdataset.error_message()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

# use_strict_dotpath() will trigger if False will prefix a root period
# for dot paths and labels when specified in function calls.
def use_strict_dotpath(on_off = True):
    return libdataset.use_strict_dotpath(True)

# is_verbose() returns true is verbose is enabled, false otherwise
def is_verbose():
    return libdataset.is_verbose()

# verbose_on() turns verboseness off
def verbose_on():
    return libdataset.verbose_on()

# verbose_off() turns verboseness on
def verbose_off():
    return libdataset.verbose_off()

# dataset_version() returns version of dataset 
# shared library semver.
def dataset_version():
    value = libdataset.dataset_version()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode() 

#
# Now write our Python idiomatic function
#

# Initializes a Dataset Collection
def init(collection_name):
    '''initialize a dataset collection with the given name'''
    if libdataset.init_collection(c_char_p(collection_name.encode('utf8'))):
        return ''
    return error_message()

# is_open checks to see if a dataset collection is already open
def is_open(collection_name):
    return libdataset.is_collection_open(c_char_p(collection_name.encode('utf8')))

# open_collection opens a dataset collection (it needs to exist)
def open_collection(collection_name):
    if libdataset.open_collection(c_char_p(collection_name.encode('utf8'))):
        return ''
    return error_message()

# close closes a dataset collection
def close_collection(collection_name):
    if libdataset.close_collection(c_char_p(collection_name.encode('utf8'))):
        return ''
    return error_message()

# close_all closes all open dataset collection
def close_all():
    if libdataset.close_all_collections():
        return ''
    return error_message()

# Has key, checks if a key is in the dataset collection
def has_key(collection_name, key):
    return libdataset.key_exists(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')))

# Create a JSON record in a Dataset Collectin
def create(collection_name, key, value):
    '''create a new JSON record in the collection based on collection name, record key and JSON string, returns True/False'''
    if isinstance(key, str) == False:
        key = f"{key}"
    if libdataset.create_object(c_char_p(collection_name.encode('utf8')),
            c_char_p(key.encode('utf8')),
            c_char_p(json.dumps(value).encode('utf8'))):
        return ''
    return error_message()
    
# Read a JSON record from a Dataset collection
def read(collection_name, key, clean_object = False):
    '''read a JSON record from a collection with the given name and record key, returns a dict and an error string'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    value = libdataset.read_object(c_char_p(collection_name.encode('utf8')), 
            c_char_p(key.encode('utf8')), clean_object)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    rval = value.decode()
    if type(rval) is str:
        if rval == "":
            return {}, error_message()
        return json.loads(rval), ''
    return {}, f"Can't read {key} from {collection_name}, {error_message()}"
    
# Read a list of JSON records from a Dataset collection
# NOTE: this provides dataset cli behavior for reading back a list
# of records effeciently ...
def read_list(collection_name, keys, clean_object = False):
    # Pack our keys as an array of string
    l = []
    for key in keys:
        if not isinstance(key, str):
            key = f"{key}"
        l.append(key)
    # Generate our JSON version of they key list
    keys_as_json = json.dumps(l)
    value = libdataset.read_object_list(c_char_p(collection_name.encode('utf-8')), c_char_p(keys_as_json.encode('utf-8')), clean_object)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    rval = value.decode()
    if isinstance(rval, str):
        if rval == "":
            return [], error_message()
        return json.loads(rval), error_message()
    return [], f"Can't read {keys} from {collection_name}, {error_message()}"



# Update a JSON record from a Dataset collection
def update(collection_name, key, value):
    '''update a JSON record from a collection with the given name, record key, JSON string returning True/False'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    if libdataset.update_object(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')), c_char_p(json.dumps(value).encode('utf8'))):
        return ''
    return error_message()

# Delete a JSON record from a Dataset collection
def delete(collection_name, key):
    '''delete a JSON record (and any attachments) from a collection with the collectin name and record key, returning True/False'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    if libdataset.delete_object(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8'))):
        return ''
    return error_message()

# Keys returns a list of keys from a collection optionally applying a filter or sort expression
def keys(collection_name):
    '''keys returns an unsorted list of keys for a collection'''
    value = libdataset.keys(c_char_p(collection_name.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)
    
# Key filter takes a list of keys and filter expression and returns a filtered list of keys
def key_filter(collection_name, keys, filter_expr):
    '''key_filter takes a list of keys (if empty or * then it uses all keys in collection) and a filter expression returning a filtered list'''
    key_list = json.dumps(keys)
    value = libdataset.key_filter(c_char_p(collection_name.encode('utf8')), c_char_p(key_list.encode('utf8')), c_char_p(filter_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)
    
# Key sort takes sort expression and an optional list of keys and returns a sorted list of keys
def key_sort(collection_name, keys, sort_expr):
    '''key_filter takes a list of keys (if empty or * then it uses all keys in collection) and a filter expression returning a filtered list'''
    key_list = json.dumps(keys)
    value = libdataset.key_sort(c_char_p(collection_name.encode('utf8')), c_char_p(key_list.encode('utf8')), c_char_p(sort_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)

# Count returns an integer of the number of keys in a collection
def count(collection_name, filter = ''):
    '''count returns an integer of the number of keys in a collection'''
    return libdataset.count_objects(c_char_p(collection_name.encode('utf8')))


#
# import_csv - import a CSV file into a collection
# syntax: COLLECTION CSV_FILENAME ID_COL
# 
# options:
#
#      use_header_row (bool)
#      overwrite (bool)
# 
# Returns: error string
def import_csv(collection_name, csv_name, id_col, use_header_row = True, overwrite = False):
    if libdataset.import_csv(c_char_p(collection_name.encode('utf8')), 
            c_char_p(csv_name.encode('utf8')), 
            c_int(id_col), use_header_row, 
            overwrite):
        return ''
    return error_message()

#
# export_csv - export collection objects to a CSV file
# syntax: COLLECTION FRAME CSV_FILENAME
# 
# Returns: error string
def export_csv(collection_name, frame_name, csv_name):
    if libdataset.export_csv(c_char_p(collection_name.encode('utf8')), 
            c_char_p(frame_name.encode('utf8')), 
            c_char_p(csv_name.encode('utf8'))):
        return ''
    return error_message()

def status(collection_name):
    return libdataset.collection_exists(collection_name.encode('utf8'))

def list(collection_name, keys = []):
    src_keys = json.dumps(keys)
    value = libdataset.list_objects(c_char_p(collection_name.encode('utf8')), c_char_p(src_keys.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    if len(value) == 0:
        return [] 
    return json.loads(value.decode()) 

def path(collection_name, key):
    value = libdataset.object_path(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    return value.decode()

def check(collection_name):
    ok = libdataset.check_collection(c_char_p(collection_name.encode('utf8')))
    return (ok == True)

def repair(collection_name):
    return libdataset.repair_collection(c_char_p(collection_name.encode('utf8')))

def attach(collection_name, key, filenames = [], semver = ''):
    if semver == '':
        semver = 'v0.0.0'
    srcFNames = json.dumps(filenames)
    if not isinstance(srcFNames, bytes):
        srcFNames = srcFNames.encode('utf8')
    return libdataset.attach(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')), c_char_p(semver.encode('utf8')), c_char_p(srcFNames))
    
def attachments(collection_name, key):
    value = libdataset.attachments(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    s = value.decode().strip(' ')
    if len(s) > 0:
        return s.split("\n")
    return ''

def detach(collection_name, key, filenames = [], semver = ''):
    '''Get attachments for a specific key.  If the version semver is not provided, it will default to the current version.  Provide [] as filenames if you want to get all attachments'''
    if semver == '':
        semver = 'v0.0.0'
    srcFNames = json.dumps(filenames)
    if not isinstance(srcFNames, bytes):
        srcFNames = srcFNames.encode('utf8')
    return libdataset.detach(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')), c_char_p(semver.encode('utf8')), c_char_p(srcFNames))

def prune(collection_name, key, filenames = [], semver = ''):
    '''Delete attachments for a specific key.  If the version semver is not provided, it will default to the current version.  Provide [] as filenames if you want to delete all attachments'''
    if semver == '':
        semver = 'v0.0.0'
    fnames = json.dumps(filenames).encode('utf8')
    return libdataset.prune(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')), c_char_p(semver.encode('utf8')), c_char_p(fnames))

def join(collection_name, key, obj = {}, overwrite = False):
    src = json.dumps(obj).encode('utf8')
    cOverwrite = c_int(0)
    if overwrite == True:
        cOverwrite = c_int(1)
    return libdataset.join_objects(c_char_p(collection_name.encode('utf8')), c_char_p(key.encode('utf8')), c_char_p(src), cOverwrite)

def clone(collection_name, keys, destination_name):
    src_keys = json.dumps(keys)
    return libdataset.clone_collection(c_char_p(collection_name.encode('utf-8')), c_char_p(src_keys.encode('utf-8')), c_char_p(destination_name.encode('utf-8')))

def clone_sample(collection_name, training_name, test_name = "", sample_size = 0):
    return libdataset.clone_sample( c_char_p(collection_name.encode('utf-8')), c_char_p(training_name.encode('utf-8')), c_char_p(test_name.encode('utf-8')), c_int(sample_size))

def frame_create(collection_name, frame_name, keys, dot_paths, labels):
    src_keys = json.dumps(keys)
    src_dot_paths = json.dumps(dot_paths)
    if len(labels) == 0 and len(dot_paths) > 0:
        for item in dot_paths:
            if item.startswith("."):
                item = item[1:]
            labels.append(item)
    src_labels = json.dumps(labels)
    return libdataset.frame_create(c_char_p(collection_name.encode('utf-8')),
        c_char_p(frame_name.encode('utf-8')),
        c_char_p(src_keys.encode('utf-8')),
        c_char_p(src_dot_paths.encode('utf-8')),
        c_char_p(src_labels.encode('utf-8')))


def has_frame(collection_name, frame_name):
    return libdataset.frame_exists(c_char_p(collection_name.encode('utf-8')),
            c_char_p(frame_name.encode('utf-8')))

def frame_keys(collection_name, frame_name):
    value = libdataset.frame_keys(c_char_p(collection_name.encode('utf-8')),
            c_char_p(frame_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return json.loads(value)


def frame(collection_name, frame_name):
    value = libdataset.frame(c_char_p(collection_name.encode('utf-8')),
            c_char_p(frame_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '' or len(value) == 0:
        return None
    return json.loads(value)

    
def frame_objects(collection_name, frame_name):
    value = libdataset.frame_objects(c_char_p(collection_name.encode('utf-8')),
            c_char_p(frame_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '' or len(value) == 0:
        return []
    return json.loads(value)

def frames(collection_name):
    value = libdataset.frames(c_char_p(collection_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '' or len(value) == 0: 
        return [] 
    return json.loads(value)

def frame_reframe(collection_name, frame_name, keys = []):
    src_keys = json.dumps(keys)
    return libdataset.frame_reframe(c_char_p(collection_name.encode('utf-8')),
        c_char_p(frame_name.encode('utf-8')),
        c_char_p(src_keys.encode('utf-8')))

def frame_refresh(collection_name, frame_name):
    return libdataset.frame_refresh(c_char_p(collection_name.encode('utf-8')))

def frame_clear(collection_name, frame_name):
    return libdataset.frame_clear(c_char_p(collection_name.encode('utf-8')), c_char_p(frame_name.encode('utf-8')))

def delete_frame(collection_name, frame_name):
    return libdataset.frame_delete(c_char_p(collection_name.encode('utf-8')), c_char_p(frame_name.encode('utf-8')))

def frame_grid(collection_name, frame_name, include_headers = True):
    header_int = 0
    if include_headers == True:
        header_int = 1
    value = libdataset.frame_grid(c_char_p(collection_name.encode('utf-8')),
            c_char_p(frame_name.encode('utf-8')),
            header_int)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '':
        return []
    return json.loads(value)

def sync_recieve_csv(collection_name, frame_name, csv_filename, overwrite = False):
    overwrite_i  = 0
    if overwrite:
        overwrite_i = 1
    return libdataset.sync_recieve_csv(
            c_char_p(collection_name.encode('utf-8')), 
            c_char_p(frame_name.encode('utf-8')), 
            c_char_p(csv_filename.encode('utf-8')), 
            c_int(overwrite_i))

def sync_send_csv(collection_name, frame_name, csv_filename, overwrite = False):
    overwrite_i = 0
    if overwrite:
        overwrite_i = 1
    return libdataset.sync_send_csv(
            c_char_p(collection_name.encode('utf-8')), 
            c_char_p(frame_name.encode('utf-8')), 
            c_char_p(csv_filename.encode('utf-8')), 
            c_int(overwrite_i))

def create_objects(collection_name, keys, default_object):
    c_name = c_char_p(collection_name.encode('utf-8'))
    keys_as_json = c_char_p(json.dumps(keys).encode('utf8'))
    object_as_json = c_char_p(json.dumps(default_object).encode('utf8'))
    return libdataset.create_objects(c_name, keys_as_json, object_as_json)

def update_objects(collection_name, keys, objects):
    c_name = c_char_p(collection_name.encode('utf-8'))
    keys_as_json = c_char_p(json.dumps(keys).encode('utf8'))
    objects_as_json = c_char_p(json.dumps(objects).encode('utf8'))
    return libdataset.update_objects(c_name, keys_as_json, objects_as_json)

def set_who(collection_name, names = []):
    c_name = c_char_p(collection_name.encode('utf-8'))
    names_as_json = c_char_p(json.dumps(names).encode('utf8'))
    return libdataset.set_who(c_name, names_as_json)

def set_what(collection_name, src = ""):
    c_name = c_char_p(collection_name.encode('utf-8'))
    c_src = c_char_p(src.encode('utf8'))
    return libdataset.set_what(c_name, c_src)

def set_when(collection_name, src = ""):
    c_name = c_char_p(collection_name.encode('utf-8'))
    c_src = c_char_p(src.encode('utf8'))
    return libdataset.set_when(c_name, c_src)

def set_where(collection_name, src = ""):
    c_name = c_char_p(collection_name.encode('utf-8'))
    c_src = c_char_p(src.encode('utf8'))
    return libdataset.set_where(c_name, c_src)

def set_version(collection_name, src = ""):
    c_name = c_char_p(collection_name.encode('utf-8'))
    c_src = c_char_p(src.encode('utf8'))
    return libdataset.set_version(c_name, c_src)

def set_contact(collection_name, src = ""):
    c_name = c_char_p(collection_name.encode('utf-8'))
    c_src = c_char_p(src.encode('utf8'))
    return libdataset.set_contact(c_name, c_src)


def get_who(collection_name):
    c_name = c_char_p(collection_name.encode('utf-8'))
    value = libdataset.get_who(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    rval = value.decode()
    if type(rval) is str:
        return json.loads(rval)
    return []

def get_what(collection_name):
    c_name = c_char_p(collection_name.encode('utf-8'))
    value = libdataset.get_what(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_where(collection_name):
    c_name = c_char_p(collection_name.encode('utf-8'))
    value = libdataset.get_where(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_when(collection_name):
    c_name = c_char_p(collection_name.encode('utf-8'))
    value = libdataset.get_when(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_version(collection_name):
    c_name = c_char_p(collection_name.encode('utf-8'))
    value = libdataset.get_version(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_contact(collection_name):
    c_name = c_char_p(collection_name.encode('utf-8'))
    value = libdataset.get_contact(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()











