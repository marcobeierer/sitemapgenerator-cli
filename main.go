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
func main() {
	log.SetFlags(0)

	tokenPath := flag.String("tokenpath", "", "path to the token file")
	flag.Parse()

	token, ok := readToken(*tokenPath)
	if !ok {
		log.Println("could not read token from file")
		return
	}

	url := flag.Arg(0)
	if url == "" {
		log.Printf("usage: %s [flags] url\n", os.Args[0])
		return
	}

	for {
		if body, contentType, stats, limitReached, ok := doRequest(url, token); ok {
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
			log.Println("request failed")
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
func doRequest(url, token string) (string, string, string, bool, bool) {
	urlBase64 := base64.URLEncoding.EncodeToString([]byte(url))

	// TODO make max_etchers and reference count as param
	req, err := http.NewRequest("GET", "https://api.marcobeierer.com/sitemap/v2/"+urlBase64+"?pdfs=1&origin_system=cli&max_fetchers=3&reference_count_threshold=5", nil)
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
