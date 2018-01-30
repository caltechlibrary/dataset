#!/bin/bash

NAME="dataset"
OS=$(uname)
case "${OS}" in
    "Darwin")
    EXT=".dylib"
    ;;
    "Windows")
    EXT=".dll"
    ;;
    *)
    EXT=".so"
    ;;
esac
go build -buildmode=c-shared -o "${NAME}${EXT}" "${NAME}.go"
