package command

import (
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"

	"github.com/manifoldco/promptui"
)

type Branch struct {
	logger  *log.Logger
	tracker issue.Tracker
	git     *git.Git
	branch  *branch.Manager
}

func NewBranch(logger *log.Logger, tracker issue.Tracker, git *git.Git, branch *branch.Manager) *Branch {
	return &Branch{
		logger,
		tracker,
		git,
		branch,
	}
}

func (command Branch) Run(issueID string, assign bool, force bool) {
	issue := command.tracker.GetIssue(issueID)
	command.RunWithIssue(issue, assign, force)
}

func (command Branch) RunWithIssue(issue *issue.Issue, assign bool, force bool) {
	branch := command.branch.Generate(issue)

	if command.git.IsBranchExist(branch.Raw) {
		command.logger.Warn("⚠️ Branch named %s already exists", branch.Raw)

		if !force {
			switchBranchPrompt := promptui.Prompt{
				Label:     "Would you like to switch to this branch?",
				IsConfirm: true,
				Default:   "y",
			}
			_, switchBranchPromptError := switchBranchPrompt.Run()

			if switchBranchPromptError != nil {
				command.logger.Fatal("The operation was %s", "canceled")
			}
		}

		// TODO: Handle error
		command.git.SwitchBranch(branch.Raw)
		command.logger.Info("✅ %s has just been checkout", branch.Raw)
	} else {
		if !force {
			command.logger.Info("🌳 Branch will be generated\nPress %s to continue, %s to cancel", "ENTER", log.ErrorStyle.Render("CTRL+C"))
			branchNamePrompt := promptui.Prompt{
				Label:     "Branch",
				Default:   branch.Raw,
				Pointer:   promptui.PipeCursor,
				AllowEdit: true,
			}
			newBranchName, err := branchNamePrompt.Run()
			if err != nil {
				command.logger.Fatal("The operation was %s", "canceled")
			}
			branch.Raw = newBranchName
		}

		command.git.CreateBranch(branch.Raw)
		command.logger.Info("✅ %s branch was created", branch.Raw)
	}

	if assign {
		if assignError := command.tracker.SelfAssignIssue(branch.IssueID); assignError != nil {
			command.logger.Debug("%v", assignError)
			command.logger.Info("❌ %s Unable to assign issue %s...", log.ErrorStyle.Render("Oups..."), issue.ID)
		} else {
			command.logger.Info("✅ Jira %s has been assigned to %s", issue.ID)
		}
	}
}
