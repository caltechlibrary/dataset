#!/usr/bin/env python3.7
import os
import sys
import shutil
import json
import csv
from pathlib import Path
from libdataset import dataset 

path_sep = "/"
if sys.platform.startswith('win'):
    path_sep = "\\"

os.makedirs('testout', 0o777, exist_ok = True)

def cleanup(c_name):
    keys = dataset.keys(c_name)
    fnames = dataset.frame_names(c_name)
    for fname in fnames:
        if dataset.delete_frame(c_name, fname) == False:
            print(f'WARNING delete_frame({c_name}, {fname}) failed, {err}')
    if keys != None and len(keys) > 0:
        for key in keys:
            if dataset.delete(c_name, key) == False:
                print(f'WARNING delete({c_name}, {key}) failed, {err}')
        fnames = dataset.frame_names(c_name)
    if fnames != None and len(fnames) > 0:
        print(f'Cleanup failed, {c_name} has following frame_names {fnames}')
        sys.exit(1)
    keys = dataset.keys(c_name)
    if keys != None and len(keys) > 0:
        print(f'Cleanup failed, {c_name} has following keys {keys}')
        sys.exit(1)
    if dataset.is_open(c_name):
        if dataset.close_collection(c_name) == False:
            err = dataset.error_message()
            print(f'Failed, close_collection({c_name}), {err}')
    elif os.path.exists(c_name):
        shutil.rmtree(c_name)
    if dataset.init(c_name, "") == False:
        err = dataset.error_message()
        print(f'Failed, init({c_name}), {err}')
        sys.exit(1)


# Setup our test collection deleting it first if neccessary
def test_setup(t, c_name, test_name):
    if os.path.exists(c_name) == False:
        t.print(f'Creating {c_name} for {test_name}')
        if dataset.init(c_name, "") == False:
            err = dataset.error_message()
            t.error(f"{test_name} Failed, could not create collection, {err}")
            sys.exit(1)
        if os.path.exists(c_name) == False:
            t.print(f"{c_name} does not exist! {test_name}")
            sys.exit(1)
    else:
        t.print(f'Using {c_name}')

    
def test_libdataset(t, c_name):
    # Clean up stale result test collections
    cleanup(c_name)

    src = '''
    [
        {"_Key": "k1", "title": "One thing", "name": "Fred"},
        {"_Key": "k2", "title": "Two things", "name": "Frieda"},
        {"_Key": "k3", "title": "Three things", "name": "Fiona"}
    ]
    '''
    
    for obj in json.loads(src):
        key = obj['_Key']
        if dataset.create(c_name, key, obj) == False:
            err = dataset.error_message()
            t.error(f"expected '', got '{err}' for dataset.create({c_name}, {key}, {obj})")
            sys.exit(1)
    
    expected_keys = [ "k1", "k2", "k3" ]
    keys = dataset.keys(c_name)
    i = 0
    for key in keys:
        if not key in expected_keys:
            t.error(f"expected {key} in {expected_keys} for {c_name}")
            sys.exit(1)
        obj, err = dataset.read(c_name, key)
        if err != '':
            t.error(f"expected '', got '{err}' for dataset.read({c_name}, {key}")
            sys.exit(1)
        obj['t_count'] = i
        i += 1
        if dataset.update(c_name, key, obj) == False:
            err = dataset.error_message()
            t.error(f"expected '', got '{err}' for dataset.update({c_name}, {key}, ...")
            sys.exit(1)
    
    f_name = "f1"
    if dataset.frame_create(c_name, f_name, keys[1:], ['._Key', '.title'], [ 'id', 'title' ]) == False:
        err = dataset.error_message()
        t.error(f"expected '', got '{err}' for dataset.frame_create({c_name}, {f_name}, ...)")
        sys.exit(1)
    
    ok = dataset.has_frame(c_name, f_name)
    if ok != True:
            t.error(f"expected 'True', got '{ok}' for dataset.has_frame({c_name}, {f_name})")
            sys.exit(1)
    
    
    expected_keys = keys[1:]
    keys = dataset.frame_keys(c_name, f_name)
    for i, expected in enumerate(expected_keys):
        key = keys[i]
        if key != expected:
            t.error(f"expected ({i}) '{expected}', got '{key}' for dataset.frame_keys({c_name}, {f_name})")
            sys.exit(1)

