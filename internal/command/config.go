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

func (cmd Config) Run(profileName string, list bool) {
	if list {
		for _, profile := range cmd.configuration.JSON.Profiles {
			switch profile.Type {
			case configuration.ProfileTypeJira:
				cmd.logger.Info("- [%s] Type %s, user %s", profile.Name, profile.Type, profile.Jira.Email)
			case configuration.ProfileTypeGithub:
				cmd.logger.Info("- [%s] Type %s, user %s", profile.Name, profile.Type, profile.Github.User)
			}
		}
		return
	}

	if cmd.profile == nil {
		cmd.logger.Info("Create new profile : %s", profileName)

		profileType := cmd.selectProfileType()
		cmd.profile = &configuration.Profile{
			Name: profileName,
			Type: profileType,
			Jira: configuration.Jira{
				Host:  "",
				Token: "",
			},
		}
	} else {
		cmd.logger.Info("Update profile : %s", profileName)
	}

	if cmd.profile.Type == configuration.ProfileTypeJira {
		cmd.createOrUpdateJiraProfile()
		cmd.profile.Github = configuration.Github{}
	} else {
		cmd.createOrUpdateGithubProfile()
		cmd.profile.Jira = configuration.Jira{}
	}

	err := cmd.configuration.AddProfile(*cmd.profile)
	if err != nil {
		cmd.logger.Fatal("❌ Unable to save configuration")
	}
	cmd.logger.Info("✅ Done!")
}

func (cmd Config) createOrUpdateJiraProfile() {
	cmd.logger.Info("Enter the Jira API URL (Example %s):", "https://jira.mycompagny.com")
	cmd.question("Endpoint", &cmd.profile.Jira.Host, false)
	if !isValidURLRegex(cmd.profile.Jira.Host) {
		cmd.logger.Fatal("%s '%s' is not a valid URL. Please make sure it's a full URL including the scheme (e.g. https://example.com)", "ERROR", cmd.profile.Jira.Host)
	}

	cmd.logger.Info("Enter the Jira Token (See %s%s):", cmd.profile.Jira.Host, "/manage-profile/security/api-tokens")
	cmd.question("Token", &cmd.profile.Jira.Token, true)

	// TODO: Refactoring
	jiraService := issue.NewJira(cmd.logger, cmd.profile, git.NewGit(cmd.logger))
	jiraUser, jiraUserError := jiraService.GetMyself()

	if jiraUserError != nil {
		cmd.logger.Debug("Unable to fetch user accound due to %v", jiraUserError)
		cmd.logger.Fatal("❌ Unable to fetch Jira account in %s.", cmd.profile.Jira.Host)
	}

	cmd.profile.Jira.Email = jiraUser.EmailAddress
	cmd.profile.Jira.AccountID = jiraUser.AccountID
	cmd.profile.Jira.UserKey = jiraUser.Key
}

func (cmd Config) createOrUpdateGithubProfile() {
	cmd.logger.Info("Enter your Github %s (Example %s):", "Username", "ealenn")
	cmd.question("User", &cmd.profile.Github.User, false)

	cmd.logger.Info("Enter your Github %s (See %s):", "Token", "https://github.com/settings/tokens")
	cmd.logger.Info("Leave %s to use only %s repositories", "blank", "public")
	cmd.question("Token", &cmd.profile.Github.Token, true)
}

func (cmd Config) question(label string, value *string, password bool) {
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
		cmd.logger.Fatal("The operation was %s", "canceled")
	}

	*value = strings.TrimSpace(result)
}

func (cmd Config) selectProfileType() configuration.ProfileType {
	prompt := promptui.Select{
		Label: "Type",
		Items: []configuration.ProfileType{configuration.ProfileTypeJira, configuration.ProfileTypeGithub},
	}

	_, result, err := prompt.Run()

	if err != nil {
		cmd.logger.Fatal("The operation was %s", "canceled")
	}

	return configuration.ProfileType(result)
}

func isValidURLRegex(url string) bool {
	re := regexp.MustCompile(`^(https?://)?([\w-]+\.)+[\w-]{2,}(/.*)?$`)
	return re.MatchString(url)
}
