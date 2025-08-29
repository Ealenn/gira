package command

import (
	"fmt"

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

func (cmd Config) Run(profileName string, list bool, remove bool) {
	profileExist := true
	if cmd.profile == nil {
		profileExist = false
		cmd.profile = &configuration.Profile{
			Name: profileName,
		}
	}

	/*
	 * List
	 */
	if list {
		for _, profile := range cmd.configuration.JSON.Profiles {
			switch profile.Type {
			case configuration.ProfileTypeJira:
				cmd.logger.Info("- [%s] Type Jira on %s", profile.Name, profile.Jira.Host)
			case configuration.ProfileTypeGithub:
				cmd.logger.Info("- [%s] Type GitHub with user %s", profile.Name, profile.Github.User)
			}
		}
		return
	}

	/*
	 * Remove
	 */
	if remove {
		if profileExist {
			if !forms.NewConfirm(cmd.logger).Ask(fmt.Sprintf("Confirm deletion of profile : %s", profileName), fmt.Sprintf("The profile %s will be deleted", profileName), forms.TypeConfirm).Confirmed {
				cmd.logger.Fatal("The operation was %s", "canceled")
			}
			if err := cmd.configuration.RemoveProfile(*cmd.profile); err != nil {
				cmd.logger.Fatal("❌ Unable to save configuration")
			}
		}

		cmd.logger.Info("✅ Done!")
		return
	}

	/*
	 * Create/Update
	 */
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