#
# test_basic(collection_name) runs tests on basic CRUD ops
# 
def test_basic(t, collection_name):
    '''test_basic(collection_name) runs tests on basic CRUD ops'''
    cleanup(collection_name)

    keys = dataset.keys(collection_name)
    if keys != None and len(keys) != 0:
        t.error(f"Something is wrong {collection_name} should be empty, got {len(keys)} keys {keys}")
        sys.exit(1)

    # Setup a test record
    key = "2488"
    value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}
    
    # We should have an empty collection, we will create our test record.
    if dataset.create(collection_name, key, value) == False:
        err = dataset.error_message()
        t.error(f"Failed, could not create record {key}")
    
    # Check to see that we have only one record
    key_count = dataset.count(collection_name)
    if key_count != 1:
        t.error(f"Failed {collection_name}, expected count to be 1, got {key_count}")
    
    # Do a minimal test to see if the record looks like it has content
    keyList = dataset.keys(collection_name)
    if keyList == None:
        t.error(f"Expected keys in keyList, not None")
    elif len(keyList) != 1:
        t.error(f"Expected one key in keyList, {keyList}")
    rec, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    for k, v in value.items():
       if type(v) != list:
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
    if dataset.update(collection_name, key, value) == False:
        err = dataset.error_message()
        t.error(f"Failed, count not update record {key}, {value}, {err}")
    rec, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    for k, v in value.items():
       if type(v) != list:
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
    cwd = f"{Path('.').resolve()}"
    expected_s = path_sep.join([cwd, collection_name, "pairtree", "24", "88", (key+".json")])
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
    if dataset.delete(collection_name, key) == False:
        err = dataset.error_message()
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
        if dataset.create(collection_name, k, v) == False:
            err = dataset.error_message()
            t.error("Failed, could not add", k, "to", collection_name, ', ', err)
    
    # Test keys, filtering keys and sorting keys
    all_keys = dataset.keys(collection_name)
    if len(all_keys) != test_count:
        t.error("Expected (a)", test_count,"keys back, got", keys)
    
    # key_filter test removed, not found in 1.1 of libdataset.go
    #filter_expr = '(eq .categories "non-fiction, memoir")'
    #keys = dataset.key_filter(collection_name, all_keys, filter_expr)
    #if len(keys) != 1:
    #    t.error("Expected (b) one key for", filter_expr, "got", keys)
    # 
    # filter_expr = '(contains .categories "novel")'
    #filtered_keys = dataset.key_filter(collection_name, all_keys, filter_expr)
    #if len(filtered_keys) != 3:
    #    t.error("Expected (c) three keys for", filter_expr, "got", keys)
    
    # key_sort test removed, not found in v1.1 of libdataset.go
    #sort_expr = '+.title'
    #keys = dataset.key_sort(collection_name, filtered_keys, sort_expr)
    #if len(keys) != 3:
    #    t.error("Expected (d) three keys for", filter_expr, "got", keys)
    #expected_keys = [ "gutenberg:21839", "gutenberg:21489", "gutenberg:2488" ]
    #for i, k in enumerate(expected_keys):
    #    if i < len(keys) and keys[i] != k:
    #        t.error("Expected (e)", k, "got", keys[i])
    

#
# test_issue32() make sure issue 32 stays fixed.
#
def test_issue32(t, collection_name):
    if dataset.create(collection_name, "k1", {"one":1}) == False:
        err = dataset.error_message()
        t.error("Failed to create k1 in", collection_name, ', ', err)
        return
    ok = dataset.has_key(collection_name, "k1")
    if ok == False:
        t.error("Failed, has_key k1 should return", True)
    ok = dataset.has_key(collection_name, "k2")
    if ok == True:
        t.error("Failed, has_key k2 should return", False)



