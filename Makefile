#
# Makefile for dataset, a Go package, command line tools and a web service.
#
# Everyday targets:   build, test, test-release, install, uninstall, clean
# Release targets:    generate, release, dist (see "Releasing" below)
#
# A plain `make`, `make build`, `make test` and `make install` need only git,
# make, go and pandoc. `cmt` is run only by `make generate` and by
# `make release`, never by the everyday targets, so a fresh clone builds with
# the committed version.go, CITATION.cff and installers.
#
PROJECT = dataset

GIT_GROUP = caltechlibrary

PROGRAMS = $(shell ls -1 cmd)

# version.go is left out: it has a rule (cmt, under generate), and listing it
# here would make every build run cmt through that rule.
PACKAGE = $(shell ls -1 *.go | grep -v 'version.go')

# The manual pages are the programs' own help output, written into docs/ by
# the build, which is where the documentation site is built from. Section 1
# has one page per program; section 5 has the three datasetd pages.
MAN_PAGES = $(shell ls -1 docs/*.1.md 2>/dev/null | sed -E 's|^docs/||; s/\.1.md/.1/g')

MAN_PAGES_5 = $(shell ls -1 docs/*.5.md 2>/dev/null | sed -E 's|^docs/||; s/\.5.md/.5/g')

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

RELEASE_HASH = $(shell git log --pretty=format:'%h' -n 1)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

OS = $(shell uname)

PREFIX = $(HOME)

ifneq ($(prefix),)
	PREFIX = $(prefix)
endif

EXT =
ifeq ($(OS), Windows)
	EXT = .exe
endif

# One entry per release platform: name:GOOS:GOARCH
PLATFORMS = Linux-x86_64:linux:amd64 Linux-aarch64:linux:arm64 \
	macOS-x86_64:darwin:amd64 macOS-arm64:darwin:arm64 \
	Windows-x86_64:windows:amd64 Windows-arm64:windows:arm64

build: $(PROGRAMS)
	@$(MAKE) pages5
	@$(MAKE) man

# Building a program also writes its manual page (the program's own -help
# output) into docs/. NOTE: unlike irdmtools, no page per help topic is
# written: dataset has no "-help topics" listing yet (see TODO.md).
$(PROGRAMS): $(PACKAGE)
	@mkdir -p bin
	go build -o "bin/$@$(EXT)" ./cmd/$@
	@mkdir -p docs
	@./bin/$@$(EXT) -help >docs/$@.1.md

# pages5 writes the datasetd section 5 pages from its help topics. It needs
# the programs built, so build runs it as a sub-make once they are.
pages5: .FORCE
	@mkdir -p docs
	@for TOPIC in api service yaml; do ./bin/datasetd$(EXT) --help $$TOPIC >docs/datasetd_$$TOPIC.5.md || exit 1; done

# man renders each docs/NAME.N.md into man/manN/NAME.N. build runs it as a
# sub-make so that it sees the pages written there: MAN_PAGES is read when the
# Makefile is read.
man: $(MAN_PAGES) $(MAN_PAGES_5)

$(MAN_PAGES): .FORCE
	@mkdir -p man/man1
	pandoc docs/$@.md --from markdown --to man -s >man/man1/$@

$(MAN_PAGES_5): .FORCE
	@mkdir -p man/man5
	pandoc docs/$@.md --from markdown --to man -s >man/man5/$@

# website builds the documentation site: the project pages, docs/ and how-to/.
website: .FORCE
	make -f website.mak
	cd docs && make -f website.mak
	cd how-to && make -f website.mak

# test runs offline: vet, then the Go tests (TestRunAPI starts a datasetd on
# localhost:8585 and takes several seconds).
test: .FORCE
	go vet ./...
	go test ./...

# test-release runs the tests of release.bash, and of release.ps1 when pwsh is
# installed. They use throwaway repositories and a stub gh, so nothing reaches
# GitHub.
test-release: .FORCE
	./release_test.bash
	@if command -v pwsh >/dev/null 2>&1; then RELEASE_IMPL=ps1 ./release_test.bash; else echo "pwsh not found: release.ps1 not tested"; fi

# check vets the code and lists, without changing them, files gofmt would
# reformat. Existing files are not reformatted as a side effect of other work.
check: .FORCE
	go vet ./...
	@UNFORMATTED="$$(gofmt -l . )"; if [ -n "$$UNFORMATTED" ]; then echo "not gofmt-formatted (left as is):"; echo "$$UNFORMATTED"; fi

# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build .FORCE
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME$(EXT) ]; then mv -v ./bin/$$FNAME$(EXT) $(PREFIX)/bin/$$FNAME$(EXT); fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@if [ ! -d $(PREFIX)/man/man1 ]; then mkdir -p $(PREFIX)/man/man1; fi
	@for MAN_PAGE in $(MAN_PAGES); do cp -v man/man1/$$MAN_PAGE $(PREFIX)/man/man1/;done
	@if [ ! -d $(PREFIX)/man/man5 ]; then mkdir -p $(PREFIX)/man/man5; fi
	@for MAN_PAGE in $(MAN_PAGES_5); do cp -v man/man5/$$MAN_PAGE $(PREFIX)/man/man5/;done
	@echo ""
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	-for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME$(EXT) ]; then rm -v $(PREFIX)/bin/$$FNAME$(EXT); fi; done
	-for MAN_PAGE in $(MAN_PAGES); do if [ -f "$(PREFIX)/man/man1/$$MAN_PAGE" ]; then rm "$(PREFIX)/man/man1/$$MAN_PAGE"; fi; done
	-for MAN_PAGE in $(MAN_PAGES_5); do if [ -f "$(PREFIX)/man/man5/$$MAN_PAGE" ]; then rm "$(PREFIX)/man/man5/$$MAN_PAGE"; fi; done

hash: .FORCE
	@git log --pretty=format:'%h' -n 1

# clean never removes version.go: it is committed, so that a fresh clone builds.
clean: .FORCE
	-if [ -d bin ]; then rm -fR bin; fi
	-if [ -d dist ]; then rm -fR dist; fi
	-if [ -d man ]; then rm -fR man; fi
	-if [ -d testout ]; then rm -fR testout; fi

status:
	git status

save:
	@git commit -am "$${msg:-Quick Save}"
	git push origin $(BRANCH)

refresh:
	git fetch origin
	git pull origin $(BRANCH)

#
# Releasing
#
# 1. Edit codemeta.json by hand: version, dateModified, datePublished,
#    releaseNotes.
# 2. make release      regenerates version.go, CITATION.cff and both installers
#                      from codemeta.json with cmt (the only place cmt runs
#                      besides an explicit `make generate`), rebuilds the
#                      programs and their manual pages from that version.go,
#                      runs the tests, builds one zip per platform in dist/ and
#                      runs release-check, which stops here if anything does
#                      not carry the version and the hash of the commit, or a
#                      zip lacks one of the programs. No libdataset zip is
#                      built: libdataset is deprecated and not released (dataset
#                      DR-0004); `make libdataset.wasm` builds it for yourself.
# 3. ./release.bash    refuses if anything beyond those files is uncommitted,
#                      asks for confirmation, commits and pushes the regenerated
#                      files ("prep for vX.Y.Z"), and creates the DRAFT GitHub
#                      release from dist/*.zip. Publish the draft from the
#                      repository's releases page (or with gh release edit).
#
# Run `make test-release` after changing release.bash, release.ps1 or
# release_test.bash.
#

# generate runs cmt over codemeta.json. Run it when the version or the
# metadata change; a plain `make` leaves the committed copies alone.
generate: version.go CITATION.cff installer.sh installer.ps1

version.go: .FORCE
	cmt codemeta.json version.go

CITATION.cff: .FORCE
	cmt codemeta.json CITATION.cff

installer.sh: .FORCE
	cmt codemeta.json installer.sh
	chmod 775 installer.sh

installer.ps1: .FORCE
	cmt codemeta.json installer.ps1

# dist builds one zip per platform in PLATFORMS, starting from an empty dist/.
# Each zip carries the programs (bin/), the manual pages (man/ and their
# sources, docs/*.N.md) and the top level documents.
dist: build .FORCE
	@rm -fR dist
	@mkdir -p dist
	@for ENTRY in $(PLATFORMS); do \
		NAME=$$(echo $$ENTRY | cut -d: -f1); GOOS=$$(echo $$ENTRY | cut -d: -f2); GOARCH=$$(echo $$ENTRY | cut -d: -f3); \
		SUFFIX=""; if [ "$$GOOS" = "windows" ]; then SUFFIX=".exe"; fi; \
		rm -fR dist/stage; mkdir -p dist/stage/bin dist/stage/docs; \
		for FNAME in $(PROGRAMS); do \
			env GOOS=$$GOOS GOARCH=$$GOARCH go build -o "dist/stage/bin/$${FNAME}$$SUFFIX" ./cmd/$$FNAME || exit 1; \
		done; \
		cp LICENSE codemeta.json CITATION.cff README.md CHANGES.md INSTALL.md INSTALL_NOTES_macOS.md INSTALL_NOTES_Windows.md dist/stage/; \
		cp docs/*.?.md dist/stage/docs/; \
		cp -R man dist/stage/; \
		(cd dist/stage && zip -qr ../$(PROJECT)-v$(VERSION)-$$NAME.zip .) || exit 1; \
		echo "dist/$(PROJECT)-v$(VERSION)-$$NAME.zip"; \
	done
	@rm -fR dist/stage

# libdataset.wasm is the WebAssembly build of the library (GOOS=wasip1). The
# library is deprecated and not part of a release; this target is for people who
# build their own (see libdataset.md).
libdataset.wasm: libdataset/libdataset.go libdataset/registry.go
	@mkdir -p dist
	GOOS=wasip1 GOARCH=wasm go build -o dist/libdataset.wasm ./libdataset/
	@echo "Built dist/libdataset.wasm ($$(ls -lh dist/libdataset.wasm | awk '{print $$5}'))"

# libdataset-dist lays out dist/dataset-vVERSION-libdataset.zip from your own
# libdataset.wasm. It is not part of release; if you run it after dist (which
# starts from an empty dist/) the zip is added beside the platform zips.
libdataset-dist: libdataset.wasm libdataset.md \
		wrappers/python/libdataset/__init__.py \
		wrappers/python/libdataset/libdataset.py \
		wrappers/python/README.md \
		wrappers/typescript/libdataset.ts \
		wrappers/typescript/README.md
	@mkdir -p dist/libdataset-staging/python/libdataset
	@mkdir -p dist/libdataset-staging/typescript
	@cp dist/libdataset.wasm dist/libdataset-staging/
	@cp wrappers/python/libdataset/*.py dist/libdataset-staging/python/libdataset/
	@cp wrappers/python/README.md dist/libdataset-staging/python/
	@cp wrappers/typescript/libdataset.ts dist/libdataset-staging/typescript/
	@cp wrappers/typescript/README.md dist/libdataset-staging/typescript/
	@cp LICENSE dist/libdataset-staging/
	@cp libdataset.md dist/libdataset-staging/
	@(cd dist/libdataset-staging && zip -qr ../$(PROJECT)-v$(VERSION)-libdataset.zip .)
	@rm -rf dist/libdataset-staging
	@echo "Built dist/$(PROJECT)-v$(VERSION)-libdataset.zip"

release-version: .FORCE
	@$(MAKE) -B generate

# release-check stops the release unless version.go, the installers,
# CITATION.cff, the manual pages in docs/ and the zip files all carry the
# version in codemeta.json and the hash of the commit being released, every
# platform zip carries every program.
release-check: .FORCE
	@HASH=$$(git log --pretty=format:%h -n 1); \
	if ! grep -q 'Version = "$(VERSION)"' version.go; then \
		echo "error: version.go does not carry version $(VERSION)" >&2; exit 1; fi; \
	if ! grep -q "ReleaseHash = \"$$HASH\"" version.go; then \
		echo "error: version.go does not carry hash $$HASH" >&2; exit 1; fi; \
	if ! grep -q '^version: $(VERSION)' CITATION.cff; then \
		echo "error: CITATION.cff does not carry version $(VERSION)" >&2; exit 1; fi; \
	for FNAME in installer.sh installer.ps1; do \
		if ! grep -q 'VERSION.*"$(VERSION)"' $$FNAME; then \
			echo "error: $$FNAME does not carry version $(VERSION)" >&2; exit 1; fi; \
	done; \
	for FNAME in docs/*.?.md; do \
		if ! head -4 $$FNAME | grep -q "version $(VERSION) $$HASH"; then \
			echo "error: $$FNAME does not carry version $(VERSION) $$HASH" >&2; exit 1; fi; \
	done; \
	COUNT=0; \
	for ENTRY in $(PLATFORMS); do \
		NAME=$$(echo $$ENTRY | cut -d: -f1); GOOS=$$(echo $$ENTRY | cut -d: -f2); ZIP=dist/$(PROJECT)-v$(VERSION)-$$NAME.zip; \
		SUFFIX=""; if [ "$$GOOS" = "windows" ]; then SUFFIX=".exe"; fi; \
		if [ ! -f $$ZIP ]; then echo "error: missing $$ZIP" >&2; exit 1; fi; \
		for FNAME in $(PROGRAMS); do \
			if ! unzip -Z1 $$ZIP | grep -qx "bin/$$FNAME$$SUFFIX"; then \
				echo "error: $$ZIP does not carry bin/$$FNAME$$SUFFIX" >&2; exit 1; fi; \
		done; \
		for PAGE in $$(unzip -Z1 $$ZIP 'docs/*.?.md'); do \
			if ! unzip -p $$ZIP $$PAGE | head -4 | grep -q "version $(VERSION) $$HASH"; then \
				echo "error: $$ZIP $$PAGE does not carry version $(VERSION) $$HASH" >&2; exit 1; fi; \
		done; \
		COUNT=$$((COUNT + 1)); \
	done; \
	echo "release-check: version.go, CITATION.cff, installers, docs/ and $$COUNT zip files carry $(VERSION) $$HASH"

# release: see "Releasing" above. It commits nothing; release.bash does that.
release: clean release-version test dist release-check
	@printf "\nrelease-check passed. Now run: ./release.bash   (or ./release.bash --dry-run first)\n\n"

.FORCE:
