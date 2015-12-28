# Sitemap Generator CLI
A command line interface for my sitemap generator written in Go (golang).

## Installation
	go get -u github.com/webguerilla/sitemapgenerator-cli
	cd $GOPATH/github.com/webguerilla/sitemapgenerator-cli
	go install

## Usage
	sitemapgenerator [flags] url

### Supported Flags
- tokenpath
	- Path to the token file

### Example
	sitemapgenerator -tokenpath token.txt https://www.marcobeierer.com

## Online Sitemap Generator
The sitemap generator is also available as online tool on [my website](https://www.marcobeierer.com/tools/sitemap-generator).
