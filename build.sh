#!/bin/sh
BIN=build/app
rm -f $BIN
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $BIN github.com/thlcodes/minrevpro/cmd/minrevpro
stat -f%z $BIN | awk '{printf "%.2f MB", $1/(1024*1024)}'
cp ./start.sh build/start.sh