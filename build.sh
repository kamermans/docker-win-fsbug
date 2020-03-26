#!/bin/bash -e

BIN="fsbug"

if ! [ -d build ]; then
    mkdir build
fi

GOOS=linux GOARCH=amd64 go build -o build/$BIN-linux .
GOOS=darwin GOARCH=amd64 go build -o build/$BIN-macos .
GOOS=windows GOARCH=amd64 go build -o build/$BIN-win.exe .
