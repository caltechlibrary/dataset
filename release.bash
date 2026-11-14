#!/bin/bash
#
# release.bash creates a draft GitHub release from dist/*.zip. Run `make
# release` first: it regenerates the version files, builds the zip files and
# checks them (release-check), so by the time this script runs the checks have
# already passed.
#
# Usage: ./release.bash [--dry-run] [--yes]
#
# The release tag and notes come from codemeta.json (version, releaseNotes).
# The script refuses to run if anything is uncommitted other than the files
# `make release` and the author's version bump change (PREP_FILES below), or if
# commits are not pushed. After you confirm, it commits those files as
# "prep for vX.Y.Z", pushes, writes a checksums file next to the zip files and
# creates a DRAFT release carrying the zips and the checksums. Publish the draft
# from the repository's releases page.
#
# Options:
#   --dry-run   run every check and print what would happen; create nothing
#   --yes       do not ask for confirmation
#   -h, --help  this text
#
# Needs: bash, git, jq, gh, and sha256sum or shasum.
#
# EXIT STATUS
#   0   draft release created (or --dry-run checks passed)
#   1   a check refused: uncommitted changes beyond the release files, unpushed
#       commits, the tag points elsewhere, the release already exists, or the
#       prompt was declined
#   2   usage: unknown option or surplus argument
#   65  codemeta.json has no usable version (missing or unexpected characters)
#   66  an input is missing: codemeta.json, dist/*.zip, or a required command
#   69  git could not push, or gh could not create the release
#   74  git could not commit the release files
#
set -euo pipefail

E_NEGATIVE=1 E_USAGE=2 E_DATA=65 E_NO_INPUT=66 E_UNAVAILABLE=69 E_IO=74

