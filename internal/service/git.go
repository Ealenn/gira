package service

import (
	"os/exec"
	"strings"

	"github.com/Ealenn/gira/internal/log"
)

type Git struct {
	logger *log.Logger
}

func NewGit(logger *log.Logger) *Git {
	return &Git{
		logger,
	}
}

func (git *Git) CreateBranch(name string) []byte {
	cmd := exec.Command("git", "checkout", "-b", name)
	output, err := cmd.CombinedOutput()

	if err != nil {
		git.logger.Fatal("%s", output)
	}

	return output
}

func (git *Git) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	response, err := cmd.CombinedOutput()

	return strings.TrimSpace(string(response)), err
}

func (git *Git) IsBranchExist(name string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", name)
	_, err := cmd.CombinedOutput()

	return err == nil
}

func (git *Git) SwitchBranch(name string) bool {
	cmd := exec.Command("git", "checkout", name)
	_, err := cmd.CombinedOutput()

	return err == nil
}
