
# Searchingo

Searchingo is an app to browse the latest publics repos on Github. This app uses the Github API.

## Table of Content

- [Run Locally](https://github.com/sihamais/searchingo#run-locally)
- [Features](https://github.com/sihamais/searchingo#features)
- [API Reference](https://github.com/sihamais/searchingo#api-reference)
- [Lessons Learned](https://github.com/sihamais/searchingo#lessons-learned)
- [Optimizations](https://github.com/sihamais/searchingo#optimizations)
- [Tech Stack](https://github.com/sihamais/searchingo#tech-stack)
- [Used Resources](https://github.com/sihamais/searchingo#used-resources)
## Run Locally

Clone the project

```bash
  git clone https://github.com/sihamais/searchingo.git
```

Go to the project directory

```bash
  cd Searchingo
```

Install dependencies

```bash
  cd cmd
  go get
```

Start the server

```bash
  cd cmd
  go run main.go
```

Open ```localhost:8080``` on your browser and then you're good to *go*.

  
## Features

- Browse the latest public repos on github
- Results sorted by number of bytes (desc) written in the chosen language
- Statistics displayed in charts
- Parallel processing of the repos

  
## API Reference

#### Get the landing page

```HTTP
  GET /
```

#### Get search results

```HTTP
  GET /search
```

| Query | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `q`      | `string` | **Required**. Name of a language |


  
## Lessons Learned

The real difficulty in this project is to find a way to process the data without overflowing the rate limit enforced by Github. The obvious approach for this project is to get the list of public repos: 
```HTTP
GET /repositories
```
However, the rate limit being fixed to 60 requests per hour, it is impossible to get every repository's informations since there are 100.

Two solutions were possible :

#### 1. Establish a github authentication service to increase the rate limit to 5000 requests per hour

This solution requires authentication from the user, which goes against the subject of the project which is supposed to allow browsing of public repositories without needing to authenticate.  

Also, this solution is a double-edged sword since the limit is not canceled and the problem persists. 

#### 2. Use the Search API offered bu Github

This solutions uses another request which does not necessarily respect the instructions of the project but displays the same result.  

This solution is not ideal either since the rate limit is inevitable. However, it allows the user to be able to browse the repositories without needing to authenticate.  

The user will probably not be able to search more than 3 times in an hour. If the page takes too much time to load, it is probably because the rate limit have been reached. My trick in this case is to switch location in the VPN if the user has one.

## Optimizations

#### 1. Number of lines in the repository
Since Github API does not provide a method to get the number of lines, I've chose to display the number of bytes of the chosen language.   
This does still demonstrate the scale of the repository. Thus, the list is sorted according to the number of bytes.

#### 2. Proportion of each repository using the pie chart
Since our metric is the number of bytes, the chart will relay on this statistics to display the proportion each repository represents.

#### 3. Parallelisation
I've used a goroutine to manage the threads and channels to send and receive repositories from the workers. The number of the workers is defined manually.





  
## Tech Stack

**Client:** HTML, CSS, VanillaJS, Charts.js

**Server:** Go, Gin

  
## Used Resources
#### Github API:
 - [List Public Repositories](https://docs.github.com/en/rest/reference/repos#list-public-repositories)
 - [Rate Limiting](https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting)
 - [Search Repositories](https://docs.github.com/en/rest/reference/search#search-repositories)
 - [List Repositories Languages](https://docs.github.com/en/rest/reference/repos#list-repository-languages)

#### Go:

 - [Gin Web Framework](https://github.com/gin-gonic/gin#gin-web-framework)
 - [A Tour of Go](https://tour.golang.org/welcome/1)

#### Charts.js:
 - [Bar Chart](https://www.chartjs.org/docs/latest/charts/bar.html)
 - [Pie Chart](https://www.chartjs.org/docs/latest/charts/doughnut.html)
