#!/usr/bin/env python3
# 
# py/dataset.go is a C shared library targetting support in Python for dataset
# 
# @author R. S. Doiel, <rsdoiel@library.caltech.edu>
#
# Copyright (c) 2018, Caltech
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
go_error_message = lib.error_message
go_error_message.restype = ctypes.c_char_p

go_use_strict_dotpath = lib.use_strict_dotpath
# Args: is 1 (true) or 0 (false)
go_use_strict_dotpath.argtypes = [ctypes.c_int]
go_use_strict_dotpath.restype = ctypes.c_int

go_version = lib.dataset_version
go_version.restype = ctypes.c_char_p

go_is_verbose = lib.is_verbose
go_is_verbose.restype = ctypes.c_int

go_verbose_on = lib.verbose_on
go_verbose_on.restype = ctypes.c_int

go_verbose_off = lib.verbose_off
go_verbose_off.restype = ctypes.c_int

go_init = lib.init_collection
# Args: collection_name (string), layout (int - 0 UNKNOWN, 1 BUCKETS, 2 PAIRTREE)
go_init.argtypes = [ctypes.c_char_p, ctypes.c_int]
# Returns: true (1), false (0)
go_init.restype = ctypes.c_int

go_create_record = lib.create_record
# Args: collection_name (string), key (string), value (JSON source)
go_create_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_create_record.restype = ctypes.c_int

go_read_record = lib.read_record
# Args: collection_name (string), key (string)
go_read_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON source)
go_read_record.restype = ctypes.c_char_p

go_update_record = lib.update_record
# Args: collection_name (string), key (string), value (JSON sourc)
go_update_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_update_record.restype = ctypes.c_int

go_delete_record = lib.delete_record
# Args: collection_name (string), key (string)
go_delete_record.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_delete_record.restype = ctypes.c_int

go_has_key = lib.has_key
# Args: collection_name (string), key (string)
go_has_key.argtypes = [ctypes.c_char_p,ctypes.c_char_p]
# Returns: true (1), false (0)
go_has_key.restype = ctypes.c_int

go_keys = lib.keys
# Args: collection_name (string), filter_expr (string), sort_expr (string)
go_keys.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON source)
go_keys.restype = ctypes.c_char_p

go_key_filter = lib.key_filter
# Args: collection_name (string), key_list (JSON array source), filter_expr (string)
go_key_filter.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON source)
go_key_filter.restype = ctypes.c_char_p

go_key_sort = lib.key_sort
# Args: collection_name (string), key_list (JSON array source), sort order (string)
go_key_sort.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON source)
go_key_sort.restype = ctypes.c_char_p

go_count = lib.count
# Args: collection_name (string)
go_count.argtypes = [ctypes.c_char_p]
# Returns: value (int)
go_count.restype = ctypes.c_int

go_indexer = lib.indexer
# Args: collection_name (string),  index_name (string), index_map_name (string), key_list (JSON array source), batch_size (int)
go_indexer.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
# Returns: true (1), false (0)
go_indexer.restype = ctypes.c_int

go_deindexer = lib.deindexer
# Args: collection_name (string), index_name (string), key_list (JSON array source), batch_size (int)
go_deindexer.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
# Returns: true (1), false (0)
go_deindexer.restype = ctypes.c_int

go_find = lib.find
# Args: index_names (JSON array source), query_string (string), options (JSON object source)
go_find.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON array source)
go_find.restype = ctypes.c_char_p

# NOTE: this diverges from cli and reflects low level dataset organization
#
# import_csv - import a CSV file into a collection
# syntax: COLLECTION CSV_FILENAME ID_COL
# 
# options that should support sensible defaults:
#
#      UseHeaderRow (bool, 1 true, 0 false)
#      Overwrite (bool, 1 true, 0 false)
# 
# Returns: true (1), false (0)
go_import_csv = lib.import_csv
go_import_csv.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.c_int, ctypes.c_int]
go_import_csv.restype = ctypes.c_int

# NOTE: this diverges from cli and uses libdataset.go bindings
#
# export_csv - export collection objects to a CSV file
# syntax examples: COLLECTION FRAME CSV_FILENAME
# 
# Returns: true (1), false (0)
go_export_csv = lib.export_csv
go_export_csv.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_export_csv.restype = ctypes.c_int

# NOTE: this diverges from the cli and uses libdataset.go bindings
# import_gsheet - import a GSheet into a collection
# syntax: COLLECTION GSHEET_ID SHEET_NAME ID_COL CELL_RANGE
# 
# options that should support sensible defaults:
#
#      UseHeaderRow
#      Overwrite
#
# Returns: true (1), false (0)
go_import_gsheet = lib.import_gsheet
go_import_gsheet.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int, ctypes.c_char_p, ctypes.c_int, ctypes.c_int]
go_import_gsheet.restype = ctypes.c_int

