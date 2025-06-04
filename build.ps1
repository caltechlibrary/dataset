#
# Build Dataset command line tools on Windows in PowerShell
#
if (Test-Path -Path .\bin) {
  del bin\*.exe
} else {
  mkdir bin
}
go build -o bin\dataset.exe cmd\dataset\dataset.go
go build -o bin\datasetd.exe cmd\datasetd\datasetd.go
go build -o bin\dsquery.exe cmd\dsquery\dsquery.go
go build -o bin\dsimporter.exe cmd\dsimporter\dsimporter.go
