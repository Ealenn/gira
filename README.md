# 🦒 Gira — Git, Jira & GitHub Issues CLI Tool  <!-- omit in toc -->

[![GitHub stars](https://img.shields.io/github/stars/Ealenn/gira?style=for-the-badge&logo=github)](https://github.com/Ealenn/gira/stargazers) 
[![GitHub issues](https://img.shields.io/github/issues/Ealenn/gira?style=for-the-badge&logo=github)](https://github.com/Ealenn/gira/issues) 
[![Github download](https://img.shields.io/github/downloads/ealenn/gira/total?style=for-the-badge&logo=github)](https://github.com/Ealenn/gira/releases) 
[![DockerHub](https://img.shields.io/docker/pulls/ealen/gira.svg?style=for-the-badge&logo=docker)](https://hub.docker.com/repository/docker/ealen/gira) 
[![License](https://img.shields.io/github/license/ealenn/gira?style=for-the-badge&logo=opensourceinitiative)](https://github.com/Ealenn/gira?tab=GPL-3.0-1-ov-file)

Gira is a powerful command-line tool that bridges your Git workflow with both **Jira** and **GitHub** issues. 
It helps you automate tasks like creating branches from issue, viewing issue details, and keeping issue tracking in sync with Git — all from your terminal.

Use Gira to speed up development workflows, reduce context switching, and streamline project tracking whether you're using Jira, GitHub, or both.

- [📦 Installation](#-installation)
  - [🔧 Automatic](#-automatic)
  - [👉 Manual](#-manual)
  - [🐳 Docker](#-docker)
- [✨ Shell Autocompletion](#-shell-autocompletion)
- [🚀 Usage](#-usage)
  - [⚙️ `config`: Configure Gira profile with accounts and tokens](#️-config-configure-gira-profile-with-accounts-and-tokens)
    - [Default Profile](#default-profile)
    - [Custom Profiles](#custom-profiles)
  - [🌱 `branch`: Create a new Git branch using issue ID (Jira or GitHub)](#-branch-create-a-new-git-branch-using-issue-id-jira-or-github)
  - [🕵️ `issue`: Show details of issue (from current branch or specified issue ID)](#️-issue-show-details-of-issue-from-current-branch-or-specified-issue-id)
  - [🌐 `open`: Open the issue in your browser](#-open-open-the-issue-in-your-browser)
  - [🥷 `ninja`: Create a new issue and branch in one go](#-ninja-create-a-new-issue-and-branch-in-one-go)

## 📦 Installation

You can use Gira either as a native binary or through Docker. Choose what fits your environment best.

### 🔧 Automatic

Automated install/update, don't forget to always verify what you're piping into bash.

The script installs downloaded binary to HOME directory by default, but it can be changed by setting DIR environment variable.

#### 🐧 Unix (Linux / MacOS) <!-- omit in toc -->

```sh
curl https://raw.githubusercontent.com/Ealenn/gira/master/install_unix.sh | bash
```

#### 🪟 Windows <!-- omit in toc -->

```sh
Invoke-RestMethod https://raw.githubusercontent.com/Ealenn/gira/master/install_windows.ps1 | Invoke-Expression
```

### 👉 Manual

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

### 🐳 Docker

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

> 💡 **Note**: The Gira Docker image is also available on GitHub Container Registry.
> 
> If your company restricts access to docker.io, you can use the GitHub-hosted image instead by replacing `ealen/gira` with `ghcr.io/ealenn/gira` in the above commands:
> ```sh
> alias gira='docker run -it --rm -v "$HOME:/root" -v "$PWD:/app" -w /app ghcr.io/ealenn/gira'
> ```

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
  branch      Create a new Git branch using issue
  completion  Generate the autocompletion script for the specified shell
  config      Configure Gira with accounts and tokens
  help        Help about any command
  issue       Show details of an issue (from current branch or specified issue ID)
  ninja       Create a new issue and associated branch in one command
  open        Open issue in web browser (from current branch or specified issue ID)
  version     Display the current Gira version and check for available updates

Flags:
  -h, --help             help for gira
  -p  --profile string   Configuration profile to use (default "default")
      --verbose          Print detailed operation logs and debug information
  -v, --version          version for gira

Use "gira [command] --help" for more information about a command.
```

🐞 Use the `--verbose` flag to display detailed exception or stack traces, which can help you better understand what went wrong.

This information is especially useful before creating a bug issue, as it provides more context for troubleshooting and reporting problems.

Example : 
```sh
❯ gira branch TEST-123 --verbose
[DEBUG] Issue TEST-123 response status 404 
[FATAL] Unable to find issue TEST-123
```

### ⚙️ `config`: Configure Gira profile with accounts and tokens

The `gira config` command sets up the Gira CLI by allowing you to configure one or more accounts, each with its own credentials. 

You can create multiple profiles to connect to different sources, such as `Jira` or `GitHub` issues, making it easy to switch between environments or accounts.

For each profile, you'll specify the source type along with the necessary credentials:

- For Jira: Provide the Jira host URL and API token.
- For GitHub: Provide optional GitHub personal access token.

This configuration is stored in your local Gira config file and enables the CLI to communicate with the appropriate service when running commands like `branch` or `issue`.

#### Default Profile

Running `gira config` with no additional arguments will set up or update your **default** profile:

```
❯ gira config
Use the arrow keys to navigate: ↓ ↑ → ← 
? Type: 
  ▸ JIRA
    GITHUB

Enter the Jira API URL : https://jira.mycompany.com
Enter the Jira Token : **********
✅ Done!
```

#### Custom Profiles

You can also configure custom profiles and use them in any Gira command by specifying the `--profile` option.

This is useful if you need to work with multiple instances or accounts.

```
❯ gira config --profile perso
Enter the Jira API URL : https://jira.personal.com
Enter the Jira Token : **********
✅ Done!
```

Once configured, you can specify the profile in other commands:

```
❯ gira branch --profile perso
```

or 

```
❯ gira branch -p perso
```

This flexibility allows you to easily manage and switch between multiple Jira or Github accounts or environments as needed.

### 🌱 `branch`: Create a new Git branch using issue ID (Jira or GitHub)

Creates a new Git branch based on issue.

The branch name is generated by combining the issue ID with a slugified version of the issue summary (e.g., "feature/ABC-123/fix-login-bug"). 

This helps enforce consistent naming conventions and improve traceability between code and issues.

#### Usage <!-- omit in toc -->
```
Usage:
  gira branch [issue] [flags]

Aliases:
  branch, checkout

Examples:
  gira branch ISSUE-123
  gira branch -a ISSUE-123

Flags:
  -a, --assign   assign the issue to the currently logged-in user after creating the Git branch
  -f, --force    disable interactive prompts and force branch creation even if checks would normally prevent it
  -h, --help     help for branch
```

#### Example <!-- omit in toc -->
```
❯ gira branch ISSUE-123
Branch feature/ISSUE-123/update-app-dependencies-to-the-latest-version will be generated
Press ENTER to continue, CTRL+C to cancel
```

### 🕵️ `issue`: Show details of issue (from current branch or specified issue ID)

Displays detailed information about an issue.

- If no issue ID is provided, the issue associated with the current Git branch is used.
- If an issue ID is specified, the command will display information for that issue.

This includes the issue key, summary, description, status, priority, assignee, and other relevant metadata.

Useful for quickly reviewing the context of your work without leaving the terminal.

#### Usage <!-- omit in toc -->
```
Usage:
  gira issue [issueId] [flags]

Examples:
  gira issue
  gira issue ABC-123

Flags:
  -h, --help   help for issue
```

#### Example <!-- omit in toc -->
```
❯ gira issue

Issue: PROJ-457
Summary: Fix 500 error when submitting user registration form
Priority: High - Status: In Progress
Assignee: Alice Martin <alice.martin@company.com>
Description:

Users receive a 500 Internal Server Error after submitting the registration form.
The issue appears to be related to missing validation on the email field when the user already exists.

Steps to reproduce:
1. Go to /register
2. Fill the form with an existing email
3. Submit

Expected: A validation message
Actual: 500 error

Refer to the backend error logs and the related ticket: DEVOPS-123.

🔗 More: https://jira.company.com/browse/PROJ-457
```

### 🌐 `open`: Open the issue in your browser

The `gira open` command quickly opens the web page for the current issue (or a specified one) in your default browser.

It works with both Jira and GitHub Issues, making it easy to jump from the terminal directly to the issue tracker for viewing, editing, or commenting.

- If no issue ID is provided, open uses the issue associated with the current Git branch.
- If an issue ID is specified, it will open that issue directly.

#### Usage <!-- omit in toc -->
```
Usage:
  gira gira open [issueId] [flags]

Examples:
  gira open
  gira open ABC-123

Flags:
  -h, --help   help for issue
```

#### Example <!-- omit in toc -->
```
❯ gira open
🔗 Opening https://jira.company.com/browse/PROJ-457
```

### 🥷 `ninja`: Create a new issue and branch in one go

The `gira ninja` command speeds up your workflow by creating a new issue (in Jira or GitHub) and immediately generating a Git branch for it, all in a single step.

This is especially useful when starting work on a brand-new task and wanting to dive straight into coding.

#### Usage <!-- omit in toc -->
```
Usage:
  gira ninja [flags]

Flags:
  -h, --help   help for issue
```

#### Example <!-- omit in toc -->
```
❯ gira ninja

Use the arrow keys to navigate: ↓ ↑ → ← 
? Type: 
  ▸ FEATURE
    BUG
```