# The files a release legitimately changes: the author's version bump and what
# `make release` regenerates. Anything else uncommitted stops the release.
PREP_FILES="codemeta.json version.go CITATION.cff installer.sh installer.ps1"
is_prep_file() { # is_prep_file PATH
	case "$1" in
	docs/*.?.md) return 0 ;;
	esac
	for F in $PREP_FILES; do [ "$1" = "$F" ] && return 0; done
	return 1
}

die() { # die CODE MESSAGE
	local code="$1"
	shift
	echo "release.bash: $*" >&2
	exit "$code"
}

usage() {
	sed -n '2,/^set -euo/p' "$0" | sed '$d' | sed 's/^# \{0,1\}//'
}

DRY_RUN=0
ASSUME_YES=0
for ARG in "$@"; do
	case "$ARG" in
	--dry-run) DRY_RUN=1 ;;
	--yes) ASSUME_YES=1 ;;
	-h | --help)
		usage
		exit 0
		;;
	-*) die "$E_USAGE" "usage: unknown option $ARG (see --help)" ;;
	*) die "$E_USAGE" "usage: unexpected argument $ARG (see --help)" ;;
	esac
done

# Required commands, all reported at once.
MISSING=""
for CMD in gh jq git; do
	command -v "$CMD" >/dev/null 2>&1 || MISSING="$MISSING $CMD"
done
if command -v sha256sum >/dev/null 2>&1; then
	SHA="sha256sum"
elif command -v shasum >/dev/null 2>&1; then
	SHA="shasum -a 256"
else
	MISSING="$MISSING sha256sum"
fi
[ -z "$MISSING" ] || die "$E_NO_INPUT" "required command(s) not found:$MISSING"

[ -f codemeta.json ] || die "$E_NO_INPUT" "codemeta.json not found in $(pwd)"
VERSION="$(jq -r '.version // empty' codemeta.json 2>/dev/null)" ||
	die "$E_DATA" "codemeta.json is not valid JSON"
[ -n "$VERSION" ] || die "$E_DATA" "codemeta.json has no version"
printf '%s' "$VERSION" | grep -qE '^[0-9a-zA-Z._-]+$' ||
	die "$E_DATA" "version contains unexpected characters: $VERSION"
TAG="v$VERSION"
NAME="$(jq -r '.name // empty' codemeta.json)"
[ -n "$NAME" ] || NAME="$(basename "$(pwd)")"

shopt -s nullglob
ZIPS=(dist/*.zip)
shopt -u nullglob
[ "${#ZIPS[@]}" -gt 0 ] || die "$E_NO_INPUT" "no dist/*.zip; run make release first"

# A release is built from what is on GitHub: clean tree, everything pushed.
git rev-parse --is-inside-work-tree >/dev/null 2>&1 || die "$E_NO_INPUT" "not inside a git work tree"
DIRTY_PREP=()
while IFS= read -r LINE; do
	[ -n "$LINE" ] || continue
	CHANGED="${LINE:3}"
	# Only a plain modification of a release file is allowed (" M" or "M ").
	case "${LINE:0:2}" in
	" M" | "M " | "MM") is_prep_file "$CHANGED" && { DIRTY_PREP+=("$CHANGED"); continue; } ;;
	esac
	die "$E_NEGATIVE" "uncommitted changes beyond the release files ($CHANGED); commit or revert them first"
done < <(git status --porcelain)
UPSTREAM="$(git rev-parse --abbrev-ref --symbolic-full-name '@{u}' 2>/dev/null)" ||
	die "$E_NEGATIVE" "commits not pushed: the branch has no upstream"
[ "$(git rev-list --count "$UPSTREAM..HEAD")" = "0" ] ||
	die "$E_NEGATIVE" "commits not pushed to $UPSTREAM; push first"

HEAD_HASH="$(git rev-parse HEAD)"
if TAG_HASH="$(git rev-parse -q --verify "refs/tags/${TAG}^{commit}" 2>/dev/null)"; then
	[ "$TAG_HASH" = "$HEAD_HASH" ] ||
		die "$E_NEGATIVE" "tag $TAG exists and points at ${TAG_HASH:0:7}, not HEAD ${HEAD_HASH:0:7}"
fi
if gh release view "$TAG" >/dev/null 2>&1; then
	die "$E_NEGATIVE" "release $TAG already exists on GitHub"
fi

CHECKSUMS="${NAME}-${TAG}-checksums.txt"
echo "tag:       $TAG (commit ${HEAD_HASH:0:7})"
echo "release:   draft"
echo "files:"
for Z in "${ZIPS[@]}"; do echo "  $Z"; done
echo "  dist/$CHECKSUMS (to be written)"
echo "notes:     $(jq -r '.releaseNotes // ""' codemeta.json | head -c 200)"

if [ "${#DIRTY_PREP[@]}" -gt 0 ]; then
	echo "would commit and push: ${DIRTY_PREP[*]}"
fi

if [ "$DRY_RUN" = "1" ]; then
	echo "dry run: nothing created"
	exit 0
fi

if [ "$ASSUME_YES" != "1" ]; then
	read -r -p "Create the draft release $TAG on GitHub? (y/N) " YES_NO || YES_NO=""
	[ "$YES_NO" = "y" ] || die "$E_NEGATIVE" "release cancelled; nothing sent to GitHub"
fi

if [ "${#DIRTY_PREP[@]}" -gt 0 ]; then
	git add -- "${DIRTY_PREP[@]}"
	git commit -q -m "prep for $TAG" || die "$E_IO" "git could not commit the release files"
	git push -q || die "$E_UNAVAILABLE" "git push failed; the release files are committed locally, nothing sent to GitHub"
	echo "committed and pushed: ${DIRTY_PREP[*]}"
fi

(cd dist && $SHA "${ZIPS[@]##*/}" >"$CHECKSUMS")
echo "checksums written to dist/$CHECKSUMS"

NOTES_FILE="$(mktemp)"
trap 'rm -f "$NOTES_FILE"' EXIT
jq -r '.releaseNotes // ""' codemeta.json >"$NOTES_FILE"

gh release create "$TAG" --draft --notes-file "$NOTES_FILE" "${ZIPS[@]}" "dist/$CHECKSUMS" ||
	die "$E_UNAVAILABLE" "gh could not create the release $TAG"

echo "draft release $TAG created; review and publish it on the repository's releases page"
