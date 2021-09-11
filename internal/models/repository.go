package models

type Repository struct {
	FullName string `json:"full_name"`
	Lang     string `json:"language"`
}

type RepositoryData struct {
	Repo  Repository `json:"full_name"`
	Lines map[string]int
}

type APIResponse struct {
	Repos []Repository `json:"items"`
}
