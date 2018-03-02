#!/bin/bash

PKGASSETS=$(which pkgassets)
if [ "$PKGASSETS" = "" ]; then
    cat <<EOT >&2
You need to have pkgassets installed and in your path.
To install pkgassets try:

    go get -u github.com/caltechlibrary/pkgassets/... 

EOT
    exit 1
fi

function buildHelp() {
    PROG="$1"
    pkgassets -o "cmd/${PROG}/assets.go" -p main \
        -ext=".md" -strip-prefix="/" -strip-suffix=".md" \
        Examples "examples/${PROG}" \
        Help "docs/${PROG}"
    git add "cmd/${PROG}/assets.go"
}

# build Help assets 
buildHelp dataset
buildHelp dsws

# build Template assets
pkgassets -o cmd/dsws/templates.go -p main Defaults defaults
git add cmd/dsws/templates.go

