#!/usr/bin/env python3
# 
# libdataet.py is a C type wrapper for our libdataset.go is a C shared.
# It is used to test our dataset functions exported from the C-Shared
# library libdataset.so, libdataset.dynlib or libdataset.dll.
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
import sys
import os
import json
from ctypes import *

# Figure out shared library extension
lib_basename = 'libdataset'
ext = '.so'
if sys.platform.startswith('win'):
    ext = '.dll'
if sys.platform.startswith('darwin'):
    ext = '.dylib'
if sys.platform.startswith('linux'):
    ext = '.so'

# Find our shared library and load it
dir_path = os.path.realpath(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))
lib_path = os.path.join(dir_path, lib_basename+ext)
libdataset = CDLL(lib_path)

# error_clear clears the error values
libdataset.error_clear.restype = None

# Setup our Go functions to be nicely wrapped
libdataset.error_message.restype = c_char_p

# Args: is 1 (true) or 0 (false)
libdataset.use_strict_dotpath.argtypes = [c_int]
libdataset.use_strict_dotpath.restype = c_bool

libdataset.is_verbose.restype = c_bool

libdataset.verbose_on.restype = c_bool

libdataset.verbose_off.restype = c_bool

libdataset.dataset_version.restype = c_char_p

# Args: collection_name (string)
libdataset.init_collection.argtypes = [c_char_p]
# Returns: true (1), false (0)
libdataset.init_collection.restype = c_bool

libdataset.is_collection_open.restype = c_bool

libdataset.open_collection.argtypes = [c_char_p]
libdataset.open_collection.restype = c_bool

libdataset.close_collection.argtypes = [c_char_p]
libdataset.close_collection.restype = c_bool

libdataset.close_all_collections.restype = c_bool

# Args: collection_name (string), key (string), value (JSON source)
libdataset.create_object.argtypes = [c_char_p, c_char_p, c_char_p]
libdataset.create_object.restype = c_bool

# Args: collection_name (string), key (string), clean_object (bool)
libdataset.read_object.argtypes = [c_char_p, c_char_p, c_bool ]
# Returns: value (JSON source)
libdataset.read_object.restype = c_char_p

# THIS IS A HACK, ctypes doesn't **easily** support undemensioned arrays
# of strings. So we will assume the array of keys has already been
# transformed into JSON before calling libdataset.read_list.
# Args: collection_name (string), keys (list of strings AS JSON!!!), clean_object (bool)
libdataset.read_object_list.argtypes = [ c_char_p, c_char_p, c_bool ]
# Returns: value (JSON source)
libdataset.read_object_list.restype = c_char_p

# Args: collection_name (string), key (string), value (JSON sourc)
libdataset.update_object.argtypes = [c_char_p, c_char_p, c_char_p ]
libdataset.update_object.restype = c_bool

# Args: collection_name (string), key (string)
libdataset.delete_object.argtypes = [c_char_p, c_char_p ]
libdataset.delete_object.restype = c_bool

# Args: collection_name (string), key (string)
libdataset.key_exists.argtypes = [c_char_p,c_char_p ]
libdataset.key_exists.restype = c_bool

# Args: collection_name (string)
libdataset.keys.argtypes = [c_char_p ]
# Returns: value (JSON source)
libdataset.keys.restype = c_char_p

# Args: collection_name (string), key_list (JSON array source), filter_expr (string)
libdataset.key_filter.argtypes = [c_char_p, c_char_p, c_char_p ]
# Returns: value (JSON source)
libdataset.key_filter.restype = c_char_p

# Args: collection_name (string), key_list (JSON array source), sort order (string)
libdataset.key_sort.argtypes = [c_char_p, c_char_p, c_char_p ]
# Returns: value (JSON source)
libdataset.key_sort.restype = c_char_p

# Args: collection_name (string)
libdataset.count_objects.argtypes = [ c_char_p ]
# Returns: value (int)
libdataset.count_objects.restype = c_int

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
# Args: collection_name (string), frame_name (string), ID column (int), use header row (bool), overwrite (bool)
libdataset.import_csv.argtypes = [ c_char_p, c_char_p, c_int, c_bool, c_bool ]
# Returns: true (1), false (0)
libdataset.import_csv.restype = c_bool

# NOTE: this diverges from cli and uses libdataset.go bindings
#
# export_csv - export collection objects to a CSV file
# syntax examples: COLLECTION FRAME CSV_FILENAME
# 
# Returns: true (1), false (0)
# Args: collection_name (strng), frame_name (string), csv_filename (string)
libdataset.export_csv.argtypes = [ c_char_p, c_char_p, c_char_p ]
# Returns: True on success, False otherwise
libdataset.export_csv.restype = c_bool


# NOTE: libdataset.sync_* diverges from cli and uses libdataset.go bindings
#
# Args: 
libdataset.sync_recieve_csv.argtypes = [ c_char_p, c_char_p, c_char_p, c_int ]
# Returns: true (1), false (0)
libdataset.sync_recieve_csv.restype = c_bool
# Args: 
libdataset.sync_send_csv.argtypes = [ c_char_p, c_char_p, c_char_p, c_int ]
# Returns: true (1), false (0)
libdataset.sync_send_csv.restype = c_bool

# Returns: true (1), false (0)
libdataset.collection_exists.restype = c_bool

# Args: collection_name (string), key list (JSON array source)
libdataset.list_objects.argtypes = [ c_char_p, c_char_p ]
# Returns: value (JSON Array of Objects source)
libdataset.list_objects.restype = c_char_p

# FIXME: for Python library only accept single return a single key's path
# Args: collection_name (string), key (string)
libdataset.object_path.argtypes = [ c_char_p, c_char_p ]
# Returns: value (string)
libdataset.object_path.restype = c_char_p

