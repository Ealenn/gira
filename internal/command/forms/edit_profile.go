package forms

import (
	"fmt"
	"regexp"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
	"github.com/charmbracelet/huh"
)

type EditProfile struct {
	logger *log.Logger
	ui     *huh.Form
}

func NewEditProfile(logger *log.Logger) *EditProfile {
	return &EditProfile{
		logger,
		nil,
	}
}

func (form EditProfile) Ask(profile *configuration.Profile) {
	/*
	* Profile
	 */
	form.ui = form.getProfileForm(profile)
	profileFormErr := form.ui.Run()

	if profileFormErr != nil {
		form.logger.Fatal("The operation was %s", "canceled")
	}

	form.ui.View()

	/*
	* Account
	 */
	form.ui = form.getAccountForm(profile)
	accountFormErr := form.ui.Run()

	if accountFormErr != nil {
		form.logger.Fatal("The operation was %s", "canceled")
	}

	form.ui.View()
}

func (form EditProfile) getProfileForm(profile *configuration.Profile) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[configuration.ProfileType]().
				Title("Type").
				Options(
					huh.Option[configuration.ProfileType]{Key: "Jira", Value: configuration.ProfileTypeJira},
					huh.Option[configuration.ProfileType]{Key: "GitHub", Value: configuration.ProfileTypeGithub},
				).
				Value(&profile.Type),
		)).WithTheme(huh.ThemeDracula())
}

func (form EditProfile) getAccountForm(profile *configuration.Profile) *huh.Form {
	var steps []*huh.Group

	switch profile.Type {
	case configuration.ProfileTypeJira:
		steps = append(steps, huh.NewGroup(
			huh.NewInput().
				Title("Host").
				Validate(func(url string) error {
					re := regexp.MustCompile(`^(https?://)?([\w-]+\.)+[\w-]{2,}(/.*)?$`)
					if !re.MatchString(url) {
						return fmt.Errorf("❌ %s (example: %s)", "The Jira host must be a valid URL", "https://your-domain.atlassian.net")
					}
					return nil
				}).
				Value(&profile.Jira.Host),
			huh.NewInput().
				Title("Token").
				EchoMode(huh.EchoModePassword).
				Value(&profile.Jira.Token),
		))
	case configuration.ProfileTypeGithub:
		steps = append(steps, huh.NewGroup(
			huh.NewInput().
				Title("User").
				Value(&profile.Github.User),
			huh.NewInput().
				Title("Token").
				EchoMode(huh.EchoModePassword).
				Value(&profile.Github.Token),
		))
	}

	return huh.NewForm(
		steps...,
	).WithTheme(huh.ThemeDracula())
}