def test_check_repair(t, c_name):
    cleanup(c_name)
    t.print("Testing status on", c_name)
    # Make sure we have a left over collection to check and repair
    dataset.init(c_name, "")
    if dataset.status(c_name) == False:
        t.error("Failed, expected dataset.status() == True, got", ok, "for", c_name)
        return

    if dataset.has_key(c_name, 'one') == False:
        dataset.create(c_name, 'one', {"one": 1})
    t.print("Testing check on", c_name)
    # Check our collection
    ok = dataset.check(c_name)
    if ok == False:
        t.error("Failed, expected check", c_name, "to return True, got", ok)
        return

    # Break and recheck our collection
    print(f"Removing {c_name}/collection.json to cause a fail")
    if os.path.exists(c_name + "/collection.json"):
        os.remove(c_name + "/collection.json")
    print(f"Testing check on (broken) {c_name}")
    ok = dataset.check(c_name)
    if ok == True:
        t.error("Failed, expected check", c_name, "to return False, got", ok)
    else:
        t.print(f"Should have see error output for broken {c_name}")

    # Repair our collection
    t.print("Testing repair on", c_name)
    if dataset.repair(c_name) == False:
        err = dataset.error_message()
        t.error("Failed, expected repair to return True, got, ", err)
    if os.path.exists(c_name + "/collection.json") == False:
        t.error(f"Failed, expected recreated {c_name}/collection.json")
 
        
def test_attachments(t, c_name):
    t.print("Testing attach, attachments, detach and prune")
    # Generate two files to attach.
    attachedNames = [ 'a1.txt', 'a2.txt' ]
    filenames = []
    for i, f_name in enumerate(attachedNames):
        fp = open(os.path.join('testout', f_name), 'w')
        fp.write(f"This is file ({i}) {f_name}\n")
        fp.close()
        filenames.append(os.path.join('testout', f_name))
        if not os.path.exists(filenames[i]):
            t.error(f'Failed to write testdata {filenames[i]}')

    ok = dataset.status(c_name)
    if ok == False:
        t.error("Failed,", c_name, "missing")
        return
    keys = dataset.keys(c_name)
    if len(keys) < 1:
        t.error("Failed,", c_name, "should have keys")
        return

    key = keys[0]
    if dataset.attach(c_name, key, filenames) == False:
        err = dataset.error_message()
        t.error("Failed, to attach files for", c_name, key, filenames, ', ', err)
        return

    l = dataset.attachments(c_name, key)
    #print(f'DEBUG attachments {l}')
    if len(l) != 2:
        t.error(f"Failed, expected two attachments for {c_name} -> {key}, got ({len(l)}) {l}")
        return
    for a_name in attachedNames:
        if not a_name in l:
            t.error(f'expected {a_name} in attached list {l}')
    # Check that attachments aren't impacted by update
    if dataset.update(c_name, key, {"testing":"update"}) == False:
        err = dataset.error_message()
        t.error("Failed, to update record", c_name, key, err)
        return
    l = dataset.attachments(c_name, key)
    if len(l) != 2:
        t.error(f"Failed, expected two attachments after update for {c_name} got ({len(l)}) {l}")
        return
    for a_name in attachedNames:
        if not a_name in l:
            t.error(f'expected {a_name} after updated in attached list {l}')

    # First try detaching one file at a time.
    for f_name in attachedNames:
        # Remove the stale files from the local folder.
        if os.path.exists(f_name):
            os.remove(f_name)
        if dataset.detach(c_name, key, [f_name]) == False:
            err = dataset.error_message()
            t.error("Failed single file detach, expected True for", c_name, key, f_name, ', ', err)

    # Test explicit filenames list detach
    if dataset.detach(c_name, key, attachedNames) == False:
        err = dataset.error_message()
        t.error("Failed detach list, expected True for", c_name, key, filenames, ', ', err)
    for f_name in attachedNames:
        if os.path.exists(f_name):
            os.remove(f_name)
        else:
            t.error("Failed detach list, expected", f_name, "to be detached from", c_name, key)

    # Test detaching all files
    if dataset.detach(c_name, key, []) == False:
        err = dataset.error_message()
        t.error("Failed detaching all, expected True for (detaching all)", c_name, key, ', ', err)
    for f_name in attachedNames:
        if os.path.exists(f_name):
            os.remove(f_name)
        else:
            t.error("Failed detaching all, expected", f_name, "for detaching all from", c_name, key)

    if dataset.prune(c_name, key, [attachedNames[0]]) == False:
        err = dataset.message()
        t.error("Failed, expected True for prune", c_name, key, [attachedNames[0]], ', ', err)
    l = dataset.attachments(c_name, key)
    if len(l) != 1:
        t.error(f"Failed, expected one {attachedNames[0]} pruned for {c_name} -> {key}, got ({len(l)}) {l}")

    if dataset.prune(c_name, key, []) == False:
        err = dataset.error_message()
        t.error("Failed, expected True for prune (all)", c_name, key, ', ', err)
    l = dataset.attachments(c_name, key)
    if len(l) != 0:
        t.error("Failed, expected zero files after prune for", c_name, key, "got", l)


