
# Compilation Notes

On Unix-like systems (e.g. Darwin, Linux) building dataset and libdataset is generally as easy as running the GNU Make command. On Windows you need to take a more Window-ish approach and run `make.bat`.

## Windows 10

+ Install Go 1.16 via the Windows' installer available from https://golang.org/downloads

+ install git and gcc via Miniconda
+ run `go get -u github.com\caltechlibrary\dataset`
+ run `go get -u github.com\caltechlibrary\cli`
+ run the make.bat 

Here's an example of what I've done after opening the "Anaconda Prompt"

```
    cd %USERPROFILE%
    mkdir go
    mkdir go\bin
    mkdir go\src
    conda install git
    conda install m2w64-gcc
    go get -u github.com\caltechlibrary\dataset
    go get -u github.com\caltechlibrary\cli
    cd go\src\caltechlibrary\dataset
    .\make.bat
    move dataset.exe "%USERPROFILE\go\bin\dataset.exe"
```

To build the DLL

```
    cd libdataset
    .\make.bat
```

Both the dataset command line exe and DLL will need to be copied
to where they available to your windows applications.

