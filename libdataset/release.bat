@echo off
REM
REM A simple batch file to build the c-shared library and
REM package the Python3 module from the Windows 10 command prompt.
REM
REM Requires: Go v1.14 or better 
REM Miniconda Python v3.8 or better.
REM Using conda: 
REM `conda install git` 
REM `conda install m2w64-gcc` 
REM `conda install mw-zip`
REM
REM Replace %VERSION_% with the version number of the release.
REM
SET VERSION_NO=0.1.11
@echo on
go build -buildmode=c-shared -o "libdataset.dll" "..\libdataset\libdataset.go"
mkdir ..\dist
copy libdataset.dll ..\dist\
cd ..\dist
copy ..\README.md .\
copy ..\LICENSE .\
copy ..\INSTALL.md .\
zip "libdataset-%VERSION_NO%-windows-amd64.zip" libdataset.dll README.md LICENSE INSTALL.md 
