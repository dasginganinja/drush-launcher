#!/usr/bin/env bash

set -eou pipefail

# Define the GitHub repository owner, repository name, and binary name
repo_owner="dasginganinja"
repo_name="drush-launcher"
binary_name="drush-launcher"

# Function to download the release asset from GitHub
download_release_asset() {
    local ASSET="$1"
    download_url="https://github.com/${repo_owner}/${repo_name}/releases/download/${latest_release}/${ASSET}"
    echo "Downloading ${latest_release} release..."
    curl -s -L -o "${ASSET}" "$download_url"

    echo "Verifying checksum"
    curl -s -L -o checksums.txt "https://github.com/${repo_owner}/${repo_name}/releases/download/${latest_release}/${binary_name}_${latest_version}_checksums.txt"
    CHECKSUM_SOURCE=$(grep "${ASSET}" checksums.txt)
    CHECKSUM_DL=$(sha256sum ${ASSET})
    if [ "$CHECKSUM_SOURCE" != "$CHECKSUM_DL" ]; then
      echo "Checksums do not match."
      rm "${ASSET}" checksums.txt
      exit 1
    fi
    rm checksums.txt
}

# Check if the user specified the target architecture and OS via environment variables
if [ "${GOOS+1}" ] && [ "${GOARCH+1}" ]; then
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
latest_version=$(echo "${latest_release}" | cut -c 2-)  # Remove the leading 'v' from the version

# Check if the latest release version is available
if [ -z "$latest_release" ]; then
    echo "Failed to fetch the latest release version."
    exit 1
fi

# Determine the binary extension based on the operating system
if [[ "$GOOS" == "windows" ]]; then
    binary_extension=".zip"
else
    binary_extension=".tar.gz"
fi

# Check if the required binary for the user's architecture exists in the latest release
binary_filename="${binary_name}_${latest_version}_${GOOS}_${GOARCH}${binary_extension}"
echo "Binary filename: " $binary_filename
release_assets=$(curl -s "https://api.github.com/repos/${repo_owner}/${repo_name}/releases/tags/${latest_release}" | jq -r '.assets[].name')
if [[ "$release_assets" == *"$binary_filename"* ]]; then
    echo "Binary for ${GOOS}-${GOARCH} architecture found in the latest release."
    download_release_asset "${binary_filename}"
    if [[ "$GOOS" == "windows" ]]; then
        unzip "${binary_filename}"
    else
        tar -zxf "${binary_filename}"
    fi
    rm "${binary_filename}"
else
    echo "Binary for ${GOOS}-${GOARCH} architecture not found in the latest release. Attempting to build locally..."
    go build -o "${binary_name}"
fi

# Make the binary executable
chmod +x "${binary_name}"

echo "The ${binary_name} binary is ready."
