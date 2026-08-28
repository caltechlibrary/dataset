#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = dataset

GIT_GROUP = caltechlibrary

RELEASE_DATE=$(shell date +'%Y-%m-%d')

RELEASE_HASH=$(shell git log --pretty=format:'%h' -n 1)

BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

# Getting the current version from codemeta.json varies depending on OS.
ifeq ($(OS), Windows_NT)
	VERSION = $(shell grep '"version":' codemeta.json | jq -r .version)
else
	VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)
endif


MAN_PAGES = dataset.1 datasetd.1 dsquery.1

MAN_PAGES_MISC = datasetd_yaml.5 datasetd_service.5 datasetd_api.5

PROGRAMS = dataset datasetd dsquery

PREFIX = $(HOME)

ifneq ($(prefix),)
        PREFIX = $(prefix)
endif

EXT =
ifeq ($(OS), Windows_NT)
	EXT = .exe
endif

EXT_WEB = .wasm

DIST_FOLDERS = bin/* man/*

build: version.go $(PROGRAMS) man CITATION.cff about.md installer.sh installer.ps1

libdataset.wasm: libdataset/libdataset.go libdataset/registry.go
	@mkdir -p dist
	GOOS=wasip1 GOARCH=wasm go build -o dist/libdataset.wasm ./libdataset/
	@echo "Built dist/libdataset.wasm ($$(ls -lh dist/libdataset.wasm | awk '{print $$5}'))"

version.go: .FORCE
	cmt codemeta.json version.go
	-git add version.go


$(PROGRAMS): cmd/*/*.go $(PACKAGE)
	@mkdir -p bin
	go build -o bin/$@$(EXT) cmd/$@/$@.go
	./bin/$@ -help >$@.1.md

man: $(MAN_PAGES) $(MAN_PAGES_LIB) $(MAN_PAGES_MISC)

$(MAN_PAGES): .FORCE
	mkdir -p man/man1
	pandoc $@.md --from markdown --to man -s >man/man1/$@

$(MAN_PAGES_MISC): $(PROGRAMS) .FORCE
	@./bin/datasetd --help api >datasetd_api.5.md
	@./bin/datasetd --help service >datasetd_service.5.md
	@./bin/datasetd --help yaml >datasetd_yaml.5.md
	mkdir -p man/man5
	pandoc $@.md --from markdown --to man -s >man/man5/$@


CITATION.cff: codemeta.json .FORCE
	cmt codemeta.json CITATION.cff

about.md: codemeta.json .FORCE
	cmt codemeta.json about.md

installer.sh: .FORCE
	cmt codemeta.json installer.sh
	chmod 775 installer.sh
	git add -f installer.sh

installer.ps1: .FORCE
	cmt codemeta.json installer.ps1
	chmod 775 installer.sh
	git add -f installer.ps1

# NOTE: on macOS you must use "mv" instead of "cp" to avoid problems
install: build
	@if [ ! -d $(PREFIX)/bin ]; then mkdir -p $(PREFIX)/bin; fi
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv -v ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"
	@echo ""
	@for FNAME in $(MAN_PAGES); do if [ -f "./man/man1/$${FNAME}" ]; then cp -v "./man/man1/$${FNAME}" "$(PREFIX)/man/man1/$${FNAME}"; fi; done
	@for FNAME in $(MAN_PAGES_LIB); do if [ -f "./man/man3/$${FNAME}" ]; then cp -v "./man/man3/$${FNAME}" "$(PREFIX)/man/man3/$${FNAME}"; fi; done
	@for FNAME in $(MAN_PAGES_MISC); do if [ -f "./man/man5/$${FNAME}" ]; then cp -v "./man/man5/$${FNAME}" "$(PREFIX)/man/man5/$${FNAME}"; fi; done
	@echo "Make sure $(PREFIX)/man is in your MANPATH"
	@echo ""

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	@echo "Removing manpages in $(PREFIX)/man"
	@for FNAME in $(MAN_PAGES); do if [ -f "$(PREFIX)/man/man1/$${FNAME}" ]; then rm -v "$(PREFIX)/man/man1/$${FNAME}"; fi; done
	@for FNAME in $(MAN_PAGES_LIB); do if [ -f "$(PREFIX)/man/man3/$${FNAME}" ]; then rm -v "$(PREFIX)/man/man3/$${FNAME}"; fi; done
	@for FNAME in $(MAN_PAGES_MISC); do if [ -f "$(PREFIX)/man/man5/$${FNAME}" ]; then rm -v "$(PREFIX)/man/man5/$${FNAME}"; fi; done