def test_join(t, collection_name):
    key = "test_join1"
    obj1 = { "one": 1}
    obj2 = { "two": 2}
    err = ''
    if dataset.has_key(collection_name, key):
        if dataset.update(collection_name, key, obj1) == False:
            err = dataset.error_message()
            t.error(f'update({collection_name}, {key}, {obj1}) failed, {err}')
    else:
        if dataset.create(collection_name, key, obj1) == False:
            err = dataset.error_message()
            t.error(f'create({collection_name}, {key}, {obj1}) failed, {err}')
    if err != '':
        t.error(f'Failed (a), could not add record for test ({collection_name}, {key}, {obj1}), {err}')
        return
    if dataset.join(collection_name, key, obj2, overwrite = False) == False:
        t.error(f'Failed (b), join({collection_name}, {key}, {obj2}, overwrite = False) -> returned False')
    obj_result, err = dataset.read(collection_name, key)
    if err != '':
        t.error(f'Unexpected error for {key} in {collection_name}, {err}')
    if obj_result.get('one') != 1:
        t.error(f'Failed (c) to join append key {key}, {obj_result}')
    if obj_result.get("two") != 2:
        t.error(f'Failed (d) to join append key {key}, {obj_result}')
    obj2['one'] = 3
    obj2['two'] = 3
    obj2['three'] = 3
    if dataset.join(collection_name, key, obj2, overwrite = True) == False:
        t.error(f'Failed (e) join({collection_name}, {key}, {obj2}, overwrite = True) -> False')
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
def test_issue43(t, c_name, csv_name):
    cleanup(c_name)
    if dataset.init(c_name) == False:
        err = dataset.error_message()
        t.error(f'Failed, need a {c_name} to run test')
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
        if dataset.create(c_name, key, row) == False:
            err = dataset.error_message()
            t.error(f"Can't add test row {key} to {c_name}, {err}")
            return

    dataset.use_strict_dotpath(False)
    # Setup frame
    frame_name = 'f1'
    keys = dataset.keys(c_name)
    if dataset.frame_create(c_name, frame_name, keys, ["._Key",".c1",".c2",".c3",".c4"], ["_Key", "c1", "c2", "c3", "c4"]) == False:
        err = dataset.error_message()
        t.error(f'frame_create({c_name}, {frame_name}, ...) failed, {err}')
        return
    if dataset.export_csv(c_name, frame_name, csv_name) == False:
        err = dataset.error_message()
        t.error(f'export_csv({c_name}, {frame_name}, {csv_name} should have emitted warnings, not error')
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
    if dataset.clone_sample(c_name, training_name, "", test_name, "", sample_size) == False:
        err = dataset.error_message()
        t.error(f"can't clone sample {c_name} size {sample_size} into {training_name}, {test_name} error {err}")

def test_frame1(t, c_name):
    cleanup(c_name)
    if dataset.init(c_name) == False:
        err = dataset.error_message()
        t.errorf(f'failed to create {c_name}, {err}')
        return

    data = [
        { "id":    "A", "one":   "one", "two":   22, "three": 3.0, "four":  ["one", "two", "three"] },
        { "id":    "B", "two":   2000, "three": 3000.1 },
        { "id": "C" },
        { "id":    "D", "one":   "ONE", "two":   20, "three": 334.1, "four":  [] }
    ]
    keys = []
    dot_paths = ["._Key", ".one", ".two", ".three", ".four"]
    labels = ["_Key", "one", "two", "three", "four"]
    for row in data:
        key = row['id']
        keys.append(key)
        if dataset.create(c_name, key, row) == False:
            err = dataset.error_message()
            t.error(f'create({c_name}, {key}, {row}) failed, {err}')
    f_name = 'f1'
    if dataset.frame_create(c_name, f_name, keys, dot_paths, labels) == False:
        err = dataset.error_message()
        t.error(f'frame_create({c_name}, {f_name}, {keys}, {dot_paths}, {labels}) -> {err}')
        return
    if dataset.frame_reframe(c_name, f_name, keys) == False:
        t.error(f'frame_reframe({c_name}, {f_name}) returned False')
    l = dataset.frame_names(c_name)
    if len(l) != 1 or l[0] != 'f1':
        t.error(f"expected one frame name, f1, got {l}")
    if dataset.delete_frame(c_name, f_name) == False:
        t.error(f'delete_frame({c_name}, {f_name}) returned False.')


