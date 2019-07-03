#!/bin/bash
set -e -u

version="$1"

if [ "$version" = "" ]; then
    version="dev"
fi

env GOOS=linux GOARCH=amd64 go build -a -v -o ./bin/linux/amd64/sitemapgenerator -ldflags "-X main.version=$version"
env GOOS=darwin GOARCH=amd64 go build -a -v -o ./bin/darwin/amd64/sitemapgenerator -ldflags "-X main.version=$version"
env GOOS=windows GOARCH=amd64 go build -a -v -o ./bin/windows/amd64/sitemapgenerator.exe -ldflags "-X main.version=$version"
