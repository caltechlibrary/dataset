@echo off
REM This is a Windows 10 Batch file for building dataset command
REM from the command prompt.
REM
REM It requires: go version 1.12.4 or better and the cli for git installed
REM
go version
echo "Getting ready to build the dataset.exe and write to .\bin"

go build -o bin\dataset.exe cmd\dataset\dataset.go

echo "You can now copy the contents of .\bin to your program directory"
