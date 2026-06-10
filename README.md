# Sitemap Generator CLI

A command line interface for my XML Sitemap Generator written in Go.

The CLI sends requests to the external sitemap generator API at `https://api.marcobeierer.com/sitemap/v2/`, which crawls the URL and generates the sitemap. 

It is free for sitemaps with up to 500 URLs. A token is required for larger sites and additional features such as image and video sitemap support. Tokens are available at [marcobeierer.com/purchase](https://www.marcobeierer.com/purchase).

## Installation

### From Source

With Go 1.16 or newer:

```
go install github.com/marcobeierer/sitemapgenerator-cli@latest
```

From a local checkout:

```
go install .
```

### Precompiled

Precompiled binaries for 64-bit Linux, macOS, and Windows systems are available in the `bin` folder of this repository. They are named `sitemapgenerator` on Linux and macOS and `sitemapgenerator.exe` on Windows.

## Usage

```
sitemapgenerator-cli <command> <url> [flags]
```

Supported commands are `run`, `download`, and `stats`. If you use a precompiled binary, replace `sitemapgenerator-cli` in the examples with `sitemapgenerator` or `sitemapgenerator.exe`.

### run

Starts or continues sitemap generation for the URL. The generated sitemap XML is written to standard output, so it can be redirected directly to a file.

```
sitemapgenerator-cli run https://www.marcobeierer.com -tokenpath token.txt > sitemap.xml
```

Supported flags:

- `-tokenpath`: path to the token file.
- `-max_fetchers`: maximum number of concurrent connections. Default: `3`.
- `-reference_count_threshold`: exclude images and videos embedded on more than this number of HTML pages. Default: `-1`.
- `-enable_index_file`: enable generation of a sitemap index file, recommended for large websites. Default: `false`.
- `-max_request_retries`: number of retries for each failed request. Default: `5`.
- `-request_retry_timeout`: timeout in seconds after a failed request. Default: `30`.
- `-sleep_time`: seconds between each update request. Default: `5`.

### download

Downloads a previously generated sitemap and sitemap index files for the URL to the given output directory.

```
sitemapgenerator-cli download https://www.marcobeierer.com -tokenpath token.txt -out_dir ./sitemaps
```

Supported flags:

- `-tokenpath`: path to the token file.
- `-out_dir`: output directory for downloaded sitemap files.

### stats

Prints sitemap generation statistics as JSON.

```
sitemapgenerator-cli stats https://www.marcobeierer.com -tokenpath token.txt
```

Supported flags:

- `-tokenpath`: path to the token file.

## Agent Skill

This repository includes an Agent Skills-compatible guide at `skills/sitemapgenerator-cli/SKILL.md` for agents that need to use this CLI. The skill documents command selection, safe token handling, common gotchas, and validation steps for `run`, `download`, and `stats`.

## Online Sitemap Generator

The sitemap generator is also available as an online tool on [my website](https://www.marcobeierer.com/tools/sitemap-generator).
