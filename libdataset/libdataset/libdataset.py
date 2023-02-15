#!/usr/bin/env python3
# 
# libdataet.py is a C type wrapper for our libdataset.go is a C shared.
# It is used to test our dataset functions exported from the C-Shared
# library libdataset.so, libdataset.dynlib or libdataset.dll.
# 
# @author Thomas E. (Tom) Morrell
# @author R. S. Doiel, <rsdoiel@caltech.edu>
#
# Copyright (c) 2023, Caltech
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
import sys, platform
import os
import json
from ctypes import *

# Figure out shared library extension
#
# NOTE: Assemble our library name based on CPU type and extension.
# Currently we're supporting libdataset for Windows and Linux on
# Intel only but for macOS we include both Intel and M1 support via
# the cpu types of "amd64" and "arm64"
#
lib_basename = 'libdataset'
cpu = 'amd64'
ext = '.so'
if sys.platform.startswith('win'):
    cpu = 'amd64'
    ext = '.dll'
if sys.platform.startswith('darwin'):
    #M1 mac uses a special dylib
    if platform.processor() == 'arm':
        cpu = 'arm64'
        ext = '.dylib'
    else:
        cpu = 'amd64'
        ext = '.dylib'
if sys.platform.startswith('linux'):
    cpu = 'amd64'
    ext = '.so'

# Find our shared library and load it
dir_path = os.path.realpath(os.path.join(os.path.dirname(os.path.realpath(__file__)), '..'))
# NOTE: we ignore the cpu type in libdataset testing.
lib_path = os.path.join(dir_path, f'{lib_basename}{ext}')
libdataset = CDLL(lib_path)

#
# Setup our Go/C-shared library wrapper
#

# NOTE: we use a horrible hack in this library. It is a royal pain
# to pass dynamic dataset structures between C and Python let alone
# between Go and C. As a result I've chosen the easy programing path
# of passing JSON source between the code spaces. This has proven
# simple, reliable and **INEFFICENT** in memory usage. I've opted for
# reliability and simplicity. RSD, 2020-03-18
#
# A someday feature would be to replace "libdataset.py" wrapper with
# a native Python implementation of the libdataset. Then you could
# could void to runtimes with indepentent memory management. I've just
# never had time to do this. RSD, 2023-02-15.

# dataset_version() returns the version number of the libdataset
# used.
#
# Return: semver (string)
libdataset.dataset_version.restype = c_char_p

# error_clear() clears the error values
#
# It takes no args and returns no value.
libdataset.error_clear.restype = None

# error_message() returns the error messages aggregated
# by previously envoked shared library functions. It clears the
# message aggregation as it returns the messages.
#
# Return: error message text (string)
libdataset.error_message.restype = c_char_p

# use_strict_dotpath() sets the state of the strict dotpath
# interpretation. Strict dot paths expect a period at the
# beginning, non strict will prefix a period automatigally.
# This is useful where you're using labels in a report
# (the period is annoying) and also in generating transformed
# object attributes (the period is useful).
#
# Args: is True (1) or False (0)
# Return: True (1) or False (0)
libdataset.use_strict_dotpath.argtypes = [c_int]
libdataset.use_strict_dotpath.restype = c_bool

# is_verbose() returns the state of the verbose flag.
#
# Returns: True (1) or False (0)
libdataset.is_verbose.restype = c_bool

# verbose_on() sets the verbose flag to True.
#
# Returns: True (1) or False (0)
libdataset.verbose_on.restype = c_bool

# verbose_off() sets the verbose flag to False
#
# Returns: True (1) or False (0)
libdataset.verbose_off.restype = c_bool

# init_collection() creates a dataset collection. If the
# dsn value is an empty string the collection will use pairtree
# storage implementation, othersize it will use a SQL store described
# by the dsn (data source name) to store the JSON documents.
#
# Args: collection_name (string), dsn (string)
# Returns: true (1), false (0)
libdataset.init_collection.argtypes = [ c_char_p, c_char_p ]
libdataset.init_collection.restype = c_bool

# is_collection_open() checks to see if the collection
# is already open and in the list of open collections.
#
# Args: collection_name (string)
# Returns: Ture (1) or False (0)
libdataset.is_collection_open.argtypes = [ c_char_p ]
libdataset.is_collection_open.restype = c_bool

# collections() returns a list of opened collections.
#
# Returns: string (names of the open collections as JSON)
libdataset.collections.restype = c_char_p

# open_collection() explicitly opens a collection and adds
# it to the open collection list. Returns True on success, 
# false otherwise.
#
# Args: collection_name (string)
# Returns: True (1) or False (0)
libdataset.open_collection.argtypes = [c_char_p]
libdataset.open_collection.restype = c_bool

