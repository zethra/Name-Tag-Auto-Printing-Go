#!/bin/bash
export GOPATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
go get github.com/kardianos/osext
go get github.com/satori/go.uuid
go get github.com/gorilla/schema
go build -o $GOPATH/bin/NameTagAutoPrinting $GOPATH/src/main.go
