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

build: $(LIB_NAME)$(EXT)

$(LIB_NAME)$(EXT): $(LIB_NAME).go ../dataset.go
	go build -buildmode=c-shared -o "$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"

test: $(LIB_NAME)$(EXT)
	python3 test_libdataset.py

clean:
	if [ -f "$(LIB_NAME).dll" ]; then rm "$(LIB_NAME).dll"; fi
	if [ -f "$(LIB_NAME).dylib" ]; then rm "$(LIB_NAME).dylib"; fi
	if [ -f "$(LIB_NAME).so" ]; then rm "$(LIB_NAME).so"; fi
	if [ -f "$(LIB_NAME).h" ]; then rm "$(LIB_NAME).h"; fi
	if [ -f "libdataset/$(LIB_NAME).dll" ]; then rm "libdataset/$(LIB_NAME).dll"; fi
	if [ -f "libdataset/$(LIB_NAME).dylib" ]; then rm "libdataset/$(LIB_NAME).dylib"; fi
	if [ -f "libdataset/$(LIB_NAME).so" ]; then rm "libdataset/$(LIB_NAME).so"; fi
	if [ -f "libdataset/$(LIB_NAME).h" ]; then rm "libdataset/$(LIB_NAME).h"; fi

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

release: $(LIB_NAME)$(EXT)
	mkdir -p ../dist/
	cp ../LICENSE ../dist/
	cp ../README.md ../dist/
	cp ../INSTALL.md ../dist/
	go build -buildmode=c-shared -o "$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"
	cp -v $(LIB_NAME)$(EXT) ../dist/
	cp -v $(LIB_NAME).h ../dist/
	cd ../dist && zip $(LIB_NAME)-$(VERSION)-$(OS)-$(ARCH).zip $(LIB_NAME)$(EXT) $(LIB_NAME).h  README.md LICENSE INSTALL.md



.FORCE:
