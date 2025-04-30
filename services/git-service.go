package services

import (
	"os/exec"
	"strings"
)

type GitService struct {
	loggerService *LoggerService
}

func NewGitService(loggerService *LoggerService) *GitService {
	return &GitService{
		loggerService,
	}
}

func (gitService *GitService) CreateBranch(name string) []byte {
	cmd := exec.Command("git", "checkout", "-b", name)
	output, err := cmd.CombinedOutput()

	if err != nil {
		gitService.loggerService.Fatal("%s", output)
	}

	return output
}

func (gitService *GitService) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	response, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(response)), err
}

func (gitService *GitService) IsBranchExist(name string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", name)
	_, err := cmd.CombinedOutput()

	return err == nil
}

func (gitService *GitService) SwitchBranch(name string) bool {
	cmd := exec.Command("git", "checkout", name)
	_, err := cmd.CombinedOutput()

	return err == nil
}
