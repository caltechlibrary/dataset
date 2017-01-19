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
	go install cmds/dataset/dataset.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css htdocs/index.md
	mkpage "content=htdocs/index.md" templates/default/index.html > htdocs/index.html
	./mk-website.bash

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

test:
	go test

clean:
	if [ -f index.html ]; then /bin/rm *.html; fi
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT_NAME)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT_NAME)-$(VERSION)-release.zip; fi

release:
	./mk-release.bash

status:
	git status

save:
	git commit -am "Quick save"
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

