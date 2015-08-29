#!/bin/bash
export GOPATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
go get github.com/kardianos/osext
go get github.com/satori/go.uuid
go build -o $GOPATH/bin/NameTagAutoPrinting.exe $GOPATH/src/main.go