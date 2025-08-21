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

func (cmd Ninja) Run(enableAI bool, force bool) {
	issueOptions := cmd.createIssueOptions(enableAI)

	cmd.logger.Log("\n------------ PREVIEW ------------")
	cmd.logger.Log("Type:%s\nTitle: %s\nDescription:%s", issueOptions.Type, issueOptions.Title, issueOptions.Description)
	cmd.logger.Log("---------------------------------\nWould you like to create this issue?")

	if askContinue := cmd.askSelect("", []string{"Yes", "No"}); askContinue == "No" {
		cmd.logger.Fatal("The operation was %s", "canceled")
	}

	issue := cmd.tracker.CreateIssue(issueOptions)
	cmd.logger.Info("Issue %s created, see %s", issue.ID, issue.URL)

	NewBranch(cmd.logger, cmd.tracker, cmd.git, cmd.branch).RunWithIssue(issue, true, enableAI, force)
}

func (cmd Ninja) createIssueOptions(enableAI bool) issue.CreateIssueOptions {
	issueTypeString := cmd.askSelect("Type", []issue.Type{issue.TypeFeature, issue.TypeBug})
	project := cmd.ask("Project")
	title := cmd.ask("Title")
	if enableAI {
		title = cmd.rewriteWithAI("Title", title)
	}
	description := cmd.ask("Description")
	if enableAI {
		description = cmd.rewriteWithAI("Description", description)
	}

	return issue.CreateIssueOptions{
		Title:       title,
		Description: description,
		Type:        issue.Type(issueTypeString),
		Project:     project,
	}
}

func (cmd Ninja) askSelect(label string, items interface{}) string {
	prompt := promptui.Select{Label: label, HideSelected: true, HideHelp: true, Items: items}
	_, str, err := prompt.Run()

	if err != nil {
		cmd.logger.Fatal("The operation was %s", "canceled")
	}

	return str
}

func (cmd Ninja) ask(label string) string {
	prompt := promptui.Prompt{Label: label, HideEntered: true}
	result, err := prompt.Run()

	if err != nil {
		cmd.logger.Fatal("The operation was %s", "canceled")
	}

	return result
}

func (cmd Ninja) rewriteWithAI(label string, text string) string {
	agent := ai.NewOpenAI(cmd.logger)
	context := fmt.Sprintf("Issue creation, this is the %s of the new issue", label)
	aiResponse, aiSuggestErr := agent.IssueRewrite(context, text)

	if aiSuggestErr == nil {
		cmd.logger.Log("AI suggestion for %s:\n%s", label, aiResponse)
		confirmPrompt := promptui.Prompt{Label: "Would you like to apply it? ", IsConfirm: true, HideEntered: true}
		confirmResult, _ := confirmPrompt.Run()

		if confirmResult == "y" {
			return aiResponse
		}
	}

	return text
}
