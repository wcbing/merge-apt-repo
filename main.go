package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/wcbing/merge-apt-repo/merge"
)

var config map[string]string

func readConfig() {
	if content, err := os.ReadFile("data/config.json"); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(content, &config); err != nil {
		log.Fatal(err)
	}
}

func main() {
	readConfig()
	if config["packages_https"] == "" {
		log.Print("packages_https and packages_http must be set in config.json")
	} else {
		if config["merge_all_https"] != "" {
			merge.MergeAll(config["merge_all_https"], config["packages_https"])
		}
		if config["merge_latest_https"] != "" {
			merge.MergeLatest(config["merge_latest_https"], config["packages_https"])
		}
	}
	if config["packages_http"] == "" {
		log.Fatal("packages_https and packages_http must be set in config.json")
	} else {
		if config["merge_all_http"] != "" {
			merge.MergeAll(config["merge_all_http"], config["packages_http"])
		}
	}
}
