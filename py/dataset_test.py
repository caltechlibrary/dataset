#!/usr/bin/env python3
import sys
import os
import shutil
import json
import dataset

#
# test_basic(collection_name) runs tests on basic CRUD ops
# 
def test_basic(t, collection_name):
    '''test_basic(collection_name) runs tests on basic CRUD ops'''
    # Setup a test record
    key = "2488"
    value = { "title": "Twenty Thousand Leagues Under the Seas: An Underwater Tour of the World", "formats": ["epub","kindle","plain text"], "authors": [{ "given": "Jules", "family": "Verne" }], "url": "https://www.gutenberg.org/ebooks/2488"}
    
    # We should have an empty collection, we will create our test record.
    ok = dataset.create(collection_name, key, value)
    if ok == False:
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
    ok = dataset.update(collection_name, key, value)
    if ok == False:
       t.error(f"Failed, count not update record {key}, {value}")
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
    expected_s = collection_name+"/aa/"+key+".json"
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
    ok = dataset.delete(collection_name, key)
    if ok == False:
        t.error("Failed, could not delete record", key)
    

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
            t.error("Failed, could not add", k, "to", collection_name)
    
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
# test_extract(collection_name) tests extracting unique values form a collection based on a dot path
#
def test_extract(t, collection_name):
    '''test_extract() tests extracting unique values form a collection based on a dot path'''
    # Test extracting the family names
    v = dataset.extract(collection_name, 'true', '.authors[:].family')
    if not isinstance(v, list):
        t.error("Failed, expected a list, got", type(v), v)
        return
    
    if len(v) != 4:
        t.error("Failed expected list to be of length 4, got", len(v))
    
    targets = [ "Austin", "Fremont", "Twain", "Verne" ]
    for s in targets:
        if s not in v:
            t.error("Failed, expected to find", s, "in", v)

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
    
    ok = dataset.indexer(collection_name, index_name, index_map_name, batch_size = 2)
    if ok == False:
        t.error("Failed to index", collection_name)
    results = dataset.find(index_name, '+family:"Verne"')
    if results["total_hits"] != 2: 
        t.print("Warning: unexpected results", json.dumps(results, indent = 4))
    hits = results["hits"]
    
    k1 = hits[0]["id"]
    ok = dataset.deindexer(collection_name, index_name, [k1])
    if ok == False:
        t.print("deindexer failed for key", k1)

#
# test_issue32() make sure issue 32 stays fixed.
#
def test_issue32(t, collection_name):
    ok = dataset.create(collection_name, "k1", {"one":1})
    if ok == False:
        t.error("Failed to create k1 in", collection_name)
        return
    ok = dataset.has_key(collection_name, "k1")
    if ok == False:
        t.error("Failed, has_key k1 should return", True)
    ok = dataset.has_key(collection_name, "k2")
    if ok == True:
        t.error("Failed, has_key k2 should return", False)

#
# test_gsheet(t, collection_name, setup_bash), if setup_bash exists run Google Sheets tests.
#
def test_gsheet(t, collection_name, setup_bash):
    '''if setup_bash exists, run Google Sheets tests'''
    if os.path.exists(setup_bash) == False:
        t.verbose_on()
        t.print("Skipping test_gsheet(", collection_name, setup_bash, ")")
        return
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
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
        t.error("Failed, could not parse CLIENT_SECRET_JSON in", setup_bash, cfg)
        return
    else:
        client_secret_name = cfg.get("client_secret_json")

    if cfg.get("spreadsheet_id") == None:
        t.error("Failed, could not parse SPREADSHEET_ID in", setup_bash)
        return
    else:
        sheet_id = cfg.get("spreadsheet_id")
    client_secret_name = "../" + client_secret_name

    ok = dataset.init(collection_name)
    if ok == False:
        t.error("Failed, could not create collection")
        return

    cnt = dataset.count(collection_name)
    if cnt != 0:
        t.error("Failed to initialize a fresh collection", collection_name)
        return

    # Setup some test data to work with.
    ok = dataset.create(collection_name, "Wilson1930",  {"additional":"Supplemental Files Information:\nGeologic Plate: Supplement 1 from \"The geology of a portion of the Repetto Hills\" (Thesis)\n","description_1":"Supplement 1 in CaltechDATA: Geologic Plate","done":"yes","identifier_1":"https://doi.org/10.22002/D1.638","key":"Wilson1930","resolver":"http://resolver.caltech.edu/CaltechTHESIS:12032009-111148185","subjects":"Repetto Hills, Coyote Pass, sandstones, shales"})
    if ok != True:
        t.error("Failed, could not create test record in", collection_name)
        return
    
    cnt = dataset.count(collection_name)
    if cnt != 1:
        t.error("Failed, should have one test record in", collection_name)
        return

    sheet_name = "Sheet1"
    cell_range = 'A1:Z'
    filter_expr = 'true'
    dot_exprs = ['.done','.key','.resolver','.subjects','.additional','.identifier_1','.description_1']
    column_names = ['Done','Key','Resolver','Subjects','Additional','Identifier 1','Description 1']
    t.print("Testing gsheet export support", sheet_id, sheet_name, cell_range, filter_expr, dot_exprs, column_names)
    ok = dataset.export_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, filter_expr, dot_exprs, column_names)
    if ok != True:
        t.error("Failed, count not export-gsheet in", collection_name)
        return

    t.print("Testing gsheet import support (should fail)", sheet_id, sheet_name, cell_range, 2, False)
    dataset.verbose_off()
    ok = dataset.import_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, id_col = 2, overwrite = False)
    if ok == True:
        t.error("Failed, should NOT be able to import-gsheet over our existing collection without overwrite = True")
        return

    t.print("Testing gsheet import support (should succeeed)", sheet_id, sheet_name, cell_range, 2, True)
    ok = dataset.import_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, id_col = 2, overwrite = True) 
    if ok == False:
        t.error("Failed, should be able to import-gsheet over our existing collection with overwrite=True")
        return

    # Check to see if this throws error correctly, i.e. should have exit code 1
    dataset.use_strict_dotpath(True)
    sheet_name="Sheet1"
    dot_exprs = ['true','.done','.key','.QT_resolver','.subjects','.additional[]','.identifier_1','.description_1']
    ok = dataset.export_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, filter_expr, dot_exprs = dot_exprs)
    if ok == True:
        t.error("Failed, export_gsheet should throw error for bad dotpath in export_gsheet")
    #dataset.verbose_on()
    dataset.use_strict_dotpath(False)
    sheet_name = "Sheet1"
    ok = dataset.export_gsheet(collection_name, client_secret_name, sheet_id, sheet_name, cell_range, filter_expr, dot_exprs = dot_exprs)
    if ok == False:
        t.error("Failed, export_gsheet should only warn of error for bad dotpath in export_gsheet")
    #dataset.verbose_off()

