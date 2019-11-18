#!/usr/bin/env python3
# 
# libdataet.py is a C type wrapper for our libdataset.go is a C shared.
# It is used to test our dataset functions exported from the C-Shared
# library libdataset.so, libdataset.dynlib or libdataset.dll.
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
import ctypes
import sys
import os
import json

# Figure out shared library extension
go_basename = 'lib/libdataset'
ext = '.so'
if sys.platform.startswith('win'):
    ext = '.dll'
if sys.platform.startswith('darwin'):
    ext = '.dylib'
if sys.platform.startswith('linux'):
    ext = '.so'

# Find our shared library and load it
dir_path = os.path.dirname(os.path.realpath(__file__))
lib = ctypes.cdll.LoadLibrary(os.path.join(dir_path, go_basename+ext))

# error_clear clears the error values
go_error_clear = lib.error_clear

# Setup our Go functions to be nicely wrapped
go_error_message = lib.error_message
go_error_message.restype = ctypes.c_char_p

go_use_strict_dotpath = lib.use_strict_dotpath
# Args: is 1 (true) or 0 (false)
go_use_strict_dotpath.argtypes = [ctypes.c_int]
go_use_strict_dotpath.restype = ctypes.c_int

go_dataset_version = lib.dataset_version
go_dataset_version.restype = ctypes.c_char_p

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

go_create_object = lib.create_object
# Args: collection_name (string), key (string), value (JSON source)
go_create_object.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_create_object.restype = ctypes.c_int

go_read_object = lib.read_object
# Args: collection_name (string), key (string), clean_object (int)
go_read_object.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
# Returns: value (JSON source)
go_read_object.restype = ctypes.c_char_p

