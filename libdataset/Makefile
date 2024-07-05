#
# Simple Makefile for building C-Shared library and Python3 module.
#

LIB_NAME = libdataset

VERSION = $(shell grep '"version":' ../codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

# x86_64 for Intel, arm64 for M1, arm7 for Pi
ARCH = $(strip $(shell uname -m))

OS = $(shell uname)

EXT = .so
ifeq ($(OS), Windows)
	EXT = .dll
	OS = windows
endif
ifeq ($(OS), Darwin)
	EXT = .dylib
	OS = macos
endif
ifeq ($(OS), Linux)
	EXT = .so
	OS = linux
endif

USE_SUDO =
ifeq ($(use_sudo), true)
	USE_SUDO = sudo
endif

build: $(LIB_NAME)$(EXT)

$(LIB_NAME)$(EXT): $(LIB_NAME).go ../*.go
	go build -buildmode=c-shared -o "$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"

test: build
	python3 test_libdataset.py

clean:
	@if [ -f "$(LIB_NAME).dll" ]; then rm "$(LIB_NAME).dll"; fi
	@if [ -f "$(LIB_NAME).dylib" ]; then rm "$(LIB_NAME).dylib"; fi
	@if [ -f "$(LIB_NAME).so" ]; then rm "$(LIB_NAME).so"; fi
	@if [ -f "$(LIB_NAME)-amd64.dll" ]; then rm "$(LIB_NAME)-amd64.dll"; fi
	@if [ -f "$(LIB_NAME)-arm64.dll" ]; then rm "$(LIB_NAME)-arm64.dll"; fi
	@if [ -f "$(LIB_NAME)-amd64.dylib" ]; then rm "$(LIB_NAME)-amd64.dylib"; fi
	@if [ -f "$(LIB_NAME)-arm64.dylib" ]; then rm "$(LIB_NAME)-arm64.dylib"; fi
	@if [ -f "$(LIB_NAME)-amd32.so" ]; then rm "$(LIB_NAME)-amd32.so"; fi
	@if [ -f "$(LIB_NAME)-amd64.so" ]; then rm "$(LIB_NAME)-amd64.so"; fi
	@if [ -f "$(LIB_NAME)-arm64.so" ]; then rm "$(LIB_NAME)-arm64.so"; fi
	@if [ -f "$(LIB_NAME).h" ]; then rm "$(LIB_NAME).h"; fi
	@if [ -f "$(LIB_NAME)-js.wasm" ]; then rm "$(LIB_NAME)-js.wasm"; fi
	@if [ -d "dist" ]; then rm -fR dist; fi
	@if [ -d "testout" ]; then rm -fR testout; fi

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

# WASM code build is experimental, Python maybe able to load WASM code via wasmer-python, https://github.com/wasmerio/wasmer-python
# This would let me avoid having at have seperate machines to build a libdataset C-shared library.
wasm: $(LIB_NAME).go
	env CGO_ENABLED=1 GOOS=js GOARCH=wasm go build -o $(LIB_NAME)-js.wasm $(LIB_NAME).go

release: js-wasm $(LIB_NAME)$(EXT)
	mkdir -p dist/man/man3
	cp ../codemeta.json dist/
	cp ../CITATION.cff dist/
	cp ../LICENSE dist/
	cp -vR ../man/man3 dist/man/
	cp README.md dist/
	cp INSTALL.md dist/
	go build -buildmode=c-shared -o "$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"
	cp -v $(LIB_NAME)$(EXT) dist/
	cp -v $(LIB_NAME).h dist/
	cd dist && zip $(LIB_NAME)-v$(VERSION)-$(OS)-$(ARCH).zip $(LIB_NAME)$(EXT) $(LIB_NAME).h $(LIB_NAME)-js.wasm codemeta.json CITATION.cff README.md LICENSE INSTALL.md


.FORCE:

