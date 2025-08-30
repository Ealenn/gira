package command

import (
	"github.com/Ealenn/gira/internal/ai"
	"github.com/Ealenn/gira/internal/branch"
	"github.com/Ealenn/gira/internal/command/forms"
	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
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

func (cmd Branch) Run(issueID string, assign bool, enableAI bool, force bool) {
	issue := cmd.tracker.GetIssue(issueID)
	cmd.RunWithIssue(issue, assign, enableAI, force)
}

func (cmd Branch) RunWithIssue(issue *issue.Issue, assign bool, enableAI, force bool) {
	generatedBranch := cmd.branch.FromIssue(issue, &branch.FromIssueOptions{})

	if enableAI {
		agent := ai.NewOpenAI(cmd.logger)
		response, err := agent.BranchNames(issue)

		if err == nil && len(response) > 0 {
			var branches []*branch.Branch
			branches = append(branches, generatedBranch)

			for _, iaGeneratedBranchName := range response {
				branches = append(branches, cmd.branch.FromIssue(issue, &branch.FromIssueOptions{
					TitleOverride: iaGeneratedBranchName,
				}))
			}

			generatedBranch = &forms.NewSelect(cmd.logger).Ask("🔎 Choose the branch name to create", "(You will be able to edit it afterwards)", branches...).Branch
		}
	}

	if !force {
		forms.NewEditBranch(cmd.logger).Ask("✒️ Tweak branch name before creating?", "", generatedBranch)
	}

	if cmd.git.IsBranchExist(generatedBranch.Raw) {
		cmd.logger.Warn("⚠️ Branch named %s already exists", generatedBranch.Raw)

		if !force {
			if !forms.NewConfirm(cmd.logger).Ask("♻️ Would you like to switch to this branch?", generatedBranch.Raw, forms.TypeYesNo).Confirmed {
				cmd.logger.Fatal("The operation was %s", "canceled")
			}
		}
		cmd.git.SwitchBranch(generatedBranch.Raw) // TODO: Handle error
		cmd.logger.Info("✅ %s has just been checkout", generatedBranch.Raw)
		return
	}

	if !force {
		if !forms.NewConfirm(cmd.logger).Ask(
			"🌳 Create this branch?",
			generatedBranch.Raw,
			forms.TypeConfirm,
		).Confirmed {
			cmd.logger.Fatal("❌ The operation was %s", "canceled")
		}
	}

	cmd.git.CreateBranch(generatedBranch.Raw)
	cmd.logger.Info("✅ %s branch was created", generatedBranch.Raw)

	if assign {
		if assignError := cmd.tracker.SelfAssignIssue(generatedBranch.IssueID); assignError != nil {
			cmd.logger.Debug("%v", assignError)
			cmd.logger.Info("❌ %s Unable to assign issue %s...", log.ErrorStyle.Render("Oups..."), issue.ID)
		} else {
			cmd.logger.Info("✅ Jira %s has been assigned", issue.ID)
		}
	}
}
