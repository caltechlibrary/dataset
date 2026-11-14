#!/bin/bash
#
# Tests for release.bash and release.ps1. Each case builds a throwaway git
# repository with a bare remote, a codemeta.json, a dist/ holding fake zips and
# a stub `gh` on PATH that records its arguments, so nothing here can reach
# GitHub.
#
# The same checks run against either script, so the two are shown to agree on
# every exit status and message the checks look at:
#
#   ./release_test.bash                      release.bash   (needs bash, git, jq)
#   RELEASE_IMPL=ps1 ./release_test.bash     release.ps1    (needs pwsh and git)
#   make test-release                        both
#
set -u

HERE="$(cd "$(dirname "$0")" && pwd)"
IMPL="${RELEASE_IMPL:-bash}"
case "$IMPL" in
bash) SCRIPT=release.bash ;;
ps1)
	SCRIPT=release.ps1
	PWSH="$(command -v pwsh)" || { echo "release_test.bash: pwsh is needed for RELEASE_IMPL=ps1" >&2; exit 66; }
	;;
*) echo "release_test.bash: RELEASE_IMPL must be bash or ps1" >&2; exit 2 ;;
esac
# run_script runs the script under test in the current directory.
run_script() {
	if [ "$IMPL" = ps1 ]; then "$PWSH" -NoProfile -File "./$SCRIPT" "$@"; else bash "./$SCRIPT" "$@"; fi
}
PASS=0
FAIL=0
ROOT="$(mktemp -d)"
trap 'rm -rf "$ROOT"' EXIT