# Setup our test collection, recreate it if necessary
def test_setup(t, collection_name):
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    ok = dataset.init(collection_name)
    if ok == False:
        t.error("Failed, could not create collection")
        return


def test_check_repair(t, collection_name):
    t.print("Testing status on", collection_name)
    # Make sure we have a left over collection to check and repair
    if os.path.exists(collection_name) == False:
        dataset.init(collection_name)
    ok = dataset.status(collection_name)
    if ok == False:
        t.error("Failed, expected dataset.status() == True, got", ok, "for", collection_name)
        return

    t.print("Testing check on", collection_name)
    # Check our collection
    ok = dataset.check(collection_name)
    if ok == False:
        t.error("Failed, expected check", collection_name, "to return True, got", ok)

    # Break and recheck our collection
    if os.path.exists(collection_name + "/collection.json"):
        os.remove(collection_name + "/collection.json")
    t.print("Testing check on (broken)", collection_name)
    ok = dataset.check(collection_name)
    if ok == True:
        t.error("Failed, expected check", collection_name, "to return False, got", ok)

    # Repair our collection
    t.print("Testing repair on", collection_name)
    ok = dataset.repair(collection_name)
    if ok == False:
        t.error("Failed, expected repair to return True, got", ok)
 
        
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
    ok = dataset.attach(collection_name, key, filenames)
    if ok == False:
        t.error("Failed, to attach files for", collection_name, key, filenames)
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
    ok = dataset.detach(collection_name, key, [filenames[1]])
    if ok == False:
        t.error("Failed, expected True for", collection_name, key, filenames[1])
    if os.path.exists(filenames[1]):
        os.remove(filenames[1])
    else:
        t.error("Failed to detch", filenames[1], "from", collection_name, key)

    # Test explicit filenames detch
    ok = dataset.detach(collection_name, key, filenames)
    if ok == False:
        t.error("Failed, expected True for", collection_name, key, filenames)

    for fname in filenames:
        if os.path.exists(fname):
            os.remove(fname)
        else:
            t.error("Failed, expected", fname, "to be detached from", collection_name, key)

    # Test detaching all files
    ok = dataset.detach(collection_name, key)
    if ok == False:
        t.error("Failed, expected True for (detaching all)", collection_name, key)
    for fname in filenames:
        if os.path.exists(fname):
            os.remove(fname)
        else:
            t.error("Failed, expected", fname, "for detaching all from", collection_name, key)

    ok = dataset.prune(collection_name, key, [filenames[0]])
    if ok == False:
        t.error("Failed, expected True for prune", collection_name, key, [filenames[0]])
    l = dataset.attachments(collection_name, key)
    if len(l) != 1:
        t.error("Failed, expected one file after prune for", collection_name, key, [filenames[0]], "got", l)

    ok = dataset.prune(collection_name, key)
    if ok == False:
        t.error("Failed, expected True for prune (all)", collection_name, key)
    l = dataset.attachments(collection_name, key)
    if len(l) != 0:
        t.error("Failed, expected zero files after prune for", collection_name, key, "got", l)


