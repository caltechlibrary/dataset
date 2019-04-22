@echo off
REM
REM A simple batch file to build the c-shared library and
REM package the Python3 module from the Windows 10 command prompt.
REM
REM Requires: Go v1.23.4 or better 
REM Miniconda Python 3.7 or better.
REM Using conda: `conda install git` `conda install m2w64-gcc`
REM

REM python3 setup.py install --user --record files.txt

go build -buildmode=c-shared -o "dataset\libdataset.dll" "..\libdataset\libdataset.go"
python setup.py sdist
echo "Check .\dist to see if this worked!"
dir dist\
echo "Ready to copy dist\dataset-v0.0.X.tar.gz to ..\dist\py3-dataset-v0.0.X-windows-amd64.tar.gz"

