#!/usr/bin/env python3
# 
# py_libdataset is a wrapper around our C-Shared library of libdataset.go
# used for testing the C-Shared library functions.
# 
# @author R. S. Doiel, <rsdoiel@library.caltech.edu>
#
# Copyright (c) 2019, Caltech
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
import ctypes

from py_libdataset.cwrapper import go_basename , go_error_clear, go_error_message , go_use_strict_dotpath , go_dataset_version , go_is_verbose , go_verbose_on , go_verbose_off , go_init , go_create_object , go_read_object , go_read_object_list , go_update_object , go_delete_object , go_key_exists , go_keys , go_key_filter , go_key_sort , go_count , go_import_csv , go_export_csv , go_import_gsheet , go_export_gsheet , go_sync_recieve_csv , go_sync_send_csv , go_sync_recieve_gsheet , go_sync_send_gsheet , go_status , go_list , go_path , go_check , go_repair , go_attach , go_attachments , go_detach , go_prune , go_join , go_clone , go_clone_sample , go_grid , go_frame_create, go_frame_keys, go_frame_objects, go_frame_exists , go_frames , go_frame_reframe , go_frame_delete , go_frame_grid , go_update_objects, go_set_who, go_get_who, go_set_what, go_get_what, go_set_where, go_get_where, go_set_when, go_get_when, go_set_version, go_get_version, go_set_contact, go_get_contact

#
# These are our Python idiomatic functions
# calling the C type wrapped functions in libdataset.py
#
def error_clear():
    go_error_clear()

def error_message():
    value = go_error_message()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode() 


def use_strict_dotpath(on_off = True):
    if on_off == True:
        go_use_strict_dotpath(1)
        return True
    go_use_strict_dotpath(0)
    return False

# is_verbose returns true is verbose is enabled, false otherwise
def is_verbose():
    ok = go_is_verbose()
    return (ok == 1)

# verbose_on turns verboseness off
def verbose_on():
    ok = go_verbose_on()
    return (ok == 1)

# verbose_off turns verboseness on
def verbose_off():
    ok = go_verbose_off()
    return (ok == 1)

# Returns version of dataset shared library
def dataset_version():
    value = go_dataset_version()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

#
# Now write our Python idiomatic function
#

def error_message():
    value = go_error_message()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode() 


def use_strict_dotpath(on_off = True):
    if on_off == True:
        go_use_strict_dotpath(1)
        return True
    go_use_strict_dotpath(0)
    return False

# is_verbose returns true is verbose is enabled, false otherwise
def is_verbose():
    ok = go_is_verbose()
    return (ok == 1)

# verbose_on turns verboseness off
def verbose_on():
    ok = go_verbose_on()
    return (ok == 1)

# verbose_off turns verboseness on
def verbose_off():
    ok = go_verbose_off()
    return (ok == 1)

# Initializes a Dataset Collection
def init(collection_name, layout = "pairtree"):
    '''initialize a dataset collection with the given name'''
    collection_layout = 0
    if layout == "buckets":
        collection_layout = 1
    elif layout == "pairtree":
        collection_layout = 2
    ok = go_init(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_int(collection_layout))
    if ok == 1:
        return ''
    return error_message()

# Has key, checks if a key is in the dataset collection
def key_exists(collection_name, key):
    ok = go_key_exists(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(key.encode('utf8')))
    return (ok == 1)

