@echo off
set GOPATH=%cd%
go build -o %GOPATH%\bin\NameTagAutoPrinting.exe %GOPATH%\src\main.go