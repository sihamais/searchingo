package repos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sihamais/searchingo/internal/models"
	"sort"
	"strings"
	"time"
)

// Get the 100 last created repositories name and language from Github API
func GetRecentRepos() ([]models.Repository, error) {

	var repos models.APIResponse

	t := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	githubUrl := fmt.Sprintf("https://api.github.com/search/repositories?q=created:>=%s+is:public&sort=author-date&per_page=100&page=1", t)

	resp, err := http.Get(githubUrl)
	if err != nil {
		return repos.Repos, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return repos.Repos, err
	}

	return repos.Repos, nil
}

// Filter a list of repositories keeping only the ones matching a lang and also return statistics about the lang
func FilterReposByLang(repos []models.Repository, searchLang string) ([]models.Repository, map[string]int) {

	var reposMatch []models.Repository
	langStats := make(map[string]int)

	for _, r := range repos {

		// Map counting the lang we found
		langStats[r.Lang]++

		if strings.EqualFold(r.Lang, searchLang) {
			reposMatch = append(reposMatch, r)
		}
	}

	return reposMatch, langStats
}

// Get the repositories from a channel, get additional data, and send the updated structure on another channel
func getRepoDataWorker(jobs <-chan models.Repository, out chan<- models.RepositoryData) {

	for repo := range jobs {

		var data models.RepositoryData

		resp, err := http.Get("https://api.github.com/repos/" + repo.FullName + "/languages")
		if err != nil {
			return
		}

		if err = json.NewDecoder(resp.Body).Decode(&data.Lines); err != nil {
			return
		}

		data.Repo = repo
		out <- data
	}
}

// Concurrently search for the additional statistics of a list of repositories and return the updated list
func GetReposData(repos []models.Repository) ([]models.RepositoryData, error) {

	var results []models.RepositoryData

	// Create two channel for communications between goroutines
	jobsLen := len(repos)
	jobs := make(chan models.Repository, jobsLen)
	out := make(chan models.RepositoryData, jobsLen)

	// You can change this
	nbCores := 4
	// Limit the number of workers to the nb of jobs OR cores on the machine
	nbWorkers := min(jobsLen, nbCores)

	// Start workers and keep count of them
	for w := 1; w <= nbWorkers; w++ {
		go getRepoDataWorker(jobs, out)
	}

	// Share the repositories with the workers and close their input channel
	for _, repo := range repos {
		jobs <- repo
	}
	close(jobs)

	// Combine the data obtained by the workers in a list
	for i := 0; i < jobsLen; i++ {
		results = append(results, <-out)
	}

	return results, nil
}

// Sort a list of repositories by number of lines
func SortReposByLines(repos []models.RepositoryData) {

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Lines[repos[i].Repo.Lang] > repos[j].Lines[repos[j].Repo.Lang]
	})

}

// Compute the total number of lines of a list of repositories
func LinesCount(repos []models.RepositoryData) int {

	linesCount := 0

	for _, r := range repos {
		linesCount += r.Lines[r.Repo.Lang]
	}
	return linesCount
}

// Compute the min of two ints
func min(i int, j int) int {
	if i < j {
		return i
	} else {
		return j
	}
}