# NOTE: this diverges from the cli and uses the libdataset.go bindings
# export_gsheet - export collection objects to a GSheet
# syntax examples: COLLECTION FRAME GSHEET_ID GSHEET_NAME CELL_RANGE
#
# Returns: true (1), false (0)
go_export_gsheet = lib.export_gsheet
go_export_gsheet.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
go_export_gsheet.restype = ctypes.c_int

go_status = lib.status
# Returns: true (1), false (0)
go_status.restype = ctypes.c_int

go_list = lib.list
# Args: collection_name (string), key list (JSON array source)
go_list.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON Array of Objects source)
go_list.restype = ctypes.c_char_p

# FIXME: for Python library only accept single return a single key's path
go_path = lib.path
# Args: collection_name (string), key (string)
go_path.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Retrusn: value (string)
go_path.restype = ctypes.c_char_p

go_check = lib.check
# Args: collection_name (string)
go_check.argtypes = [ctypes.c_char_p]
# Returns: true (1), false (0)
go_check.restype = ctypes.c_int

go_repair = lib.repair
# Args: collection_name (string)
go_repair.argtypes = [ctypes.c_char_p]
# Returns: true (1), false (0)
go_repair.restype = ctypes.c_int

go_attach = lib.attach
# Args: collection_name (string), key (string), filepath (string)
go_attach.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_attach.restype = ctypes.c_int

go_attachments = lib.attachments
# Args: collection_name (string), key (string)
go_attachments.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_attachments.restype = ctypes.c_char_p

go_detach = lib.detach
# Args: collection_name (string), key (string), filepaths (string)
go_detach.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_detach.restype = ctypes.c_int

go_prune = lib.prune
# Args: collection_name (string), key (string), filepaths (string)
go_prune.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_prune.restype = ctypes.c_int

go_join = lib.join
# Args: collection_name (string), key (string), value (JSON source), overwrite (1: true, 0: false)
go_join.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
# Returns: true (1), false (0)
go_join.restype = ctypes.c_int

go_clone = lib.clone
# Args: collection_name (string), new_collection_name (string), ????
go_clone.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_clone.restype = ctypes.c_int

go_clone_sample = lib.clone_sample
# Args: collection_name (string), new_sample_collection_name (string), new_rest_collection_name (string), sample size ????
go_clone_sample.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int ]
# Returns: true (1), false (0)
go_clone_sample.restype = ctypes.c_int

go_grid = lib.grid
# Args: collection_name (string), keys??? (JSON source), dotpaths???? (JSON source)
go_grid.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON 2D array source)
go_grid.restype = ctypes.c_char_p