def test_s3(t):
    aws_sdk_load_config = os.getenv("AWS_SDK_LOAD_CONFIG", "")
    collection_name = os.getenv("DATASET", "")
    if aws_sdk_load_config != "1" or collection_name[0:5] != "s3://":
        t.verbose_on()
        t.print("Skipping test_s3(), missing environment AWS_SDK_LOAD_CONFIG and DATASET")
        return
    
    ok = dataset.status(collection_name)
    if ok == False:
        t.print("Missing", collection_name, "attempting to initialize", collection_name)
        ok = dataset.init(collection_name)
        if ok == False:
            t.error("Aborting, couldn't initialize", collection_name)
            return
    else:
        t.print("Using collection initialized as", collection_name)

    collection_name = os.getenv("DATASET")
    record = { "one": 1 }
    key = "s3t1"
    ok = dataset.create(collection_name, key, record)
    if ok == False:
        t.error("Failed to create record", collection_name, key, record)
    record2, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    if record2.get("one") != 1:
        t.error("Failed, read", collection_name, key, record2)
    record["two"] = 2
    ok = dataset.update(collection_name, key, record)
    if ok == False:
        t.error("Failed to update record", collection_name, key, record)
    record2, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    if record2.get("one") != 1:
        t.error("Failed, 2nd read", collection_name, key, record2)
    if record2.get("two") != 2:
        t.error("Failed, 2nd read", collection_name, key, record2)
    ok = dataset.delete(collection_name, key)
    if ok == False:
        t.error("Failed to delete record", collection_name, key, record)
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
    if ok == True:
        ok = dataset.update(collection_nane, key, obj1)
    else:
        ok = dataset.create(collection_name, key, obj1)
    if ok == False:
        t.error("Failed, could not add record for test", collection, key, obj1)
        return
    ok = dataset.join(collection_name, key, "append", obj2)
    if ok == False:
        t.error("Failed, join for", collection_name, key, "append", obj2)
    obj_result, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    if obj_result.get("one") != 1:
        t.error("Failed to join append key", key, obj_result)
    if obj_result.get("two") != 2:
        t.error("Failed to join append key", key, obj_result)
    obj2["one"] = 3
    obj2["two"] = 3
    obj2["three"] = 3
    ok = dataset.join(collection_name, key, "overwrite", obj2)
    if ok == False:
        t.error("Failed to join overwrite", collection_name, key, "overwrite", obj2)
    obj_result, err = dataset.read(collection_name, key)
    if err != "":
        t.error(f"Unexpected error for {key} in {collection_name}, {err}")
    for k in obj_result:
        if k != "_Key" and obj_result[k] != 3:
            t.error("Failed to update value in join overwrite", k, obj_result)
    ok = dataset.join(collection_name, key, "fred and mary", obj2)
    if ok == True:
        t.error("Failed, expected error for join type 'fred and mary'")
    
#
# test_issue43() When exporting records to a table using
# use_srict_dotpath(True), the rows are getting miss aligned.
#
def test_issue43(t, collection_name, csv_name):
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    if os.path.exists(csv_name):
        os.remove(csv_name)
    ok = dataset.init(collection_name)
    if ok == False:
        t.error(f"Failed, need a {collection_name} to run test")
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
        ok = dataset.create(collection_name, key, row)
        if ok == False:
            t.error(f"Can't add test row {key} to {collection_name}")
            return
    #dataset.verbose_on()
    dataset.use_strict_dotpath(False)
    ok = dataset.export_csv(collection_name, csv_name, "true", ["._Key",".c1",".c2",".c3",".c4"])
    if ok == False:
       t.error(f"csv_export({collection_name}, {csv_name} should have emitted warnings, not error")
       return
    with open(csv_name, mode = "r", encoding = "utf-8") as f:
        rows = f.read()

    for row in rows.split("\n"):
        if len(row) > 0:
            cells = row.split(",")
            if len(cells) < 5:
                t.error(f"row error {csv_name} for {cells}")


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

    # Pre-test check
    error_count = 0
    ok = True
    dataset.verbose_off()

    collection_name = "test_collection.ds"
    test_runner = TestRunner(os.path.basename(__file__))
    test_runner.add(test_setup, [collection_name])
    test_runner.add(test_basic, [collection_name])
    test_runner.add(test_keys, [collection_name])
    test_runner.add(test_extract, [collection_name])
    test_runner.add(test_search, [collection_name, "test_index_map.json", "test_index.bleve"])
    test_runner.add(test_issue32, [collection_name])
    test_runner.add(test_attachments, [collection_name])
    test_runner.add(test_join, [collection_name])
    test_runner.add(test_check_repair, ["test_check_and_repair.ds"])
    test_runner.add(test_gsheet, ["test_gsheet.ds", "../etc/test_gsheet.bash"])
    test_runner.add(test_issue43,["test_issue43.ds", "test_issue43.csv"])
    test_runner.add(test_s3)
    test_runner.run()

