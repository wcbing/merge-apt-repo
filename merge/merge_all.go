package merge

import (
	"sync"
)

func MergeAll(repoListFile, packagesFile string) {
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
			results <- getRemotePackages(r)
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
