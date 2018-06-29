#!/bin/bash

function cleanUpHTML() {
	findfile -s ".html" . | while read P; do
		rm "$P"
	done
}

function FindNavMD() {
	DNAME="$1"
	if [ -f "${DNAME}/nav.md" ]; then
		echo "${DNAME}/nav.md"
		return
	fi
	DNAME=$(dirname "${DNAME}")
	FindNavMD "${DNAME}"
}

# Cleanup stale HTML files
cleanUpHTML

# Look through files and build new site
mkpage "nav=nav.md" "content=markdown:$(cat LICENSE)" page.tmpl >license.html
mkpage "nav=nav.md" "content=markdown:$(cat INSTALL.md)" page.tmpl >install.html
git add -f license.html install.html
findfile -s ".md" . | while read P; do
	DNAME=$(dirname "$P")
	if [[ "${DNAME:0:3}" = "etc" || "${DNAME:0:5}" = "demos" || ${DNAME:0:8} = "examples" ]]; then
		echo "Skipping $DNAME"
	else
		FNAME=$(basename "$P")
		case "$FNAME" in
		"README.md")
			if [ ! -f "${DNAME}/index.md" ]; then
				HTML_NAME="${DNAME}/index.html"
			else
				HTML_NAME="${DNAME}/README.html"
			fi
			;;
        "INSTALL.md")
            HTML_NAME="${DNAME}/install.html"
            ;;
		*)
			HTML_NAME=$(echo "$P" | sed -E 's/.md$/.html/g')
			;;
		esac
		if [ "${DNAME:0:4}" != "dist" ] && [ "${FNAME}" != "nav.md" ]; then
			NAV=$(FindNavMD "$DNAME")
			echo "Building $HTML_NAME from $DNAME/$FNAME and $NAV"
			mkpage "nav=$NAV" "content=${DNAME}/${FNAME}" page.tmpl >"${HTML_NAME}"
			git add -f "${HTML_NAME}"
		else
			echo "Skipping $P"
		fi
	fi
done

# Cleanup some files we don't want to publish
if [ -f "NOTES.html" ]; then
	git rm -f "NOTES.html"
fi
if [ -f "DEVELOPMENT-NOTES.html" ]; then
	git rm -f "DEVELOPMENT-NOTES.html"
fi
if [ -f "README.html" ]; then
	git rm -f "README.html"
fi
if [ -f "TODO.html" ]; then
	git rm -f "TODO.html"
fi
if [ -f "RELEASES.html" ]; then
	git rm -f "RELEASES.html"
fi
