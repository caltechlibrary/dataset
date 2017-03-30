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

echo "Checking software..."
softwareCheck mkpage
echo "Generating website"
MakePage nav.md README.md index.html
MakePage nav.md INSTALL.md install.html
MakePage nav.md "markdown:$(cat LICENSE)" license.html

# Build utility docs pages
GDD=$(which godocdown)
if [ "$GDD" != "" ]; then
    read -p "Overwrite docs/package.md from source code? Y/N " Y_OR_N
    if [ "$Y_OR_N" = "Y" ] || [ "$Y_OR_N" = "y" ]; then
        godocdown . > docs/package.md
    fi
fi
for FNAME in index package dataset; do
	MakePage "docs/nav.md" "docs/${FNAME}.md" "docs/${FNAME}.html"
done
