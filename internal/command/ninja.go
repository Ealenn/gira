package command

import (
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
	"github.com/manifoldco/promptui"
)

type Ninja struct {
	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager
}

func NewNinja(logger *log.Logger, tracker issue.Tracker, git *git.Git, branch *branch.Manager) *Ninja {
	return &Ninja{
		logger,
		tracker,
		git,
		branch,
	}
}

func (command Ninja) Run(force bool) {
	issueOptions := command.createIssueOptions()

	issue := command.tracker.CreateIssue(issueOptions)
	command.logger.Info("Issue %s created, see %s", issue.ID, issue.URL)

	NewBranch(command.logger, command.tracker, command.git, command.branch).RunWithIssue(issue, true, force)
}

func (command Ninja) createIssueOptions() issue.CreateIssueOptions {
	typeSelect := promptui.Select{Label: "Type", Items: []issue.Type{issue.TypeFeature, issue.TypeBug}}
	_, issueTypeString, err := typeSelect.Run()
	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	projectPrompt := promptui.Prompt{Label: "Project", HideEntered: true}
	project, err := projectPrompt.Run()
	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	titlePrompt := promptui.Prompt{Label: "Title", HideEntered: true}
	title, err := titlePrompt.Run()
	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	descriptionPrompt := promptui.Prompt{Label: "Description", HideEntered: true}
	description, err := descriptionPrompt.Run()
	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	return issue.CreateIssueOptions{
		Title:       title,
		Description: description,
		Type:        issue.Type(issueTypeString),
		Project:     project,
	}
}
