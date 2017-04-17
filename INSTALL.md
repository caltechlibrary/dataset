
# Installation

This is generalized instructions for a release.  For deployment suggestions see NOTES.md

## Compiled version

*dataset* is a set of command line programs run from a shell like Bash. It allows you to organize JSON documents
into a collection by unique filename. 

Compiled versions are available for Mac OS X (amd64 processor), Linux (amd64), Windows (amd64) and Rapsberry Pi (both ARM6 and ARM7)

VERSION_NUMBER is a [symantic version number](http://semver.org/) (e.g. v0.1.2)

### Mac OS X

1. Download **dataset-VERSION_NUMBER-release.zip** from https://github.com/caltechlibrary/dataset/releases/latest
2. Open a finder window, find and unzip **dataset-VERSION_NUMBER-release.zip**
3. Look in the unziped folder and find the files in *dist/macosx-amd64/*
4. Drag (or copy) *dataset* to a "bin" directory in your path
5. Open and "Terminal" and run `dataset -h` to confirm you were successful

### Windows

1. Download **dataset-VERSION_NUMBER-release.zip** from https://github.com/caltechlibrary/dataset/releases/latest
2. Open the file manager find and unzip **dataset-VERSION_NUMBER-release.zip**
3. Look in the unziped folder and find the files in *dist/windows-amd64/*
4. Drag (or copy) *dataset.exe* to a directory in your path
5. Open Bash and and run `dataset -h` to confirm you were successful

### Linux

1. Download **dataset-VERSION_NUMBER-release.zip** from https://github.com/caltechlibrary/dataset/releases/latest
2. Find and unzip **dataset-VERSION_NUMBER-release.zip**
3. In the unziped directory and find the files in *dist/linux-amd64/*
4. Copy *dataset* to a "bin" directory (e.g. cp ~/Downloads/dataset-VERSION_NUMBER-release/dist/linux-amd64/dataset ~/bin/)
5. From the shell prompt run `dataset -h` to confirm you were successful

### Raspberry Pi

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM 7 support).

1. Download **dataset-VERSION_NUMBER-release.zip** from https://github.com/caltechlibrary/dataset/releases/latest
2. Find and unzip **dataset-VERSION_NUMBER-release.zip**
3. In the unziped directory and find the files in *dist/rasbian-arm7/*
4. Copy *dataset* to a "bin" directory (e.g. cp ~/Downloads/dataset-VERSION_NUMBER-release/dist/rasbian-arm7/dataset ~/bin/)
5. From the shell prompt run `dataset -h` to confirm you were successful


## Compiling from source

_dataset_ is "go gettable".  Use the "go get" command to download the dependant packages
as well as _dataset_'s source code.

```shell
    go get -u github.com/caltechlibrary/dataset
```

And then compile

```shell
    cd $GOPATH/src/github.com/caltechlibrary/dataset
    make
    make test
    make install
```


