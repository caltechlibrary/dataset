#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = dataset

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PKGASSETS = $(shell which pkgassets)

PROJECT_LIST = dataset 

build: $(PROJECT_LIST)

dataset: bin/dataset bin/dsindexer bin/dsfind bin/dsws

dataset.go: assets.go

assets.go: 
	pkgassets -p dataset -o assets.go SiteDefaults defaults


bin/dataset: dataset.go attachments.go cmds/dataset/dataset.go
	go build -o bin/dataset cmds/dataset/dataset.go

bin/dsindexer: dataset.go search.go cmds/dsindexer/dsindexer.go
	go build -o bin/dsindexer cmds/dsindexer/dsindexer.go

bin/dsfind: dataset.go search.go formats.go cmds/dsfind/dsfind.go
	go build -o bin/dsfind cmds/dsfind/dsfind.go
	
bin/dsws: dataset.go search.go assets.go formats.go cmds/dsws/dsws.go
	go build -o bin/dsws cmds/dsws/dsws.go

install: $(PROJECT_LIST)
	env GOBIN=$(GOPATH)/bin go install cmds/dataset/dataset.go
	env GOBIN=$(GOPATH)/bin go install cmds/dsindexer/dsindexer.go
	env GOBIN=$(GOPATH)/bin go install cmds/dsfind/dsfind.go
	env GOBIN=$(GOPATH)/bin go install cmds/dsws/dsws.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	./mk-website.bash

test:
	go test

format:
	gofmt -w assets.go
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
	golint assets.go
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
	if [ "$(PKGASSETS)" != "" ] && [ -f assets.go ]; then rm assets.go; fi
	if [ -f index.html ]; then rm *.html; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then rm $(PROJECT)-$(VERSION)-release.zip; fi

dist/linux-amd64:
	env  GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/dataset cmds/dataset/dataset.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/dsindexer cmds/dsindexer/dsindexer.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/dsfind cmds/dsfind/dsfind.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/dsws cmds/dsws/dsws.go

dist/windows-amd64:
	env  GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/dataset.exe cmds/dataset/dataset.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/dsindexer.exe cmds/dsindexer/dsindexer.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/dsfind.exe cmds/dsfind/dsfind.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/dsws.exe cmds/dsws/dsws.go

dist/macosx-amd64:
	env  GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/dataset cmds/dataset/dataset.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/dsindexer cmds/dsindexer/dsindexer.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/dsfind cmds/dsfind/dsfind.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/dsws cmds/dsws/dsws.go

dist/raspbian-arm7:
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/dataset cmds/dataset/dataset.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/dsindexer cmds/dsindexer/dsindexer.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/dsfind cmds/dsfind/dsfind.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/dsws cmds/dsws/dsws.go

release: dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v docs/dataset.md dist/
	cp -v docs/dsindexer.md dist/
	cp -v docs/dsfind.md dist/
	cp -v docs/dsws.md dist/
	cp -v docs/operations.md dist/
	cp -v docs/file-system-layout.md dist/
	cp -v docs/package.md dist/
	cp -v how-to/import-csv-rows-as-json-documents.md dist/
	cp -v how-to/use-dataset-with-s3.md dist/
	zip -r $(PROJECT)-$(VERSION)-release.zip dist/*


status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	./mk-website.bash NO
	./publish.bash

