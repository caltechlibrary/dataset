#!/usr/bin/env python3
import sys
import os
import shutil
import json
import dataset

#
# test_basic(collection_name) runs tests on basic CRUD ops
# 
def test_basic(collection_name):
    '''test_basic(collection_name) runs tests on basic CRUD ops'''
    error_count = 0
    # Setup a test record
    key = "2488"
    value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}
    
    # We should have an empty collection, we will create our test record.
    ok = dataset.create(collection_name, key, value)
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
    rec = dataset.read(collection_name, key)
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
    ok = dataset.update(collection_name, key, value)
    if ok == False:
       print("Failed, count not update record", key, value)
       error_count += 1
    rec = dataset.read(collection_name, key)
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
    
    # Test path to record
    expected_s = collection_name+"/aa/"+key+".json"
    expected_l = len(expected_s)
    p = dataset.path(collection_name, key)
    if len(p) != expected_l:
        print("Failed, expected length", expected_l, "got", len(p))
        error_count += 1
    if p != expected_s:
        print("Failed, expected", expected_s, "got", p)
        error_count += 1

    # Test listing records
    l = dataset.list(collection_name, [key])
    if len(l) != 1:
        print("Failed, list should return an array of one record, got", l)
        error_count += 1
        return error_count

    # test deleting a record
    ok = dataset.delete(collection_name, key)
    if ok == False:
        print("Failed, could not delete record", key)
        error_count += 1
    # test_base() done
    if error_count > 0:
        print("Test failed")
    return error_count
    

#
# test_keys(collection_name) test getting, filter and sorting keys
#
def test_keys(collection_name):
    '''test_keys(collection_name) test getting, filter and sorting keys'''
    error_count = 0
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
    test_count = len(test_records)
    
    for k in test_records:
        v = test_records[k]
        ok = dataset.create(collection_name, k, v)
        if ok == False:
            print("Failed, could not add", k, "to", collection_name)
            error_count += 1
    
    # Test keys, filtering keys and sorting keys
    keys = dataset.keys(collection_name)
    if len(keys) != test_count:
        print("Expected", test_count,"keys back, got", keys)
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
            error_count += 1
        i += 1

    # test_keys() done, return error count
    if error_count > 0:
        print("Test failed")
    return error_count
    
#
# test_extract(collection_name) tests extracting unique values form a collection based on a dot path
#
def test_extract(collection_name):
    '''test_extract() tests extracting unique values form a collection based on a dot path'''
    error_count = 0
    # Test extracting the family names
    v = dataset.extract(collection_name, 'true', '.authors[:].family')
    if not isinstance(v, list):
        print("Failed, expected a list, got", type(v), v)
        error_count += 1
        return error_count
    
    if len(v) != 4:
        print("Failed expected list to be of length 4, got", len(v))
        error_count += 1
    
    targets = [ "Austin", "Fremont", "Twain", "Verne" ]
    for s in targets:
        if s not in v:
            print("Failed, expected to find", s, "in", v)
            error_count += 1
    # test_extract() done, return error count
    if error_count > 0:
        print("Test failed")
    return error_count
    

#
# test_search(collection_name, index_map_name, index_name) tests indexer, deindexer and find funcitons
#
def test_search(collection_name, index_map_name, index_name):
    '''test indexer, deindexer and find functions'''
    error_count = 0
    dataset.verbose_on()
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
    
    ok = dataset.indexer(collection_name, index_name, index_map_name, batch_size = 2)
    if ok == False:
        print("Failed to index", collection_name)
        error_count += 1
    results = dataset.find(index_name, '+family:"Verne"')
    if results["total_hits"] != 2: 
        print("Warning: unexpected results", json.dumps(results, indent = 4))
    hits = results["hits"]
    
    k1 = hits[0]["id"]
    ok = dataset.deindexer(collection_name, index_name, [k1])
    if ok == False:
        print("deindexer failed for key", k1)
    # test_search(), done
    if error_count > 0:
        print("Test failed")
    return error_count 
    