# THIS IS A HACK, ctypes doesn't **easily** support undemensioned arrays
# of strings. So we will assume the array of keys has already been
# transformed into JSON before calling go_read_list.
go_read_object_list = lib.read_object_list
# Args: collection_name (string), keys (list of strings AS JSON!!!), clean_object (int)
go_read_object_list.argtypes = [ ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
# Returns: value (JSON source)
go_read_object_list.restype = ctypes.c_char_p

go_update_object = lib.update_object
# Args: collection_name (string), key (string), value (JSON sourc)
go_update_object.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_update_object.restype = ctypes.c_int

go_delete_object = lib.delete_object
# Args: collection_name (string), key (string)
go_delete_object.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_delete_object.restype = ctypes.c_int

go_key_exists = lib.key_exists
# Args: collection_name (string), key (string)
go_key_exists.argtypes = [ctypes.c_char_p,ctypes.c_char_p]
# Returns: true (1), false (0)
go_key_exists.restype = ctypes.c_int

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

# NOTE: go_sync_* diverges from cli in that it separates the functions
# specifically for CSV files and GSheets.
#
# Returns: true (1), false (0)
go_sync_recieve_csv = lib.sync_recieve_csv
go_sync_recieve_csv.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
go_sync_recieve_csv.restype = ctypes.c_int

go_sync_send_csv = lib.sync_send_csv
go_sync_send_csv.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
go_sync_send_csv.restype = ctypes.c_int

go_sync_recieve_gsheet = lib.sync_recieve_gsheet
go_sync_recieve_gsheet.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
go_sync_recieve_gsheet.restype = ctypes.c_int

go_sync_send_gsheet = lib.sync_send_gsheet
go_sync_send_gsheet.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
go_sync_send_gsheet.restype = ctypes.c_int

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
# Args: collection_name (string), key (string), semver (string), filenames (string)
go_attach.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_attach.restype = ctypes.c_int

go_attachments = lib.attachments
# Args: collection_name (string), key (string)
go_attachments.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_attachments.restype = ctypes.c_char_p

go_detach = lib.detach
# Args: collection_name (string), key (string), semver (string), basename (string)
go_detach.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_detach.restype = ctypes.c_int

go_prune = lib.prune
# Args: collection_name (string), key (string), semver (string) basename (string)
go_prune.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
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

go_frame_create = lib.frame_create
# Args: collection_name (string), frame_name (string), keys (JSON source), dotpaths (JSON source), labels (JSON source)
go_frame_create.argtypes = [ctypes.c_char_p, ctypes.c_char_p,  ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_frame_create.restype = ctypes.c_int

go_frame_exists = lib.frame_exists
# Args: collection_name (string), fame_name (string)
go_frame_exists.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_frame_exists.restype = ctypes.c_int

go_frame_keys = lib.frame_keys
# Args: collection_name (string), fame_name (string)
go_frame_keys.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_frame_keys.restype = ctypes.c_char_p

go_frame_objects = lib.frame_objects
# Args: collection_name (string), fame_name (string)
go_frame_objects.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_frame_objects.restype = ctypes.c_char_p


go_frames = lib.frames
# Args: collection_name)
go_frames.argtypes = [ctypes.c_char_p]
# Returns: frame names (JSON Array Source)
go_frames.restype = ctypes.c_char_p

go_frame_refresh = lib.frame_refresh
# Args: collection_name (string), frame_name (string), keys??? (JSON source)
go_frame_refresh.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_frame_refresh.restype = ctypes.c_int

go_frame_reframe = lib.frame_reframe
# Args: collection_name (string), frame_name (string), keys??? (JSON source)
go_frame_reframe.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p]
# Returns: value (JSON object source)
go_frame_reframe.restype = ctypes.c_int

go_frame_delete = lib.frame_delete
# Args: collection_name (string), frame_name (string)
go_frame_delete.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_frame_delete.restype = ctypes.c_int

go_frame_clear = lib.frame_clear
# Args: collection_name (string), frame_name (string)
go_frame_clear.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
# Returns: true (1), false (0)
go_frame_clear.restype = ctypes.c_int

go_frame_grid = lib.frame_grid
# Args: collection_name (string), frame_name (string), include header (int)
go_frame_grid.argtypes = [ctypes.c_char_p, ctypes.c_char_p, ctypes.c_int]
# Returns: frame names (JSON Array Source)
go_frame_grid.restype = ctypes.c_char_p

go_make_objects = lib.make_objects
# Args: collection_name (string), keys_as_json (string), object_as_json (string)
go_make_objects.argtypes = [ ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 is has errors
go_make_objects.restype = ctypes.c_int

go_update_objects = lib.update_objects
# Args: collection_name (string), keys_as_json (string), objects_as_json (string)
go_update_objects.argtypes = [ ctypes.c_char_p, ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 is has errors
go_update_objects.restype = ctypes.c_int

go_set_who = lib.set_who
# Args: collection_name (string), names_as_json (string)
go_set_who.argtypes = [ ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 has error(s)
go_set_who.restype = ctypes.c_int

go_set_what = lib.set_what
# Args: collection_name (string), what value (string)
go_set_what.argtypes = [ ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 has error(s)
go_set_what.restype = ctypes.c_int

go_set_when = lib.set_when
# Args: collection_name (string), when value (string)
go_set_when.argtypes = [ ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 has error(s)
go_set_when.restype = ctypes.c_int

go_set_where = lib.set_where
# Args: collection_name (string), where value (string)
go_set_where.argtypes = [ ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 has error(s)
go_set_where.restype = ctypes.c_int

go_set_version = lib.set_version
# Args: collection_name (string), version value (string)
go_set_version.argtypes = [ ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 has error(s)
go_set_version.restype = ctypes.c_int

go_set_contact = lib.set_contact
# Args: collection_name (string), contact value (string)
go_set_contact.argtypes = [ ctypes.c_char_p, ctypes.c_char_p ]
# Returns: int 0 is no error, 1 has error(s)
go_set_contact.restype = ctypes.c_int

go_get_who = lib.get_who
# Args: collection_name (string)
go_get_who.argtypes = [ ctypes.c_char_p ]
# Returns: frame names (JSON Array Source)
go_get_who.restype = ctypes.c_char_p

go_get_what = lib.get_what
# Args: collection_name (string)
go_get_what.argtypes = [ ctypes.c_char_p ]
# Returns: frame names (JSON Array Source)
go_get_what.restype = ctypes.c_char_p

go_get_where = lib.get_where
# Args: collection_name (string)
go_get_where.argtypes = [ ctypes.c_char_p ]
# Returns: frame names (JSON Array Source)
go_get_where.restype = ctypes.c_char_p

go_get_when = lib.get_when
# Args: collection_name (string)
go_get_when.argtypes = [ ctypes.c_char_p ]
# Returns: frame names (JSON Array Source)
go_get_when.restype = ctypes.c_char_p

go_get_version = lib.get_version
# Args: collection_name (string)
go_get_version.argtypes = [ ctypes.c_char_p ]
# Returns: frame names (JSON Array Source)
go_get_version.restype = ctypes.c_char_p

go_get_contact = lib.get_contact
# Args: collection_name (string)
go_get_contact.argtypes = [ ctypes.c_char_p ]
# Returns: frame names (JSON Array Source)
go_get_contact.restype = ctypes.c_char_p

