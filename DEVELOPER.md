
# Compilation Notes

On Unix-like systems (e.g. Darwin, Linux) building dataset, libdataset and python is generally as easy as running the GNU Make command. On Windows you need to take a more Window-ish approach.

## Windows 10

+ Install Go 1.12.x via the Windows' installer available from https://golang.org/downloads
+ Install Miniconda (Python 3.7) from https://docs.conda.io/en/latest/miniconda.html

Now we ready to finish our setup via the "Anaconda Prompt".

+ install git and gcc
+ run `go get -u github.com\caltechlibrary\dataset`
+ run `go get -u github.com\caltechlibrary\cli`
+ run the make.bat and python setup.py install

Here's an example of what I've done after openning the "Anaconda Prompt"

```
    cd %USERPROFILE%
    mkdir go
    mkdir go\src
    conda install git
    conda install m2w64-gcc
    go get -u github.com\caltechlibrary\dataset
    go get -u github.com\caltechlibrary\cli
    cd go\src\caltechlibrary\dataset
    .\make.bat
    move dataset.exe "%USERPROFILE\go\bin\dataset.exe"
    cd py
    .\make.bat
    python setup.py install
    cd %USERPROFILE%
```

Both the dataset command line should now work from the Anaconda Prompt
as well as the dataset Python module being available to Python.

