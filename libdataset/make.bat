@echo off
REM
REM A simple batch file to build the c-shared library and
REM package the Python3 module from the Windows 10 command prompt.
REM
REM Requires: Go v1.23.1 or better 
REM Miniconda Python 3.7 or better.
REM Using conda: `conda install git` `conda install m2w64-gcc`
REM
echo Running: go build -buildmode=c-shared -o "libdataset.dll" "..\libdataset\libdataset.go"
go build -buildmode=c-shared -o "libdataset.dll" "..\libdataset\libdataset.go"
