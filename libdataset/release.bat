@echo off
REM
REM A simple batch file to build the c-shared library and
REM package the Python3 module from the Windows 10 command prompt.
REM
REM Requires: Go v1.23.4 or better 
REM Miniconda Python 3.7 or better.
REM Using conda: `conda install git` `conda install m2w64-gcc` `conda install tar gz`
REM
@echo on
go build -buildmode=c-shared -o "libdataset.dll" "..\libdataset\libdataset.go"
mkdir ..\dist
copy libdataset.dll ..\dist\
cd ..\dist
copy ..\README.md .\
copy ..\LICENSE .\
copy ..\INSTALL.md .\
tar zcvf libdataset-VERSION_NO-windows-amd64.tar.gz libdataset.dll README.md LICENSE INSTALL.md 