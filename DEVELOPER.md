
# Compilation Notes

On Unix-like systems (e.g. Darwin, Linux, Windows with the Linux subsystem enabled) building dataset and datasetd is generally as easy as running the GNU Make command. On Windows without the Linux subsystem you need to take a more Window-ish approach and run `make.bat`.

## Windows 11

+ Install Go 1.21.1 via the Windows' installer available from https://golang.org/downloads
+ Install git
+ Run `go get -u github.com\caltechlibrary\dataset`
+ Change into the dataset directory
+ Run the make.bat 

Here's an example of what I've done after opening a command window

```
    cd %USERPROFILE%
    mkdir go
    mkdir go\bin
    mkdir go\src
    go get -u github.com\caltechlibrary\dataset
    cd go\src\caltechlibrary\dataset
    .\make.bat
    move dataset.exe "%USERPROFILE\go\bin\dataset.exe"
```

The dataset command line exe will likely need to be copied
to where your windows command line applications at located.

