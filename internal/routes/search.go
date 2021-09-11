package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sihamais/searchingo/internal/models"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {

	searchLang := c.Query("q")

	//  http.GET timeout -> creer http.Client(Timeout: 10).Get()
	repos, err := GetRecentRepos()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusRequestTimeout, gin.H{"Error": "Fetch Repositories"})
		return
	}

	reposMatch := FilterReposByLang(repos, searchLang)

	results, err := GetReposData(reposMatch)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusRequestTimeout, gin.H{"Error": "Fetch Repositories Data"})
		return
	}

	SortReposByLines(results)
	
	linesCount := LinesCount(results, searchLang)

	c.HTML(http.StatusOK, "search.tmpl", gin.H{"repos": results, "search": searchLang, "lines": linesCount})
}

func GetRecentRepos() ([]models.Repository, error) {

	var repos models.APIResponse

	resp, err := http.Get("https://api.github.com/search/repositories?q=created:>2021-09-10+is:public&sort=author-date&per_page=100&page=1")
	if err != nil {
		return repos.Repos, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return repos.Repos, err
	}

	return repos.Repos, nil
}

func FilterReposByLang(repos []models.Repository, searchLang string) []models.Repository {

	var reposMatch []models.Repository

	for _, r := range repos {
		if strings.EqualFold(r.Lang, searchLang) {
			reposMatch = append(reposMatch, r)

		}
	}

	return reposMatch
}

func GetRepoLang(repo models.Repository) (models.RepositoryData, error) {

	var data models.RepositoryData

	resp, err := http.Get("https://api.github.com/repos/" + repo.FullName + "/languages")
	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data.Lines); err != nil {
		return data, err
	}

	data.Repo = repo

	return data, nil
}

func GetReposData(repos []models.Repository) ([]models.RepositoryData, error) {

	var results []models.RepositoryData

	for _, repo := range repos {
		data, err := GetRepoLang(repo)
		if err != nil {
			return results, err
		}
		results = append(results, data)
	}

	return results, nil
}

func SortReposByLines(repos []models.RepositoryData) {

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Lines[repos[i].Repo.Lang] > repos[j].Lines[repos[j].Repo.Lang]
	})

}

func LinesCount(repos []models.RepositoryData, searchLang string) int {

    linesCount := 0

    for _, r := range repos {
        linesCount += r.Lines[r.Repo.Lang]
    }
    return linesCount
}