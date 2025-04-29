#!/bin/bash

# Allow specifying a custom destination directory (default: ~/.local/bin)
DIR="${DIR:-"$HOME/.local/bin"}"

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

# Download and install
echo "Downloading ${GITHUB_FILE} from ${GITHUB_URL}..."
curl -L -o gira "$GITHUB_URL" || { echo "Download failed"; exit 1; }

chmod +x gira
mv gira "$DIR/gira"

echo "✅ gira installed to $DIR"

# Optional: suggest adding to PATH if not already in it
if [[ ":$PATH:" != *":$DIR:"* ]]; then
  echo "ℹ️  To use 'gira' everywhere, add this line to your shell config:"
  echo "    export PATH=\"\$PATH:$DIR\""
fi
