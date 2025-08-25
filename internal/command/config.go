package command

import (
	"github.com/Ealenn/gira/internal/command/forms"
	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
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
	/*
	 * List
	 */
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

	/*
	 * Create/Update
	 */
	profileExist := true
	if cmd.profile == nil {
		profileExist = false
		cmd.profile = &configuration.Profile{
			Name: profileName,
		}
	}
	forms.NewEditProfile(cmd.logger).Ask(cmd.profile)

	if !profileExist {
		cmd.logger.Info("Create new profile : %s...", profileName)
	} else {
		cmd.logger.Info("Update profile : %s...", profileName)
	}

	err := cmd.configuration.AddProfile(*cmd.profile)
	if err != nil {
		cmd.logger.Fatal("❌ Unable to save configuration")
	}
	cmd.logger.Info("✅ Done!")
}
