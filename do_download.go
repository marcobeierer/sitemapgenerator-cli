package main

import (
	"fmt"
	"log"
	"os"

	securejoin "github.com/cyphar/filepath-securejoin"
)

// TODO be careful to not expose sitemapgenerator files when vendored...
func doDownload(urlBase64, token, outDir string) {
	if outDir == "" {
		log.Fatalln("no out dir provided")
	}

	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	outPath, err := securejoin.SecureJoin(currentPath, outDir)
	if err != nil {
		log.Fatalln(err)
	}

	err = downloadFile(urlBase64, outPath+"/sitemap.xml", token)
	if err != nil {
		log.Fatalln(err)
	}

	stats, err := getStats(urlBase64, token)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < stats.SitemapIndexCount; i++ {
		format := "sitemap.%0" + fmt.Sprintf("%d", stats.SitemapIndexNumberOfDigits) + "d.xml"
		filename := fmt.Sprintf(format, i)

		err = downloadFile(urlBase64, outPath+"/"+filename, token)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
