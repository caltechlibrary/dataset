#!/bin/bash

function softwareCheck() {
    for CMD in $@; do
        APP=$(which $CMD)
        if [ "$APP" = "" ]; then
            echo "Skipping, missing $CMD"
            exit 1;
        fi 
    done
}

function mkPage () {
    nav="$1"
    content="$2"
    html="$3"

    echo "Rendering $html"
    mkpage \
        "nav=$nav" \
        "content=$content" \
        page.tmpl > $html
    git add $html
}

echo "Checking software..."
softwareCheck mkpage
echo "Generating website"
mkPage nav.md README.md index.html
mkPage nav.md INSTALL.md install.html
mkPage nav.md "markdown:$(cat LICENSE)" license.html

# Build utility docs pages
for FNAME in dataset; do
    mkPage nav.md $FNAME.md $FNAME.html
done
