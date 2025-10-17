#!/bin/bash

# cogmit installation script
# Usage: curl -fsSL https://raw.githubusercontent.com/nicoaudy/cogmit/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="nicoaudy/cogmit"
BINARY_NAME="cogmit"
INSTALL_DIR="/usr/local/bin"

# Get the latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Detect platform and architecture
detect_platform() {
    local os arch

    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          echo -e "${RED}❌ Unsupported operating system: $(uname -s)${NC}" >&2; exit 1 ;;
    esac

    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        armv7l)         arch="arm64" ;; # Fallback for some ARM systems
        *)              echo -e "${RED}❌ Unsupported architecture: $(uname -m)${NC}" >&2; exit 1 ;;
    esac

    echo "${os}-${arch}"
}

# Download and install binary
install_binary() {
    local version=$1
    local platform=$2
    local download_url="https://github.com/${REPO}/releases/download/${version}/cogmit-${platform}"

    # Add .exe extension for Windows
    if [ "${platform#windows-}" != "$platform" ]; then
        download_url="${download_url}.exe"
    fi

    echo -e "${BLUE}📥 Downloading cogmit ${version} for ${platform}...${NC}"

    # Create temp directory
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    # Download binary
    if ! curl -fsSL "$download_url" -o "$BINARY_NAME"; then
        echo -e "${RED}❌ Failed to download binary${NC}" >&2
        exit 1
    fi

    # Make executable
    chmod +x "$BINARY_NAME"

    # Check if we can write to install directory
    if [ -w "$INSTALL_DIR" ]; then
        # We have write permissions
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        # Need sudo
        echo -e "${YELLOW}🔐 Installing to ${INSTALL_DIR} requires sudo permissions${NC}"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    fi

    # Cleanup
    cd /
    rm -rf "$temp_dir"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local installed_version=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        echo -e "${GREEN}✅ cogmit installed successfully!${NC}"
        echo -e "${GREEN}   Version: ${installed_version}${NC}"
        echo -e "${GREEN}   Location: $(which $BINARY_NAME)${NC}"
        echo ""
        echo -e "${BLUE}🚀 Next steps:${NC}"
        echo -e "   1. Run: ${YELLOW}cogmit setup${NC} to configure your settings"
        echo -e "   2. Run: ${YELLOW}cogmit generate${NC} to generate commit messages"
        echo -e "   3. Run: ${YELLOW}cogmit --help${NC} to see all options"
    else
        echo -e "${RED}❌ Installation failed - cogmit not found in PATH${NC}" >&2
        exit 1
    fi
}

# Main installation function
main() {
    echo -e "${BLUE}🤖 Installing cogmit...${NC}"
    echo ""

    # Get latest version
    echo -e "${BLUE}🔍 Checking for latest version...${NC}"
    local version=$(get_latest_version)
    if [ -z "$version" ]; then
        echo -e "${RED}❌ Failed to get latest version${NC}" >&2
        exit 1
    fi
    echo -e "${GREEN}   Latest version: ${version}${NC}"

    # Detect platform
    echo -e "${BLUE}🔍 Detecting platform...${NC}"
    local platform=$(detect_platform)
    echo -e "${GREEN}   Platform: ${platform}${NC}"

    # Install binary
    install_binary "$version" "$platform"

    # Verify installation
    verify_installation
}

# Run main function
main "$@"
