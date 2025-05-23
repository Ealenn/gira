#!/usr/bin/env sh

set -eu
printf '\n'

BOLD="$(tput bold 2>/dev/null || printf '')"
GREY="$(tput setaf 0 2>/dev/null || printf '')"
UNDERLINE="$(tput smul 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
BLUE="$(tput setaf 4 2>/dev/null || printf '')"
MAGENTA="$(tput setaf 5 2>/dev/null || printf '')"
NO_COLOR="$(tput sgr0 2>/dev/null || printf '')"

info() {
  printf '%s\n' "${BOLD}${GREY}>${NO_COLOR} $*"
}

warn() {
  printf '%s\n' "${YELLOW}! $*${NO_COLOR}"
}

error() {
  printf '%s\n' "${RED}x $*${NO_COLOR}" >&2
}

completed() {
  printf '%s\n' "${GREEN}✓${NO_COLOR} $*"
}

# Verify curl is installed
if ! command -v curl >/dev/null 2>&1; then
  error "curl is not installed. Please install curl and try again."
  exit 1
fi

# Allow specifying a custom destination directory
if [ -n "${DIR:-}" ]; then
  # User provided DIR
  :
elif [ "$(uname -s)" = "Darwin" ]; then
  DIR="$HOME/Library/Application Support/gira/bin"
else
  DIR="$HOME/.local/bin"
fi

# Ensure install directory exists
mkdir -p "$DIR"

# Detect architecture
ARCH=$(uname -m)
case $ARCH in
    i386|i686) ARCH="386" ;;
    x86_64) ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

# Construct binary filename
GITHUB_FILE="gira-${OS}-${ARCH}"

# Get latest release tag from GitHub API
GITHUB_LATEST_VERSION=$(curl -s https://api.github.com/repos/Ealenn/gira/releases/latest | grep tag_name | cut -d '"' -f 4)

# Build the download URL
GITHUB_URL="https://github.com/Ealenn/gira/releases/download/${GITHUB_LATEST_VERSION}/${GITHUB_FILE}"
info "Downloading ${GITHUB_FILE} from ${GITHUB_URL}"
printf '\n'

# Download and install
if ! curl --fail --silent -L -o gira "$GITHUB_URL"; then
  error "Failed to download the Gira binary. Please verify that the release exists for your platform."
  warn "More information: https://github.com/Ealenn/gira"
  exit 1
fi

chmod +x gira
mv gira "$DIR/gira"

completed "Gira $GITHUB_LATEST_VERSION installed to $DIR"
printf '\n'

# Suggest adding autocompletion
info "Gira supports autocompletion for major shells like Bash, Zsh, Fish, and PowerShell."
for s in "bash" "zsh" "fish"
do
  config_file="~/.${s}rc"
  config_cmd="eval \"\$(gira completion ${s})\""

  case ${s} in
    fish )
      # shellcheck disable=SC2088
      config_file="~/.config/fish/config.fish"
      ;;
  esac

  printf "  %s\n  Add the following to the end of %s:\n\n\t%s\n\n" \
    "${BOLD}${UNDERLINE}${s}${NO_COLOR}" \
    "${BOLD}${config_file}${NO_COLOR}" \
    "${config_cmd}"
done
printf "\n"

# Optional: suggest adding to PATH if not already in it
case ":$PATH:" in
  *":$DIR:"*) ;;
  *)
    warn "To use 'gira' command everywhere, add this line to your shell config:"
    info "    export PATH=\"\$PATH:$DIR\""
    printf '\n'
    ;;
esac
