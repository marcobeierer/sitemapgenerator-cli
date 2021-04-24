package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// log.Println() writes to stderr by default
// fmt.Println() writes to stdout by default
// log.Fatal is equivalent to log.Print() followed by a call to os.Exit(1).
func main() {
	log.SetFlags(0)

	tokenPath := flag.String("tokenpath", "", "path to the token file")
	referenceCountThreshold := flag.Int64("reference_count_threshold", -1, "With the reference count threshold you can define that images and videos that are embedded on more than the selected number of HTML pages are excluded from the sitemap.")
	maxFetchers := flag.Int64("max_fetchers", 3, "Number of the maximal concurrent connections.")

	maxRequestRetries := flag.Int64("max_request_retries", 5, "Number of retries for each failed request")
	requestRetryTimeoutInSeconds := flag.Int64("request_retry_timeout", 30, "Timeout in seconds after a failed request")
	sleepTimeInSeconds := flag.Int64("sleep_time", 5, "Seconds between each update request")

	flag.Parse()

	token, ok := readToken(*tokenPath)
	if !ok {
		log.Fatalln("could not read token from file")
		return
	}

	url := flag.Arg(0)
	if url == "" {
		log.Fatalf("usage: %s [flags] url\n", os.Args[0])
		return
	}

	retriesCount := int64(0)
	for {
		if body, statusCode, contentType, stats, limitReached, ok := doRequest(url, token, *maxFetchers, *referenceCountThreshold); ok {
			retriesCount = 0 // always reset retries count on a successfull request

			if contentType == "application/xml" {
				if stats != "" {
					log.Println(stats)
				}
				if limitReached {
					log.Println("WARNING: the URL limit was reached and the sitemap probably is not complete")
				}
				fmt.Println(body)
				return
			} else {
				log.Println(body) // stats are just set in final request, before, stats are in body
			}
		} else if statusCode == 0 && retriesCount < *maxRequestRetries {
			// do up to three retries if request fails
			// the easiest way to simulate retries is to add an invalid port the sitemap generator API URL (api.marcobeierer.com) below
			retriesCount++

			// sleep a little longer if there was an error, might be a refused connection due to too much requests in short time
			time.Sleep(time.Duration(*requestRetryTimeoutInSeconds) * time.Second)

			// don't `continue` because we want to sleep anyway
		} else {
			if retriesCount > 0 {
				log.Fatalln("multiple request failed, abort sitemap generation")
			} else {
				log.Fatalln("request failed, abort sitemap generation")
			}
			return
		}
		time.Sleep(time.Duration(*sleepTimeInSeconds) * time.Second)
	}
}

func readToken(tokenPath string) (string, bool) {
	if tokenPath == "" {
		return "", true
	}

	bytes, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Println(err)
		return "", false
	}

	return fmt.Sprintf("%s", bytes), true
}

// returns body, statusCode, contentType, stats (as unparsed json) limitReached, and bool if successful
func doRequest(url, token string, maxFetchers, referenceCountThreshold int64) (string, int, string, string, bool, bool) {
	urlBase64 := base64.URLEncoding.EncodeToString([]byte(url))

	// TODO max_fetchers as param
	requestURL := fmt.Sprintf("https://api.marcobeierer.com/sitemap/v2/%s?pdfs=1&origin_system=cli&max_fetchers=%d&reference_count_threshold=%d", urlBase64, maxFetchers, referenceCountThreshold)
	//requestURL := fmt.Sprintf("http://marco-desktop:9999/sitemap/v2/%s?pdfs=1&origin_system=cli&max_fetchers=%d&reference_count_threshold=%d&enable_index_file=1", urlBase64, maxFetchers, referenceCountThreshold)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		// err could just be invalid method or URL parse error
		log.Println(err)
		return "", -1, "", "", false, false // -1 because it doesn't make sense to retry in these cases
	}

	if token != "" {
		token = strings.TrimSuffix(token, "\n")
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", 0, "", "", false, false // 0 because we may retry to connect, err could for example be `connection refused`
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("content-type")
	stats := resp.Header.Get("X-Stats")
	limitReached := resp.Header.Get("X-Limit-Reached") == "1"

	if resp.StatusCode != http.StatusOK {
		log.Printf("got status code %d, expected 200\n", resp.StatusCode)
		return "", resp.StatusCode, contentType, stats, limitReached, false
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", resp.StatusCode, contentType, stats, limitReached, false
	}

	return string(bytes), resp.StatusCode, contentType, stats, limitReached, true
}
