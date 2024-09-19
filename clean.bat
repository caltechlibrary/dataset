@echo off
REM This is a Windows 10 Batch file for building dataset command
REM from the command prompt.
REM
REM It requires: go version 1.23.1 or better and the cli for git installed
REM
DEL /S bin
RMDIR /S bin
DEL *.exe
DEL libdataset\*.dll
