package command

import (
	"bufio"
	"os"

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
	branch  *branch.BranchManager
}

func NewBranch(logger *log.Logger, tracker issue.Tracker, git *git.Git, branch *branch.BranchManager) *Branch {
	return &Branch{
		logger,
		tracker,
		git,
		branch,
	}
}

func (command Branch) Run(issueID string, assign bool, force bool) {
	issue := command.tracker.GetIssue(issueID)
	branch := command.branch.Generate(issue)

	if command.git.IsBranchExist(branch.Raw) {
		command.logger.Warn("⚠️ Branch named %s already exists", branch.Raw)

		if !force {
			command.logger.Info("Would you like to switch to this branch?\nPress %s to continue, %s to cancel", "ENTER", log.ErrorStyle.Render("CTRL+C"))
			bufio.NewReader(os.Stdin).ReadLine()
		}

		command.git.SwitchBranch(branch.Raw)
		command.logger.Info("✅ %s has just been checkout", branch.Raw)
	} else {
		if !force {
			command.logger.Info("🌳 Branch will be generated\nPress %s to continue, %s to cancel", "ENTER", log.ErrorStyle.Render("CTRL+C"))
			prompt := promptui.Prompt{
				Label:     "Branch",
				Default:   branch.Raw,
				Pointer:   promptui.PipeCursor,
				AllowEdit: true,
			}
			newBranchName, err := prompt.Run()
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
			command.logger.Info("❌ %s Unable to assign issue %s...", log.ErrorStyle.Render("Oups..."), issueID)
		} else {
			command.logger.Info("✅ Jira %s has been assigned to %s", issueID)
		}
	}
}
