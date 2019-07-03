# Sitemap Generator CLI
A command line interface for my XML Sitemap Generator written in Go (golang).

## Installation

### From Source
	go get -u github.com/marcobeierer/sitemapgenerator-cli
	cd $GOPATH/github.com/marcobeierer/sitemapgenerator-cli
	go install

### Precompiled
You can find precompiled binaries for 64 bit Linux, MacOS and Windows systems in the `bin` folder of this repository.

## Usage
	sitemapgenerator [flags] url

The sitemap is written to the standard output. It is thus possible to redirect the output directly to a file.

### Supported Flags
- tokenpath
	- Path to the token file
- max\_fetchers
	- Number of the maximal concurrent connections.
- reference\_count\_threshold
	- With the reference count threshold you can define that images and videos that are embedded on more than the selected number of HTML pages are excluded from the sitemap.

### Example
	sitemapgenerator -tokenpath token.txt https://www.marcobeierer.com > sitemap.xml

## Online Sitemap Generator
The sitemap generator is also available as online tool on [my website](https://www.marcobeierer.com/tools/sitemap-generator).

## Where do I get a Token?
You can use the Sitemap Generator for websites with up to 500 URL for free. If your website has more URLs, you can [purchase a token on my website](https://www.marcobeierer.com/purchase).
