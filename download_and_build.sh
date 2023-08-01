#!/bin/bash

# Define the GitHub repository owner, repository name, and binary name
repo_owner="dasginganinja"
repo_name="drush-launcher"
binary_name="drush-launcher"

# Function to download the release asset from GitHub
download_release_asset() {
    download_url="https://github.com/${repo_owner}/${repo_name}/releases/download/${latest_release}/${binary_name}-${GOOS}-${GOARCH}"
    echo "Downloading ${latest_release} release..."
    curl -L -o "${binary_name}" "$download_url"
}

# Check if the user specified the target architecture and OS via environment variables
if [ -n "$GOOS" ] && [ -n "$GOARCH" ]; then
    echo "Using user-specified target architecture: ${GOOS}-${GOARCH}"
else
    # Detect the operating system
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        GOOS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        GOOS="darwin"
    elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
        GOOS="windows"
    else
        echo "Unsupported operating system: $OSTYPE"
        exit 1
    fi

    # Detect the architecture
    if [[ "$GOOS" == "windows" ]]; then
        GOARCH="$PROCESSOR_ARCHITECTURE"
    else
        GOARCH=$(uname -m)
    fi
fi

# Fetch the latest release version using GitHub API
latest_release=$(curl -s "https://api.github.com/repos/${repo_owner}/${repo_name}/releases/latest" | jq -r '.tag_name')

# Check if the latest release version is available
if [ -z "$latest_release" ]; then
    echo "Failed to fetch the latest release version."
    exit 1
fi

# Check if the required binary for the user's architecture exists in the latest release
binary_filename="${binary_name}-${GOOS}-${GOARCH}"
release_assets=$(curl -s "https://api.github.com/repos/${repo_owner}/${repo_name}/releases/tags/${latest_release}" | jq -r '.assets[].name')
if [[ "$release_assets" == *"$binary_filename"* ]]; then
    echo "Binary for ${GOOS}-${GOARCH} architecture found in the latest release."
    download_release_asset
else
    echo "Binary for ${GOOS}-${GOARCH} architecture not found in the latest release. Attempting to build locally..."
    go build -o "${binary_name}"
fi

# Make the binary executable
chmod +x "${binary_name}"

echo "The ${binary_name} binary is ready."
