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


dataset$(EXT): bin/dataset$(EXT) bin/dsws$(EXT)

cmd/dataset/assets.go:
	pkgassets -o cmd/dataset/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples examples/dataset Help docs/dataset
	git add cmd/dataset/assets.go

cmd/dsws/assets.go:
	pkgassets -o cmd/dsws/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples examples/dsws Help docs/dsws 
	git add cmd/dsws/assets.go

cmd/dsws/templates.go:
	pkgassets -o cmd/dsws/templates.go -p main Defaults defaults
	git add cmd/dsws/templates.go

bin/dataset$(EXT): dataset.go attachments.go repair.go sort.go gsheet/gsheet.go cmd/dataset/dataset.go cmd/dataset/assets.go
	go build -o bin/dataset$(EXT) cmd/dataset/dataset.go cmd/dataset/assets.go

bin/dsws$(EXT): dataset.go search.go formats.go cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go
	go build -o bin/dsws$(EXT) cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go

build: $(PROJECT_LIST)

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/dataset/dataset.go cmd/dataset/assets.go
	env GOBIN=$(GOPATH)/bin go install cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go

python:
	cd py && $(MAKE)

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	bash mk-website.bash

test: bin/dataset$(EXT) bin/dsws$(EXT)
	go test
	cd gsheet && go test -client-secret="../etc/client_secret.json" -spreadsheet-id="1y23sLVy4rfL2U81kYhOYG6x3dTxnexqJcVBasIsyEx8"
	bash test_cmd.bash
	cd py && $(MAKE) test

format:
	gofmt -w dataset.go
	gofmt -w dataset_test.go
	gofmt -w attachments.go
	gofmt -w attachments_test.go
	gofmt -w search.go
	gofmt -w search_test.go
	gofmt -w formats.go
	gofmt -w cmd/dataset/dataset.go

lint:
	golint dataset.go
	golint dataset_test.go
	golint attachments.go
	golint attachments_test.go
	golint search.go
	golint search_test.go
	golint formats.go
	golint cmd/dataset/dataset.go

clean: 
	if [ "$(PKGASSETS)" != "" ]; then bash rebuild-assets.bash; fi
	if [ -f index.html ]; then rm *.html; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d testdata ]; then rm -fR testdata; fi
	cd py && $(MAKE) clean

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dataset cmd/dataset/dataset.go cmd/dataset/assets.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dsws cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dataset.exe cmd/dataset/dataset.go cmd/dataset/assets.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dsws.exe cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dataset cmd/dataset/dataset.go cmd/dataset/assets.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dsws cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dataset cmd/dataset/dataset.go cmd/dataset/assets.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dsws cmd/dsws/dsws.go cmd/dsws/assets.go cmd/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

distribute_docs:
	rm -fR dist
	mkdir -p dist/how-to
	mkdir -p dist/docs
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v docs/*.md dist/docs/
	cp -v how-to/*.md dist/how-to/
	cp -vR demos dist/
	bash package-versions.bash > dist/package-versions.txt

update_patch_no:
	./update_version.py --yes

release: dataset.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7
	cd py && $(MAKE) release

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	bash mk-website.bash
	bash publish.bash

