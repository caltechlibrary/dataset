#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = dataset

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PKGASSETS = $(shell which pkgassets)

PROJECT_LIST = dataset


dataset: bin/dataset bin/dsindexer bin/dsfind bin/dsws

cmds/dataset/assets.go:
	pkgassets -o cmds/dataset/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples examples/dataset Help docs/dataset
	git add cmds/dataset/assets.go

cmds/dsindexer/assets.go:
	pkgassets -o cmds/dsindexer/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples examples/dsindexer Help docs/dsindexer
	git add cmds/dsindexer/assets.go

cmds/dsfind/assets.go:
	pkgassets -o cmds/dsfind/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples examples/dsfind Help docs/dsfind
	git add cmds/dsfind/assets.go

cmds/dsws/assets.go:
	pkgassets -o cmds/dsws/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Examples examples/dsws Help docs/dsws 
	git add cmds/dsws/assets.go

cmds/dsws/templates.go:
	pkgassets -o cmds/dsws/templates.go -p main Defaults defaults
	git add cmds/dsws/templates.go


bin/dataset: dataset.go attachments.go repair.go cmds/dataset/dataset.go cmds/dataset/assets.go
	go build -o bin/dataset cmds/dataset/dataset.go cmds/dataset/assets.go

bin/dsindexer: dataset.go search.go cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go
	go build -o bin/dsindexer cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go

bin/dsfind: dataset.go search.go formats.go cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	go build -o bin/dsfind cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	
bin/dsws: dataset.go search.go formats.go cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go
	go build -o bin/dsws cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go

build: $(PROJECT_LIST)

install: 
	env GOBIN=$(GOPATH)/bin go install cmds/dataset/dataset.go cmds/dataset/assets.go
	env GOBIN=$(GOPATH)/bin go install cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go
	env GOBIN=$(GOPATH)/bin go install cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	env GOBIN=$(GOPATH)/bin go install cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	./mk-website.bash

test:
	go test
	cd gsheets && go test

format:
	gofmt -w dataset.go
	gofmt -w dataset_test.go
	gofmt -w attachments.go
	gofmt -w attachments_test.go
	gofmt -w search.go
	gofmt -w search_test.go
	gofmt -w formats.go
	gofmt -w cmds/dataset/dataset.go
	gofmt -w cmds/dsindexer/dsindexer.go
	gofmt -w cmds/dsfind/dsfind.go

lint:
	golint dataset.go
	golint dataset_test.go
	golint attachments.go
	golint attachments_test.go
	golint search.go
	golint search_test.go
	golint formats.go
	golint cmds/dataset/dataset.go
	golint cmds/dsindexer/dsindexer.go
	golint cmds/dsfind/dsfind.go

clean: 
	if [ "$(PKGASSETS)" != "" ]; then rm cmds/*/assets.go; fi
	if [ "$(PKGASSETS)" != "" ]; then rm cmds/*/templates.go; fi
	if [ -f index.html ]; then rm *.html; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dataset cmds/dataset/dataset.go cmds/dataset/assets.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dsindexer cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dsfind cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/dsws cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dataset.exe cmds/dataset/dataset.go cmds/dataset/assets.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dsindexer.exe cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dsfind.exe cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/dsws.exe cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dataset cmds/dataset/dataset.go cmds/dataset/assets.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dsindexer cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dsfind cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/dsws cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dataset cmds/dataset/dataset.go cmds/dataset/assets.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dsindexer cmds/dsindexer/dsindexer.go cmds/dsindexer/assets.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dsfind cmds/dsfind/dsfind.go cmds/dsfind/assets.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/dsws cmds/dsws/dsws.go cmds/dsws/assets.go cmds/dsws/templates.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* demos/*
	rm -fR dist/bin

distribute_docs:
	rm -fR dist
	mkdir -p dist/how-to
	mkdir -p dist/docs
	mkdir -p dist/demos
	mkdir -p dist/examples
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v docs/*.md dist/docs/
	cp -v how-to/*.md dist/how-to/
	cp -v examples/*.md dist/examples/
	cp -vR demos dist/

release: dataset.go assets.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	./mk-website.bash NO
	./publish.bash

