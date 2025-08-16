package command

import (
	"github.com/Ealenn/gira/internal/ai"
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

func (command Branch) Run(issueID string, assign bool, enableAI bool, force bool) {
	issue := command.tracker.GetIssue(issueID)
	command.RunWithIssue(issue, assign, enableAI, force)
}

func (command Branch) RunWithIssue(issue *issue.Issue, assign bool, enableAI, force bool) {
	generatedBranch := command.branch.FromIssue(issue, &branch.FromIssueOptions{})

	if enableAI {
		agent := ai.NewOpenAI(command.logger)
		response, err := agent.BranchNames(issue)

		if err == nil && len(response) > 0 {
			var branches []*branch.Branch
			branches = append(branches, generatedBranch)

			for _, iaGeneratedBranchName := range response {
				branches = append(branches, command.branch.FromIssue(issue, &branch.FromIssueOptions{
					TitleOverride: iaGeneratedBranchName,
				}))
			}

			generatedBranch = command.selectBranchName(branches)
		}
	}

	if command.git.IsBranchExist(generatedBranch.Raw) {
		command.logger.Warn("⚠️ Branch named %s already exists", generatedBranch.Raw)

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
		command.git.SwitchBranch(generatedBranch.Raw)
		command.logger.Info("✅ %s has just been checkout", generatedBranch.Raw)
	} else {
		if !force {
			command.logger.Info("🌳 Branch will be generated\nPress %s to continue, %s to cancel", "ENTER", log.ErrorStyle.Render("CTRL+C"))
			branchNamePrompt := promptui.Prompt{
				Label:     "Branch",
				Default:   generatedBranch.Raw,
				Pointer:   promptui.PipeCursor,
				AllowEdit: true,
			}
			newBranchName, err := branchNamePrompt.Run()
			if err != nil {
				command.logger.Fatal("The operation was %s", "canceled")
			}
			generatedBranch.Raw = newBranchName
		}

		command.git.CreateBranch(generatedBranch.Raw)
		command.logger.Info("✅ %s branch was created", generatedBranch.Raw)
	}

	if assign {
		if assignError := command.tracker.SelfAssignIssue(generatedBranch.IssueID); assignError != nil {
			command.logger.Debug("%v", assignError)
			command.logger.Info("❌ %s Unable to assign issue %s...", log.ErrorStyle.Render("Oups..."), issue.ID)
		} else {
			command.logger.Info("✅ Jira %s has been assigned to %s", issue.ID)
		}
	}
}

func (command *Branch) selectBranchName(branches []*branch.Branch) *branch.Branch {
	items := make([]string, len(branches))
	for index, branch := range branches {
		items[index] = branch.Raw
	}

	typeSelect := promptui.Select{Label: "Branch name", Items: items}
	index, _, err := typeSelect.Run()

	if err != nil {
		command.logger.Fatal("The operation was %s", "canceled")
	}

	return branches[index]
}
