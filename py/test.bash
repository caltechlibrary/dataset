#!/bin/bash
if [[ -d "TestCollection.ds" ]]; then
    rm -fR "TestCollection.ds"
fi
python3 dataset.py "TestCollection.ds"
