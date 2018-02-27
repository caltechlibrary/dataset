#!/usr/bin/env python3
import sys

print("Starting dataset_test.py")
import dataset

# Pre-test check
error_count = 0
ok = True
dataset.verbose_off()
#dataset.verbose_on() # DEBUG

if len(sys.argv) > 1:
    collection_name = sys.argv[1]
    #print("Initializing", collection_name)
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

# Test extracting the family names
v = dataset.extract(collection_name, 'true', '.authors[:].family')
if not isinstance(v, list):
    print("Failed, expected a list, got", type(v), v)
    error_count += 1
    sys.exit(1)

if len(v) != 1:
    printf("Failed expected list to be of length 1, got", len(v))
    error_count += 1
    sys.exit(1)

if "Verne" not in v:
    print("Failed, expected a list of family_names with Verne", v)
    error_count += 1
    sys.exit(1)
print("OK, extract works for true .authors[:].family", v)

# Finally test deleting a record
ok = dataset.delete_record(collection_name, key)
if ok == False:
    print("Failed, could not delete record", key)
    error_count += 1

# Test count after delete
cnt = dataset.count(collection_name)
if cnt != 0:
    print("Failed, expected zero records, got", cnt)

if error_count > 0:
    print("Failed", error_count, "tests")
    sys.exit(1)
print("Success!")