# new_repo makes $ROOT/<name>/work with a clean tree pushed to a bare remote,
# version 1.2.3, one fake zip and a stub gh. It prints the work directory.
new_repo() {
	local name="$1" base="$ROOT/$1"
	mkdir -p "$base/bin"
	git init -q --bare "$base/remote.git"
	git clone -q "$base/remote.git" "$base/work" 2>/dev/null
	(
		cd "$base/work"
		git config user.email test@example.edu
		git config user.name Test
		jq -n '{name:"demo", version:"1.2.3", releaseNotes:"Notes with `backticks` and $dollars"}' >codemeta.json
		printf '/dist/\n' >.gitignore
		printf 'package demo\n' >version.go
		cp "$HERE/$SCRIPT" .
		git add . && git commit -q -m init && git push -q -u origin HEAD 2>/dev/null
		mkdir dist
		printf 'zip-a' >dist/demo-v1.2.3-Linux-x86_64.zip
		printf 'zip-b' >dist/demo-v1.2.3-macOS-arm64.zip
	)
	cat >"$base/bin/gh" <<'STUB'
#!/bin/bash
# Stub gh: log the arguments, and the notes file's content for release create.
LOG="$(dirname "$0")/gh.log"
echo "gh $*" >>"$LOG"
if [ "$1 $2" = "release view" ]; then
	[ "${GH_RELEASE_EXISTS:-}" = "1" ] && exit 0
	exit 1
fi
if [ "$1 $2" = "release create" ]; then
	[ "${GH_FAIL:-}" = "1" ] && exit 1
	echo "STATE clean=$(git status --porcelain | wc -l | tr -d ' ') ahead=$(git rev-list --count '@{u}..HEAD')" >>"$LOG"
	while [ $# -gt 0 ]; do
		if [ "$1" = "--notes-file" ] || [ "$1" = "-F" ]; then echo "NOTESPATH: $2" >>"$LOG"; echo "NOTES: $(cat "$2")" >>"$LOG"; fi
		shift
	done
fi
exit 0
STUB
	chmod +x "$base/bin/gh"
	echo "$base/work"
}

# run_release DIR ARGS... runs the script under test in DIR with the stub gh first on
# PATH, answering any prompt with "n", and sets $CODE and $OUT.
run_release() {
	local dir="$1"
	shift
	OUT="$(cd "$dir" && PATH="$(dirname "$dir")/bin:$PATH" run_script "$@" </dev/null 2>&1)"
	CODE=$?
}

check() { # check NAME EXPECTED_CODE [PATTERN]
	local name="$1" want="$2" pat="${3:-}"
	if [ "$CODE" != "$want" ]; then
		echo "FAIL $name: exit $CODE, want $want"
		echo "$OUT" | sed 's/^/     | /'
		FAIL=$((FAIL + 1))
	elif [ -n "$pat" ] && ! printf '%s' "$OUT" | grep -q -- "$pat"; then
		echo "FAIL $name: output does not contain '$pat'"
		echo "$OUT" | sed 's/^/     | /'
		FAIL=$((FAIL + 1))
	else
		echo "ok   $name"
		PASS=$((PASS + 1))
	fi
}

sha256() { if command -v sha256sum >/dev/null 2>&1; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi; }

gh_log() { cat "$(dirname "$1")/bin/gh.log" 2>/dev/null; }

# --- usage: exit 2 --------------------------------------------------------
W="$(new_repo usage)"
run_release "$W" --bogus;           check "bad flag"        2 "usage: unknown option --bogus"
check "message carries the script name, not the exit code" 2 "^release\.[a-z0-9]*: usage: unknown option"
run_release "$W" --dry-run extra;   check "surplus argument" 2 "usage: unexpected argument extra"
run_release "$W" --help;            check "--help"          0 "EXIT STATUS"
run_release "$W" -h;                check "-h"              0 "EXIT STATUS"

# --- data and input: 65, 66 -----------------------------------------------
W="$(new_repo badtag)"
(cd "$W" && jq '.version="1.2.3; rm -rf /"' codemeta.json >c && mv c codemeta.json && git commit -q -am bad && git push -q)
run_release "$W" --dry-run;         check "bad tag characters" 65 "version"

W="$(new_repo nover)"
(cd "$W" && jq 'del(.version)' codemeta.json >c && mv c codemeta.json && git commit -q -am nover && git push -q)
run_release "$W" --dry-run;         check "no version"      65 "no version"

W="$(new_repo badjson)"
(cd "$W" && echo '{' >codemeta.json && git commit -q -am badjson && git push -q)
run_release "$W" --dry-run;         check "codemeta.json is not JSON" 65 "not valid JSON"

W="$(new_repo noname)"
(cd "$W" && jq 'del(.name)' codemeta.json >c && mv c codemeta.json && git commit -q -am noname && git push -q)
run_release "$W" --dry-run;         check "no name uses the directory name" 0 "work-v1.2.3-checksums.txt"

W="$(new_repo longnotes)"
(cd "$W" && jq --arg n "$(printf 'a%.0s' $(seq 1 200))ZZZ" '.releaseNotes=$n' codemeta.json >c && mv c codemeta.json && git commit -q -am notes && git push -q)
run_release "$W" --dry-run;         check "notes are shown, 200 characters of them" 0 "^notes:     a\{200\}$"
if printf '%s' "$OUT" | grep -q ZZZ; then echo "FAIL notes are cut at 200 characters"; FAIL=$((FAIL + 1)); else echo "ok   notes are cut at 200 characters"; PASS=$((PASS + 1)); fi

W="$(new_repo nocodemeta)"
(cd "$W" && git rm -q codemeta.json && git commit -q -m x && git push -q)
run_release "$W" --dry-run;         check "no codemeta.json" 66 "codemeta.json"

W="$(new_repo nodist)"
rm -f "$W"/dist/*.zip
run_release "$W" --dry-run;         check "no dist zips"    66 "dist"

W="$(new_repo nogh)"
rm -f "$(dirname "$W")/bin/gh"
OUT="$(cd "$W" && PATH="/usr/bin:/bin" run_script --dry-run </dev/null 2>&1)"; CODE=$?
check "gh missing" 66 "gh"

W="$(new_repo nogit)"
rm -rf "$W/.git"
run_release "$W" --dry-run;         check "not inside a git work tree" 66 "not inside a git work tree"

# --- state forbids: 1 -----------------------------------------------------
W="$(new_repo dirty)"
(cd "$W" && echo '# change' >>.gitignore)
run_release "$W" --dry-run;         check "dirty tree"      1 "uncommitted"

# The files `make release` regenerates may be dirty: the script commits them.
W="$(new_repo prepdry)"
(cd "$W" && echo '// regenerated' >>version.go && jq '.releaseNotes="new notes"' codemeta.json >c && mv c codemeta.json)
run_release "$W" --dry-run;         check "release files dirty, dry run" 0 "would commit"
if git -C "$W" status --porcelain | grep -q version.go; then echo "ok   dry run leaves the files uncommitted"; PASS=$((PASS + 1)); else echo "FAIL dry run committed"; FAIL=$((FAIL + 1)); fi

# Each file a release regenerates may be dirty, and the generated man pages
# (docs/*.?.md) too; other docs and other names may not.
for F in CITATION.cff installer.sh installer.ps1 docs/demo.1.md docs/sub/demo.2.md; do
	W="$(new_repo "prep-$(echo "$F" | tr '/.' '__')")"
	(cd "$W" && mkdir -p "$(dirname "$F")" && echo first >"$F" && git add -A && git commit -q -m "add $F" && git push -q && echo changed >>"$F")
	run_release "$W" --dry-run;         check "release file $F may be dirty" 0 "would commit and push: .*$F"
done
for F in docs/guide.md docs/demo.12.md docs/demo.md docs/demo.1.MD CITATION.CFF README.md Makefile; do
	W="$(new_repo "stray-$(echo "$F" | tr '/.' '__')")"
	(cd "$W" && mkdir -p "$(dirname "$F")" && echo first >"$F" && git add -A && git commit -q -m "add $F" && git push -q && echo changed >>"$F")
	run_release "$W" --dry-run;         check "$F may not be dirty" 1 "uncommitted"
done

# A change that is staged, or staged and then changed again, is still a plain
# modification of a release file.
W="$(new_repo staged)"
(cd "$W" && echo '// regenerated' >>version.go && git add version.go)
run_release "$W" --dry-run;         check "release file staged" 0 "would commit and push: version.go"
W="$(new_repo stagedmore)"
(cd "$W" && echo '// regenerated' >>version.go && git add version.go && echo '// more' >>version.go)
run_release "$W" --dry-run;         check "release file staged then changed again" 0 "would commit and push: version.go"

W="$(new_repo untracked)"
(cd "$W" && echo 'cff-version: 1.2.0' >CITATION.cff)
run_release "$W" --dry-run;         check "untracked release file is refused" 1 "uncommitted"

W="$(new_repo prepdirty)"
(cd "$W" && echo '// regenerated' >>version.go && echo '# stray' >>.gitignore)
run_release "$W" --dry-run;         check "stray dirty file beside release files" 1 "uncommitted"

W="$(new_repo noupstream)"
(cd "$W" && git branch -q --unset-upstream)
run_release "$W" --dry-run;         check "branch with no upstream" 1 "not pushed"

W="$(new_repo unpushed)"
(cd "$W" && echo x >f && git add f && git commit -q -m local)
run_release "$W" --dry-run;         check "unpushed commit" 1 "not pushed"

W="$(new_repo exists)"
GH_RELEASE_EXISTS=1 run_release "$W" --dry-run; check "release already exists" 1 "already exists"

W="$(new_repo tagelsewhere)"
(cd "$W" && git tag v1.2.3 && echo y >g && git add g && git commit -q -m next && git push -q)
run_release "$W" --dry-run;         check "tag points elsewhere" 1 "tag"

W="$(new_repo tagathead)"
(cd "$W" && git tag v1.2.3)
run_release "$W" --dry-run;         check "tag already at HEAD is fine" 0 "v1.2.3"

# --- dry run: 0, creates nothing ------------------------------------------
W="$(new_repo dry)"
run_release "$W" --dry-run;         check "dry run succeeds" 0 "v1.2.3"
check "dry run shows the tag and commit" 0 "^tag:       v1.2.3 (commit [0-9a-f]\{7\})"
check "dry run says it is a draft" 0 "^release:   draft"
check "dry run lists the files" 0 "^files:"
check "dry run lists each zip" 0 "^  dist/demo-v1.2.3-Linux-x86_64.zip"
check "dry run lists the checksums file" 0 "^  dist/demo-v1.2.3-checksums.txt (to be written)"
check "dry run says nothing was created" 0 "^dry run: nothing created"
if printf '%s' "$OUT" | grep -q "would commit"; then echo "FAIL a clean tree reports files to commit"; FAIL=$((FAIL + 1)); else echo "ok   a clean tree reports no files to commit"; PASS=$((PASS + 1)); fi
if gh_log "$W" | grep -q "release create"; then echo "FAIL dry run called gh release create"; FAIL=$((FAIL + 1)); else echo "ok   dry run creates no release"; PASS=$((PASS + 1)); fi
if ls "$W"/dist/*checksums* >/dev/null 2>&1; then echo "FAIL dry run wrote a checksums file"; FAIL=$((FAIL + 1)); else echo "ok   dry run writes no checksums"; PASS=$((PASS + 1)); fi

# --- prompt declined: 1, creates nothing ----------------------------------
W="$(new_repo declined)"
run_release "$W";                   check "prompt declined" 1 "release cancelled; nothing sent to GitHub"
# bash's read -p prints its prompt only on a terminal, so this one is ps1 only
[ "$IMPL" != ps1 ] || check "the prompt names the tag" 1 "Create the draft release v1.2.3 on GitHub"
if gh_log "$W" | grep -q "release create"; then echo "FAIL declined run called gh release create"; FAIL=$((FAIL + 1)); else echo "ok   declined run creates no release"; PASS=$((PASS + 1)); fi

# The answer is read from standard input, and only a lower case y confirms.
# run_answer DIR ANSWER runs the script with ANSWER typed at the prompt.
run_answer() {
	OUT="$(cd "$1" && PATH="$(dirname "$1")/bin:$PATH" run_script 2>&1 <<<"$2")"
	CODE=$?
}
W="$(new_repo answer-upper)"
run_answer "$W" "Y";                check "an upper case Y does not confirm" 1 "release cancelled"
W="$(new_repo answer-lower)"
run_answer "$W" "y";                check "a lower case y confirms" 0 "draft release v1.2.3 created"
if gh_log "$W" | grep -q "release create"; then echo "ok   a confirmed run creates the release"; PASS=$((PASS + 1)); else echo "FAIL a confirmed run did not call gh release create"; FAIL=$((FAIL + 1)); fi

# --- real run with --yes --------------------------------------------------
W="$(new_repo real)"
run_release "$W" --yes;             check "release with --yes" 0 "draft"
check "the checksums file is announced" 0 "^checksums written to dist/demo-v1.2.3-checksums.txt"
check "the closing message names the draft" 0 "^draft release v1.2.3 created; review and publish it"
LOG="$(gh_log "$W")"
NOTES_PATH="$(printf '%s\n' "$LOG" | sed -n 's/^NOTESPATH: //p')"
if [ -n "$NOTES_PATH" ] && [ ! -e "$NOTES_PATH" ]; then echo "ok   the temporary notes file is removed"; PASS=$((PASS + 1)); else echo "FAIL the temporary notes file '$NOTES_PATH' was left behind"; FAIL=$((FAIL + 1)); fi
for want in "release create v1.2.3" "--draft" "demo-v1.2.3-Linux-x86_64.zip" "demo-v1.2.3-macOS-arm64.zip" "demo-v1.2.3-checksums.txt"; do
	if printf '%s' "$LOG" | grep -q -- "$want"; then echo "ok   gh create carries $want"; PASS=$((PASS + 1)); else echo "FAIL gh create missing $want"; echo "$LOG" | sed 's/^/     | /'; FAIL=$((FAIL + 1)); fi
done
if printf '%s' "$LOG" | grep -q 'NOTES: Notes with `backticks` and $dollars'; then echo "ok   release notes pass through unmangled"; PASS=$((PASS + 1)); else echo "FAIL release notes mangled"; echo "$LOG" | sed 's/^/     | /'; FAIL=$((FAIL + 1)); fi
CS="$W/dist/demo-v1.2.3-checksums.txt"
if [ -f "$CS" ] && grep -q "demo-v1.2.3-Linux-x86_64.zip" "$CS" && ! grep -q "dist/" "$CS"; then echo "ok   checksums file lists zips by bare name"; PASS=$((PASS + 1)); else echo "FAIL checksums file wrong"; FAIL=$((FAIL + 1)); fi

# the checksums are the real SHA-256 of each zip, lower case, two spaces, bare name
for z in "$W"/dist/*.zip; do
	want="$(sha256 "$z")  $(basename "$z")"
	if grep -qxF "$want" "$CS"; then echo "ok   checksum line for $(basename "$z")"; PASS=$((PASS + 1)); else echo "FAIL checksum line for $(basename "$z"): want '$want'"; sed 's/^/     | /' "$CS"; FAIL=$((FAIL + 1)); fi
done
if [ "$(wc -l <"$CS" | tr -d ' ')" = "2" ]; then echo "ok   one checksum line per zip"; PASS=$((PASS + 1)); else echo "FAIL checksum file has the wrong number of lines"; FAIL=$((FAIL + 1)); fi

# --- real run commits and pushes the release files before calling gh -----
W="$(new_repo prepreal)"
(cd "$W" && echo '// regenerated' >>version.go)
run_release "$W" --yes;             check "release with dirty release files" 0 "draft"
check "the prep commit is announced" 0 "^committed and pushed: version.go"
if gh_log "$W" | grep -q "STATE clean=0 ahead=0"; then echo "ok   tree was clean and pushed when gh create ran"; PASS=$((PASS + 1)); else echo "FAIL tree not clean and pushed at gh create"; gh_log "$W" | sed 's/^/     | /'; FAIL=$((FAIL + 1)); fi
if git -C "$(dirname "$W")/remote.git" log --format=%s | head -1 | grep -q "prep for v1.2.3"; then echo "ok   remote has the prep commit"; PASS=$((PASS + 1)); else echo "FAIL remote lacks the prep commit"; FAIL=$((FAIL + 1)); fi

W="$(new_repo commitfail)"
(cd "$W" && echo '// regenerated' >>version.go && printf '#!/bin/sh\nexit 1\n' >.git/hooks/pre-commit && chmod +x .git/hooks/pre-commit)
run_release "$W" --yes;             check "commit fails" 74 "commit"
if gh_log "$W" | grep -q "release create"; then echo "FAIL gh create ran after a failed commit"; FAIL=$((FAIL + 1)); else echo "ok   no release after a failed commit"; PASS=$((PASS + 1)); fi

W="$(new_repo pushfail)"
(cd "$W" && echo '// regenerated' >>version.go && git remote set-url origin /nonexistent/remote.git)
run_release "$W" --yes;             check "push fails" 69 "push"
if gh_log "$W" | grep -q "release create"; then echo "FAIL gh create ran after a failed push"; FAIL=$((FAIL + 1)); else echo "ok   no release after a failed push"; PASS=$((PASS + 1)); fi

# --- gh failure: 69 -------------------------------------------------------
W="$(new_repo ghfail)"
GH_FAIL=1 run_release "$W" --yes;   check "gh release create fails" 69 "gh"

echo
echo "$PASS passed, $FAIL failed"
[ "$FAIL" = 0 ]
