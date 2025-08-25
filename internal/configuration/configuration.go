package configuration

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/Ealenn/gira/internal/log"
)

type Configuration struct {
	logger *log.Logger
	JSON   JSONConfiguration
	Path   string
}

func New(logger *log.Logger) *Configuration {
	homeDirPath, homeDirPathError := os.UserHomeDir()
	if homeDirPathError != nil {
		logger.Fatal("unable to find home directory %v", homeDirPathError)
	}
	configurationFilePath := filepath.Join(homeDirPath, ".gira")

	var jsonConfiguration JSONConfiguration

	if _, statError := os.Stat(configurationFilePath); statError != nil {
		fileContent, createConfigurationError := updateConfiguration(configurationFilePath, JSONConfiguration{
			Profiles:         []Profile{},
			LastVersionCheck: 0,
		})

		if createConfigurationError != nil {
			logger.Fatal("unable to create configuration due to %v", createConfigurationError)
		}

		jsonConfiguration = *fileContent
	} else {
		fileContent, readConfigurationError := readConfiguration(configurationFilePath)
		if readConfigurationError != nil {
			logger.Fatal("unable to read configuration in %s due to %v", configurationFilePath, readConfigurationError)
		}

		jsonConfiguration = *fileContent
	}

	return &Configuration{
		logger: logger,
		JSON:   jsonConfiguration,
		Path:   configurationFilePath,
	}
}

func (configuration *Configuration) RemoveProfile(profile Profile) error {
	for index, jsonProfile := range configuration.JSON.Profiles {
		if jsonProfile.Name == profile.Name {
			configuration.JSON.Profiles = slices.Delete(configuration.JSON.Profiles, index, index+1)
			break
		}
	}

	_, err := updateConfiguration(configuration.Path, configuration.JSON)
	return err
}

func (configuration *Configuration) AddProfile(newProfile Profile) error {
	found := false
	for index, jsonProfile := range configuration.JSON.Profiles {
		if jsonProfile.Name == newProfile.Name {
			configuration.JSON.Profiles[index] = newProfile
			found = true
			break
		}
	}

	if !found {
		configuration.JSON.Profiles = append(configuration.JSON.Profiles, newProfile)
	}

	_, err := updateConfiguration(configuration.Path, configuration.JSON)
	return err
}

func (configuration *Configuration) GetProfile(name string) *Profile {
	for _, element := range configuration.JSON.Profiles {
		if element.Name == name {
			return &element
		}
	}

	return nil
}

func (configuration *Configuration) IsValid(profile *Profile) bool {
	if profile.Type == ProfileTypeJira {
		parsedJiraHost, parsedJiraHostError := url.ParseRequestURI(profile.Jira.Host)
		if parsedJiraHostError != nil {
			return false
		}
		if parsedJiraHost.Scheme != "http" && parsedJiraHost.Scheme != "https" {
			return false
		}
		if len(profile.Jira.Token) < 2 {
			return false
		}
	}

	return true
}

func (configuration *Configuration) VersionChecked() {
	configuration.JSON.LastVersionCheck = time.Now().Unix()
	updateConfiguration(configuration.Path, configuration.JSON)
}

func readConfiguration(path string) (*JSONConfiguration, error) {
	rawFileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fileContent JSONConfiguration
	if err := json.Unmarshal(rawFileContent, &fileContent); err != nil {
		return nil, err
	}

	return &fileContent, nil
}

func updateConfiguration(path string, jsonConfiguration JSONConfiguration) (*JSONConfiguration, error) {
	jsonFileContent, err := json.Marshal(jsonConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new json configuration : %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create json configuration file : %v", err)
	}
	defer file.Close()

	if _, err := file.Write(jsonFileContent); err != nil {
		return nil, err
	}

	return &jsonConfiguration, nil
}
