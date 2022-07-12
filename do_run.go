package main

import (
	"log"
	"time"
)

func doRun(urlBase64, token string, maxFetchers, referenceCountThreshold int64, enableIndexFile bool, maxRequestRetries, requestRetryTimeoutInSeconds, sleepTimeInSeconds int64) {
	retriesCount := int64(0)
	for {
		if body, statusCode, contentType, stats, limitReached, ok := doRequest(urlBase64, token, maxFetchers, referenceCountThreshold, enableIndexFile); ok {
			retriesCount = 0 // always reset retries count on a successful request

			if contentType == "application/xml" {
				if stats != "" {
					log.Println(stats)
				}
				if limitReached {
					log.Println("WARNING: the URL limit was reached and the sitemap probably is not complete")
				}

				return
			} else {
				log.Println(body) // stats are just set in final request, before, stats are in body
			}
		} else if statusCode == 0 && retriesCount < maxRequestRetries {
			// do up to three retries if request fails
			// the easiest way to simulate retries is to add an invalid port the sitemap generator API URL (api.marcobeierer.com) below
			retriesCount++

			// sleep a little longer if there was an error, might be a refused connection due to too much requests in short time
			time.Sleep(time.Duration(requestRetryTimeoutInSeconds) * time.Second)

			// don't `continue` because we want to sleep anyway
		} else {
			if retriesCount > 0 {
				log.Fatalln("multiple request failed, abort sitemap generation")
			} else {
				log.Fatalln("request failed, abort sitemap generation")
			}
			return
		}
		time.Sleep(time.Duration(sleepTimeInSeconds) * time.Second)
	}
}