# close_collection() closes a previously opened collection.
# It removes it from the open collections list. Returns True
# on success, False otherwise.
#
# Args: collection_name (string)
# Returns: True (1) or False (0)
libdataset.close_collection.argtypes = [c_char_p]
libdataset.close_collection.restype = c_bool

# close_all_collections closes all opened collections.
# The open collection list is cleared.
#
# Returns: True (1) or False (0)
libdataset.close_all_collections.restype = c_bool

# create_object() creates a JSON object in a collection.
#
# Args: collection_name (string), key (string), value (JSON source)
# Returns: True (1) or False (0)
libdataset.create_object.argtypes = [c_char_p, c_char_p, c_char_p]
libdataset.create_object.restype = c_bool

# read_object() retrieves a JSON object from a collection.
#
# Args: collection_name (string), key (string)
# Returns: value (JSON source)
libdataset.read_object.argtypes = [ c_char_p, c_char_p ]
libdataset.read_object.restype = c_char_p

# read_object_version retrieves an object from the collection
# using the key an semver provides.
#
# Args: collection_name (string), key (string), semver (string)
# Returns: value (JSON source)
libdataset.read_object_version.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.read_object_version.restype = c_char_p

# update_object() updates an object in the collection given a key
# and new object. If versioning is enabled the new version will use
# the incremented semver value indicated in the versioning setting.
#
# Args: collection_name (string), key (string), value (JSON sourc)
# Returns: value (JSON source)
libdataset.update_object.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.update_object.restype = c_bool

# delete_object will remove an object form a colleciton (including
# all versions of the object if versioning has been enabled).
# If you are versioning your collection normally you avoid deleting
# and replace the object tomb stone object. Returns true on 
# successful removal, false otherwise.
#
# Args: collection_name (string), key (string)
# Returns: True (1) or False (0)
libdataset.delete_object.argtypes = [ c_char_p, c_char_p ]
libdataset.delete_object.restype = c_bool

# has_key() tests for a key in a collection.
#
# Args: collection_name(string), key (string)
# Returns: (bool)
libdataset.has_key.argtypes = [ c_char_p, c_char_p ]
libdataset.has_key.restype = c_bool

# keys() returns a list of all keys in a collection.
#
# Args: collection_name (string)
# Returns: value (JSON source)
libdataset.keys.argtypes = [ c_char_p ]
libdataset.keys.restype = c_char_p

# count_objects() returns the number of objects in a collection.
#
# Args: collection_name (string)
# Returns: value (int)
libdataset.count_objects.argtypes = [ c_char_p ]
libdataset.count_objects.restype = c_int

# NOTE: import_csv, export_csv, sync_* diverges from cli and
# reflects the low level dataset organization. 

# import_csv - import a CSV file into a collection. The column
# id must be specified. Returns true on success, false if errors
# encountered.
#
# Args: collection_name (string), frame_name (string), ID column (int),
#       use header row (bool), overwrite (bool)
# Returns: True (1), False (0)
libdataset.import_csv.argtypes = [ c_char_p, c_char_p, c_int, c_bool, c_bool ]
libdataset.import_csv.restype = c_bool

# export_csv - export collection objects to a CSV file using a frame.
# 
# Args: collection_name (strng), frame_name (string), csv_filename (string)
# Returns: True (1), False (0)
libdataset.export_csv.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.export_csv.restype = c_bool

# sync_receive_csv() retrieves data from a CSV file and updates a 
# collection using a frame to map columns to attributes. Returns 
# True on success, False if error is encountered.
#
# Args: collection_name (string), frame_name (string), 
#       csv_filename (string), overwrite (bool)
# Returns: True (1), False (0)
libdataset.sync_recieve_csv.argtypes = [ c_char_p, c_char_p, c_char_p, c_int ]
libdataset.sync_recieve_csv.restype = c_bool

# sync_send_csv() updates a CSV file based on the objects in a collection
# using a frame. The frame is used to map object attributes to columns.
# Returns True on success, False if errors encountered.
#
# Args: collection_name (string), frame_name (string),
#       csv_filename (string), ovewrite (bool)
# Returns: True (1), False (0)
libdataset.sync_send_csv.argtypes = [ c_char_p, c_char_p, c_char_p, c_int ]
libdataset.sync_send_csv.restype = c_bool

# collection_exists() returns True if a collection exists, False otherwise.
#
# NOTE: collection_exists() will be renamed has_collection() in a coming 
# release of libdataset.
#
# Returns: True (1), False (0)
libdataset.collection_exists.restype = c_bool

# list_objects() returns a list of objects for a list of keys as
# JSON. Returns a JSON list of object as source.
#
# Args: collection_name (string), key list (JSON array source)
# Returns: value (JSON Array of Objects source)
libdataset.list_objects.argtypes = [ c_char_p, c_char_p ]
libdataset.list_objects.restype = c_char_p

