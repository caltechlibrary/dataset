#
# Simple Makefile for building C-Shared library 
# on a Linux amd64 Linux cross compiled to
# darwin amd64 and windows amd64. Based on the 
# Thomas Pöchtrager's tools setup at 
#
# + https://github.com/tpoechtrager/wclang
# + https://github.com/tpoechtrager/osxcross
#
#
PROJECT = dataset

LIB_NAME = libdataset

VERSION = $(shell grep -m 1 'Version =' ../$(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

ARCH = x86_64

OS = $(shell uname)

# Darwin cross compile tools
DARWIN_TOOLS = $(HOME)/src/github.com/tpoechtrager/osxcross/target/bin

# Windows cross compile tools
WINDOWS_TOOLS = $(HOME)/src/github.com/tpoechtrager/wclang/target/bin

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

build: darwin-amd64 windows-amd64

darwin-amd64: 
	env PATH="$(DARWIN_TOOLS):$(PATH)" CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC=o64-clang go build -buildmode=c-shared -o "$(LIB_NAME).dynlib" "$(LIB_NAME).go"

windows-amd64: 
	env PATH="$(WINDOWS_TOOLS):$(PATH)" CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=w64-clang go build -buildmode=c-shared -o "$(LIB_NAME).dll" "$(LIB_NAME).go"

clean:
	if [ -f "$(LIB_NAME).so" ]; then rm "$(LIB_NAME).so"; fi
	if [ -f "$(LIB_NAME).dll" ]; then rm "$(LIB_NAME).dll"; fi
	if [ -f "$(LIB_NAME).dynlib" ]; then rm "$(LIB_NAME).dynlib"; fi
	if [ -f "$(LIB_NAME).h" ]; then rm "$(LIB_NAME).h"; fi

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

release: release-darwin-amd64 release-windows-amd64

release-darwin-amd64: darwin-amd64 FORCE
	mkdir -p ../dist/
	cp ../LICENSE ../dist/
	cp ../README.md ../dist/
	cp ../INSTALL.md ../dist/
	cp -v $(LIB_NAME).dynlib ../dist/
	cp -v $(LIB_NAME).h ../dist/
	cd ../dist && tar zcvf $(LIB_NAME)-$(VERSION)-darwin-amd64.tar.gz $(LIB_NAME).dynlib $(LIB_NAME).h  README.md LICENSE INSTALL.md

release-windows-amd64: windows-amd64 FORCE
	mkdir -p ../dist/
	cp ../LICENSE ../dist/
	cp ../README.md ../dist/
	cp ../INSTALL.md ../dist/
	cp -v $(LIB_NAME).dll ../dist/
	cp -v $(LIB_NAME).h ../dist/
	cd ../dist && tar zcvf $(LIB_NAME)-$(VERSION)-windows-amd64.tar.gz $(LIB_NAME).dll $(LIB_NAME).h  README.md LICENSE INSTALL.md

FORCE:
