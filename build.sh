#!/bin/bash
export GOPATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
go build -o $GOPATH/bin/NameTagAutoPrinting.exe $GOPATH/src/main.go