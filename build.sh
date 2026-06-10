#!/bin/bash
set -e -u

version="$1"

if [ "$version" = "" ]; then
    version="dev"
fi

mkdir -p ./bin
rm -f ./bin/sitemapgenerator-linux-amd64 ./bin/sitemapgenerator-darwin-amd64 ./bin/sitemapgenerator-windows-amd64.exe

env GOOS=linux GOARCH=amd64 go build -a -v -o ./bin/sitemapgenerator-linux-amd64 -ldflags "-X main.version=$version"
env GOOS=darwin GOARCH=amd64 go build -a -v -o ./bin/sitemapgenerator-darwin-amd64 -ldflags "-X main.version=$version"
env GOOS=windows GOARCH=amd64 go build -a -v -o ./bin/sitemapgenerator-windows-amd64.exe -ldflags "-X main.version=$version"