#
# test_issue32() make sure issue 32 stays fixed.
#
def test_issue32(collection_name):
    error_count = 0
    ok = dataset.create(collection_name, "k1", {"one":1})
    if ok == False:
        print("Failed to create k1 in", collection_name)
        error_count += 1
        return error_count
    ok = dataset.has_key(collection_name, "k1")
    if ok == False:
        print("Failed, has_key k1 should return", True)
        error_count += 1
    ok = dataset.has_key(collection_name, "k2")
    if ok == True:
        print("Failed, has_key k2 should return", False)
        error_count += 1
    # test_issue32() done
    if error_count > 0:
        print("Test failed")
    return error_count
    
#
# test_gsheet(collection_name, setup_bash), if setup_bash exists run Google Sheets tests.
#
def test_gsheet(collection_name, setup_bash):
    '''if setup_bash exists, run Google Sheets tests'''
    if os.path.exists(setup_bash) == False:
        print("Skipping test_gsheet(", collection_name, setup_bash, ")")
        return 0
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    error_count = 0
    cfg = {}
    # read the environment settings from fname, turn into object.
    with open(setup_bash) as f:
        lines = f.readlines()
        for line in lines:
            if "export " in line:
                k,v = line.strip("export ").split("=", 2)
                k = k.strip("'\"\n ").lower()
                v = v.strip("'\"\n ")
                cfg[k] = v
    
    client_secret_name = ""
    sheet_id = ""
    if cfg.get("client_secret_json") == None:
        print("Failed, could not parse CLIENT_SECRET_JSON in", setup_bash, cfg)
        error_count += 1
        return error_count
    else:
        client_secret_name = cfg.get("client_secret_json")

    if cfg.get("spreadsheet_id") == None:
        print("Failed, could not parse SPREADSHEET_ID in", setup_bash)
        error_count += 1
        return error_count
    else:
        sheet_id = cfg.get("spreadsheet_id")
    client_secret_name = "../" + client_secret_name

    ok = dataset.init_collection(collection_name)
    if ok == False:
        print("Failed, could not create collection")
        error_count += 1
        return error_count

    cnt = dataset.count(collection_name)
    if cnt != 0:
        print("Failed to initialize a fresh collection", collection_name)
        error_count += 1
        return error_count

    # Setup some test data to work with.
    ok = dataset.create(collection_name, "Wilson1930",  {"additional":"Supplemental Files Information:\nGeologic Plate: Supplement 1 from \"The geology of a portion of the Repetto Hills\" (Thesis)\n","description_1":"Supplement 1 in CaltechDATA: Geologic Plate","done":"yes","identifier_1":"https://doi.org/10.22002/D1.638","key":"Wilson1930","resolver":"http://resolver.caltech.edu/CaltechTHESIS:12032009-111148185","subjects":"Repetto Hills, Coyote Pass, sandstones, shales"})
    if ok != True:
        print("Failed, could not create test record in", collection_name)
        error_count += 1
        return error_count
    
    cnt = dataset.count(collection_name)
    if cnt != 1:
        print("Failed, should have one test record in", collection_name)
        error_count += 1
        return error_count

    sheet_name = "Sheet1"
    cell_range = 'A1:Z'
    filter_expr = 'true'
    dot_exprs = ['.done','.key','.resolver','.subjects','.additional','.identifier_1','.description_1']
    column_names = ['Done','Key','Resolver','Subjects','Additional','Identifier 1','Description 1']
    print("Testing gsheet export support", sheet_id, sheet_name, cell_range, filter_expr, dot_exprs, column_names)
    ok = dataset.export_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, filter_expr, dot_exprs, column_names)
    if ok != True:
        print("Failed, count not export-gsheet in", collection_name)
        error_count += 1
        return error_count

    print("Testing gsheet import support (should fail)", sheet_id, sheet_name, cell_range, 2, False)
    dataset.verbose_off()
    ok = dataset.import_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, id_col = 2, overwrite = False)
    if ok == True:
        print("Failed, should NOT be able to import-gsheet over our existing collection without overwrite = True")
        error_count += 1
        return error_count

    print("Testing gsheet import support (should succeeed)", sheet_id, sheet_name, cell_range, 2, True)
    ok = dataset.import_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, id_col = 2, overwrite = True) 
    if ok == False:
        print("Failed, should be able to import-gsheet over our existing collection with overwrite=True")
        error_count += 1
        return error_count

    # Check to see if this throws error correctly, i.e. should have exit code 1
    sheet_name="Sheet2"
    dot_exprs = ['true','.done','.key','.QT_resolver','.subjects','.additional[]','.identifier_1','.description_1']
    ok = dataset.export_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, filter_expr, dot_exprs = dot_exprs)
    if ok == True:
        print("Failed, export_gsheet should throw error for bad dotpath in export_gsheet")
        error_count += 1
        return error_count

    # test_gsheet() done
    return error_count


