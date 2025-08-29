package forms

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/log"
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
				Description("See https://support.atlassian.com/organization-administration/docs/understand-user-api-tokens/").
				EchoMode(huh.EchoModePassword).
				Value(&profile.Jira.Token),
		), huh.NewGroup(
			huh.NewInput().
				Title("Dashboard ID").
				Description("Optional: Used to fetch issues, keep empty if you don't use 'dash' command").
				Validate(func(s string) error {
					if _, err := strconv.Atoi(s); s != "" && err != nil {
						return fmt.Errorf("❌ %s (example: %s)", "Please enter a valid Dashboard ID", "1234")
					}
					return nil
				}).
				Value(&profile.Jira.Board),
			huh.NewInput().
				Title("JQL").
				Description("Optional: Used to filter issues on 'dash' command").
				Value(&profile.Jira.JQL),
		))
	case configuration.ProfileTypeGithub:
		steps = append(steps, huh.NewGroup(
			huh.NewInput().
				Title("User").
				Value(&profile.Github.User),
			huh.NewInput().
				Title("Token").
				Description("See https://github.com/settings/tokens").
				EchoMode(huh.EchoModePassword).
				Value(&profile.Github.Token),
		))
	}

	return huh.NewForm(
		steps...,
	).WithTheme(huh.ThemeDracula())
}