# Args: collection_name (string)
libdataset.check_collection.argtypes = [ c_char_p ]
# Returns: true (1), false (0)
libdataset.check_collection.restype = c_bool

# Args: collection_name (string)
libdataset.repair_collection.argtypes = [ c_char_p ]
# Returns: true (1), false (0)
libdataset.repair_collection.restype = c_bool

# Args: collection_name (string), key (string), semver (string), filenames (string)
libdataset.attach.argtypes = [ c_char_p, c_char_p, c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.attach.restype = c_bool

# Args: collection_name (string), key (string)
libdataset.attachments.argtypes = [ c_char_p, c_char_p ]
libdataset.attachments.restype = c_char_p

# Args: collection_name (string), key (string), semver (string), basename (string)
libdataset.detach.argtypes = [ c_char_p, c_char_p, c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.detach.restype = c_bool

# Args: collection_name (string), key (string), semver (string) basename (string)
libdataset.prune.argtypes = [ c_char_p, c_char_p, c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.prune.restype = c_bool

# Args: collection_name (string), key (string), value (JSON source), overwrite (1: true, 0: false)
libdataset.join_objects.argtypes = [ c_char_p, c_char_p, c_char_p, c_bool ]
# Returns: true (1), false (0)
libdataset.join_objects.restype = c_bool

# Args: collection_name (string), new_collection_name (string), ????
libdataset.clone_collection.argtypes = [ c_char_p, c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.clone_collection.restype = c_bool

# Args: collection_name (string), new_sample_collection_name (string), new_rest_collection_name (string), sample size (int)
libdataset.clone_sample.argtypes = [ c_char_p, c_char_p, c_char_p, c_int ]
# Returns: true (1), false (0)
libdataset.clone_sample.restype = c_bool

# Args: collection_name (string), frame_name (string), keys (JSON source), dotpaths (JSON source), labels (JSON source)
libdataset.frame_create.argtypes = [ c_char_p, c_char_p,  c_char_p, c_char_p, c_char_p ]
# Returns: value (JSON object source)
libdataset.frame_create.restype = c_bool

# Args: collection_name (string), frame_name (string)
libdataset.frame_exists.argtypes = [ c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.frame_exists.restype = c_bool

# Args: collection_name (string), frame_name (string)
libdataset.frame_keys.argtypes = [ c_char_p, c_char_p ]
# Returns: value (JSON object source)
libdataset.frame_keys.restype = c_char_p

# Args: collection_name (string), frame_name (string)
libdataset.frame.argtypes = [ c_char_p, c_char_p ]
# Returns: value (JSON object source)
libdataset.frame.restype = c_char_p

# Args: collection_name (string), frame_name (string)
libdataset.frame_objects.argtypes = [ c_char_p, c_char_p ]
# Returns: value (JSON object source)
libdataset.frame_objects.restype = c_char_p

# Args: collection_name)
libdataset.frames.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.frames.restype = c_char_p

# Args: collection_name (string), frame_name (string)
libdataset.frame_refresh.argtypes = [ c_char_p, c_char_p]
# Returns: value (JSON object source)
libdataset.frame_refresh.restype = c_bool

# Args: collection_name (string), frame_name (string), keys (JSON source)
libdataset.frame_reframe.argtypes = [ c_char_p, c_char_p, c_char_p ]
# Returns: value (JSON object source)
libdataset.frame_reframe.restype = c_bool

# Args: collection_name (string), frame_name (string)
libdataset.frame_delete.argtypes = [ c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.frame_delete.restype = c_bool

# Args: collection_name (string), frame_name (string)
libdataset.frame_clear.argtypes = [ c_char_p, c_char_p ]
# Returns: true (1), false (0)
libdataset.frame_clear.restype = c_bool

# Args: collection_name (string), frame_name (string), include header (bool)
libdataset.frame_grid.argtypes = [ c_char_p, c_char_p, c_bool ]
# Returns: frame names (JSON Array Source)
libdataset.frame_grid.restype = c_char_p

# Args: collection_name (string), keys_as_json (string), object_as_json (string)
libdataset.create_objects.argtypes = [ c_char_p, c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.create_objects.restype = c_bool

# Args: collection_name (string), keys_as_json (string), objects_as_json (string)
libdataset.update_objects.argtypes = [ c_char_p, c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.update_objects.restype = c_bool

# Args: collection_name (string), name (string)
libdataset.set_who.argtypes = [ c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.set_who.restype = c_bool

# Args: collection_name (string), what value (string)
libdataset.set_what.argtypes = [ c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.set_what.restype = c_bool

# Args: collection_name (string), when value (string)
libdataset.set_when.argtypes = [ c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.set_when.restype = c_bool

# Args: collection_name (string), where value (string)
libdataset.set_where.argtypes = [ c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.set_where.restype = c_bool

# Args: collection_name (string), version value (string)
libdataset.set_version.argtypes = [ c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.set_version.restype = c_bool

# Args: collection_name (string), contact value (string)
libdataset.set_contact.argtypes = [ c_char_p, c_char_p ]
# Returns: True (1) success, False (0) if there are errors
libdataset.set_contact.restype = c_bool

# Args: collection_name (string)
libdataset.get_who.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.get_who.restype = c_char_p

# Args: collection_name (string)
libdataset.get_what.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.get_what.restype = c_char_p

# Args: collection_name (string)
libdataset.get_where.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.get_where.restype = c_char_p

# Args: collection_name (string)
libdataset.get_when.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.get_when.restype = c_char_p

# Args: collection_name (string)
libdataset.get_version.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.get_version.restype = c_char_p

# Args: collection_name (string)
libdataset.get_contact.argtypes = [ c_char_p ]
# Returns: frame names (JSON Array Source)
libdataset.get_contact.restype = c_char_p
