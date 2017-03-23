#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = dataset

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROJECT_LIST = dataset

build: package $(PROJECT_LIST)

package: dataset.go
	go build

dataset: bin/dataset

bin/dataset: dataset.go cmds/dataset/dataset.go
	go build -o bin/dataset cmds/dataset/dataset.go

install: $(PROJECT_LIST)
	env GOBIN=$(GOPATH)/bin go install cmds/dataset/dataset.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	./mk-website.bash

test:
	go test

format:
	goimports -w dataset.go
	goimports -w dataset_test.go
	goimports -w cmds/dataset/dataset.go
	gofmt -w dataset.go
	gofmt -w dataset_test.go
	gofmt -w cmds/dataset/dataset.go

lint:
	golint dataset.go
	golint dataset_test.go
	golint cmds/dataset/dataset.go

clean:
	if [ -f index.html ]; then /bin/rm *.html; fi
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

dist/linux-amd64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/dataset cmds/dataset/dataset.go

dist/windows-amd64:
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/dataset cmds/dataset/dataset.go

dist/macosx-amd64:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/dataset cmds/dataset/dataset.go

dist/raspbian-arm7:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/dataset cmds/dataset/dataset.go


release: dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v dataset.md dist/
	zip -r $(PROJECT)-$(VERSION)-release.zip dist/*


status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

