#!/usr/bin/env python3
import sys
import os
import shutil
import json
import dataset
import csv

#
# test_basic(collection_name) runs tests on basic CRUD ops
# 
def test_basic(t, collection_name):
    '''test_basic(collection_name) runs tests on basic CRUD ops'''
    # Setup a test record
    key = "2488"
    value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}
    
    # We should have an empty collection, we will create our test record.
    err = dataset.create(collection_name, key, value)
    if err != '':
        t.error(f"Failed, could not create record {key}")
    
    # Check to see that we have only one record
    key_count = dataset.count(collection_name)
    if key_count != 1:
        t.error(f"Failed, expected count to be 1, got {key_count}")
    
    # Do a minimal test to see if the record looks like it has content
    keyList = dataset.keys(collection_name)
    rec, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    for k, v in value.items():
       if not isinstance(v, list):
            if k in rec and rec[k] == v:
                t.print("OK, found", k, " -> ", v)
            else:
                t.error(f"epxected {rec[k]} got {v}")
       else:
            if k == "formats" or k == "authors":
                t.print("OK, expected lists for", k, " -> ", v)
            else:
                t.error(f"Failed, expected {k} with list v, got {v}")
    
    # Test updating record
    value["verified"] = True
    err = dataset.update(collection_name, key, value)
    if err != '':
       t.error(f"Failed, count not update record {key}, {value}, {err}")
    rec, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    for k, v in value.items():
       if not isinstance(v, list):
           if k in rec and rec[k] == v:
               t.print("OK, found", k, " -> ", v)
           else:
               t.error("expected {rec[k]} got {v} for key {k}")
       else:
           if k == "formats" or k == "authors":
               t.print("OK, expected lists for", k, " -> ", v)
           else:
               t.error("Failed, expected {k} with a list for v, got {v}")
    
    # Test path to record
    expected_s = "/".join([collection_name, "pairtree", "24", "88", (key+".json")])
    expected_l = len(expected_s)
    p = dataset.path(collection_name, key)
    if len(p) != expected_l:
        t.error("Failed, expected length", expected_l, "got", len(p))
    if p != expected_s:
        t.error("Failed, expected", expected_s, "got", p)

    # Test listing records
    l = dataset.list(collection_name, [key])
    if len(l) != 1:
        t.error("Failed, list should return an array of one record, got", l)
        return

    # test deleting a record
    err = dataset.delete(collection_name, key)
    if err != '':
        t.error("Failed, could not delete record", key, ", ", err)
    

