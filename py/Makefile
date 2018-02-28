#
# Simple Makefile for building C-Shared library and Python3 module.
#
PROJECT = dataset

LIB_NAME = libdataset

VERSION = $(shell grep -m 1 'Version =' ../$(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

ARCH = x86_64

OS = $(shell uname)

EXT = .so
ifeq ($(OS), Windows)
	EXT = .dll
	ARCH = x86_64
	os = windows
endif
ifeq ($(OS), Darwin)
	EXT = .dylib
	ARCH = $(shell arch)
	OS = maxosx
endif
ifeq ($(OS), Linux)
	EXT = .so
	ARCH = $(shell arch)
	OS = linux
endif
ifeq ($(ARCH), i386)
	ARCH = amd64
endif
ifeq ($(ARCH), x86_64)
	ARCH = amd64
endif


$(LIB_NAME)$(EXT): $(LIB_NAME).go
	go build -buildmode=c-shared -o "$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"

test: $(LIB_NAME)$(EXT)
	python3 dataset_test.py "test_collection.ds"

clean:
	if [ -f "$(LIB_NAME)$(EXT)" ]; then rm "$(LIB_NAME)$(EXT)"; fi
	if [ -f "$(LIB_NAME).h" ]; then rm "$(LIB_NAME).h"; fi
	if [ -d "test_collection.ds" ]; then rm -fR "test_collection.ds"; fi	
	if [ -d "__pycache__" ]; then rm -fR "__pycache__"; fi	

release: $(LIB_NAME)$(EXT)
	go build -buildmode=c-shared -o "$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"
	mkdir -p ../dist
	zip -r "../dist/py3-$(LIB_NAME)-$(VERSION)-$(OS)-$(ARCH).zip" "*.py" "$(LIB_NAME)$(EXT)" "README.md"

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

