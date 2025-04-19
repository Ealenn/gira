package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type Configuration struct {
	JSON    JSONConfiguration
	Path    string
	isDebug bool
}

type JSONConfiguration struct {
	JiraHost      string `json:"JiraHost"`
	JiraToken     string `json:"JiraToken"`
	JiraAccountID string `json:"JiraAccountID"`
}

func New() *Configuration {
	homeDirPath, homeDirPathError := os.UserHomeDir()
	if homeDirPathError != nil {
		log.Fatalf("unable to find home directory %v", homeDirPathError)
	}
	configurationFilePath := filepath.Join(homeDirPath, ".gira")

	if _, statError := os.Stat(configurationFilePath); statError != nil {
		configuration, createConfigurationError := createConfiguration(Configuration{
			JSON: JSONConfiguration{
				JiraHost:      "",
				JiraToken:     "",
				JiraAccountID: "",
			},
			Path:    configurationFilePath,
			isDebug: os.Getenv("DEBUG") == "TRUE",
		})

		if createConfigurationError != nil {
			log.Fatalf("unable to create configuration due to %v", createConfigurationError)
		}
		return configuration
	}

	fileContent, readConfigurationError := readConfiguration(configurationFilePath)
	if readConfigurationError != nil {
		log.Fatalf("unable to read configuration in %s due to %v", configurationFilePath, readConfigurationError)
	}

	return &Configuration{
		JSON: JSONConfiguration{
			JiraHost:      fileContent.JSON.JiraHost,
			JiraToken:     fileContent.JSON.JiraToken,
			JiraAccountID: fileContent.JSON.JiraAccountID,
		},
		Path:    configurationFilePath,
		isDebug: os.Getenv("DEBUG") == "TRUE",
	}
}

func (configuration *Configuration) IsValid() bool {
	parsedJiraHost, parsedJiraHostError := url.ParseRequestURI(configuration.JSON.JiraHost)
	if parsedJiraHostError != nil {
		return false
	}

	if parsedJiraHost.Scheme != "http" && parsedJiraHost.Scheme != "https" {
		return false
	}

	if len(configuration.JSON.JiraToken) < 5 {
		return false
	}

	if len(configuration.JSON.JiraAccountID) < 5 {
		return false
	}

	return true
}

func readConfiguration(configurationFilePath string) (*Configuration, error) {
	rawFileContent, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return nil, err
	}

	var fileContent JSONConfiguration
	if err := json.Unmarshal(rawFileContent, &fileContent); err != nil {
		return nil, err
	}

	return &Configuration{
		JSON: fileContent,
		Path: configurationFilePath,
	}, nil
}

func createConfiguration(configuration Configuration) (*Configuration, error) {
	jsonFileContent, err := json.Marshal(configuration.JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to create json configuration : %v", err)
	}

	file, err := os.Create(configuration.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := file.Write(jsonFileContent); err != nil {
		return nil, err
	}

	return &configuration, nil
}
