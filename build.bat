@echo off
REM This is a Windows 10 Batch file for building dataset command
REM from the command prompt.
REM
REM It requires: go version 1.23.1 or better and the cli for git installed
REM
go version
echo Getting ready to build the dataset.exe

SET PROJECT=dataset
echo Release info for %PROJECT%
echo Displaying version number from codemeta.json
@REM jq-windows-amd64 -r .version codemeta.json
jq -r .version codemeta.json
echo Enter the version number you want to release as.
SET /P DS_VERSION=
echo Displaying current hash using git log --pretty="%h" -n 1
git log --pretty="%%h" -n 1
echo Enter hash to use for release
SET /P RELEASE_HASH=
echo Today's date
date /T
echo Enter date for release
SET /P RELEASE_DATE=
echo Generating version.go using Pandoc
echo '' | pandoc --from t2t --to plain ^
                --metadata-file codemeta.json ^
                --metadata package=%PROJECT% ^
                --metadata version=%DS_VERSION% ^
                --metadata release_date=%RELEASE_DATE% ^
                --metadata release_hash=%RELEASE_HASH% ^
                --template codemeta-version-go.tmpl ^
                LICENSE >version.go
IF NOT EXIST bin MKDIR bin

echo Compiling bin\dataset.exe
go build -o bin\dataset.exe "cmd\dataset\dataset.go"
echo Compiling bin\datasetd.exe
go build -o bin\datasetd.exe "cmd\datasetd\datasetd.go"
echo Compiling bin\dsquery.exe
go build -o bin\dsquery.exe "cmd\dsquery\dsquery.go"

echo Checking compile should see version number of dataset
.\bin\dataset.exe -version
echo Checking compile should see version number of datasetd
.\bin\datasetd.exe -version
echo Checking compile should see version number of dsquery
.\bin\dsquery.exe -version

echo If OK, you can now copy the compiled programs to %USERPROFILE%\bin
echo.
echo       copy bin\dataset.exe %USERPROFILE%\bin
echo       copy bin\datasetd.exe %USERPROFILE%\bin
echo       copy bin\dsquery.exe %USERPROFILE%\bin
echo.
@echo on