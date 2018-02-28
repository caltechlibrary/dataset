#!/usr/bin/env python3
import sys
import os
import shutil

print("Starting dataset_test.py")
import dataset

# Pre-test check
error_count = 0
ok = True
dataset.verbose_off()

if len(sys.argv) > 1:
    collection_name = sys.argv[1]
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    ok = dataset.init_collection(collection_name)
    if ok == False:
        print("Failed, could not create collection")
        error_count += 1
else:
    print("To run tests provide a collection name for testing,", sys.argv[0], '"test_collection.ds"')
    sys.exit(1)


# Setup a test record
key = "2488"
value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}

#
# Ready to run tests
#

# We should have an empty collection, we will create our test record.
ok = dataset.create_record(collection_name, key, value)
if ok == False:
    print("Failed, could not create record",key)
    error_count += 1

# Check to see that we have only one record
key_count = dataset.count(collection_name)
if key_count != 1:
    print("Failed, expected count to be 1, got", key_count)
    error_count += 1

# Do a minimal test to see if the record looks like it has content
keyList = dataset.keys(collection_name)
rec = dataset.read_record(collection_name, key)
for k, v in value.items():
   if not isinstance(v, list):
        if k in rec and rec[k] == v:
            print("OK, found", k, " -> ", v)
   else:
        if k == "formats" or k == "authors":
            print("OK, expected lists for", k, " -> ", v)
        else:
            print("Failed, expected", k, "with v",v)
            error_count += 1

# Test updating record
value["verified"] = True
ok = dataset.update_record(collection_name, key, value)
if ok == False:
   print("Failed, count not update record", key, value)
   error_count += 1
rec = dataset.read_record(collection_name, key)
for k, v in value.items():
   if not isinstance(v, list):
       if k in rec and rec[k] == v:
           print("OK, found", k, " -> ", v)
   else:
       if k == "formats" or k == "authors":
           print("OK, expected lists for", k, " -> ", v)
       else:
           print("Failed, expected", k, "with v",v)
           error_count += 1

# test deleting a record
ok = dataset.delete_record(collection_name, key)
if ok == False:
    print("Failed, could not delete record", key)
    error_count += 1

# Test count after delete
key_list = dataset.keys(collection_name)
cnt = dataset.count(collection_name)
if cnt != 0:
    print("Failed, expected zero records, got", cnt, key_list)
    error_count += 1

#
# Generate multiple records for collection for testing keys and extract
#
test_records = {
    "gutenberg:21489": {"title": "The Secret of the Island", "formats": ["epub","kindle", "plain text", "html"], "authors": [{"given": "Jules", "family": "Verne"}], "url": "http://www.gutenberg.org/ebooks/21489", "categories": "fiction, novel"},
    "gutenberg:2488": { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488", "categories": "fiction, novel"},
    "gutenberg:21839": { "title": "Sense and Sensibility", "formats": ["epub", "kindle", "plain text"], "authors": [{"given": "Jane", "family": "Austin"}], "url": "http://www.gutenberg.org/ebooks/21839", "categories": "fiction, novel" },
    "gutenberg:3186": {"title": "The Mysterious Stranger, and Other Stories", "formats": ["epub","kindle", "plain text", "html"], "authors": [{ "given": "Mark", "family": "Twain"}], "url": "http://www.gutenberg.org/ebooks/3186", "categories": "fiction, short story"},
    "hathi:uc1321060001561131": { "title": "A year of American travel - Narrative of personal experience", "formats": ["pdf"], "authors": [{"given": "Jessie Benton", "family": "Fremont"}], "url": "https://babel.hathitrust.org/cgi/pt?id=uc1.32106000561131;view=1up;seq=9", "categories": "non-fiction, memoir" }
}
test_record_count = len(test_records)

for k in test_records:
    v = test_records[k]
    ok = dataset.create_record(collection_name, k, v)
    if ok == False:
        print("Failed, could not add", k, "to", collection_name)
        error_count += 1

# Test keys, filtering keys and sorting keys
keys = dataset.keys(collection_name)
if len(keys) != test_record_count:
    print("Expected", test_record_count,"keys back, got", keys)
    error_count += 1

dataset.verbose_on()
filter_expr = '(eq .categories "non-fiction, memoir")'
keys = dataset.keys(collection_name, filter_expr)
if len(keys) != 1:
    print("Expected one key for", filter_expr, "got", keys)
    error_count += 1

filter_expr = '(contains .categories "novel")'
keys = dataset.keys(collection_name, filter_expr)
if len(keys) != 3:
    print("Expected three keys for", filter_expr, "got", keys)
    error_count += 1

sort_expr = '+.title'
filter_expr = '(contains .categories "novel")'
keys = dataset.keys(collection_name, filter_expr, sort_expr)
if len(keys) != 3:
    print("Expected three keys for", filter_expr, "got", keys)
    error_count += 1
i = 0
expected_keys = ["gutenberg:21839", "gutenberg:21489", "gutenberg:2488"]
for k in expected_keys:
    if i < len(keys) and keys[i] != k:
        print("Expected", k, "got", keys[i])
    i += 1


# Test extracting the family names
v = dataset.extract(collection_name, 'true', '.authors[:].family')
if not isinstance(v, list):
    print("Failed, expected a list, got", type(v), v)
    error_count += 1
    sys.exit(1)

if len(v) != 4:
    print("Failed expected list to be of length 4, got", len(v))
    error_count += 1

targets = [ "Austin", "Fremont", "Twain", "Verne" ]
for s in targets:
    if s not in v:
        print("Failed, expected to find", s, "in", v)
        error_count += 1

# Wrap up tests
if error_count > 0:
    print("Failed", error_count, "tests")
    sys.exit(1)
print("Success!")
