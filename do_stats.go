package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func doStats(urlBase64, token string) {
	stats, err := getStats(urlBase64, token)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := json.MarshalIndent(stats, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(data))
}
