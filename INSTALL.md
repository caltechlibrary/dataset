Installation
============

__dataset__ is a set of command line programs run from a shell like Bash.
It is designed for single user, single process management of a JSON
object documents as a collection where JSON documents are referenced
by a unique identifier or key.  __datasetd__ is a web service which
serves a similar purpose but is intended for supporting multi-user
and multi-processes.

This is generalized instructions for a release.  For deployment suggestions
see NOTES.md

Quick install with curl
-----------------------

There is an experimental installer.sh script that can be run with the
following command to install lastest table release. This may work for
macOS, Linux and if you're using Windows with the Unix subsystem. This
would be run from your shell (e.g. Terminal on macOS).

~~~
curl https://caltechlibrary.github.io/dataset/installer.sh | sh
~~~

This will install dataset and datasetd in your `$HOME/bin` directory.

Compiled version
----------------

Compiled versions are available for macOS (Intel and M1), Linux (Intel),
Windows (Intel and ARM64) and Raspberry Pi (ARM).

VERSION_NUMBER is a [semantic version number](http://semver.org/) (e.g. v2.0.0)


For all the released version go to the project page on GitHub and click
latest release

>    https://github.com/caltechlibrary/dataset/releases/latest


| Platform         | Zip Filename                             |
|------------------|------------------------------------------|
| Windows (Intel)  | dataset-VERSION_NUMBER-windows-amd64.zip |
| Windows (ARM 64) | dataset-VERSION_NUMBER-windows-arm64.zip |
| macOS (Intel)    | dataset-VERSION_NUMBER-macos-amd64.zip   |
| macOS (M1)       | dataset-VERSION_NUMBER-macos-arm64.zip   |
| Linux (Intel)    | dataset-VERSION_NUMBER-linux-amd64.zip   |
| Raspberry Pi OS (ARM7) | dataset-VERSION_NUMBER-raspbian-arm7.zip |


The basic recipe
----------------

- Find the Zip file listed matching the architecture you're running and download it
    - (e.g. if you're on a Windows 10 laptop/Surface with a amd64 style CPU you'd choose the Zip file with "windows-amd64" in the name).
- Download the zip file and unzip the file.
- Copy the contents of the folder named "bin" to a folder that is in your path
    - (e.g. "bin" in your "HOME" directory is common).
- Adjust your PATH if needed
- Test


### macOS

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to "bin" folder in HOME folder (or another folder in your PATH)
4. Make sure the new location in in our path
5. Test

Here's an example of the commands run in the Terminal App after
downloading the zip file.

```shell
    cd Downloads/
    unzip dataset-*-macos-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    dataset -version
```

### Windows

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to the "bin" directory in your "HOME" directory (or a folder in your path)
4. Test

#### Intel Hardware

Here's an example of the commands run in from the Bash shell on Windows 11
after downloading the zip file (assume Linux Subsystem for Windows).

```shell
    cd Downloads/
    unzip dataset-*-windows-amd64.zip
    mkdir -p $HOME/bin
    mv -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    dataset -version
```

#### ARM64 Hardware

Here's an example of the commands run in from the Bash shell on Windows 11
after downloading the zip file (assumes Linux Subsystem for Windows).

```shell
    cd Downloads/
    unzip dataset-*-windows-arm64.zip
    mkdir -p $HOME/bin
    mv -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    dataset -version
```


### Linux

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to the "bin" directory in your "HOME" directory.
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip dataset-*-linux-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    dataset -version
```


### Raspberry Pi OS

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM 7
support). Testing is done on Raspberry Pi 4 B devices using 32bit
Raspberry Pi OS.

1. Download the zip file
2. Unzip the zip file
3. Copy the executable to `$HOME/bin` (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip dataset-*-raspbian-arm7.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    dataset -version
```


## Compiling from source

Clone the repository and then compile.

```shell
    cd
    git clone https://github.com/caltechlibrary/dataset
    cd dataset
    make
    # Add any missing dependencies you might need in your Go environment
    make test
    make install
```

### Requirements

- Go version 1.19.2 or better
- Pandoc version 2.19.2 or better
- GNU Make
- Common POSIX/Unix utilities, e.g. cat, sed, grep

### Windows compilation

The tool chain to compile on Windows make several assumptions.

1. You're using Anaconda shell and have the C tool chain installed for
   cgo to work
2. GNU Make, cat, grep and sed
3. You have the latest Go installed

Since I don't assume a POSIX shell environment on windows I have made
batch files to perform some of what Make under Linux and macOS would do.

- make.bat builds our application and depends on go and jq commands
- release.bat builds a release, will prompt for version
- clean.bat removes executable and temp files


Compilation assumes [go](https://github.com/golang/go) v1.19.2 or better.