def test_frame2(t, c_name):
    cleanup(c_name)

    data = [
        { "id": "A", "nameIdentifiers": [
                {
                    "nameIdentifier": "0000-000X-XXXX-XXXX",
                    "nameIdentifierScheme": "ORCID",
                    "schemeURI": "http://orcid.org/"
                },
                {
                    "nameIdentifier": "H-XXXX-XXXX",
                    "nameIdentifierScheme": "ResearcherID",
                    "schemeURI": "http://www.researcherid.com/rid/"
                }], "two":   22, "three": 3.0, "four":  ["one", "two", "three"] },
        { "id": "B", "two":   2000, "three": 3000.1 },
        { "id": "C" },
        { "id": "D", "nameIdentifiers": [
                {
                    "nameIdentifier": "0000-000X-XXXX-XXXX",
                    "nameIdentifierScheme": "ORCID",
                    "schemeURI": "http://orcid.org/"
                }], "two":   20, "three": 334.1, "four":  [] }
    ]
    keys = []
    saved_keys = []
    dot_paths = [".id", ".nameIdentifiers", ".nameIdentifiers[:].nameIdentifier", ".two", ".three", ".four"]
    labels = ["id", "nameIdentifiers", "nameIdentifier", "two", "three", "four"]
    for row in data:
        key = row['id']
        saved_keys.append(key.lower())
        if dataset.create(c_name, key, row) == False:
            err = dataset.error_message()
            t.error(f'create({c_name}, {key}, {row}) failed, {err}')
    keys = dataset.keys(c_name)
    for i, key in enumerate(saved_keys):
        if not dataset.has_key(c_name, key):
            t.error(f'expected has_key True ({i}) {key} in {c_name}')
        if not key in keys:
            t.error(f'expected key ({i}) {key} in {c_name}')
    f_name = 'f1'
    if dataset.frame_create(c_name, f_name, keys, dot_paths, labels) == False:
        err = dataset.error_message()
        t.error(f'frame_create({c_name}, {f_name}, {keys}, ...) failed, {err}')
        return
    f = dataset.frame(c_name, f_name)
    if f == None:
        t.error(f'after frame_create(), frame({c_name}, {f_name}) returned None, expected frame data')
        return
    if dataset.frame_reframe(c_name, f_name, keys) == False:
        t.error(f'frame_reframe({c_name}, {f_name}) returned False')
        return
    f = dataset.frame(c_name, f_name)
    if f == None:
        t.error(f'after frame_reframe(), frame({c_name}, {f_name}) returned None, expected frame data')
        return
    if len(f['keys']) == 0:
        t.error(f'missing keys after reframe, frame({c_name}, {f_name}) -> {f}')
        return

    l = dataset.frame_names(c_name)
    if len(l) != 1 or l[0] != 'f1':
        t.error(f"expected one frame name, f1, got {l}")
        return
    object_result = dataset.frame_objects(c_name, f_name)
    if len(object_result) != 4:
        t.error(f'frame_objects({c_name}, {f_name}), expected 4 got {len(object_result)} -> {object_result}')
        return
    count_nameId = 0
    count_nameIdObj = 0
    for i, obj in enumerate(object_result):
        if 'id' not in obj:
            t.error(f'Did not get id in object {i} -> {ob}')
        if 'nameIdentifiers' in obj:
            count_nameId += 1
            for idv in obj['nameIdentifiers']:
                if 'nameIdentifier' not in idv:
                    t.error('Missing part of object (#{i})')
        if 'nameIdentifier' in obj:
            count_nameIdObj += 1
            if "0000-000X-XXXX-XXXX" not in obj['nameIdentifier']:
                t.error('Missing object in complex dot path')
    if count_nameId != 2:
        t.error(f"Incorrect number of nameIdentifiers elements, expected 2, got {count_nameId}")
        return
    if count_nameIdObj != 2:
        t.error(f"Incorrect number of nameIdentifier elements, expected 2, got {count_nameIdObj}")
        return

