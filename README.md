# Sitemap Generator CLI
A command line interface for my XML sitemap generator written in Go (golang).

## Installation
	go get -u github.com/webguerilla/sitemapgenerator-cli
	cd $GOPATH/github.com/webguerilla/sitemapgenerator-cli
	go install

## Usage
	sitemapgenerator [flags] url

The sitemap is written to the standard output. It is thus possible to redirect the output directly to a file.

### Supported Flags
- tokenpath
	- Path to the token file

### Example
	sitemapgenerator -tokenpath token.txt https://www.marcobeierer.com > sitemap.xml

## Online Sitemap Generator
The sitemap generator is also available as online tool on [my website](https://www.marcobeierer.com/tools/sitemap-generator).
