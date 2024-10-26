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
//	    "only_latest": only has the latest version or not
//	    "repo": repo url, end with "/"
//	    "xxx_path": repo xxx Packages file path, start with no "/"
//	}

type repo struct {
	Name       string `json:"name"`
	OnlyLatest bool   `json:"only_latest"`
	Repo       string `json:"repo"`
	MixPath    string `json:"mix_path"` // 用于混合了多个平台包的扁平仓库
	Amd64Path  string `json:"amd64_path"`
	Arm64Path  string `json:"arm64_path"`
	AllPath    string `json:"all_path"`
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
func getRemotePackages(repoUrl, filePath string) (content []byte) {

	// get packages file content
	fileUrl := repoUrl + filePath
	if resp, err := http.Get(fileUrl); err != nil {
		log.Print(err)
	} else if resp.StatusCode != http.StatusOK {
		log.Printf("GetError: %s returned status %s\n", fileUrl, resp.Status)
	} else {
		content, _ = io.ReadAll(resp.Body)
	}

	// complete the two newlines if the ending is less than two newlines
	// 结尾不足两个换行符的话，补全两个换行符
	if !bytes.HasSuffix(content, []byte("\n\n")) {
		content = append(content, []byte("\n")...)
	}
	return bytes.Replace(content, []byte("Filename: "), []byte("Filename: "+repoUrl), -1)
}

// save packages info to file
func savePackages(packagesFile string, packages []byte) {
	// check dir exists
	packageDir := packagesFile[:strings.LastIndex(packagesFile, "/")]
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
	if f, err := os.OpenFile(packagesFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Print(err)
	} else if _, err := f.Write(packages); err != nil {
		f.Close()
		log.Print(err)
	} else if err := f.Close(); err != nil {
		log.Print(err)
	}
}
