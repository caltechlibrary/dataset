@echo off
REM
REM A simple batch file to build the c-shared library and
REM package the Python3 module from the Windows 10 command prompt.
REM
REM Requires: Go v1.23.1 or better 
REM Miniconda Python v3.8 or better.
REM Using conda: 
REM `conda install git` 
REM `conda install m2w64-gcc` 
REM `conda install mw-zip`
REM
REM Replace %VERSION_% with the version number of the release.
REM
echo Displaying version number from codemeta.json
@REM jq-windows-amd64 -r .version ..\codemeta.json
jq -r .version ..\codemeta.json
echo Enter the version number you want to release as without "v" prefix.
SET /P DS_VERSION=
IF [%DS_VERSION%] == [] SET DS_VERSION=0.0.0
echo Using Version number %DS_VERSION%
@echo on
IF NOT EXIST libdataset.dll go build -buildmode=c-shared -o "libdataset.dll" "libdataset.go"
IF NOT EXIST dist MKDIR dist
COPY libdataset.dll dist\
COPY ..\codemeta.json dist\
COPY ..\CITATION.cff dist\
COPY ..\README.md dist\
COPY ..\LICENSE dist\
COPY ..\INSTALL.md dist\
CD dist
zip "libdataset-v%DS_VERSION%-windows-amd64.zip" libdataset.dll README.md LICENSE INSTALL.md
CD ..
DIR dist\*.zip
