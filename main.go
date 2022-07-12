package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	http.DefaultClient.Timeout = 30 * time.Second
}

// log.Println() writes to stderr by default
// fmt.Println() writes to stdout by default
// log.Fatal is equivalent to log.Print() followed by a call to os.Exit(1).
func main() {
	log.SetFlags(log.Lshortfile)

	// TODO can tokenPath be streamlined? define globally and read token just once?
	var tokenPath string

	runFlagSet := flag.NewFlagSet("run", flag.ExitOnError)
	runFlagSet.StringVar(&tokenPath, "tokenpath", "", "path to the token file")

	referenceCountThreshold := runFlagSet.Int64("reference_count_threshold", -1, "With the reference count threshold you can define that images and videos that are embedded on more than the selected number of HTML pages are excluded from the sitemap.")
	maxFetchers := runFlagSet.Int64("max_fetchers", 3, "Number of the maximal concurrent connections.")
	enableIndexFile := runFlagSet.Bool("enable_index_file", false, "Enable generation of a sitemap index file, recommended for large websites.")

	maxRequestRetries := runFlagSet.Int64("max_request_retries", 5, "Number of retries for each failed request")
	requestRetryTimeoutInSeconds := runFlagSet.Int64("request_retry_timeout", 30, "Timeout in seconds after a failed request")
	sleepTimeInSeconds := runFlagSet.Int64("sleep_time", 5, "Seconds between each update request")

	downloadFlagSet := flag.NewFlagSet("download", flag.ExitOnError)
	downloadFlagSet.StringVar(&tokenPath, "tokenpath", "", "path to the token file")
	outDir := downloadFlagSet.String("out_dir", "", "TODO")

	statsFlagSet := flag.NewFlagSet("stats", flag.ExitOnError)
	statsFlagSet.StringVar(&tokenPath, "tokenpath", "", "path to the token file")

	if len(os.Args) < 3 {
		log.Fatalf("usage: %s <command> <url> [flags]\n", os.Args[0])
	}

	command := os.Args[1]

	url := os.Args[2]
	urlBase64 := base64.URLEncoding.EncodeToString([]byte(url))

	switch command {
	case "stats":
		err := statsFlagSet.Parse(os.Args[3:])
		if err != nil {
			log.Fatalln(err)
		}

		token, ok := readToken(tokenPath)
		if !ok {
			log.Fatalln("could not read token from file")
		}

		doStats(urlBase64, token)

	case "download":
		err := downloadFlagSet.Parse(os.Args[3:])
		if err != nil {
			log.Fatalln(err)
		}

		token, ok := readToken(tokenPath)
		if !ok {
			log.Fatalln("could not read token from file")
		}

		doDownload(urlBase64, token, *outDir)

	case "run":
		err := runFlagSet.Parse(os.Args[3:])
		if err != nil {
			log.Fatalln(err)
		}

		token, ok := readToken(tokenPath)
		if !ok {
			log.Fatalln("could not read token from file")
		}

		doRun(urlBase64, token, *maxFetchers, *referenceCountThreshold, *enableIndexFile, *maxRequestRetries, *requestRetryTimeoutInSeconds, *sleepTimeInSeconds)

	default:
		log.Fatalln("command not supported, supported commands are run, stats and download")
	}
}
