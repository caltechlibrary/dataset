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

if [ -f "${NAME}.h" ]; then
    rm "${NAME}.h" 
fi
if [ -f "${NAME}${EXT}" ]; then
    rm "${NAME}${EXT}"
fi
