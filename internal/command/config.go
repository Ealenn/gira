package command

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
	"github.com/Ealenn/gira/internal/service"
	"github.com/manifoldco/promptui"

	"github.com/charmbracelet/x/term"
)

type Config struct {
	logger        *log.Logger
	configuration *configuration.Configuration
	profile       *configuration.Profile
}

func NewConfig(logger *log.Logger, configuration *configuration.Configuration, profile *configuration.Profile) *Config {
	return &Config{
		logger:        logger,
		configuration: configuration,
		profile:       profile,
	}
}

func (command Config) Run(profileName string, list bool) {
	if list {
		for _, profile := range command.configuration.JSON.Profiles {
			switch profile.Type {
			case configuration.ProfileTypeJira:
				command.logger.Info("- [%s] Type %s, user %s", profile.Name, profile.Type, profile.Jira.Email)
			case configuration.ProfileTypeGithub:
				command.logger.Info("- [%s] Type %s, user %s", profile.Name, profile.Type, profile.Github.User)
			}
		}
		return
	}

	if command.profile == nil {
		command.logger.Info("Create new profile : %s", profileName)

		profileType := command.selectProfileType()
		command.profile = &configuration.Profile{
			Name: profileName,
			Type: profileType,
			Jira: configuration.Jira{
				Host:  "",
				Token: "",
			},
		}
	} else {
		command.logger.Info("Update profile : %s", profileName)
	}

	if command.profile.Type == configuration.ProfileTypeJira {
		command.createOrUpdateJiraProfile()
	} else {
		command.createOrUpdateGithubProfile()
	}

	err := command.configuration.AddProfile(*command.profile)
	if err != nil {
		command.logger.Fatal("❌ Unable to save configuration")
	}
	command.logger.Info("✅ Done!")
}

func (command Config) createOrUpdateJiraProfile() {
	reader := bufio.NewReader(os.Stdin)

	// Jira Endpoint
	command.logger.Info("Enter the Jira API URL (Example %s):", "https://jira.mycompagny.com")
	if command.profile.Jira.Host != "" {
		command.logger.Info("[%s]", command.profile.Jira.Host)
	}
	inputJiraHost, _ := reader.ReadString('\n')
	inputJiraHost = strings.TrimSpace(inputJiraHost)
	if inputJiraHost == "" {
		inputJiraHost = command.profile.Jira.Host
	}
	if !isValidURLRegex(inputJiraHost) {
		command.logger.Fatal("%s '%s' is not a valid URL. Please make sure it's a full URL including the scheme (e.g. https://example.com)", "ERROR", inputJiraHost)
	}
	command.profile.Jira.Host = inputJiraHost

	// Jira Token
	command.logger.Info("Enter the Jira Token (See %s%s):", inputJiraHost, "/manage-profile/security/api-tokens")
	if command.profile.Jira.Token != "" {
		command.logger.Info("[Token already defined. Press %s to continue without making any changes.]", "ENTER")
	}
	inputJiraTokenBytes, _ := term.ReadPassword(os.Stdin.Fd())
	inputJiraToken := strings.TrimSpace(string(inputJiraTokenBytes))
	if inputJiraToken == "" {
		inputJiraToken = command.profile.Jira.Token
	}
	command.profile.Jira.Token = inputJiraToken

	// Jira Profile
	jiraService := service.NewJira(command.logger, command.profile)
	jiraUser, jiraUserError := jiraService.GetMyself()

	if jiraUserError != nil {
		command.logger.Debug("Unable to fetch user accound due to %v", jiraUserError)
		command.logger.Fatal("❌ Unable to fetch Jira account in %s.", inputJiraHost)
	}

	command.profile.Jira.Email = jiraUser.EmailAddress
	command.profile.Jira.AccountID = jiraUser.AccountID
	command.profile.Jira.UserKey = jiraUser.Key
}

func (command Config) createOrUpdateGithubProfile() {
	reader := bufio.NewReader(os.Stdin)

	// Github Username
	command.logger.Info("Enter your Github %s (Example %s):", "Username", "ealenn")
	if command.profile.Github.User != "" {
		command.logger.Info("[%s]", command.profile.Github.User)
	}
	inputGithubUser, _ := reader.ReadString('\n')
	inputGithubUser = strings.TrimSpace(inputGithubUser)
	if inputGithubUser == "" {
		inputGithubUser = command.profile.Github.User
	}

	command.profile.Github.User = inputGithubUser
}

func (command Config) selectProfileType() string {
	prompt := promptui.Select{
		Label: "Type",
		Items: []string{configuration.ProfileTypeJira, configuration.ProfileTypeGithub},
	}

	_, result, err := prompt.Run()

	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	return result
}

func isValidURLRegex(url string) bool {
	re := regexp.MustCompile(`^(https?://)?([\w-]+\.)+[\w-]{2,}(/.*)?$`)
	return re.MatchString(url)
}
