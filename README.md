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
    - [⚙️ From Source (Download the Binary)](#️-from-source-download-the-binary)
    - [🐳 From Docker (No Install Required)](#-from-docker-no-install-required)
  - [✨ Shell Autocompletion](#-shell-autocompletion)
  - [🚀 Usage](#-usage)
    - [🌱 branch: Creating a Branch from Jira](#-branch-creating-a-branch-from-jira)

## 📦 Installation

You can use Gira either as a native binary or through Docker. Choose what fits your environment best.

### ⚙️ From Source (Download the Binary)

You can download the latest release directly from GitHub:

```sh
curl -sSL https://github.com/Ealenn/gira/releases/latest/download/gira-linux-amd64 -o /usr/local/bin/gira 
chmod +x /usr/local/bin/gira
```

> This will place the gira binary in your system path for global use.

### 🐳 From Docker (No Install Required) 

If you prefer using Docker, you can run Gira directly without installing it:

```sh
docker run --rm -v "$HOME:/root" -v "$PWD:/app" -w /app ealen/gira
```

To make it easier to use Gira like a native CLI, add this alias to your shell config:

```sh
# In your ~/.bashrc or ~/.zshrc
alias gira='docker run --rm -v "$HOME:/root" -v "$PWD:/app" -w /app ealen/gira'
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
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help   help for gira

Use "gira [command] --help" for more information about a command.
```

### 🌱 branch: Creating a Branch from Jira

```
❯ gira branch TICKET-123
Branch feature/TICKET-123/update-app-dependencies-to-the-latest-version will be generated
Press ENTER to continue, CTRL+C to cancel
```
