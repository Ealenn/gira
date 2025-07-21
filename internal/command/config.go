package command

import (
	"regexp"
	"strings"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"

	"github.com/manifoldco/promptui"
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
		command.profile.Github = configuration.Github{}
	} else {
		command.createOrUpdateGithubProfile()
		command.profile.Jira = configuration.Jira{}
	}

	err := command.configuration.AddProfile(*command.profile)
	if err != nil {
		command.logger.Fatal("❌ Unable to save configuration")
	}
	command.logger.Info("✅ Done!")
}

func (command Config) createOrUpdateJiraProfile() {
	command.logger.Info("Enter the Jira API URL (Example %s):", "https://jira.mycompagny.com")
	command.question("Endpoint", &command.profile.Jira.Host, false)
	if !isValidURLRegex(command.profile.Jira.Host) {
		command.logger.Fatal("%s '%s' is not a valid URL. Please make sure it's a full URL including the scheme (e.g. https://example.com)", "ERROR", command.profile.Jira.Host)
	}

	command.logger.Info("Enter the Jira Token (See %s%s):", command.profile.Jira.Host, "/manage-profile/security/api-tokens")
	command.question("Token", &command.profile.Jira.Token, true)

	// TODO: Refactoring
	jiraService := issue.NewJira(command.logger, command.profile, git.NewGit(command.logger))
	jiraUser, jiraUserError := jiraService.GetMyself()

	if jiraUserError != nil {
		command.logger.Debug("Unable to fetch user accound due to %v", jiraUserError)
		command.logger.Fatal("❌ Unable to fetch Jira account in %s.", command.profile.Jira.Host)
	}

	command.profile.Jira.Email = jiraUser.EmailAddress
	command.profile.Jira.AccountID = jiraUser.AccountID
	command.profile.Jira.UserKey = jiraUser.Key
}

func (command Config) createOrUpdateGithubProfile() {
	command.logger.Info("Enter your Github %s (Example %s):", "Username", "ealenn")
	command.question("User", &command.profile.Github.User, false)

	command.logger.Info("Enter your Github %s (See %s):", "Token", "https://github.com/settings/tokens")
	command.logger.Info("Leave %s to use only %s repositories", "blank", "public")
	command.question("Token", &command.profile.Github.Token, true)
}

func (command Config) question(label string, value *string, password bool) {
	var mask rune
	if password {
		mask = '*'
	}

	prompt := promptui.Prompt{
		Label:     label,
		AllowEdit: true,
		Default:   *value,
		Mask:      mask,
		Pointer:   promptui.PipeCursor,
	}

	result, err := prompt.Run()

	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	*value = strings.TrimSpace(result)
}

func (command Config) selectProfileType() configuration.ProfileType {
	prompt := promptui.Select{
		Label: "Type",
		Items: []configuration.ProfileType{configuration.ProfileTypeJira, configuration.ProfileTypeGithub},
	}

	_, result, err := prompt.Run()

	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	return configuration.ProfileType(result)
}

func isValidURLRegex(url string) bool {
	re := regexp.MustCompile(`^(https?://)?([\w-]+\.)+[\w-]{2,}(/.*)?$`)
	return re.MatchString(url)
}
