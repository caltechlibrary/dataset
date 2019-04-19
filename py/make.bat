@echo off
REM
REM A simple batch file to build the c-shared library and
REM package the Python3 module from the Windows 10 command prompt.
REM
REM Requires: Go v1.23.4 or better 
REM Miniconda Python 3.7 or better.
REM Cygwin installed for gcc.
REM


REM python3 setup.py install --user --record files.txt

go build -buildmode=c-shared -o "libdataset/libdataset.dll" "../libdataset/libdataset.go
python setup.py sdist
echo "Check .\dist to see if this worked!"
dir dist\
