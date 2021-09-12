package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sihamais/searchingo/internal/repos"
)

// Return the search page with the search results
func Search(c *gin.Context) {

	searchLang := c.Query("q")

	// Get the 100 last public repositories from Github
	recentRepos, err := repos.GetRecentRepos()
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{"Error": "Fetch Repositories"})
		return
	}

	// Only do the next search on repos matching the lang searched
	reposMatch, langStats := repos.FilterReposByLang(recentRepos, searchLang)

	// Concurent search of the repos lang data
	results, err := repos.GetReposData(reposMatch)
	if err != nil {
		c.JSON(http.StatusRequestTimeout, gin.H{"Error": "Fetch Repositories Data"})
		return
	}

	// Sort the list of repos by their number of lines just obtained
	repos.SortReposByLines(results)

	// Compute the total number of lines of all the repos
	linesCount := repos.LinesCount(results)

	// Return the searchpage to the client with the data
	c.HTML(http.StatusOK, "search.tmpl", gin.H{"repos": results, "search": searchLang, "lines": linesCount, "langStats": langStats})
}