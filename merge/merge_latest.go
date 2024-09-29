package merge

import (
	"bytes"
	"regexp"
	"sync"

	"github.com/hashicorp/go-version"
)

func getLatest(x64DebPackages []byte) (latestInfo []byte) {
	packageVersion := make(map[string]string)
	packageInfo := make(map[string][]byte)

	// 分组
	x64DebPackages = bytes.Replace(x64DebPackages, []byte("Package: "), []byte("{{start}}Package: "), -1)
	infoList := bytes.Split(x64DebPackages, []byte("{{start}}"))[1:]

	findName := regexp.MustCompile("Package: (.+)")
	findVersion := regexp.MustCompile("Version: (.+)")

	for _, v := range infoList {
		name := string(findName.FindSubmatch(v)[1])
		nversion := string(findVersion.FindSubmatch(v)[1])
		n1, _ := version.NewVersion(nversion)
		n2, _ := version.NewVersion(packageVersion[name])
		if packageVersion[name] == "" || n1.GreaterThan(n2) {
			packageVersion[name] = nversion
			packageInfo[name] = v
		}
		// fmt.Println(string(name), string(version))
	}
	for _, v := range packageInfo {
		latestInfo = append(latestInfo, v...)
	}
	return
}

func MergeLatest(repoListFile, packagesFile string) {
	repoList := readRepoList(repoListFile)
	if repoList == nil {
		return
	}

	var wg sync.WaitGroup
	results := make(chan []byte, len(repoList))

	for _, r := range repoList {
		wg.Add(1)
		go func(r repo) {
			defer wg.Done()
			x64DebPackages := getRemotePackages(r)
			results <- getLatest(x64DebPackages)
		}(r)
	}
	wg.Wait()
	close(results)
	var packages []byte
	for result := range results {
		packages = append(packages, result...)
	}

	savePackages(packagesFile, packages)
}