go_frame = lib.frame
# Args: collection_name (string), frame_name (string), keys??? (JSON source), dotpaths???? (JSON source)
go_frame.argtypes = [ctypes.c_char_p, ctypes.c_char_p,  ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_frame.restype = ctypes.c_char_p

go_has_frame = lib.has_frame
# Args: collection_name (string), fame_name (string)
go_has_frame.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_has_frame.restype = ctypes.c_int

go_frames = lib.frames
# Args: collection_name)
go_frames.argtypes = [ctypes.c_char_p]
# Returns: frame names (JSON Array Source)
go_frames.restype = ctypes.c_char_p

# FIXME: This changed in cli...
go_reframe = lib.reframe
# Args: collection_name (string), frame_name (string), keys??? (JSON source)
go_reframe.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_reframe.restype = ctypes.c_int

go_frame_labels = lib.frame_labels
# Args: collection_name (string), frame_name (string), label_values (JSON array of string source)
go_frame_labels.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_frame_labels.restype = ctypes.c_int


go_delete_frame = lib.delete_frame
# Args: collection_name (string), frame_name (string)
go_delete_frame.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_delete_frame.restype = ctypes.c_int


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

# Returns version of dataset shared library
def version():
    value = go_version()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')

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

# Returns version of dataset shared library
def version():
    value = go_version()
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    return value.decode() 

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
def has_key(collection_name, key):
    ok = go_has_key(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(key.encode('utf8')))
    return (ok == 1)

# Create a JSON record in a Dataset Collectin
def create(collection_name, key, value):
    '''create a new JSON record in the collection based on collection name, record key and JSON string, returns True/False'''
    if isinstance(key, str) == False:
        key = f"{key}"
    ok = go_create_record(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(key.encode('utf8')), 
            ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if ok == 1:
        return ''
    return error_message()
    
# Read a JSON record from a Dataset collection
def read(collection_name, key):
    '''read a JSON record from a collection with the given name and record key, returns a dict and an error string'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    value = go_read_record(ctypes.c_char_p(collection_name.encode('utf8')), 
            ctypes.c_char_p(key.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    rval = value.decode()
    if type(rval) is str:
        if rval == "":
            return {}, error_message()
        return json.loads(rval), ''
    return {}, f"Can't read {key} from {collection_name}, {error_message()}"
    

# Update a JSON record from a Dataset collection
def update(collection_name, key, value):
    '''update a JSON record from a collection with the given name, record key, JSON string returning True/False'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    ok = go_update_record(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(json.dumps(value).encode('utf8')))
    if ok == 1:
        return ''
    return error_message()

# Delete a JSON record from a Dataset collection
def delete(collection_name, key):
    '''delete a JSON record (and any attachments) from a collection with the collectin name and record key, returning True/False'''
    if not isinstance(key, str) == True:
        key = f"{key}"
    ok = go_delete_record(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')))
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


# Indexer takes a collection name, an index name, an index map file name, and an optional keylist 
# and creates/updates a Bleve index on disc.
def indexer(collection_name, index_name, index_map_name, key_list = [], batch_size = 0):
    '''indexes a collection given a collection name, bleve index name, index map filename, and optional key list'''
    key_list_src = json.dumps(key_list)
    ok = go_indexer(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(index_name.encode('utf8')), ctypes.c_char_p(index_map_name.encode('utf8')), ctypes.c_char_p(key_list_src.encode('utf8')), ctypes.c_int(batch_size))
    if ok == 1:
        return ''
    return error_message()

# Deindexer takes a collection name, an index name, key list and optional batch size deleting the provided keys from 
# the index.
def deindexer(collection_name, index_name, key_list, batch_size = 0):
    '''indexes a collection given a collection name, bleve index name, index map filename, and optional key list'''
    key_list_src = json.dumps(key_list).encode('utf8')
    ok = go_deindexer(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(index_name.encode('utf8')), ctypes.c_char_p(key_list_src), ctypes.c_int(batch_size))
    if ok == 1:
        return ''
    return error_message()

# Find takes an index name, query string an optional options dict and returns a search result
def find(index_names, query_string, options = {}):
    '''Find takes an index name, query string an optional options dict and returns a search result'''
    option_src = json.dumps(options)
    err = ''
    value = go_find(ctypes.c_char_p(index_names.encode('utf8')), ctypes.c_char_p(query_string.encode('utf8')), ctypes.c_char_p(option_src.encode('utf8')))
    if not isinstance(value, bytes):
        value = value.encode('utf8')
    rval = value.decode()
    err = error_message()
    if rval == "":
        return {}, err
    return json.loads(rval), err


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
            ctyles.c_int(i_overwrite))
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

def attach(collection_name, key, filenames = []):
    srcFNames = json.dumps(filenames).encode('utf8')
    ok = go_attach(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(srcFNames))
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

def detach(collection_name, key, filenames = []):
    fnames = json.dumps(filenames).encode('utf8')
    ok = go_detach(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(fnames))
    if ok == 1:
        return ''
    return error_message()

def prune(collection_name, key, filenames = []):
    fnames = json.dumps(filenames).encode('utf8')
    ok = go_prune(ctypes.c_char_p(collection_name.encode('utf8')), ctypes.c_char_p(key.encode('utf8')), ctypes.c_char_p(fnames))
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

def frame(collection_name, frame_name, keys = [], dot_paths = []):
    src_keys = json.dumps(keys)
    src_dot_paths = json.dumps(dot_paths)
    value = go_frame(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')),
        ctypes.c_char_p(src_keys.encode('utf-8')),
        ctypes.c_char_p(src_dot_paths.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '':
        return [], error_message()
    return json.loads(value), ''

def has_frame(collection_name, frame_name):
    ok = go_has_frame(ctypes.c_char_p(collection_name.encode('utf-8')),
            ctypes.c_char_p(frame_name.encode('utf-8')))
    if ok == 1:
        return True
    return False

def frames(collection_name):
    value = go_frames(ctypes.c_char_p(collection_name.encode('utf-8')))
    if not isinstance(value, bytes):
        value = value.encode('utf-8')
    if value == None or value.strip() == '':
        return []
    return json.loads(value)

def reframe(collection_name, frame_name, keys = []):
    src_keys = json.dumps(keys)
    ok = go_reframe(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')),
        ctypes.c_char_p(src_keys.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def frame_labels(collection_name, frame_name, labels):
    src_labels = json.dumps(labels)
    ok = go_frame_labels(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')),
        ctypes.c_char_p(src_labels.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

def delete_frame(collection_name, frame_name):
    ok = go_delete_frame(ctypes.c_char_p(collection_name.encode('utf-8')),
        ctypes.c_char_p(frame_name.encode('utf-8')))
    if ok == 1:
        return ''
    return error_message()

