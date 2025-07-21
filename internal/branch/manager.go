package branch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Ealenn/gira/internal/git"
	"github.com/Ealenn/gira/internal/issue"
	"github.com/Ealenn/gira/internal/log"
)

type BranchManager struct {
	logger  *log.Logger
	git     *git.Git
	tracker issue.Tracker
}

func NewBranchManager(logger *log.Logger, git *git.Git, tracker issue.Tracker) *BranchManager {
	return &BranchManager{
		logger,
		git,
		tracker,
	}
}

func (branch *BranchManager) GetCurrentBranch() *Branch {
	currentBranch, currentBranchError := branch.git.CurrentBranch()
	if currentBranchError != nil {
		branch.logger.Fatal("❌ Unable to check current branch")
	}

	branch.logger.Debug("🔎 Current branch %s", currentBranch)
	branchNameParts := strings.Split(currentBranch, `/`)
	if len(branchNameParts) < 2 {
		branch.logger.Fatal("❌ Unable to find issue in branch name %s", currentBranch)
	}

	return &Branch{
		Type:    branch.getBranchType([]string{branchNameParts[0]}),
		IssueID: branchNameParts[1],
		Title:   branchNameParts[2],
	}
}

func (branch *BranchManager) Generate(issue *issue.Issue) *Branch {
	branchTitle := strings.ToLower(strings.TrimSpace(issue.Title))
	branchTitle = regexp.MustCompile(`^\[[^\]]+\]\s*`).ReplaceAllString(branchTitle, "")
	branchTitle = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(branchTitle, "")
	branchTitle = strings.Join(strings.Fields(branchTitle), " ")
	branchTitle = strings.ReplaceAll(branchTitle, " ", "-")

	branchType := branch.getBranchType(issue.Types)

	branchRaw := fmt.Sprintf("%s/%s/%s", branchType, issue.ID, branchTitle)

	return &Branch{
		Type:    branchType,
		IssueID: issue.ID,
		Title:   branchTitle,
		Raw:     branchRaw,
	}
}

func (branch *BranchManager) getBranchType(issueTypes []string) BranchType {
	for _, issueType := range issueTypes {
		switch strings.ToLower(issueType) {
		case "bug":
			return BranchType(Bug)
		case "tasks":
		case "feature":
		case "enhancement":
			return BranchType(Feature)
		}
	}

	return BranchType(Feature)
}
