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
	OS = macosx
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

USE_SUDO =
ifeq ($(use_sudo), true)
	USE_SUDO = sudo
endif

build: $(PROJECT)/$(LIB_NAME)$(EXT)

install: $(PROJECT)/$(LIB_NAME)$(EXT)
	$(USE_SUDO) python3 setup.py install --record files.txt


$(PROJECT)/$(LIB_NAME)$(EXT): ../$(LIB_NAME)/$(LIB_NAME).go
	go build -buildmode=c-shared -o "$(PROJECT)/$(LIB_NAME)$(EXT)" "../$(LIB_NAME)/$(LIB_NAME).go"

clean:
	if [ -f "MANIFEST" ]; then $(USE_SUDO) rm MANIFEST; fi
	if [ -f "$(PROJECT)/$(LIB_NAME)$(EXT)" ]; then $(USE_SUDO) rm "$(PROJECT)/$(LIB_NAME)$(EXT)"; fi
	if [ -f "$(PROJECT)/$(LIB_NAME).h" ]; then $(USE_SUDO) rm "$(PROJECT)/$(LIB_NAME).h"; fi
	if [ -f "test_index_map.json" ]; then $(USE_SUDO) rm "test_index_map.json"; fi
	if [ -d "dist" ]; then $(USE_SUDO) rm -fR dist; fi
	if [ -d "build" ]; then $(USE_SUDO) rm -fR build; fi
	if [ -d "test_collection.ds" ]; then $(USE_SUDO) rm -fR "test_collection.ds"; fi	
	if [ -d "test_gsheet.ds" ]; then $(USE_SUDO) rm -fR "test_gsheet.ds"; fi	
	if [ -d "test_issue43.ds" ]; then $(USE_SUDO) rm -fR "test_issue43.ds"; fi	
	if [ -f "test_issue43.csv" ]; then $(USE_SUDO) rm -fR "test_issue43.csv"; fi	
	if [ -d "test_check_and_repair.ds" ]; then $(USE_SUDO) rm -fR "test_check_and_repair.ds"; fi
	if [ -d "test_index.bleve" ]; then $(USE_SUDO) rm -fR "test_index.bleve"; fi
	if [ -d "__pycache__" ]; then $(USE_SUDO) rm -fR "__pycache__"; fi	
	if [ -d "$(PROJECT)/__pycache__" ]; then $(USE_SUDO) rm -fR "$(PROJECT)/__pycache__"; fi	
	if [ -d "$(PROEJCT).egg-info" ]; then $(USE_SUDO) rm -fR "$(PROJECT).egg-info"; fi

test: clean $(PROJECT)/$(LIB_NAME)$(EXT)
	python3 dataset_test.py

release: clean $(PROJECT)/$(LIB_NAME)$(EXT)
	go build -buildmode=c-shared -o "$(PROJECT)/$(LIB_NAME)$(EXT)" "../$(LIB_NAME)/$(LIB_NAME).go"
	python3 setup.py sdist
	mkdir -p ../dist
	cp -v dist/$(PROJECT)-$(VERSION).tar.gz "../dist/py3-$(PROJECT)-$(VERSION)-$(OS)-$(ARCH).tar.gz"

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

