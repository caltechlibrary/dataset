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
	output="$3"

	echo "Rendering $output"
	mkpage \
		"nav=$nav" \
		"content=$content" \
		page.tmpl >"$output"
	git add "$output"
}

function MakeSubPages() {
		SUBDIR="$1"
    findfile -s .md "${SUBDIR}" | while read FNAME; do
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

    finddir -depth 2  "${ASSET_FOLDER}" | sort | while read DNAME; do
        T="$(basename "${DNAME}")"
        # Generate nav and topdics for folder
				echo "Scanning for topics in: ${ASSET_FOLDER}/${T}"
				cat <<EOT >"${ASSET_FOLDER}/${T}/topics.md"

# Topics

EOT

				findfile -s .md "${ASSET_FOLDER}/${T}" | sort | while read FNAME; do
						LABEL="$(basename "$FNAME" ".md")"
						case "${LABEL}" in
							"nav")
							;;
							"commands")
							;;
							"description")
							;;
							"usage")
							;;
							"index")
							;;
							"topics")
							;;
							*)
							if [ -f "${ASSET_FOLDER}/${T}/${LABEL}.md" ]; then
									echo "+ [${LABEL}](${LABEL}.html)"
							fi
							;;
						esac
				done >>"${ASSET_FOLDER}/${T}/topics.md"

				# Generating "Here" value for folder
				HERE="${DNAME}"
				if [ "${HERE}" = '.' ] && [ "$ASSET_FOLDER" = "docs" ]; then
						HERE="Documentation"
				elif [ "${HERE}" = '.' ]; then
						HERE="${ASSET_FOLDER}"
				fi

				C="$(wc -l "${ASSET_FOLDER}/${T}/topics.md" | cut -d\  -f 1)"
				if [ "${C}" != "3" ] ; then
					echo "Creating nav.md with topic links for: ${ASSET_FOLDER}/${T}"
	        cat <<EOT >"${ASSET_FOLDER}/${T}/nav.md"
+ [Home](/)
+ [Up](../)
+ [${HERE}](./)
+ [topics](topics.html)
EOT
				else
					echo "Creating nav.md without topics for asset: ${ASSET_FOLDER}/${T}"
	        cat <<EOT >"${ASSET_FOLDER}/${T}/nav.md"
+ [Home](/)
+ [Up](../)
+ [${HERE}](./)
EOT
				fi

				# Generate the index of topics for the cmd described in asset folder
        MakeSubPages "${ASSET_FOLDER}/${T}"

    done
    echo "" >>"${ASSET_FOLDER}/index.md"
}

echo "Checking software..."
softwareCheck mkpage finddir findfile
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