#
# test_keys(t, collection_name) test getting, filter and sorting keys
#
def test_keys(t, collection_name):
    '''test_keys(collection_name) test getting, filter and sorting keys'''
    # Test count after delete
    key_list = dataset.keys(collection_name)
    cnt = dataset.count(collection_name)
    if cnt != 0:
        t.error("Failed, expected zero records, got", cnt, key_list)
    
    #
    # Generate multiple records for collection for testing keys
    #
    test_records = {
        "gutenberg:21489": {"title": "The Secret of the Island", "formats": ["epub","kindle", "plain text", "html"], "authors": [{"given": "Jules", "family": "Verne"}], "url": "http://www.gutenberg.org/ebooks/21489", "categories": "fiction, novel"},
        "gutenberg:2488": { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488", "categories": "fiction, novel"},
        "gutenberg:21839": { "title": "Sense and Sensibility", "formats": ["epub", "kindle", "plain text"], "authors": [{"given": "Jane", "family": "Austin"}], "url": "http://www.gutenberg.org/ebooks/21839", "categories": "fiction, novel" },
        "gutenberg:3186": {"title": "The Mysterious Stranger, and Other Stories", "formats": ["epub","kindle", "plain text", "html"], "authors": [{ "given": "Mark", "family": "Twain"}], "url": "http://www.gutenberg.org/ebooks/3186", "categories": "fiction, short story"},
        "hathi:uc1321060001561131": { "title": "A year of American travel - Narrative of personal experience", "formats": ["pdf"], "authors": [{"given": "Jessie Benton", "family": "Fremont"}], "url": "https://babel.hathitrust.org/cgi/pt?id=uc1.32106000561131;view=1up;seq=9", "categories": "non-fiction, memoir" }
    }
    test_count = len(test_records)
    
    for k in test_records:
        v = test_records[k]
        err = dataset.create(collection_name, k, v)
        if err != '':
            t.error("Failed, could not add", k, "to", collection_name, ', ', err)
    
    # Test keys, filtering keys and sorting keys
    keys = dataset.keys(collection_name)
    if len(keys) != test_count:
        t.error("Expected", test_count,"keys back, got", keys)
    
    #dataset.verbose_on()
    filter_expr = '(eq .categories "non-fiction, memoir")'
    keys = dataset.keys(collection_name, filter_expr)
    if len(keys) != 1:
        t.error("Expected one key for", filter_expr, "got", keys)
    
    filter_expr = '(contains .categories "novel")'
    keys = dataset.keys(collection_name, filter_expr)
    if len(keys) != 3:
        t.error("Expected three keys for", filter_expr, "got", keys)
    
    sort_expr = '+.title'
    filter_expr = '(contains .categories "novel")'
    keys = dataset.keys(collection_name, filter_expr, sort_expr)
    if len(keys) != 3:
        t.error("Expected three keys for", filter_expr, "got", keys)
    i = 0
    expected_keys = ["gutenberg:21839", "gutenberg:21489", "gutenberg:2488"]
    for k in expected_keys:
        if i < len(keys) and keys[i] != k:
            t.error("Expected", k, "got", keys[i])
        i += 1
    

#
# test_search(t, collection_name, index_map_name, index_name) tests indexer, deindexer and find funcitons
#
def test_search(t, collection_name, index_map_name, index_name):
    '''test indexer, deindexer and find functions'''
    #dataset.verbose_on()
    if os.path.exists(index_name):
        shutil.rmtree(index_name)
    if os.path.exists(index_map_name):
        os.remove(index_map_name)
    
    index_map = {
        "title": {
            "object_path": ".title"
        },
        "family": {
            "object_path": ".authors[:].family"
        },
        "categories": {
            "object_path": ".categories"
        }
    }
    with open(index_map_name, 'w') as outfile:
         json.dump(index_map, outfile, indent = 4, ensure_ascii = False)
    
    err = dataset.indexer(collection_name, index_name, index_map_name, batch_size = 2)
    if err != '':
        t.error("Failed to index", collection_name, ', ', err)
    results, err = dataset.find(index_name, '+family:"Verne"')
    if err != '':
        t.error("Find failed for ", collection_name + ", "+err)
        return
    if results["total_hits"] != 2: 
        t.print("Warning: unexpected results", json.dumps(results, indent = 4))
    hits = results["hits"]
    
    k1 = hits[0]["id"]
    err = dataset.deindexer(collection_name, index_name, [k1])
    if err != '':
        t.print("deindexer failed for key", k1, ', ', err)

#
# test_issue32() make sure issue 32 stays fixed.
#
def test_issue32(t, collection_name):
    err = dataset.create(collection_name, "k1", {"one":1})
    if err != '':
        t.error("Failed to create k1 in", collection_name, ', ', err)
        return
    ok = dataset.has_key(collection_name, "k1")
    if ok == False:
        t.error("Failed, has_key k1 should return", True)
    ok = dataset.has_key(collection_name, "k2")
    if ok == True:
        t.error("Failed, has_key k2 should return", False)

# Setup our test collection, recreate it if necessary
def test_setup(t, collection_name):
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    err = dataset.init(collection_name, "pairtree")
    if err != '':
        t.error("Failed, could not create collection, ", err)
        return


def test_check_repair(t, collection_name):
    t.print("Testing status on", collection_name)
    # Make sure we have a left over collection to check and repair
    if os.path.exists(collection_name) == True:
        shutil.rmtree(collection_name)
    dataset.init(collection_name, "pairtree")
    ok = dataset.status(collection_name)
    if ok == False:
        t.error("Failed, expected dataset.status() == True, got", ok, "for", collection_name)
        return

    if dataset.has_key(collection_name, 'one') == False:
        dataset.create(collection_name, 'one', {"one": 1})
    t.print("Testing check on", collection_name)
    # Check our collection
    ok = dataset.check(collection_name)
    if ok == False:
        t.error("Failed, expected check", collection_name, "to return True, got", ok)
        return

    # Break and recheck our collection
    print(f"Removing {collection_name}/collection.json to cause a fail")
    if os.path.exists(collection_name + "/collection.json"):
        os.remove(collection_name + "/collection.json")
    print(f"Testing check on (broken) {collection_name}")
    ok = dataset.check(collection_name)
    if ok == True:
        t.error("Failed, expected check", collection_name, "to return False, got", ok)
    else:
        t.print(f"Should have see error output for broken {collection_name}")

    # Repair our collection
    t.print("Testing repair on", collection_name)
    err = dataset.repair(collection_name)
    if err != '':
        t.error("Failed, expected repair to return True, got, ", err)
    if os.path.exists(collection_name + "/collection.json") == False:
        t.error(f"Failed, expected recreated {collection_name}/collection.json")
 
        
def test_attachments(t, collection_name):
    t.print("Testing attach, attachments, detach and prune")
    # Generate two files to attach.
    with open('a1.txt', 'w') as text_file:
        text_file.write('This is file a1')
    with open('a2.txt', 'w') as text_file:
        text_file.write('This is file a2')
    filenames = ['a1.txt','a2.txt']

    ok = dataset.status(collection_name)
    if ok == False:
        t.error("Failed,", collection_name, "missing")
        return
    keys = dataset.keys(collection_name)
    if len(keys) < 1:
        t.error("Failed,", collection_name, "should have keys")
        return

    key = keys[0]
    err = dataset.attach(collection_name, key, filenames)
    if err != '':
        t.error("Failed, to attach files for", collection_name, key, filenames, ', ', err)
        return

    l = dataset.attachments(collection_name, key)
    if len(l) != 2:
        t.error("Failed, expected two attachments for", collection_name, key, "got", l)
        return

    if os.path.exists(filenames[0]):
        os.remove(filenames[0])
    if os.path.exists(filenames[1]):
        os.remove(filenames[1])

    # First try detaching one file.
    err = dataset.detach(collection_name, key, [filenames[1]])
    if err != '':
        t.error("Failed, expected True for", collection_name, key, filenames[1], ', ', err)
    if os.path.exists(filenames[1]):
        os.remove(filenames[1])
    else:
        t.error("Failed to detch", filenames[1], "from", collection_name, key)

    # Test explicit filenames detch
    err = dataset.detach(collection_name, key, filenames)
    if err != '':
        t.error("Failed, expected True for", collection_name, key, filenames, ', ', err)

    for fname in filenames:
        if os.path.exists(fname):
            os.remove(fname)
        else:
            t.error("Failed, expected", fname, "to be detached from", collection_name, key)

    # Test detaching all files
    err = dataset.detach(collection_name, key)
    if err != '':
        t.error("Failed, expected True for (detaching all)", collection_name, key, ', ', err)
    for fname in filenames:
        if os.path.exists(fname):
            os.remove(fname)
        else:
            t.error("Failed, expected", fname, "for detaching all from", collection_name, key)

    err = dataset.prune(collection_name, key, [filenames[0]])
    if err != '':
        t.error("Failed, expected True for prune", collection_name, key, [filenames[0]], ', ', err)
    l = dataset.attachments(collection_name, key)
    if len(l) != 1:
        t.error("Failed, expected one file after prune for", collection_name, key, [filenames[0]], "got", l)

    err = dataset.prune(collection_name, key)
    if err != '':
        t.error("Failed, expected True for prune (all)", collection_name, key, ', ', err)
    l = dataset.attachments(collection_name, key)
    if len(l) != 0:
        t.error("Failed, expected zero files after prune for", collection_name, key, "got", l)


def test_s3(t):
    collection_name = os.getenv("DATASET", "")
    if collection_name == "":
        t.print("Skipping test_s3(), missing environment S3 DATASET value to test with")
        return
    if collection_name[0:5] != "s3://":
        t.verbose_on()
        t.print("Skipping test_s3(), missing environment S3 DATASET value to test with")
        return
    
    ok = dataset.status(collection_name)
    if ok == False:
        t.print("Missing", collection_name, "attempting to initialize", collection_name)
        err = dataset.init(collection_name, "pairtree")
        if err != '':
            t.error("Aborting, couldn't initialize", collection_name, ', ', err)
            return
    else:
        t.print("Using collection initialized as", collection_name)

    collection_name = os.getenv("DATASET")
    record = { "one": 1 }
    key = "s3t1"
    err = dataset.create(collection_name, key, record)
    if err != '':
        t.error("Failed to create record", collection_name, key, record, ', ', err)
    record2, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    if record2.get("one") != 1:
        t.error("Failed, read", collection_name, key, record2)
    record["two"] = 2
    err = dataset.update(collection_name, key, record)
    if err != '':
        t.error("Failed to update record", collection_name, key, record, ', ', err)
    record2, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    if record2.get("one") != 1:
        t.error("Failed, 2nd read", collection_name, key, record2)
    if record2.get("two") != 2:
        t.error("Failed, 2nd read", collection_name, key, record2)
    err = dataset.delete(collection_name, key)
    if err != '':
        t.error("Failed to delete record", collection_name, key, record, ', ', err)
    ok = dataset.has_key(collection_name, key)
    if ok == True:
        t.error("Failed, delete should have removed key", collection_name, key)


def test_join(t, collection_name):
    key = "test_join1"
    obj1 = { "one": 1}
    obj2 = { "two": 2}
    ok = dataset.status(collection_name)
    if ok == False:
        t.error("Failed, collection status is False,", collection_name)
        return
    ok = dataset.has_key(collection_name, key)
    err = ''
    if ok == True:
        err = dataset.update(collection_nane, key, obj1)
    else:
        err = dataset.create(collection_name, key, obj1)
    if err != '':
        t.error(f'Failed, could not add record for test ({collection_name}, {key}, {obj1}), {err}')
        return
    err = dataset.join(collection_name, key, obj2, overwrite = False)
    if err != '':
        t.error(f'Failed, join for {collection_name}, {key}, {obj2}, overwrite = False -> {err}')
    obj_result, err = dataset.read(collection_name, key)
    if err != '':
        t.error(f'Unexpected error for {key} in {collection_name}, {err}')
    if obj_result.get('one') != 1:
        t.error(f'Failed to join append key {key}, {obj_result}')
    if obj_result.get("two") != 2:
        t.error(f'Failed to join append key {key}, {obj_result}')
    obj2['one'] = 3
    obj2['two'] = 3
    obj2['three'] = 3
    err = dataset.join(collection_name, key, obj2, overwrite = True)
    if err != '':
        t.error(f'Failed to join overwrite {collection_name}, {key}, {obj2}, overwrite = True -> {err}')
    obj_result, err = dataset.read(collection_name, key)
    if err != '':
        t.error(f'Unexpected error for {key} in {collection_name}, {err}')
    for k in obj_result:
        if k != '_Key' and obj_result[k] != 3:
            t.error('Failed to update value in join overwrite', k, obj_result)
    
#
# test_issue43() When exporting records to a table using
# use_srict_dotpath(True), the rows are getting miss aligned.
#
def test_issue43(t, collection_name, csv_name):
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    if os.path.exists(csv_name):
        os.remove(csv_name)
    err = dataset.init(collection_name, "pairtree")
    if err != '':
        t.error(f'Failed, need a {collection_name} to run test')
        return
    table = {
            "r1": {
                "c1": "one",
                "c2": "two",
                "c3": "three",
                "c4": "four"
                },
            "r2": {
                "c1": "one",
                "c3": "three",
                "c4": "four"
                },

            "r3": {
                "c1": "one",
                "c2": "two",
                "c4": "four"
                },
            "r4": {
                "c1": "one",
                "c2": "two",
                "c3": "three"
                },
            "r5": {
                "c1": "one",
                "c2": "two",
                "c3": "three",
                "c4": "four"
                }
            }
    for key in table:
        row = table[key]
        err = dataset.create(collection_name, key, row)
        if err != '':
            t.error(f"Can't add test row {key} to {collection_name}")
            return

    dataset.use_strict_dotpath(False)
    # Setup frame
    frame_name = 'f1'
    keys = dataset.keys(collection_name)
    (f1, err) = dataset.frame(collection_name, frame_name, keys, 
        ["._Key",".c1",".c2",".c3",".c4"])
    if err != '':
        t.error(err)
        return
    err = dataset.export_csv(collection_name, frame_name, csv_name)
    if err != '':
       t.error(f'export_csv({collection_name}, {frame_name}, {csv_name} should have emitted warnings, not error')
       return
    with open(csv_name, mode = 'r', encoding = 'utf-8') as f:
        rows = f.read()

    for row in rows.split('\n'):
        if len(row) > 0:
            cells = row.split(',')
            if len(cells) < 5:
                t.error(f'row error {csv_name} for {cells}')


def test_clone_sample(t, c_name, sample_size, training_name, test_name):
    if os.path.exists(training_name):
        shutil.rmtree(training_name)
    if os.path.exists(test_name):
        shutil.rmtree(test_name)
    err = dataset.clone_sample(c_name, training_name, test_name, sample_size)
    if err != '':
        t.error(f"can't clone sample {c_name} size {sample_size} into {training_name}, {test_name} error {err}")

def test_grid(t, c_name):
    if os.path.exists(c_name):
        shutil.rmtree(c_name)
    err = dataset.init(c_name, "pairtree")
    if err != '':
        t.error(err)
        return
    data = [
        { "id":    "A", "one":   "one", "two":   22, "three": 3.0, "four":  ["one", "two", "three"] },
        { "id":    "B", "two":   2000, "three": 3000.1 },
        { "id": "C" },
        { "id":    "D", "one":   "ONE", "two":   20, "three": 334.1, "four":  [] }
    ]
    keys = []
    dot_paths = ["._Key", ".one", ".two", ".three", ".four"]
    for row in data:
        key = row['id']
        keys.append(key)
        err = dataset.create(c_name, key, row)
    (g, err) = dataset.grid(c_name, keys, dot_paths)
    if err != '':
        t.error(err)

def test_frame(t, c_name):
    if os.path.exists(c_name):
        shutil.rmtree(c_name)
    err = dataset.init(c_name, "pairtree")
    if err != '':
        t.error(err)
        return
    data = [
        { "id":    "A", "one":   "one", "two":   22, "three": 3.0, "four":  ["one", "two", "three"] },
        { "id":    "B", "two":   2000, "three": 3000.1 },
        { "id": "C" },
        { "id":    "D", "one":   "ONE", "two":   20, "three": 334.1, "four":  [] }
    ]
    keys = []
    dot_paths = ["._Key", ".one", ".two", ".three", ".four"]
    for row in data:
        key = row['id']
        keys.append(key)
        err = dataset.create(c_name, key, row)
    f_name = 'f1'
    (g, err) = dataset.frame(c_name, f_name, keys, dot_paths)
    if err != '':
        t.error(err)
    err = dataset.reframe(c_name, f_name)
    if err != '':
        t.error(err)
    l = dataset.frames(c_name)
    if len(l) != 1 or l[0] != 'f1':
        t.error(f"expected one frame name, f1, got {l}")
    labels = ['Column A', 'Column B', 'Column C', 'Column D', 'Column E']
    err = dataset.frame_labels(c_name, f_name, labels)
    if err != '':
        t.error(err)
    if err != '':
        t.error(err)
    err = dataset.delete_frame(c_name, f_name)
    if err != '':
        t.error(err)

#
# test_sync_csv (issue 80) - add tests for sync_send_csv, sync_recieve_csv
#
def test_sync_csv(t, c_name):
    # Setup test collection
    if os.path.exists(c_name):
        shutil.rmtree(c_name)
    err = dataset.init(c_name, "pairtree")
    if err != '':
        t.error(err)
        return

    # Setup test CSV instance
    t_data = [
            { "_Key": "one", "value": 1 },
            { "_Key": "two", "value": 2 },
            { "_Key": "three", "value": 3  }
    ]
    csv_name = c_name.strip(".ds") + ".csv"
    if os.path.exists(csv_name):
        os.remove(csv_name)
    with open(csv_name, 'w') as csvfile:
        csv_writer = csv.DictWriter(csvfile, fieldnames = ["_Key", "value" ])
        csv_writer.writeheader()
        for obj in t_data:
            csv_writer.writerow(obj)
        
    # Import CSV into collection
    dataset.import_csv(c_name, csv_name, 1)
    for key in [ "one", "two", "three" ]:
        if dataset.has_key(c_name, key) == False:
            t.error(f"expected has_key({key}) == True, got False")
    if dataset.has_key(c_name, "five") == True:
        t.error(f"expected has_key('five') == False, got True")
    err = dataset.create(c_name, "five", {"value": 5})
    if err != "":
        t.error(err)
        return

    # Setup frame
    frame_name = 'test_sync'
    keys = dataset.keys(c_name)
    (frame, err) = dataset.frame(c_name, frame_name, keys, ["._Key", ".value"] )
    if err != '':
        t.error(err)
        return

    #NOTE: Tests for sync_send_csv and sync_receive_csv
    err = dataset.sync_send_csv(c_name, frame_name, csv_name)
    if err != '':
        t.error(err)
        return
    with open(csv_name) as fp:
        src = fp.read()
        if 'five' not in src:
            t.error(f"expected 'five' in src, got {src}")

    # Now remove "five" from collection
    err = dataset.delete(c_name, "five")
    if err != '':
        t.error(err)
        return
    if dataset.has_key(c_name, "five") == True:
        t.error(f"expected has_key(five) == False, got True")
        return
    err = dataset.sync_recieve_csv(c_name, frame_name, csv_name)
    if err != '':
        t.error(err)
        return
    if dataset.has_key(c_name, "five") == False:
        t.error(f"expected has_key(five) == True, got False")
        return


#
# test_basic_gsheet(t, collection_name) - test if basic GSheet
# import and export work.
#
def test_basic_gsheet(t, collection_name):
    '''if credentials.json and gsheet_test.json exist, run gsheet tests'''
    client_secret_name = "credentials.json"
    gsheet_test = "gsheet_test.json"
    for filename in [ client_secret_name, gsheet_test ]:
        if os.path.exists(filename) == False:
            t.print(f"Skipping test_gsheet({collection_name}), no {filename} found")
            return
    
    sheet_id = ""
    sheet_name = 'Sheet1'
    cell_range = 'A1:Z100'
    id_col = 1
    if os.path.exists('gsheet_test.json') == False:
        t.error("Can't find gsheet_test.json, skipping")
        return
    with open('gsheet_test.json') as f:
        src = f.read()
        obj = json.loads(src)
        if 'gsheet_id' in obj:
            sheet_id = obj['gsheet_id']

    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    err = dataset.init(collection_name, "pairtree")
    if err != '':
        t.error("Failed, could not create collection, ", err)
        return

    cnt = dataset.count(collection_name)
    if cnt != 0:
        t.error("Failed to initialize a fresh collection", collection_name)
        return

    # Setup some test data to work with.
    err = dataset.create(collection_name, "Wilson1930",  {"additional":"Supplemental Files Information:\nGeologic Plate: Supplement 1 from \"The geology of a portion of the Repetto Hills\" (Thesis)\n","description_1":"Supplement 1 in CaltechDATA: Geologic Plate","done":"yes","identifier_1":"https://doi.org/10.22002/D1.638","key":"Wilson1930","resolver":"http://resolver.caltech.edu/CaltechTHESIS:12032009-111148185","subjects":"Repetto Hills, Coyote Pass, sandstones, shales"})
    if err != '':
        t.error("Failed, could not create test record in", collection_name, ', ', err)
        return
    
    cnt = dataset.count(collection_name)
    if cnt != 1:
        t.error("Failed, should have one test record in", collection_name)
        return

    frame_name = 'f1'
    filter_expr = 'true'
    dot_exprs = ['.done','.key','.resolver','.subjects','.additional','.identifier_1','.description_1']
    column_names = ['_Key', 'Done','Key','Resolver','Subjects','Additional','Identifier 1','Description 1']
    keys = dataset.keys(collection_name)

    if dataset.has_frame(collection_name, frame_name):
        dataset.delete_frame(collection_name, frame_name)
    (f1, err) = dataset.frame(collection_name, frame_name, keys, dot_exprs)
    if err != '':
        t.error(err)
        return
    dataset.frame_labels(collection_name, frame_name, column_names)

    t.print("Testing gsheet export support", collection_name, frame_name, sheet_id, sheet_name, cell_range)
    err = dataset.export_gsheet(collection_name, frame_name, sheet_id, sheet_name, cell_range)
    if err != '':
        t.error("Failed, count not export-gsheet in", collection_name, ', ', err)
        return

    # Let's change our collection and see if overwrite works
    err = dataset.update(collection_name, "Wilson1930", {"Done": "No"})
    if err != '':
        t.error(f"Updated obj before import est {err}")
        return

    t.print("Testing gsheet import support (should fail)", sheet_id, sheet_name, id_col, cell_range, False)
    err = dataset.import_gsheet(collection_name, sheet_id, sheet_name, id_col, cell_range, overwrite = False)
    if err != '':
        t.error(err)
        return
   
    (obj, err) = dataset.read(collection_name, "Wilson1930")
    if err != "":
        t.error(err)
        return

    if "Done" in obj and obj['Done'] == "Yes":
        t.error(f"should NOT import-gsheet over record with overwrite = False")
        return

    t.print("Testing gsheet import support (should succeeed)", sheet_id, sheet_name, id_col, cell_range, True)
    err = dataset.import_gsheet(collection_name, sheet_id, sheet_name, id_col, cell_range, overwrite = True) 
    if err != '':
        t.error(err)
    (obj, err) = dataset.read(collection_name, "Wilson1930")
    if "Done" in obj and obj['Done'] == "No":
        t.error(f"should be able to {obj} import with overwrite = True")
        return
    

#
# test_sync_gsheet (issue 80) - add tests for 
# sync_send_gsheet and sync_recieve_gsheet 
#
def test_sync_gsheet(t, c_name):
    # Setup test collection
    if os.path.exists(c_name):
        shutil.rmtree(c_name)
    err = dataset.init(c_name, "pairtree")
    if err != '':
        t.error(err)
        return

    # Setup test CSV instance
    t_data = [
            { "_Key": "one", "value": 1 },
            { "_Key": "two", "value": 2 },
            { "_Key": "three", "value": 3  }
    ]

    gsheet_id = ""
    if os.path.exists('gsheet_test.json') == False:
        t.error("Can't find gsheet_test.json, skipping")
        return
    with open('gsheet_test.json') as f:
        src = f.read()
        obj = json.loads(src)
        if 'gsheet_id' in obj:
            gsheet_id = obj['gsheet_id']
    gsheet_name = 'test_sheet'
    cell_range = 'A1:Z100'

    for obj in t_data:
        err = dataset.create(c_name, obj['_Key'], {"value": obj['value']})
        if err != '':
            t.error(f"expected create to succeed, got {err}")
            return

    f_name = 'test_sync'
    frame, err = dataset.frame(c_name, f_name, ['one', 'two', 'three'], ['_Key', 'value'])
    if err != '':
        t.error(err)
        return

    err = dataset.export_gsheet(c_name, f_name, gsheet_id, gsheet_name, cell_range)
    #err = dataset.export_gsheet(c_name, f_name, gsheet_id, gsheet_name)
    if err != '':
        t.error(err)
        return

#
# Test harness
#
class ATest:
    def __init__(self, test_name, verbose = False):
        self._test_name = test_name
        self._error_count = 0
        self._verbose = False

    def test_name(self):
        return self._test_name

    def is_verbose(self):
        return self._verbose

    def verbose_on(self):
        self._verbose = True

    def verbose_off(self):
        self.verbose = False

    def print(self, *msg):
        if self._verbose == True:
            print(*msg)

    def error(self, *msg):
        fn_name = self._test_name
        self._error_count += 1
        print(f"\t{fn_name}", *msg)

    def error_count(self):
        return self._error_count

class TestRunner:
    def __init__(self, set_name, verbose = False):
        self._set_name = set_name
        self._tests = []
        self._error_count = 0
        self._verbose = verbose

    def add(self, fn, params = []):
        self._tests.append((fn, params))

    def run(self):
        for test in self._tests:
            fn_name = test[0].__name__
            t = ATest(fn_name, self._verbose)
            fn, params = test[0], test[1]
            fn(t, *params)
            error_count = t.error_count()
            if error_count > 0:
                print(f"\t\t{fn_name} failed, {error_count} errors found")
            else:
                print(f"\t\t{fn_name} OK")
            self._error_count += error_count
        error_count = self._error_count
        set_name = self._set_name
        if error_count > 0:
            print(f"Failed {set_name}, {error_count} total errors found")
            sys.exit(1)
        print("PASS")
        print("Ok", __file__)
        sys.exit(0)

#
# Main processing
#
if __name__ == "__main__":
    print("Starting dataset_test.py")
    print("Testing dataset version", dataset.version())
    run_gsheet_tests = False
    if "gsheet" in sys.argv or "gsheets" in sys.argv:
        if os.path.exists("gsheet_test.json") == True:
            run_gsheet_tests = True
        else:
            print("""
ERROR: Can't find gsheet_test.json and credentials.json 
needed run the 'gsheet' tests.

The gsheet testing needs to be setup manually with
the dataset command line tool before running the
Python 'gsheet' test sequence.

""")
            sys.exit(1)

    else:
        print("NOTICE: Skipped 'gsheet' tests")

    # Pre-test check
    error_count = 0
    ok = True
    dataset.verbose_off()

    collection_name = "test_collection.ds"
    test_runner = TestRunner(os.path.basename(__file__))
    test_runner.add(test_setup, [collection_name])
    test_runner.add(test_basic, [collection_name])
    test_runner.add(test_keys, [collection_name])
    test_runner.add(test_search, [collection_name, "test_index_map.json", "test_index.bleve"])
    test_runner.add(test_issue32, [collection_name])
    test_runner.add(test_attachments, [collection_name])
    test_runner.add(test_join, [collection_name])
    test_runner.add(test_check_repair, ["test_check_and_repair.ds"])
    test_runner.add(test_issue43,["test_issue43.ds", "test_issue43.csv"])
    test_runner.add(test_s3)
    test_runner.add(test_clone_sample, ["test_collection.ds", 5, "test_training.ds", "test_test.ds"])
    test_runner.add(test_grid, ["test_grid.ds"])
    test_runner.add(test_frame, ["test_frame.ds"])

    test_runner.add(test_sync_csv, ["test_sync_csv.ds"])
    if run_gsheet_tests == True:
        test_runner.add(test_basic_gsheet, ["test_gsheet.ds"])
        test_runner.add(test_sync_gsheet, ["test_sync_gsheet.ds"])

    test_runner.run()

