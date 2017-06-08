#!/bin/bash

echo "Started $(date)"
if [ ! -f dir-listing.txt ]; then
    echo "Generating the dir-listing.txt for CaltechAUTHORS...."
    aws s3 ls --recursive s3://dataset.library.caltech.edu/CaltechAUTHORS/ > dir-listing.txt
else
    echo "Using existing creation of dir-listing.txt"
fi

echo "Generating collection-fixed.json"
echo -n '{"version":"v0.0.3","name":"CaltechAUTHORS","buckets":[' > collection-fixed.json
cat dir-listing.txt | cut -d/ -f 2 | sort -u | while read ITEM; do
    name=$(basename ${ITEM} ".json")
    if [ "${name}" = "${ITEM}" ]; then
        echo -n  "\"${ITEM}\","; 
    fi
done | sed -E 's/,$//' >> collection-fixed.json
echo -n '],"keymap":{' >> collection-fixed.json
cat dir-listing.txt | cut -d/ -f 2,3 | sed -E 's/.json//;s/.tar//' | sort -u | while read ITEM; do 
    if [ "${ITEM}" != "" ]; then
        ky=$(echo "${ITEM}" | cut -d/ -f 1)
        bucket=$(echo "${ITEM}" | cut -d/ -f 2)
        if [ "${ky}" != "" ] && [ "${bucket}" != "" ]; then
            echo -n  "\"${ky}\":\"${bucket}\","; 
        fi
    fi
done | sed -E 's/,$//' >> collection-fixed.json
echo -n '},"select_lists":["keys"],"index_defs":null}' >> collection-fixed.json

echo "Generating keys-fixed.json"
echo -n '[' > keys-fixed.json
cat dir-listing.txt | cut -d/ -f 3 | sed -E 's/.json//;s/.tar//' | sort -u | while read ITEM; do 
    if [ "${ITEM}" != "" ]; then
        echo -n  "\"${ITEM}\","; 
    fi
done | sed -E 's/,$//' >> keys-fixed.json
echo -n ']' >> keys-fixed.json
echo "Done! $(date)"