#
# test_import_csv
#
def test_import_csv(t, c_name):
    csv_name = c_name.strip(".ds") + ".csv"
    if os.path.exists(csv_name):
        os.remove(csv_name)
    if os.path.exists(c_name):
        shutil.rmtree(c_name)
    dataset.init(c_name)
    with open(csv_name, 'w') as csvfile:
        csv_writer = csv.DictWriter(csvfile, fieldnames, [ 'name', 'email', 'id' ])
        csv_writer.writeheader()
        csv_writer.writer_row({'name': 'Gandolf', 'email': 'gtw@middleearth.example.edu', 'id': 'gtw'})
    err = dataset.import_csv(c_name, csv_name, False, True)

  

#
# test_sync_csv (issue 80) - add tests for sync_send_csv, sync_recieve_csv
#
def test_sync_csv(t, c_name):
    cleanup(c_name)

    # Setup test CSV instance
    t_data = [
            { "key": "one",   "value": 1 },
            { "key": "two",   "value": 2 },
            { "key": "three", "value": 3 }
    ]
    csv_name = c_name.strip(".ds") + ".csv"
    if os.path.exists(csv_name):
        os.remove(csv_name)
    with open(csv_name, 'w') as csvfile:
        csv_writer = csv.DictWriter(csvfile, fieldnames = [ "key", "value" ])
        csv_writer.writeheader()
        for obj in t_data:
            csv_writer.writerow(obj)
        
    # Import CSV into collection
    dataset.import_csv(c_name, csv_name, True, True)
    for key in [ "one", "two", "three" ]:
        if dataset.has_key(c_name, key) == False:
            t.error(f"expected has_key({key}) == True, got False")
    if dataset.has_key(c_name, "five") == True:
        t.error(f"expected has_key('five') == False, got True")
    if dataset.create(c_name, "five", {"key": "five", "value": 5}) == False:
        err = dataset.error_message()
        t.error(err)
        return

    # Setup frame
    frame_name = 'test_sync'
    keys = dataset.keys(c_name)
    if len(keys) != 4:
        t.error(f'expected 4 keys, got {keys}')
    if dataset.frame_create(c_name, frame_name, keys, [ ".key", ".value" ], [ "key", "value" ]) == False:
        err = dataset.error_message()
        t.error(f'frame_create({c_name}, {frame_name}, {keys}, ...) failed, {err}')
        return

    #NOTE: Tests for sync_send_csv and sync_receive_csv
    if dataset.sync_send_csv(c_name, frame_name, csv_name) == False:
        err = dataset.error_message()
        t.error(f'sync_send_csv({c_name}, {frame_name}, {csv_name}) failed, {err}')
        return
    with open(csv_name) as fp:
        src = fp.read()
        if 'five' not in src:
            t.error(f"expected 'five' in src, got {src}")

    # Now remove "five" from collection
    if dataset.delete(c_name, "five") == False:
        err = dataset.error_message()
        t.error(f'delete({c_name}, "five") failed, {err}')
        return
    if dataset.has_key(c_name, "five") == True:
        t.error(f"expected has_key(five) == False, got True")
        return
    if dataset.sync_recieve_csv(c_name, frame_name, csv_name) == False:
        err = dataset.error_message()
        t.error(f'sync_receive_csv({c_name}, {frame_name}, {csv_name}) failed, {err}')
        return
    if dataset.has_key(c_name, "five") == False:
        t.error(f"expected has_key(five) == True, got False")
        return

