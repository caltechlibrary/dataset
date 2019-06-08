#!/usr/bin/env python3

#
# This is an example Python script that migrates pre-v0.0.62 attachments in a tarball
# two dataset v0.0.62 attachment scheme. It doesn't make an assumption of which py_dataset
# version is available other than to read/update the JSON object. The attachment process
# is done by the dataset cli installed which must be v0.0.62 or later.
#
import os
import sys
import tarfile
from subprocess import Popen, PIPE, run
from py_dataset import dataset

def reattach(c_name, key, semver, f_name):
    cmd = [ "dataset", "attach", c_name, key, semver, f_name ]
    with Popen(cmd, stdout = PIPE, stderr = PIPE) as proc:
        err = proc.stderr.read().strip().decode('utf-8')
        if err != '':
            print(f"Failed {' '.join(cmd)} {err}")
            sys.exit(1)

def migrate_attachment(c_name, key):
    obj, err = dataset.read(c_name, key)
    obj_path = dataset.path(c_name, key).replace(key + ".json", "")
    tarball = os.path.join(obj_path, key + ".tar")
    if os.path.exists(tarball):
        tar = tarfile.open(tarball)
        tar.extractall()
        tar.close()
        files = os.listdir()
        # Prune _Attachment from object and resave
        if "_Attachments" in obj:
            del obj["_Attachments"]
            err = dataset.update(c_name, key, obj)
            if err != "":
                print(f"Can't remove _Attachments metadata, {err}")
                sys.exit(1)
        for fname in files:
            print(".", end = "")
            reattach(c_name, key, "v0.0.0", fname)
            os.remove(fname)
        # NOTE: if all re-attached then we need to remove tarball too
        os.remove(tarball)
        sys.stdout.flush()


if len(sys.argv) == 1:
    app = os.path.basename(sys.argv[0])
    print(f"USAGE: {app} DATASET_NAME", end = "\n\n")
    print("Converts attachments in a dataset from tarballs to v0.0.62 attachment scheme", end = "\n\n")
    sys.exit(0)

if not os.path.exists("tmp-attachment-migration"):
    os.mkdir("tmp-attachment-migration")
os.chdir("tmp-attachment-migration")
print(f"Working directory for migration is {os.getcwd()}")
for c_name in sys.argv:
    keys = dataset.keys(os.path.join("..", c_name))
    tot = len(keys)
    print(f"Ready to process {tot} objects")
    for i, key in enumerate(keys):
        migrate_attachment(os.path.join("..", c_name), key)
        if (i % 500) == 0:
            print(f"\n{i+1} of {tot} processed")
    print()
    print(f"Procssing {c_name} complete")
