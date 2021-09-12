package models

// Basic repository struct
type Repository struct {
	FullName string `json:"full_name"`
	Lang     string `json:"language"`
}

// Repository struct with additional data
type RepositoryData struct {
	Repo  Repository `json:"full_name"`
	Lines map[string]int
}

// Struct to decode Github API Response
type APIResponse struct {
	Repos []Repository `json:"items"`
}
