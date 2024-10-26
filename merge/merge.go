package merge

import (
	"bytes"
	"regexp"
	"sync"

	"github.com/hashicorp/go-version"
)

func getLatest(debPackages []byte) (latestInfo []byte) {
	packageVersion := make(map[string]string)
	packageInfo := make(map[string][]byte)

	// 将每个包的信息分割开，存放到 infoList 中
	debPackages = bytes.Replace(debPackages, []byte("Package: "), []byte("{{start}}Package: "), -1)
	infoList := bytes.Split(debPackages, []byte("{{start}}"))[1:]

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

func Merge(repoListFile, packagesFile string) {
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
			debPackages := []byte{}
			// 获取 Repo 中 Amd64 包信息
			if r.Amd64Path != "" {
				debPackages = append(debPackages, getRemotePackages(r.Repo, r.Amd64Path)...)
			}
			// 获取扁平 Repo 中包信息
			if r.MixPath != "" {
				debPackages = append(debPackages, getRemotePackages(r.Repo, r.MixPath)...)
			}
			// // 获取 Repo 中 Arm64 包信息
			// if r.Arm64Path != "" {
			// 	debPackages = append(debPackages, getRemotePackages(r.Repo, r.Arm64Path)...)
			// }

			// 判断是否需要筛选最新版本
			if r.OnlyLatest {
				results <- debPackages
			} else {
				results <- getLatest(debPackages)
			}
		}(r)
	}

	wg.Wait()
	close(results)

	// 将所有结果合并保存
	var packages []byte
	for result := range results {
		packages = append(packages, result...)
	}
	savePackages(packagesFile, packages)
}
