package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
func doRequest(urlBase64, token string, maxFetchers, referenceCountThreshold int64, enableIndexFile bool) (string, int, string, string, bool, bool) {
	requestURL := fmt.Sprintf("https://api.marcobeierer.com/sitemap/v2/%s?pdfs=1&origin_system=cli&max_fetchers=%d&reference_count_threshold=%d&enable_index_file=%t", urlBase64, maxFetchers, referenceCountThreshold, enableIndexFile)
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

// filename is sitemap.xml or sitemap.000001.xml, etc.
func downloadFile(urlBase64, filepathx, token string) error {
	requestURL := fmt.Sprintf("https://api.marcobeierer.com/sitemap/v2/%s/%s", urlBase64, filepath.Base(filepathx))

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		// err could just be invalid method or URL parse error
		log.Println(err)
		return err
	}

	if token != "" {
		token = strings.TrimSuffix(token, "\n")
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepathx)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func getStats(urlBase64, token string) (*Stats, error) {
	requestURL := fmt.Sprintf("https://api.marcobeierer.com/sitemap/v2/%s/stats", urlBase64)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		// err could just be invalid method or URL parse error
		log.Println(err)
		return nil, err
	}

	if token != "" {
		token = strings.TrimSuffix(token, "\n")
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	stats := Stats{}

	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &stats, nil
}
