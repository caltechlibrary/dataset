@echo off
REM This is a Windows 10 Batch file for building dataset command
REM from the command prompt.
REM
REM It requires: go version 1.12.4 or better and the cli for git installed
REM
go version
echo "Getting ready to build the dataset.exe"

go build -o dataset.exe "cmd\dataset\dataset.go" "cmd\dataset\assets.go"

echo "Checking compile should see version number of dataset"
.\dataset.exe -version

echo "If OK, you can now copy the dataset.exe to %USERPROFILE%\go\bin"
echo ""
echo "      copy dataset.exe %USERPROFILE%\go\bin
echo ""
