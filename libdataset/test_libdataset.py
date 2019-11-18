from libdataset import dataset 
import os
import sys
import shutil
import json

# Alias this to dataset
#dataset = libdataset

# Clean up stale result test collections
c_name = "t1.ds"
if os.path.exists(c_name):
    shutil.rmtree(c_name)

err = dataset.init(c_name)
if err != "":
    print(f"expected '', got '{err}' for dataset.init({c_name})")
    sys.exit(1)

src = '''
[
    {"_Key": "k1", "title": "One thing", "name": "Fred"},
    {"_Key": "k2", "title": "Two things", "name": "Frieda"},
    {"_Key": "k3", "title": "Three things", "name": "Fiona"}
]
'''

for obj in json.loads(src):
    key = obj['_Key']
    err = dataset.create(c_name, key, obj)
    if err != '':
        print(f"expected '', got '{err}' for dataset.create({c_name}, {_key}, {obj})")
        sys.exit(1)

expected_keys = [ "k1", "k2", "k3" ]
keys = dataset.keys(c_name)
i = 0
for key in keys:
    if not key in expected_keys:
        print(f"expected {key} in {expected_keys} for {c_name}")
        sys.exit(1)
    obj, err = dataset.read(c_name, key)
    if err != '':
        print(f"expected '', got '{err}' for dataset.read({c_name}, {key}")
        sys.exit(1)
    obj['t_count'] = i
    i += 1
    err = dataset.update(c_name, key, obj)
    if err != '':
        print(f"expected '', got '{err}' for dataset.update({c_name}, {key}, ...")
        sys.exit(1)

f_name = "f1"
err = dataset.frame_create(c_name, f_name, keys[1:], ['._Key', '.title'], [ 'id', 'title' ])
if err != '':
        print(f"expected '', got '{err}' for dataset.frame_create({c_name}, {f_name}, ...)")
        sys.exit(1)

ok = dataset.frame_exists(c_name, f_name)
if ok != True:
        print(f"expected 'True', got '{ok}' for dataset.frame_exists({c_name}, {f_name})")
        sys.exit(1)


expected_keys = keys[1:]
keys = dataset.frame_keys(c_name, f_name)
for i, expected in enumerate(expected_keys):
    key = keys[i]
    if key != expected:
        print(f"expected ({i}) '{expected}', got '{key}' for dataset.frame_keys({c_name}, {f_name})")
        sys.exit(1)



print('Success!')
shutil.rmtree("t1.ds")