website: .FORCE
	make -f website.mak
	cd docs && make -f website.mak
	cd how-to && make -f website.mak

check: .FORCE
	go vet *.go
	cd api && go vet *.go
	cd cli && go vet *.go
	cd config && go vet *.go
	cd cmd/dataset && go vet *.go
	cd cmd/datasetd && go vet *.go
	cd dotpath && go vet *.go
	cd dsv1 && go vet *.go
	cd dsv1/tbl && go vet *.go
	cd pairtree && go vet *.go
	cd ptstore && go vet *.go
	cd semver && go vet *.go
	cd sqlstore && go vet *.go
	cd texts && go vet *.go

test: clean
	go test

cleanweb:
	@if [ -f index.html ]; then rm *.html; fi

clean:
	go clean
	@if [ -d bin ]; then rm -fR bin; fi
	@if [ -d dist ]; then rm -fR dist; fi
	@if [ -d testout ]; then rm -fR testout; fi
	@if [ -d semver/testout ]; then rm -fR semver/testout; fi
	@if [ -d dotpath/testout ]; then rm -fR dotpath/testout; fi
	@if [ -d pairtree/testout ]; then rm -fR pairtree/testout; fi
	@if [ -d ptstore/testout ]; then rm -fR ptstore/testout; fi
	@if [ -d sqlstore/testout ]; then rm -fR sqlstore/testout; fi
	@if [ -d texts/testout ]; then rm -fR texts/testout; fi
	@if [ -d api/testout ]; then rm -fR api/testout; fi
	@if [ -d cli/testout ]; then rm -fR cli/testout; fi
	-go clean -r

dist/Linux-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-x86_64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/Linux-aarch64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-aarch64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macOS-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-x86_64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/macOS-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/Windows-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-x86_64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/Windows-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=arm64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-arm64.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

# Raspberry Pi OS (32 bit), as reported by Raspberry Pi 3B+
dist/Linux-armv7l:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/$$FNAME cmd/$$FNAME/*.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-armv7l.zip LICENSE codemeta.json CITATION.cff *.md $(DIST_FOLDERS)
	@rm -fR dist/bin

dist/libdataset: dist/libdataset.wasm libdataset.md \
		wrappers/python/libdataset/__init__.py \
		wrappers/python/libdataset/libdataset.py \
		wrappers/python/README.md \
		wrappers/typescript/libdataset.ts \
		wrappers/typescript/README.md
	@mkdir -p dist/libdataset-staging/python/libdataset
	@mkdir -p dist/libdataset-staging/typescript
	@cp dist/libdataset.wasm           dist/libdataset-staging/
	@cp wrappers/python/libdataset/*.py dist/libdataset-staging/python/libdataset/
	@cp wrappers/python/README.md       dist/libdataset-staging/python/
	@cp wrappers/typescript/libdataset.ts dist/libdataset-staging/typescript/
	@cp wrappers/typescript/README.md   dist/libdataset-staging/typescript/
	@cp LICENSE                         dist/libdataset-staging/
	@cp libdataset.md                   dist/libdataset-staging/
	@cd dist/libdataset-staging && zip -r ../$(PROJECT)-v$(VERSION)-libdataset.zip .
	@rm -rf dist/libdataset-staging
	@echo "Built dist/$(PROJECT)-v$(VERSION)-libdataset.zip"

	
distribute_docs:
	if [ -d dist ]; then rm -fR dist; fi
	mkdir -p dist
	cp -v codemeta.json dist/
	cp -v CITATION.cff dist/
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v INSTALL_NOTES_*.md dist/
	cp installer.sh dist/
	cp installer.ps1 dist/
	cp -vR man dist/

update_version:
	$(EDITOR) codemeta.json
	codemeta2cff

release: .FORCE clean build version.go CITATION.cff man website distribute_docs dist/Linux-x86_64 dist/Linux-aarch64 dist/Linux-armv7l dist/Windows-x86_64 dist/Windows-arm64 dist/macOS-x86_64 dist/macOS-arm64 libdataset.wasm dist/libdataset
	@printf "\nReady to run\n\n\t./release.bash\n\n"

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

loghash: .FORCE
	git log --pretty=format:'%h' -n 1

.FORCE:
