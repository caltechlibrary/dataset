@echo off
REM This is a Windows 10 Batch file for building dataset command
REM from the command prompt.
REM
REM It requires: go version 1.12.4 or better and the cli for git installed
REM
go version
echo Getting ready to build the dataset.exe

echo Using jq to extract version string from codemeta.json
jq .version codemeta.json > version.txt
git log --pretty=format:'%h' -n 1 >hash.txt
date >release_date.txt
SET /P DS_VERSION= < version.txt
SET PROJECT=dataset
SET RELEASE_DATE= < release_date.txt
SET RELEASE_HASH= < hash.txt
echo Generating version.go using Pandoc and jq
echo '' | pandoc --from t2t --to plain \
                --metadata-file codemeta.json \
                --metadata package=%PROJECT% \
                --metadata version=%DS_VERSION% \
                --metadata release_date=%RELEASE_DATE% \
                --metadata release_hash=%RELEASE_HASH% \
                --template codemeta-version-go.tmpl \
                LICENSE >version.go
DEL version.txt
DEL release_date.txt
DEL hash.txt
echo Building version: %DS_VERSION%
echo package datasdet > version.go
echo.  >> version.go
echo // Version of package >> version.go
echo const Version = %DS_VERSION% >> version.go
MKDIR bin
go build -o bin\dataset.exe "cmd\dataset\dataset.go" "cmd\dataset\assets.go"

echo Checking compile should see version number of dataset
.\bin\dataset.exe -version

echo If OK, you can now copy the dataset.exe to %USERPROFILE%\go\bin
echo.
echo       copy bin\dataset.exe %USERPROFILE%\go\bin
echo.
