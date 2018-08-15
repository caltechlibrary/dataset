#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = dataset

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PKGASSETS = $(shell which pkgassets)

PROJECT_LIST = dataset

OS = $(shell uname)

EXT = 
ifeq ($(OS), Windows)
	EXT = .exe
endif


dataset$(EXT): bin/dataset$(EXT)

cmd/dataset/assets.go:
	pkgassets -o cmd/dataset/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples how-to Help docs/dataset
	git add cmd/dataset/assets.go

bin/dataset$(EXT): dataset.go pairtree.go buckets.go attachments.go grid.go frame.go repair.go sort.go gsheet/gsheet.go cmd/dataset/dataset.go cmd/dataset/assets.go
	go build -o bin/dataset$(EXT) cmd/dataset/dataset.go cmd/dataset/assets.go

build: $(PROJECT_LIST) python

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/dataset/dataset.go cmd/dataset/assets.go
	if [ "$(OS)" != "Windows" ]; then mkdir -p $(GOPATH)/man/man1; fi
	if [ "$(OS)" != "Windows" ]; then $(GOPATH)/bin/dataset -generate-manpage | nroff -Tutf8 -man > $(GOPATH)/man/man1/dataset.1; fi

python:
	cd py && $(MAKE)

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	bash mk-website.bash

test: clean bin/dataset$(EXT)
	go test
	bash test_cmd.bash
	cd py && $(MAKE) test

clean: 
	if [ "$(PKGASSETS)" != "" ]; then bash rebuild-assets.bash; fi
	if [ -f index.html ]; then rm *.html; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d man ]; then rm -fR man; fi
	if [ -d testdata ]; then rm -fR testdata; fi
	cd py && $(MAKE) clean

man: build
	mkdir -p man/man1
	bin/dataset -generate-manpage | nroff -Tutf8 -man > man/man1/dataset.1

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dataset cmd/dataset/dataset.go cmd/dataset/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dataset.exe cmd/dataset/dataset.go cmd/dataset/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dataset cmd/dataset/dataset.go cmd/dataset/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dataset cmd/dataset/dataset.go cmd/dataset/assets.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

distribute_docs:
	if [ -d dist ]; then rm -fR dist; fi
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	bash package-versions.bash > dist/package-versions.txt

update_version:
	./update_version.py --yes

release: clean dataset.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7 python
	cd py && $(MAKE) release

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	bash mk-website.bash
	bash publish.bash

