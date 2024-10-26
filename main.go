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
		log.Print("packages_https not set in config.json")
	} else {
		if config["repo_list_https"] != "" {
			merge.Merge(config["repo_list_https"], config["packages_https"])
		} else {
			log.Print("repo_list_https not set in config.json")
		}
	}
	if config["packages_http"] == "" {
		log.Print("packages_http not set in config.json")
	} else {
		if config["repo_list_http"] != "" {
			merge.Merge(config["repo_list_http"], config["packages_http"])
		} else {
			log.Print("repo_list_http not set in config.json")
		}
	}
}
