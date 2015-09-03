@echo off
set GOPATH=%cd%
go get github.com/kardianos/osext
go get github.com/satori/go.uuid
go get github.com/gorilla/schema
go build -o %GOPATH%\bin\NameTagAutoPrinting.exe %GOPATH%\src\main.go