# object_path() returns the file system path to a JSON object in a
# collection if the collection uses a pairtree store. If the collection
# uses SQL store then an empty string is returned. This is because SQL
# stored JSON objects are not directly accessible from disk (e.g. the
# MySQL/PostgreSQL JSON store maybe accessed over the network).
#
# Args: collection_name (string), key (string)
# Returns: value (string)
libdataset.object_path.argtypes = [ c_char_p, c_char_p ]
libdataset.object_path.restype = c_char_p

# check_collection() checks a collection for structural errors.
# This can be slow for larger collections. Returns True if everything
# is OK, False of there are errors.
#
# Args: collection_name (string)
# Returns: True (1), False (0)
libdataset.check_collection.argtypes = [ c_char_p ]
libdataset.check_collection.restype = c_bool

# repair_collection runs the analyzer over a collection and repairs JSON
# objects and attachment discovered having a problem. Also is
# useful for upgrading a collection between dataset releases. This
# process can be slow. Returns True on successful repair, False if
# errors encountered.
#
# Args: collection_name (string)
# Returns: true (1), false (0)
libdataset.repair_collection.argtypes = [ c_char_p ]
libdataset.repair_collection.restype = c_bool

# attach() adds a file to a JSON object record. If the collection
# has versioning set then the object will be added with an appropriate
# version number. NOTE: the object is copied in full and is not a
# delta of previous objects. This can take up allot of disk space!
# A better approach to versioning would take advantage of more effecient
# storage options (e.g. ZFS file system festures, or a version control
# system like Git or Subversion). Returns True on successful attachment
# and False if errors encountered.
#
# Args: collection_name (string), key (string), filenames (string)
# Returns: true (1), false (0)
libdataset.attach.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.attach.restype = c_bool

# attachments() returns a list the files attached to a JSON object record
# as a JSON array.
#
# Args: collection_name (string), key (string)
# Return: string (JSON list of basenames)
libdataset.attachments.argtypes = [ c_char_p, c_char_p ]
libdataset.attachments.restype = c_char_p

# detach() retrieves a file from an JSON object record copying it
# out using the basename to the current working directory. It returns
# True if the file is successfully copied out, False if an error is
# encountered.
#
# Args: collection_name (string), key (string), basename (string)
# Returns: true (1), false (0)
libdataset.detach.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.detach.restype = c_bool

# detach_version() will detatch a specific version of a file from
# a JSON object in the collection. It needs the key, semver and basename
# of the file.
#
# Args: collection_name (string), key (string), semver (string), basename (string)
# Returns: True (1), False (0)
libdataset.detach_version.argtypes = [ c_char_p, c_char_p, c_char_p, c_char_p ]
libdataset.detach_version.restype = c_bool

# prune removes an attachment from an object. NOTE: it removes all versions
# of the attachment if versioning is enabled for the collection. If you
# are using versioning and need to "remove" a file, place a tomb stone
# file there instead of using prune.
#
# Args: collection_name (string), key (string), basename (string)
# Returns: true (1), false (0)
libdataset.prune.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.prune.restype = c_bool

# join_objects() joins a new object with an existing object in a collection.
# The overwrite parameters determines if attributes are overwritten if they
# are found in the collection or skipped (if overwrite is False).
#
# Args: collection_name (string), key (string), value (JSON source), overwrite (1: true, 0: false)
# Returns: True (1), False (0)
libdataset.join_objects.argtypes = [ c_char_p, c_char_p, c_char_p, c_bool ]
libdataset.join_objects.restype = c_bool

# clone_collection() takes collection and creates a copy of it. The
# collection created will be a the type indicated by the dsn (data source
# name) value. If it is an empty string then it'll use a pairtree store
# otherwise it'll be the SQL store indicated in the dsn. Returns True
# on success, False if errors encountered.
#
# Args: collection_name (string), new_collection_name (string), dsn (string)
# Returns: True (1), False (0)
libdataset.clone_collection.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.clone_collection.restype = c_bool

# clone_sample() generates a random sample of objects split between 
# training and test collections. The dsn for training and testing collection
# will set the storage type for that specific collection. The if the
# data source names are an empty string then a pairtree store will be used.
# 
# Args: collection_name (string), training_collection_name (string), training_dsn (string), test_collection_name (string), test_dsn (string), sample size (int)
# Returns: True (1), False (0)
libdataset.clone_sample.argtypes = [ c_char_p, c_char_p, c_char_p, c_char_p, c_char_p, c_int ]
libdataset.clone_sample.restype = c_bool

# frame() returns the full metadata and contents of a frame.
#
# Args: collection_name (string), frame_name (string)
# Returns: value (JSON object source)
libdataset.frame.argtypes = [ c_char_p, c_char_p ]
libdataset.frame.restype = c_char_p

