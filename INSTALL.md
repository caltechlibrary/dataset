Installation for development of **dataset**
===========================================

**dataset** The Dataset Project provides tools for working with collections of JSON documents easily. It uses a simple key and object pair to organize JSON documents into a collection. It supports SQL querying of the objects stored in a collection.

It is suitable for temporary storage of JSON objects in data processing pipelines as well as a persistent storage mechanism for collections of JSON objects.

The Dataset Project provides command line programs and a web service for working with JSON objects as a collection or individual objects. As such it is well suited for data science projects as well as building web applications that work with metadata.

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the following command to install latest table release. This may work for macOS, Linux and if you’re using Windows with the Unix subsystem. This would be run from your shell (e.g. Terminal on macOS).

~~~shell
curl https://caltechlibrary.github.io/dataset/installer.sh | sh
~~~

This will install the programs included in dataset in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command below.

~~~ps1
irm https://caltechlibrary.github.io/dataset/installer.ps1 | iex
~~~

If you are running macOS please see, [INSTALL_NOTES_macOS.md](INSTALL_NOTES_macOS.md) for details on dealing with unsigned executables.

If you are running Windows please see, [INSTALL_NOTES_Windows.md](INSTALL_NOTES_Windows.md) for details on dealing with unsigned executables.

Installing from source
----------------------

### Required software

- Golang &gt;&#x3D; 1.24.5
- CMTools &gt;&#x3D; 0.0.35

### Steps

1. git clone https://github.com/caltechlibrary/dataset
2. Change directory into the `dataset` directory
3. Make to build, test and install

~~~shell
git clone https://github.com/caltechlibrary/dataset
cd dataset
make
make test
make install
~~~

