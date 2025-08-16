package command

import (
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/Ealenn/gira/internal/ai"
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
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

func (command Ninja) Run(enableAI bool, force bool) {
	issueOptions := command.createIssueOptions(enableAI)

	command.logger.Log("Type:%s\nTitle: %s\nDescription:%s", issueOptions.Type, issueOptions.Title, issueOptions.Description)
	if askContinue := command.askSelect("Continue", []string{"Yes", "No"}); askContinue == "No" {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	issue := command.tracker.CreateIssue(issueOptions)
	command.logger.Info("Issue %s created, see %s", issue.ID, issue.URL)

	NewBranch(command.logger, command.tracker, command.git, command.branch).RunWithIssue(issue, true, enableAI, force)
}

func (command Ninja) createIssueOptions(enableAI bool) issue.CreateIssueOptions {
	issueTypeString := command.askSelect("Type", []issue.Type{issue.TypeFeature, issue.TypeBug})
	project := command.ask("Project")
	title := command.ask("Title")
	if enableAI {
		title = command.rewriteWithAI("Title", title)
	}
	description := command.ask("Description")
	if enableAI {
		description = command.rewriteWithAI("Description", description)
	}

	return issue.CreateIssueOptions{
		Title:       title,
		Description: description,
		Type:        issue.Type(issueTypeString),
		Project:     project,
	}
}

func (command Ninja) askSelect(label string, items interface{}) string {
	prompt := promptui.Select{Label: label, HideSelected: true, Items: items}
	_, str, err := prompt.Run()

	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	return str
}

func (command Ninja) ask(label string) string {
	prompt := promptui.Prompt{Label: label, HideEntered: true}
	result, err := prompt.Run()

	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	return result
}

func (command Ninja) rewriteWithAI(label string, text string) string {
	agent := ai.NewOpenAI(command.logger)
	context := fmt.Sprintf("Issue creation, this is the %s of the new issue", label)
	aiResponse, aiSuggestErr := agent.Rewrite(context, text)

	if aiSuggestErr == nil {
		command.logger.Log("AI suggestion for %s:\n%s", label, aiResponse)
		confirmPrompt := promptui.Prompt{Label: "Would you like to apply it? ", IsConfirm: true, HideEntered: true}
		confirmResult, _ := confirmPrompt.Run()

		if confirmResult == "y" {
			return aiResponse
		}
	}

	return text
}
