#!/usr/bin/env python3

#from distutils.core import setup
#from site import getsitepackages
#site_package_location = os.path.join(getsitepackages()[0], "dataset")

from setuptools import setup, find_packages

import sys
import os
import shutil
import json

readme_md = "README.md"
readme_txt = "README.txt"

def read(fname):
    with open(fname, mode = "r", encoding = "utf-8") as f:
        src = f.read()
    return src

codemeta_json = "codemeta.json"
# If we're running sdist make sure our local codemeta.json is up to date!
if "sdist" in sys.argv:
    # Project Metadata and README
    shutil.copyfile(os.path.join("..", codemeta_json),  codemeta_json)
    shutil.copyfile(os.path.join("..", readme_md),  readme_txt)

# Let's pickup as much metadata as we need from codemeta.json
with open(codemeta_json, mode = "r", encoding = "utf-8") as f:
    src = f.read()
    meta = json.loads(src)

# Let's make our symvar string
version = meta["version"]
#version = meta["version"]

# Now we need to pull and format our author, author_email strings.
author = ""
author_email = ""
for obj in meta["author"]:
    given = obj["givenName"]
    family = obj["familyName"]
    email = obj["email"]
    if len(author) == 0:
        author = f"{given} {family}"
    else:
        author = author + f", {given} {family}"
    if len(author_email) == 0:
        author_email = f"{email}"
    else:
        author_email = author_email + f", {email}"

# Setup for our Go based shared library as a "data_file" since Python doesn't grok Go.
platform = os.uname().sysname
shared_library_name = "libdataset.so"
OS_Classifier = "Operating System :: POSIX :: Linux"
if platform.startswith("Darwin"):
    shared_library_name = "libdataset.dylib"
    platform = "Mac OS X"
    OS_Classifier = "Operating System :: MacOS :: MacOS X"
elif platform.startswith("Win"):
    shared_library_name = "libdataset.dll"
    platform = "Windows"
    OS_Classifier = "Operating System :: Microsoft :: Windows :: Windows 10"
        
if os.path.exists(os.path.join("dataset", shared_library_name)) == False:
    print(f"Missing compiled shared library {shared_library_name} in dataset module")
    sys.exit(1)

# Now that we know everything configure out setup
setup(name = "dataset",
    version = version,
    description = "A python module for managing with JSON docs on local disc, in cloud storage",
    long_description = read(readme_txt),
    author = author,
    author_email = author_email,
    url = "https://caltechlibrary.github.io/dataset",
    download_url = "https://github.com/caltechlibrary/dataset/latest/releases",
    license = meta["license"],
    packages = find_packages(exclude=["*.tests", "*.tests.*", "tests.*", "tests", "*_test.py"]),
    package_data = {
        '': [ '*.txt', '*.so', '*.dll', '*.dylib'],
    },
    platforms = [platform],
    keywords = ["JSON", "CSV", "data science", "storage"],
    include_package_data = True,
    classifiers = [
        "Development Status :: Alpha",
        "Environment :: Console",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Programming Language :: Other",
        "Programming Language :: Go",
        "Intended Audience :: Science/Research",
        "Topic :: Scientific/Engineering",
        "License :: OSI Approved :: BSD License",
        OS_Classifier
    ]
)
