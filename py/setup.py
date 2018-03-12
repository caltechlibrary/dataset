#!/usr/bin/env python3
from distutils.core import setup

#FIXME: I probably want to use setuptools rather than distutils, my compiled library isn't getting copied to the right place.

import sys
import os
import shutil
import json

codemeta_json = "codemeta.json"
# If we're running sdist make sure our local codemeta.json is up to date!
if "sdist" in sys.argv:
    # Project Metadata 
    shutil.copyfile(os.path.join("..", codemeta_json),  codemeta_json)

# Let's pickup as much metadata as we need from codemeta.json
with open(codemeta_json, mode = "r", encoding = "utf-8") as f:
    src = f.read()
    meta = json.loads(src)

# Let's make our symvar string
version = "v"+meta["version"]

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
if platform.startswith("Darwin"):
    shared_library_name = "libdataset.dylib"
    platform = "Mac OS X"
elif platform.startswith("Win"):
    shared_library_name = "libdataset.dll"
    platform = "Windows"


# Now that we know everything configure out setup
setup(name = "dataset",
    version = version,
    description = "A python module for managing with JSON docs on local disc, in cloud storage",
    long_description = """This module wraps the functionality available from the Go base command line tool developed
at Caltech Library called dataset. The module, like the tool, supports working 
with collections of JSON documents on your local disc as well as in the cloud 
(e.g. AWS S3, Google Cloud Storage). It can also import/export to CSV files or 
Google Spreadsheets. In addition to managing JSON documents dataset also
supports full search and indexing.""",
    author = author,
    author_email = author_email,
    url = "https://caltechlibrary.github.io/dataset",
    download_url = "https://github.com/caltechlibrary/dataset/latest/releases",
    license = meta["license"],
    packages = ["dataset"],
    data_files = [
        ("", [shared_library_name, codemeta_json])
    ],
    platforms = [platform],
    keywords = ["JSON", "CSV", "data science", "storage"],
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
        "Operating System :: MacOS :: MacOS X",
        #"Operating System :: POSIX :: Linux",
        "Natural Language :: English",
    ]
)