#
# test_issue12() https://github.com/caltechlibrary/py_dataset/issues/12
# delete_frame() returns True but frame metadata still in memory.
#
def test_issue12(t, c_name):
    src = '''[
{"id": "1", "c1": 1, "c2": 2, "c3": 3 },
{"id": "2", "c1": 2, "c2": 2, "c3": 3 },
{"id": "3", "c1": 3, "c2": 3, "c3": 3 },
{"id": "4", "c1": 1, "c2": 1, "c3": 1 },
{"id": "5", "c1": 6, "c2": 6, "c3": 6 }
]'''
    cleanup(c_name)
    if dataset.init(c_name, "") == False:
        err = dataset.error_message()
        t.error(f'failed to create {c_name} -> {err}')
        return
    if dataset.status(c_name) == False:
        t.error(f'failed to find {c_name} after init')
        return 

    objects = json.loads(src)
    for obj in objects:
        key = obj['id']
        if dataset.has_key(c_name, key):
            dataset.update(c_name, key, obj)
        else:            
            dataset.create(c_name, key, obj)
    f_names = dataset.frame_names(c_name)
    for f_name in f_names:
        if not dataset.delete_frame(c_name, f_name):
            err = dataset.error_message()
            t.error(f'Failed to delete {f_name} from {c_name} -> "{err}"')
    f_name = 'issue12'
    dot_paths = [ ".c1", "c3" ]
    labels = [ ".col1", "col3" ]
    keys = dataset.keys(c_name)
    if not dataset.frame_create(c_name, f_name, keys, dot_paths, labels):
        err = dataset.error_message()
        t.error(f'failed to create {f_name} from {c_name} -> "{err}"')
        return
    if not dataset.has_frame(c_name, f_name):
        err = dataset.error_message()
        t.error(f'expected frame {f_name} to exists, {err}')
        return
    f_keys = dataset.frame_keys(c_name, f_name)
    if len(f_keys) == 0:
        err = dataset.error_message()
        t.error(f'expected keys in {f_name}, got zero, {err}')
        return
    f_objects = dataset.frame_objects(c_name, f_name)
    if len(f_objects) == 0:
        err = dataset.error_message()
        t.error(f'expected objects in {f_name}, got zero, {err}')
        return
    # Note test frame_clear should remove keys/objects but leave frame ...
    if not dataset.frame_clear(c_name, f_name):
        err = dataset.error_message()
        t.error(f'expected to clear frame {f_name} in {c_name}, {err}')
    else:
        f_objects = dataset.frame_objects(c_name, f_name)
        if len(f_objects) != 0:
            t.error(f'frame_clear({c_name}, {f_name}) should have removed objects!')
    if not dataset.delete_frame(c_name, f_name):
        err = dataset.error_message()
        t.error(f'expected to delete {f_name} in {c_name}, {err}')


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
        print(f"{fn_name}: ", *msg)

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
                #return
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
    app_name = os.path.basename(sys.argv[0])
    print(f"Setting up {app_name}")

    # Pre-test check
    error_count = 0
    ok = True

    print(f'Starting {app_name}')
    test_runner = TestRunner(os.path.basename(__file__), True)
    c_name = os.path.join('testout', 'test_collection.ds')
    test_runner.add(test_setup, [ c_name, 'test_setup' ])
    test_runner.add(test_libdataset, [ c_name ])
    test_runner.add(test_basic, [ c_name ])
    test_runner.add(test_keys, [ c_name ])
    test_runner.add(test_issue32, [ c_name ])
    test_runner.add(test_attachments, [ c_name ])
    test_runner.add(test_join, [ c_name ])
    test_runner.add(test_issue43, [ os.path.join('testout', "test_issue43.ds"), os.path.join('testout', "test_issue43.csv") ])
    test_runner.add(test_clone_sample, [ c_name, 5, os.path.join('testout', "test_training.ds"), os.path.join('testout', "test_test.ds") ])
    test_runner.add(test_frame1, [ os.path.join('testout', "test_frame1.ds") ])
    test_runner.add(test_frame2, [ os.path.join('testout', "test_frame2.ds") ])
    test_runner.add(test_sync_csv, [ os.path.join('testout', "test_sync_csv.ds") ])
    test_runner.add(test_check_repair, [ os.path.join('testout', "test_check_and_repair.ds") ])
    test_runner.add(test_issue12, [ os.path.join('testout', 'test_issue12.ds') ])
    test_runner.run()