# frame_create() generates a new data frame given a collection name,
# frame name, keys, dot paths and labels. The keys, dot paths and labels
# need to be JSON encoded.
#
# Args: collection_name (string), frame_name (string), keys (JSON source), dotpaths (JSON source), labels (JSON source)
# Returns: value (JSON object source)
libdataset.frame_create.argtypes = [ c_char_p, c_char_p,  c_char_p, c_char_p, c_char_p ]
libdataset.frame_create.restype = c_bool

# has_frame() checks to see if a frame name has already been defined.
#
# Args: collection_name (string), frame_name (string)
# Returns: True (1), False (0)
libdataset.has_frame.argtypes = [ c_char_p, c_char_p ]
libdataset.has_frame.restype = c_bool

# frame_keys() returns a list of keys as JSON as defined in the
# data frame.
#
# Args: collection_name (string), frame_name (string)
# Returns: value (JSON object source)
libdataset.frame_keys.argtypes = [ c_char_p, c_char_p ]
libdataset.frame_keys.restype = c_char_p

# frame_objects() returns a list of objects as JSON currently defined
# in the data frame. The returned objects are JSON encoded.
#
# Args: collection_name (string), frame_name (string)
# Returns: value (JSON object source)
libdataset.frame_objects.argtypes = [ c_char_p, c_char_p ]
libdataset.frame_objects.restype = c_char_p

# frame_names() returns a list of frames defined for the collection
# JSON encoded.
#
# Args: collection_name (string)
# Returns: frame names (JSON Array Source)
libdataset.frame_names.argtypes = [ c_char_p ]
libdataset.frame_names.restype = c_char_p

# frame_refresh() updates the objects in a data frame based on
# the current state of the collection. Any objects removed from
# the collection will be removed from the frame.
#
# Args: collection_name (string), frame_name (string)
# Returns: value (JSON object source)
libdataset.frame_refresh.argtypes = [ c_char_p, c_char_p]
libdataset.frame_refresh.restype = c_bool

# frame_reframe() replaces the object list in a data frame. Objects
# not in the new list of keys will be removed from the frame.
#
# Args: collection_name (string), frame_name (string), keys (JSON source)
# Returns: value (JSON object source)
libdataset.frame_reframe.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.frame_reframe.restype = c_bool

# frame_delete() removes a frame from a collection. Returns True
# if delete is successful, False if there were errors.
#
# Args: collection_name (string), frame_name (string)
# Returns: True (1), False (0)
libdataset.frame_delete.argtypes = [ c_char_p, c_char_p ]
libdataset.frame_delete.restype = c_bool

# frame_clear() removes all objects from a frame leaving
# the definition in place.
#
# Args: collection_name (string), frame_name (string)
# Returns: True (1), False (0)
libdataset.frame_clear.argtypes = [ c_char_p, c_char_p ]
libdataset.frame_clear.restype = c_bool

# frame_grid returns a 2D JSON array of frame data JSON encoded.
# (depreciated, this will go away in a future version of libdataset).
#
# Args: collection_name (string), frame_name (string), include header (bool)
# Returns: frame names (JSON Array Source)
libdataset.frame_grid.argtypes = [ c_char_p, c_char_p, c_bool ]
libdataset.frame_grid.restype = c_char_p

# create_objects() generates a batch of objects in a collection,
# used for testing libdataset. (depreciated, will go away in a
# a future version of libdataset)
#
# Args: collection_name (string), keys_as_json (string), object_as_json (string)
# Returns: True (1) success, False (0) if there are errors
libdataset.create_objects.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.create_objects.restype = c_bool

# update_objects()  updates a set of objects in a collections, used in
# testing a counter part to make_objects. (depreciated, this will go
# away in a future version of libdataset)
#
# Args: collection_name (string), keys_as_json (string), objects_as_json (string)
# Returns: True (1) success, False (0) if there are errors
libdataset.update_objects.argtypes = [ c_char_p, c_char_p, c_char_p ]
libdataset.update_objects.restype = c_bool

# set_versioning() sets the versioning status for a collection. The
# accepted values are "", "patch", "minor", "major". An empty string
# disables versioning (collection are not versioned by default) the
# other strings indicate the semver value incremented on a new version.
# Returns True if successful, False if an error encountered.
#
# Args: collection_name (string), versioning_setting (string)
# Returns: True (1), False 0)
libdataset.set_versioning.argtypes = [ c_char_p, c_char_p ]
libdataset.set_versioning.restype = c_bool

# get_versioning() returns the current setting of versioning for a
# collection. The values are "" (versioning disabled), "patch",
# "minor" and "major" for semver values to be incremented on
# create, update, and attach.
#
# Args: collection_name (string)
# Returns: value (string)
libdataset.get_versioning.argtypes = [ c_char_p ]
libdataset.get_versioning.restype = c_bool