# Create a JSON record in a Dataset Collectin
def create(collection_name, key, value):
    '''create a new JSON record in the collection based on collection name, record key and JSON string, returns True/False'''
    if isinstance(key, str) == False:
        key = f"{key}"
    ok = go_create_object(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(key.encode('utf8')), 
            ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if ok == 1:
        return ''
    return error_message()
    
# Read a JSON record from a Dataset collection
def read(collection_name, key, clean_object = False):
    '''read a JSON record from a collection with the given name and record key, returns a dict and an error string'''
    clean_object_int = ctypes.c_int(0)
    if clean_object == True:
        clean_object_int = ctypes.c_int(1)
    if not isinstance(key, str) == True:
        key = f"{key}"
    value = go_read_object(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(key.encode('utf8')), clean_object_int)
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
    clean_object_int = ctypes.c_int(0)
    if clean_object == True:
        clean_object_int = ctypes.c_int(1)
    l = []
    for key in keys:
        if not isinstance(key, str):
            key = f"{key}"
        l.append(key)
    # Generate our JSON version of they key list
    keys_as_json = json.dumps(l)
    value = go_read_object_list(ctypes.c_char_p(collection_name.encode('utf-8')), ctypes.c_char_p(keys_as_json.encode('utf-8')), clean_object_int)
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
    ok = go_update_object(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if ok == 1:
        return ''
    return error_message()

# Delete a JSON record from a Dataset collection
def delete(collection_name, key):
    '''delete a JSON record (and any attachments) from a collection with the collectin name and record key, returning True/False'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    ok = go_delete_object(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if ok == 1:
        return ''
    return error_message()

# Keys returns a list of keys from a collection optionally applying a filter or sort expression
def keys(collection_name, filter_expr = "", sort_expr = ""):
    '''keys returns a list of keys, optionally apply a filter and sort expression'''
    value = go_keys(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(filter_expr.encode('utf8')), ctypes.c_char_p(sort_expr.encode('utf8')))
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
    value = go_key_filter(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key_list.encode('utf8')), ctypes.c_char_p(filter_expr.encode('utf8')))
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
    value = go_key_sort(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key_list.encode('utf8')), ctypes.c_char_p(sort_expr.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    if rval == "":
        return []
    return json.loads(rval)

# Count returns an integer of the number of keys in a collection
def count(collection_name, filter = ''):
    '''count returns an integer of the number of keys in a collection'''
    return go_count(ctypes.c_char_p(collection_name.encode('utf8')))


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
    if use_header_row == True:
        i_use_header_row = 1
    else:
        i_use_header_row = 0
    if overwrite == True:
        i_overwrite = 1
    else:
        i_overwrite = 0
    ok = go_import_csv(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(csv_name.encode('utf8')), 
            ctypes.c_int(id_col), ctypes.c_int(i_use_header_row), 
            ctypes.c_int(i_overwrite))
    if ok == 1:
        return ''
    return error_message()

#
# export_csv - export collection objects to a CSV file
# syntax: COLLECTION FRAME CSV_FILENAME
# 
# Returns: error string
def export_csv(collection_name, frame_name, csv_name):
    ok = go_export_csv(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(frame_name.encode('utf8')), 
            ctypes.c_char_p(csv_name.encode('utf8')))
    if ok == 1:
        return ''
    return error_message()

# import_gsheet - import a GSheet into a collection
# syntax: COLLECTION GSHEET_ID SHEET_NAME ID_COL CELL_RANGE
# 
# options:
#
#      UseHeaderRow (bool)
#      Overwrite (bool)
#
# Returns: error string
def import_gsheet(collection_name, sheet_id, sheet_name, id_col, cell_range, use_header_row = True, overwrite = True):
    if use_header_row == True:
        i_use_header_row = 1
    else:
        i_use_header_row = 0
    if overwrite == True:
        i_overwrite = 1
    else:
        i_overwrite = 0

    if isinstance(id_col, str):
        id_col = int(id_col)
    ok = go_import_gsheet(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(sheet_id.encode('utf8')), 
            ctypes.c_char_p(sheet_name.encode('utf8')), 
            ctypes.c_int(id_col), 
            ctypes.c_char_p(cell_range.encode('utf8')), 
            ctypes.c_int(i_use_header_row), ctypes.c_int(i_overwrite))
    if ok == 1:
        return ''
    return error_message()

# export_gsheet - export collection objects to a GSheet
# syntax: COLLECTION FRAME GSHEET_ID GSHEET_NAME CELL_RANGE
# 
# Returns: error string
def export_gsheet(collection_name, frame_name, sheet_id, sheet_name, cell_range):
    ok = go_export_gsheet(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(frame_name.encode('utf8')), 
            ctypes.c_char_p(sheet_id.encode('utf8')), 
            ctypes.c_char_p(sheet_name.encode('utf8')), 
            ctypes.c_char_p(cell_range.encode('utf8')))
    if ok == 1:
        return ''
    return error_message()

def status(collection_name):
    ok = go_status(collection_name.encode('utf8'))
    return (ok == 1)

def list(collection_name, keys = []):
    src_keys = json.dumps(keys)
    value = go_list(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(src_keys.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    if len(value) == 0:
        return [] 
    return json.loads(value.decode()) 

def path(collection_name, key):
    value = go_path(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    return value.decode()

def check(collection_name):
    ok = go_check(ctypes.c_char_p(collection_name.encode('utf8')))
    return (ok == True)

def repair(collection_name):
    ok = go_repair(ctypes.c_char_p(collection_name.encode('utf8')))
    if ok == 1:
        return ''
    return error_message()

def attach(collection_name, key, filenames = [], semver = ''):
    if semver == '':
        semver = 'v0.0.0'
    srcFNames = json.dumps(filenames)
    if not isinstance(srcFNames, bytes):
        srcFNames = srcFNames.encode('utf8')
    ok = go_attach(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(semver.encode('utf8')), ctypes.c_char_p(srcFNames))
    if ok == 1:
        return ''
    return error_message()
    
def attachments(collection_name, key):
    value = go_attachments(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
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
    ok = go_detach(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(semver.encode('utf8')), ctypes.c_char_p(srcFNames))
    if ok == 1:
        return ''
    return error_message()

def prune(collection_name, key, filenames = [], semver = ''):
    '''Delete attachments for a specific key.  If the version semver is not provided, it will default to the current version.  Provide [] as filenames if you want to delete all attachments'''
    if semver == '':
        semver = 'v0.0.0'
    fnames = json.dumps(filenames).encode('utf8')
    ok = go_prune(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(semver.encode('utf8')), ctypes.c_char_p(fnames))
    if ok == 1:
        return ''
    return error_message()

def join(collection_name, key, obj = {}, overwrite = False):
    src = json.dumps(obj).encode('utf8')
    cOverwrite = ctypes.c_int(0)
    if overwrite == True:
        cOverwrite = ctypes.c_int(1)
    ok = go_join(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(src), cOverwrite)
    if ok == 1:
        return ''
    return error_message()

def clone(collection_name, keys, destination_name):
    src_keys = json.dumps(keys)
    ok = go_clone(ctypes.c_char_p(collection_name.encode('utf-8')), ctypes.c_char_p(src_keys.encode('utf-8')), ctypes.c_char_p(destination_name.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def clone_sample(collection_name, training_name, test_name = "", sample_size = 0):
    ok = go_clone_sample( ctypes.c_char_p(collection_name.encode('utf-8')), ctypes.c_char_p(training_name.encode('utf-8')), ctypes.c_char_p(test_name.encode('utf-8')), ctypes.c_int(sample_size))
    if ok == 1:
        return ''
    return error_message()

def grid(collection_name, keys, dot_paths):
    src_keys = json.dumps(keys)
    src_dot_paths = json.dumps(dot_paths)
    value = go_grid(ctypes.c_char_p(collection_name.encode('utf-8')), ctypes.c_char_p(src_keys.encode('utf-8')), ctypes.c_char_p(src_dot_paths.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    if value == None or value.strip() == "":
        return [], error_message()
    return json.loads(value), ''

def frame_create(collection_name, frame_name, keys, dot_paths, labels):
    src_keys = json.dumps(keys)
    src_dot_paths = json.dumps(dot_paths)
    if len(labels) == 0 and len(dot_paths) > 0:
        for item in dot_paths:
            if item.startswith("."):
                item = item[1:]
            labels.append(item)
    src_labels = json.dumps(labels)
    ok = go_frame_create(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')),
        ctypes.c_char_p(src_keys.encode('utf-8')),
        ctypes.c_char_p(src_dot_paths.encode('utf-8')),
        ctypes.c_char_p(src_labels.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()


def frame_exists(collection_name, frame_name):
    ok = go_frame_exists(ctypes.c_char_p(collection_name.encode('utf-8')),
            ctypes.c_char_p(frame_name.encode('utf-8')))
    if ok == 1:
        return True
    return False

def frame_keys(collection_name, frame_name):
    value = go_frame_keys(ctypes.c_char_p(collection_name.encode('utf-8')),
            ctypes.c_char_p(frame_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return json.loads(value)

def frame_objects(collection_name, frame_name):
    value = go_frame_objects(ctypes.c_char_p(collection_name.encode('utf-8')),
            ctypes.c_char_p(frame_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '' or len(value) == 0:
        return {}, error_message()
    return json.loads(value), ''

def frames(collection_name):
    value = go_frames(ctypes.c_char_p(collection_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '' or len(value) == 0: 
        return [] 
    return json.loads(value)

def frame_reframe(collection_name, frame_name, keys = []):
    src_keys = json.dumps(keys)
    ok = go_frame_reframe(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')),
        ctypes.c_char_p(src_keys.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def frame_refresh(collection_name, frame_name, keys = []):
    src_keys = json.dumps(keys)
    ok = go_frame_refresh(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')),
        ctypes.c_char_p(src_keys.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def frame_clear(collection_name, frame_name):
    ok = go_frame_clear(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def frame_delete(collection_name, frame_name):
    ok = go_frame_delete(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def frame_grid(collection_name, frame_name, include_headers = True):
    header_int = 0
    if include_headers == True:
        header_int = 1
    value = go_frame_grid(ctypes.c_char_p(collection_name.encode('utf-8')),
            ctypes.c_char_p(frame_name.encode('utf-8')),
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
    ok = go_sync_recieve_csv(
            ctypes.c_char_p(collection_name.encode('utf-8')), 
            ctypes.c_char_p(frame_name.encode('utf-8')), 
            ctypes.c_char_p(csv_filename.encode('utf-8')), 
            ctypes.c_int(overwrite_i))
    if ok == 1:
        return ''
    return error_message()


def sync_send_csv(collection_name, frame_name, csv_filename, overwrite = False):
    overwrite_i = 0
    if overwrite:
        overwrite_i = 1
    ok = go_sync_send_csv(
            ctypes.c_char_p(collection_name.encode('utf-8')), 
            ctypes.c_char_p(frame_name.encode('utf-8')), 
            ctypes.c_char_p(csv_filename.encode('utf-8')), 
            ctypes.c_int(overwrite_i))
    if ok == 1:
        return ''
    return error_message()


def sync_recieve_gsheet(collection_name, frame_name, gsheet_id, gsheet_name, cell_range = "A1:ZZ", overwrite = False):
    overwrite_i  = 0
    if overwrite == True:
        overwrite_i = 1
    ok = go_sync_recieve_gsheet(
            ctypes.c_char_p(collection_name.encode('utf-8')), 
            ctypes.c_char_p(frame_name.encode('utf-8')), 
            ctypes.c_char_p(gsheet_id.encode('utf-8')), 
            ctypes.c_char_p(gsheet_name.encode('utf-8')), 
            ctypes.c_char_p(cell_range.encode('utf-8')), 
            ctypes.c_int(overwrite_i))
    if ok == 1:
        return ''
    return error_message()


def sync_send_gsheet(collection_name, frame_name, gsheet_id, gsheet_name, cell_range = "A1:ZZ", overwrite = False):
    overwrite_i = 0
    if overwrite == True:
        overwrite_i = 1
    ok = go_sync_send_gsheet(
            ctypes.c_char_p(collection_name.encode('utf-8')), 
            ctypes.c_char_p(frame_name.encode('utf-8')), 
            ctypes.c_char_p(gsheet_id.encode('utf-8')), 
            ctypes.c_char_p(gsheet_name.encode('utf-8')), 
            ctypes.c_char_p(cell_range.encode('utf-8')), 
            ctypes.c_int(overwrite_i))
    if ok == 1:
        return ''
    return error_message()

def make_objects(collection_name, keys, default_object):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    keys_as_json = ctypes.c_char_p(json.dumps(keys).encode('utf8'))
    object_as_json = ctypes.c_char_p(json.dumps(default_object).encode('utf8'))
    ok = go_make_objects(c_name, keys_as_json, object_as_json)
    if ok == 1:
        return ''
    return error_message()

def update_objects(collection_name, keys, objects):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    keys_as_json = ctypes.c_char_p(json.dumps(keys).encode('utf8'))
    objects_as_json = ctypes.c_char_p(json.dumps(objects).encode('utf8'))
    ok = go_update_objects(c_name, keys_as_json, objects_as_json)
    if ok == 1:
        return ''
    return error_message()

def set_who(collection_name, names = []):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    names_as_json = ctypes.c_char_p(json.dumps(names).encode('utf8'))
    ok = go_set_who(c_name, names_as_json)
    if ok == 1:
        return ''
    return error_message()

def set_what(collection_name, src = ""):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    c_src = ctypes.c_char_p(src.encode('utf8'))
    ok = go_set_what(c_name, c_src)
    if ok == 1:
        return ''
    return error_message()

def set_when(collection_name, src = ""):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    c_src = ctypes.c_char_p(src.encode('utf8'))
    ok = go_set_when(c_name, c_src)
    if ok == 1:
        return ''
    return error_message()

def set_where(collection_name, src = ""):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    c_src = ctypes.c_char_p(src.encode('utf8'))
    ok = go_set_where(c_name, c_src)
    if ok == 1:
        return ''
    return error_message()

def set_version(collection_name, src = ""):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    c_src = ctypes.c_char_p(src.encode('utf8'))
    ok = go_set_version(c_name, c_src)
    if ok == 1:
        return ''
    return error_message()

def set_contact(collection_name, src = ""):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    c_src = ctypes.c_char_p(src.encode('utf8'))
    ok = go_set_contact(c_name, c_src)
    if ok == 1:
        return ''
    return error_message()


def get_who(collection_name):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    value = go_get_who(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    rval = value.decode()
    if type(rval) is str:
        return json.loads(rval)
    return []

def get_what(collection_name):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    value = go_get_what(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_where(collection_name):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    value = go_get_where(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_when(collection_name):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    value = go_get_when(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_version(collection_name):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    value = go_get_version(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()

def get_contact(collection_name):
    c_name = ctypes.c_char_p(collection_name.encode('utf-8'))
    value = go_get_contact(c_name)
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode()











