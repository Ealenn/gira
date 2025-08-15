package browser

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/Ealenn/gira/internal/log"
)

type Browser struct {
	logger *log.Logger
}

func NewBrowser(logger *log.Logger) *Browser {
	return &Browser{
		logger,
	}
}

func (browser Browser) Open(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		if browser.isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	}

	if len(args) > 1 {
		// args[0] is used for 'start' command argument, to prevent issues with URLs starting with a quote
		args = append(args[:1], append([]string{""}, args[1:]...)...)
	}

	if err := exec.Command(cmd, args...).Start(); err != nil {
		browser.logger.Debug("Open browser exception %v", err)
		browser.logger.Fatal("Unable to open link %s", url)
	}
}

func (browser Browser) isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
