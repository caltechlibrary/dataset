#!/bin/bash

#
# Crawl the included Caltech Library package and display their versions in GOPATH
#
function crawl_code() {
grep "${1}" $(findfile -s go .) | cut -d \" -f 2 | sort -u | while read PNAME; do 
    echo -n "$PNAME -- ";
    V=$(grep 'Version = `' "$GOPATH/src/$PNAME/$(basename $PNAME).go" | cut -d \` -f 2)
    if [ "$V" = "" ]; then
        echo "Unknown"
    else
        echo "$V"
    fi
done
}

echo "The following package versions were used in this release"
crawl_code "github.com/caltechlibrary/"
crawl_code "github.com/rsdoiel/"
