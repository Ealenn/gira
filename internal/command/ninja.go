package command

import (
	"fmt"

	"github.com/Ealenn/gira/internal/ai"
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/command/forms"
	"github.com/Ealenn/gira/internal/configuration"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type Ninja struct {
	profile *configuration.Profile
	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager
}

func NewNinja(logger *log.Logger, profile *configuration.Profile, tracker issue.Tracker, git *git.Git, branch *branch.Manager) *Ninja {
	return &Ninja{
		profile,
		logger,
		tracker,
		git,
		branch,
	}
}

func (cmd Ninja) Run(enableAI bool, force bool) {
	agent := ai.NewOpenAI(cmd.logger)
	options := forms.NewCreateIssue(cmd.logger).Ask(cmd.profile.Type == configuration.ProfileTypeJira)

	if enableAI {
		titleSuggestion, titleSuggestionErr := agent.IssueRewrite("Issue creation, this is the Title of the new issue", options.Title)
		if titleSuggestionErr == nil && forms.NewConfirm(cmd.logger).Ask(
			"🤖 Title suggestion", titleSuggestion, forms.TypeApply,
		).Confirmed {
			options.Title = titleSuggestion
		}

		descriptionSuggestion, descriptionSuggestionErr := agent.IssueRewrite("Issue creation, this is the Title of the new issue", options.Title)
		if descriptionSuggestionErr == nil && forms.NewConfirm(cmd.logger).Ask(
			"🤖 Description suggestion", descriptionSuggestion, forms.TypeApply,
		).Confirmed {
			options.Title = descriptionSuggestion
		}
	}

	if !force {
		if !forms.NewConfirm(cmd.logger).Ask(
			"Would you like to create this issue?",
			fmt.Sprintf("Type:%s\nTitle: %s\nDescription:\n%s", options.Type, options.Title, options.Description),
			forms.TypeYesNo,
		).Confirmed {
			cmd.logger.Fatal("The operation was %s", "canceled")
		}
	}

	issue := cmd.tracker.CreateIssue(issue.CreateIssueOptions{
		Type:        options.Type,
		Project:     options.Project,
		Title:       options.Title,
		Description: options.Description,
	})
	cmd.logger.Info("Issue %s created, see %s", issue.ID, issue.URL)

	NewBranch(cmd.logger, cmd.tracker, cmd.git, cmd.branch).RunWithIssue(issue, true, enableAI, force)
}
