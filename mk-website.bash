#!/bin/bash

function softwareCheck() {
	for CMD in "$@"; do
		APP=$(which "$CMD")
		if [ "$APP" = "" ]; then
			echo "Skipping, missing $CMD"
			exit 1
		fi
	done
}

function MakePage() {
	nav="$1"
	content="$2"
	html="$3"

	echo "Rendering $html"
	mkpage \
		"nav=$nav" \
		"content=$content" \
		page.tmpl >"$html"
	git add "$html"
}

function MakeSubPages() {
    SUBDIR="${1}"
    find "${SUBDIR}" -type f | grep -E '\.md$' | while read FNAME; do
        FNAME="$(basename "${FNAME}" ".md")"
        if [ -f "${SUBDIR}/${FNAME}.md" ] && [ "$FNAME" != "nav" ]; then
	        MakePage "${SUBDIR}/nav.md" "${SUBDIR}/${FNAME}.md" "${SUBDIR}/${FNAME}.html"
        fi
    done
}

function MakeSubSubPages() {
    ASSET_FOLDER="${1}"
    if [ ! -f "${ASSET_FOLDER}/index.md" ]; then
        echo "Creating index.md for asset: ${ASSET_FOLDER}"
        # Create an index.md for the asset list
        cat <<EOT >"${ASSET_FOLDER}/index.md"

# ${ASSET_FOLDER}

EOT

        echo "Creating nav.md for asset: ${ASSET_FOLDER}"
        # Create some nav
        cat <<EOT >"${ASSET_FOLDER}/nav.md"
+ [Home](/)
+ [Up](../)
+ [${ASSET_FOLDER}](./)
EOT

    fi

    find "${ASSET_FOLDER}" -type d -depth 1 | sort | while read DNAME; do
        T="$(basename "${DNAME}")"
        echo "Creating nav.md for asset: ${ASSET_FOLDER}/${T}"
        cat <<EOT >"${ASSET_FOLDER}/${T}/nav.md"
+ [Home](/)
+ [Up](../)
+ [${T}](./)
EOT
        # Generate the index of topics for the cmd described in asset folder
        MakeSubPages "${ASSET_FOLDER}/${T}"
    done
    echo "" >>"${ASSET_FOLDER}/index.md"
}

echo "Checking software..."
softwareCheck mkpage
echo "Generating website"
MakePage nav.md README.md index.html
MakePage nav.md INSTALL.md install.html
MakePage nav.md "markdown:$(cat LICENSE)" license.html

# Build example pages
MakeSubSubPages examples
MakeSubPages examples

# Build utility docs pages
MakeSubSubPages docs
MakeSubPages docs

# Build how-to pages
MakeSubPages how-to
