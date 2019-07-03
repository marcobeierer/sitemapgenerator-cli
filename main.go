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

	for {
		if body, contentType, stats, limitReached, ok := doRequest(url, token, *maxFetchers, *referenceCountThreshold); ok {
			if contentType == "application/xml" {
				if stats != "" {
					log.Println(stats)
				}
				if limitReached {
					log.Println("WARNING: the URL limit was reached and the sitemap probably is not complete")
				}
				fmt.Println(body)
				return
			}
		} else {
			log.Fatalln("request failed")
			return
		}
		time.Sleep(1 * time.Second)
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

// returns body, contentType, stats (as unparsed json) limitReached, and bool if successful
func doRequest(url, token string, maxFetchers, referenceCountThreshold int64) (string, string, string, bool, bool) {
	urlBase64 := base64.URLEncoding.EncodeToString([]byte(url))

	// TODO max_fetchers as param
	requestURL := fmt.Sprintf("https://api.marcobeierer.com/sitemap/v2/%s?pdfs=1&origin_system=cli&max_fetchers=%d&reference_count_threshold=%d", urlBase64, maxFetchers, referenceCountThreshold)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Println(err)
		return "", "", "", false, false
	}

	if token != "" {
		token = strings.TrimSuffix(token, "\n")
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return "", "", "", false, false
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("content-type")
	stats := resp.Header.Get("X-Stats")
	limitReached := resp.Header.Get("X-Limit-Reached") == "1"

	if resp.StatusCode != http.StatusOK {
		log.Printf("got status code %d, expected 200\n", resp.StatusCode)
		return "", contentType, stats, limitReached, false
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return "", contentType, stats, limitReached, false
	}

	return string(bytes), contentType, stats, limitReached, true
}
