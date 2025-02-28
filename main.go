package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Repositories struct {
	Total             *int   `json:"total_count,omitempty"`
	IncompleteResults *bool  `json:"incomplete_results,omitempty"`
	Items             []Repo `json:"items,omitempty"`
}

type Repo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	Issues      int       `json:"open_issues_count"`
	Created     time.Time `json:"created_at"`
	Updated     time.Time `json:"updated_at"`
	URL         string    `json:"html_url"`
}

var result Repositories
var languages = []string{"go", "java", "kotlin"}

func main() {
	now := time.Now()
	readme, err := os.OpenFile("README.md", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	readme.WriteString("# Top Projects by Stars\n")
	readme.WriteString("\n\n")

	readme.WriteString(fmt.Sprintf("Updated at: %v \n", now.Format(time.DateTime)))

	for _, lang := range languages {
		result := getGithubResult(lang)
		writeResultToReadme(strings.Title(lang), result.Items, readme)
	}
}

func getGithubResult(lang string) Repositories {
	apiURL := "https://api.github.com/search/repositories?q=language:" + lang + "&sort=stars&order=desc&per_page=100"
	resp, err := http.Get(apiURL)

	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
	}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&result); err != nil {
		log.Println(err)
	}
	return result
}

func writeResultToReadme(lang string, result []Repo, readme *os.File) {
	_, _ = readme.WriteString(fmt.Sprintf(`
## Top %s Projects

|    | Project Name | Stars | Forks | Open Issues | Description |
| -- | ------------ | ----- | ----- | ----------- | ----------- |
`, lang))

	for i, repo := range result {
		_, _ = readme.WriteString(fmt.Sprintf("| %d | [%s](%s) | %d | %d | %d | %s |\n", i+1, repo.Name, repo.URL, repo.Stars, repo.Forks, repo.Issues, repo.Description))
	}
}
