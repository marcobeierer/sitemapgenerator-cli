---
name: sitemapgenerator-cli
description: Use when invoking the Sitemap Generator CLI, `sitemapgenerator-cli`, `sitemapgenerator`, or `sitemapgenerator.exe` to run crawls, download generated XML sitemap files, or inspect sitemap generation stats.
compatibility: Requires this Go CLI or a precompiled binary and network access to https://api.marcobeierer.com/sitemap/v2/.
---

# Sitemap Generator CLI

Use this skill when a user asks an agent to generate, download, inspect, or troubleshoot XML sitemaps with this repository's CLI.

## Quick Workflow

1. Confirm the target URL is intentional and authorized for crawling.
2. Choose the command:
   - `run` starts or continues sitemap generation and writes final XML to stdout.
   - `download` downloads a previously generated `sitemap.xml` plus sitemap index files into an output directory.
   - `stats` prints generation statistics as JSON.
3. Choose the executable:
   - From a local checkout, prefer `go run . <command> <url> [flags]`.
   - If installed from source, use `sitemapgenerator-cli <command> <url> [flags]`.
   - If using precompiled binaries, use `bin/linux/amd64/sitemapgenerator`, `bin/darwin/amd64/sitemapgenerator`, or `bin/windows/amd64/sitemapgenerator.exe` as appropriate.
4. Put flags after the URL. The CLI expects `<command>` first, `<url>` second, then command flags.
5. Keep token files private. Pass token file paths with `-tokenpath`; do not print token contents.

## Command Recipes

Generate a sitemap and save the final XML:

```bash
go run . run https://www.example.com -tokenpath token.txt > sitemap.xml
```

Run without a token for a small/free sitemap:

```bash
go run . run https://www.example.com > sitemap.xml
```

Generate a large sitemap with index files enabled:

```bash
go run . run https://www.example.com -tokenpath token.txt -enable_index_file -max_fetchers 3 > sitemap.xml
```

Download previously generated sitemap files:

```bash
mkdir -p sitemaps
go run . download https://www.example.com -tokenpath token.txt -out_dir ./sitemaps
```

Read generation stats:

```bash
go run . stats https://www.example.com -tokenpath token.txt
```

## Flags

Shared flag:

- `-tokenpath`: path to a token file. Empty is accepted, but larger sites and saved/downloadable generated files require a token.

`run` flags:

- `-max_fetchers`: maximum concurrent connections. Default: `3`.
- `-reference_count_threshold`: exclude images and videos embedded on more than this number of HTML pages. Default: `-1`.
- `-enable_index_file`: enable sitemap index generation. Default: `false`.
- `-max_request_retries`: retries for failed requests. Default: `5`.
- `-request_retry_timeout`: seconds to wait after a failed request. Default: `30`.
- `-sleep_time`: seconds between generation status polls. Default: `5`.

`download` flags:

- `-out_dir`: output directory for downloaded sitemap files. Create it before running `download`.

## Gotchas

- The user passes the raw URL. Do not pre-encode it; the CLI URL-safe base64 encodes it internally.
- `run` writes the final sitemap XML to stdout. Progress, stats, warnings, and errors are logged to stderr.
- During `run`, non-XML API responses are status/progress output. The command keeps polling until the API returns `application/xml`.
- Without a token, `run` can output the final sitemap directly, but the sitemap is not saved on the server for later `download`.
- `download` writes `sitemap.xml` first, then uses `stats` to discover indexed sitemap filenames such as `sitemap.000000.xml`.
- `download` creates files but not missing directories. Ensure `-out_dir` exists before running it.
- `download` does not validate HTTP status or content type before writing files. After download, verify the files are non-empty XML before treating them as valid.
- The CLI uses the external API endpoint in `config.go`; commands require network access and can take a while on large sites.

## Validation Loop

After running commands:

1. Check the process exit status.
2. For `run`, verify the redirected output file exists, is non-empty, and starts with XML.
3. For `download`, verify `sitemap.xml` exists in `-out_dir`; if stats show index files, verify each expected indexed sitemap file exists and is XML.
4. For `stats`, verify stdout is valid JSON before using the values.
5. If validation fails, inspect stderr first; it carries CLI logs and API error status messages.