# Setup our test collection, recreate it if necessary
def test_setup(collection_name):
    error_count = 0
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    ok = dataset.init_collection(collection_name)
    if ok == False:
        print("Failed, could not create collection")
        error_count += 1
        return error_count

    # test_gsheet() done, return error count
    return error_count

def test_check_repair(collection_name):
    error_count = 0
    print("Testing status on", collection_name)
    # Make sure we have a left over collection to check and repair
    ok = dataset.status(collection_name)
    if ok == False:
        print("Failed, expected dataset.status() == True, got", ok, "for", collection_name)
        error_count += 1
        return error_count

    print("Testing check on", collection_name)
    # Check our collection
    ok = dataset.check(collection_name)
    if ok == False:
        print("Failed, expected check", collection_name, "to return True, got", ok)
        error_count += 1

    # Break and recheck our collection
    if os.path.exists(collection_name + "/collection.json"):
        os.remove(collection_name + "/collection.json")
    print("Testing check on (broken)", collection_name)
    ok = dataset.check(collection_name)
    if ok == True:
        print("Failed, expected check", collection_name, "to return False, got", ok)
        error_count += 1

    # Repair our collection
    print("Testing repair on", collection_name)
    ok = dataset.repair(collection_name)
    if ok == False:
        print("Failed, expected repair to return True, got", ok)
        error_count += 1
    return error_count
        
def test_attachments(collection_name, filenames):
    print("Testing attach, attachments, detach and prune")
    error_count = 0
    ok = dataset.status(collection_name)
    if ok == False:
        print("Failed,", collection_name, "missing")
        error_count += 1
        return error_count
    keys = dataset.keys(collection_name)
    if len(keys) < 1:
        print("Failed,", collection_name, "should have keys")
        error_count += 1
        return error_count

    key = keys[0]
    ok = dataset.attach(collection_name, key, filenames)
    if ok == False:
        print("Failed, to attach files for", collection_name, key, filenames)
        error_count += 1
        return error_count

    return error_count

#
# Main processing
#
print("Starting dataset_test.py")
print("Testing dataset version", dataset.version())

# Pre-test check
error_count = 0
ok = True
dataset.verbose_off()

collection_name = "test_collection.ds"
error_count += test_setup(collection_name)
error_count += test_basic(collection_name)
error_count += test_keys(collection_name)
error_count += test_extract(collection_name)
error_count += test_search(collection_name, "test_index_map.json", "test_index.bleve")
error_count += test_issue32(collection_name)
error_count += test_gsheet("test_gsheet.ds", "../etc/test_gsheet.bash")
error_count += test_check_repair("test_gsheet.ds")
error_count += test_attachments(collection_name, ["README.md", "Makefile"])

print("Tests completed")

# Wrap up tests
if error_count > 0:
    print("Failed", error_count, "test(s)")
    sys.exit(1)
print("Success!")

