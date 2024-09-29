package merge

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// repo info json format:
//	{
//	    "name": repo name
//	    "repo": repo url, end with "/"
//	    "amd64_path": repo amd64 Packages file path, start with no "/"
//	}

type repo struct {
	Name      string `json:"name"`
	Repo      string `json:"repo"`
	Amd64Path string `json:"amd64_path"`
}

// read repo info
func readRepoList(repoListFile string) (repoList []repo) {
	if content, err := os.ReadFile(repoListFile); err != nil {
		log.Print(err)
	} else if err := json.Unmarshal(content, &repoList); err != nil {
		log.Print(err)
	}
	return
}

// get the packages info from remote repo
func getRemotePackages(r repo) (content []byte) {

	// get amd64 packages info
	amd64Url := r.Repo + r.Amd64Path
	if resp, err := http.Get(amd64Url); err != nil {
		log.Print(err)
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("GetError: %s returned status %s\n", amd64Url, resp.Status)
	} else {
		content, _ = io.ReadAll(resp.Body)
	}

	// 结尾不足两个换行符的话，补全两个换行符
	if !bytes.HasSuffix(content, []byte("\n\n")) {
		content = append(content, []byte("\n")...)
	}
	return bytes.Replace(content, []byte("Filename: "), []byte("Filename: "+r.Repo), -1)
}

// save packages info to file
func savePackages(packagesFile string, packages []byte) {
	// check dir exists
	packageDir := packagesFile[:len(packagesFile)-len(packagesFile[strings.LastIndex(packagesFile, "/"):])]
	if _, err := os.Stat(packageDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(packageDir, os.ModePerm); err != nil {
				log.Print(err)
				return
			}
		} else {
			log.Print(err)
			return
		}
	}
	// append to the file
	f, err := os.OpenFile(packagesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
	}
	if _, err := f.Write(packages); err != nil {
		f.Close()
		log.Print(err)
	}
	if err := f.Close(); err != nil {
		log.Print(err)
	}
}
