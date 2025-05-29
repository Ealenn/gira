package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Ealenn/gira/internal/log"
)

var client = &http.Client{Timeout: 10 * time.Second}

type GitHub struct {
	logger *log.Logger
}

func NewGitHub(logger *log.Logger) *GitHub {
	return &GitHub{
		logger,
	}
}

func (github *GitHub) GetLatestRelease() (*GithubLatestReleaseResponse, error) {
	githubReleaseResponse := &GithubLatestReleaseResponse{}
	err := getJSON("https://api.github.com/repos/Ealenn/gira/releases/latest", githubReleaseResponse)

	if err != nil {
		return nil, err
	}

	return githubReleaseResponse, nil
}

func getJSON(url string, target interface{}) error {
	response, err := client.Get(url)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

type GithubLatestReleaseResponse struct {
	URL         string    `json:"url"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
}
