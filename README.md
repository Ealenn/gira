# 🦒 Gira — Git & Jira CLI Tool

[![GitHub stars](https://img.shields.io/github/stars/Ealenn/gira?style=for-the-badge&logo=github)](https://github.com/Ealenn/gira/stargazers) 
[![GitHub issues](https://img.shields.io/github/issues/Ealenn/gira?style=for-the-badge&logo=github)](https://github.com/Ealenn/gira/issues) 
[![Github download](https://img.shields.io/github/downloads/ealenn/gira/total?style=for-the-badge&logo=github)](https://github.com/Ealenn/gira/releases) 
[![DockerHub](https://img.shields.io/docker/pulls/ealen/gira.svg?style=for-the-badge&logo=docker)](https://hub.docker.com/repository/docker/ealen/gira) 
[![License](https://img.shields.io/github/license/ealenn/gira?style=for-the-badge&logo=opensourceinitiative)](https://github.com/Ealenn/gira?tab=GPL-3.0-1-ov-file)

Gira is a simple and powerful command-line tool that bridges your Git workflow with Jira. It helps you automate common tasks like creating branches from Jira issues, and updating or closing issues — all from your terminal.

Use Gira to speed up development workflows, reduce copy-pasting from Jira to Git, and keep your issue tracking in sync with your commits.

- [🦒 Gira — Git \& Jira CLI Tool](#-gira--git--jira-cli-tool)
  - [📦 Installation](#-installation)
    - [🔧 Automatic](#-automatic)
      - [🐧 Linux / MacOS](#-linux--macos)
      - [🪟 Windows](#-windows)
    - [👉 Manual (Download the Binary)](#-manual-download-the-binary)
    - [🐳 From Docker (No Install Required)](#-from-docker-no-install-required)
  - [✨ Shell Autocompletion](#-shell-autocompletion)
  - [🚀 Usage](#-usage)
    - [⚙️ `configuration`: Configure Gira with Jira account and API token](#️-configuration-configure-gira-with-jira-account-and-api-token)
    - [🌱 `branch`: Create a new Git branch using Jira issue ID](#-branch-create-a-new-git-branch-using-jira-issue-id)
      - [Usage](#usage)
      - [Example](#example)

## 📦 Installation

You can use Gira either as a native binary or through Docker. Choose what fits your environment best.

### 🔧 Automatic

Automated install/update, don't forget to always verify what you're piping into bash.

The script installs downloaded binary to HOME directory by default, but it can be changed by setting DIR environment variable.

#### 🐧 Linux / MacOS

```sh
curl https://raw.githubusercontent.com/Ealenn/gira/master/install_linux.sh | bash
```

#### 🪟 Windows

```sh
Invoke-RestMethod https://raw.githubusercontent.com/Ealenn/gira/master/install_windows.ps1 | Invoke-Expression
```

### 👉 Manual (Download the Binary)

|Platform|Download Link|Hash|
|--------|-------------|----|
|macOS (Intel)|[gira-darwin-amd64](https://github.com/Ealenn/gira/releases/latest/download/gira-darwin-amd64)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-darwin-amd64.md5)|
|macOS (ARM)|[gira-darwin-arm64](https://github.com/Ealenn/gira/releases/latest/download/gira-darwin-arm64)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-darwin-arm64.md5)|
|Linux (x86)|[gira-linux-386](https://github.com/Ealenn/gira/releases/latest/download/gira-linux-386)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-linux-386.md5)|
|Linux (AMD64)|[gira-linux-amd64](https://github.com/Ealenn/gira/releases/latest/download/gira-linux-amd64)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-linux-amd64.md5)|
|Linux (ARM64)|[gira-linux-arm64](https://github.com/Ealenn/gira/releases/latest/download/gira-linux-arm64)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-linux-arm64.md5)|
|Windows (x86)|[gira-windows-386.exe](https://github.com/Ealenn/gira/releases/latest/download/gira-windows-386.exe)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-windows-386.exe.md5)|
|Windows (AMD64)|[gira-windows-amd64.exe](https://github.com/Ealenn/gira/releases/latest/download/gira-windows-amd64.exe)|[md5](https://github.com/Ealenn/gira/releases/latest/download/gira-windows-amd64.exe.md5)|

You can download the latest release directly from GitHub:

```sh
curl -sSL https://github.com/Ealenn/gira/releases/latest/download/gira-linux-amd64 -o /usr/local/bin/gira 
chmod +x /usr/local/bin/gira
```

> This will place the gira binary in your system path for global use.

### 🐳 From Docker (No Install Required) 

If you prefer using Docker, you can run Gira directly without installing it:

```sh
docker run -it --rm -v "$HOME:/root" -v "$PWD:/app" -w /app ealen/gira
```

To make it easier to use Gira like a native CLI, add this alias to your shell config:

```sh
# In your ~/.bashrc or ~/.zshrc
alias gira='docker run -it --rm -v "$HOME:/root" -v "$PWD:/app" -w /app ealen/gira'
```

After reloading your shell, you'll be able to run gira from anywhere.

## ✨ Shell Autocompletion

Gira supports autocompletion for major shells like Bash, Zsh, Fish, and PowerShell.

Generate the completion script with:

```sh
$ gira completion [shell]
```

Supported shells: `bash` – `zsh` – `fish` – `powershell`

Example with ZSH :

```sh
# in your ~/.zshrc
eval "$(gira completion zsh)"
```

> This enables tab completion for Gira commands and flags in your shell.

## 🚀 Usage

```
Usage:
  gira [command]

Available Commands:
  branch      Create a new Git branch using Jira issue ID.
  configure   Configure Gira with Jira account and API token
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Display the current Gira version and check for available updates

Flags:
  -h, --help      help for gira
  -v, --version   version for gira

Use "gira [command] --help" for more information about a command.
```

Use the DEBUG environment variable to display detailed exception stack traces.

```sh
❯ DEBUG=1 gira branch TEST-123
[DEBUG] Issue TEST-123 response status 404 
[FATAL] Unable to find Jira TEST-123
```

### ⚙️ `configuration`: Configure Gira with Jira account and API token

Configures the Gira CLI by setting up the Jira account credentials, including the Jira host URL, email, and API token.
This command updates the configuration file to enable communication with the Jira instance for subsequent commands like 'branch'.
Ensure you have a valid Jira API token from your Atlassian account before running this command.

```
❯ gira configure
Enter the Jira API URL (Example https://jira.mycompagny.com): https://jira.mycompagny.com
Enter the Jira Token (See /manage-profile/security/api-tokens): **********
✅ Done!
```

### 🌱 `branch`: Create a new Git branch using Jira issue ID

Creates a new Git branch based on Jira issue.

The branch name is generated by combining the Jira issue ID with a slugified version of the issue summary (e.g., "feature/ABC-123/fix-login-bug"). 

This helps enforce consistent naming conventions and improve traceability between code and Jira issues.

#### Usage 
```
Usage:
  gira branch [Jira Issue ID] [flags]

Examples:
gira branch ISSUE-123
gira branch -a ISSUE-123

Aliases:
  branch, checkout

Flags:
  -a, --assignIssue   assign the issue to the currently logged-in Jira user after creating the Git branch
  -h, --help          help for branch
```

#### Example
```
❯ gira branch ISSUE-123
Branch feature/ISSUE-123/update-app-dependencies-to-the-latest-version will be generated
Press ENTER to continue, CTRL+C to cancel
```
