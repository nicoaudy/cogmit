# 🚀 Releasing cogmit

This document explains how to create releases for cogmit.

## 📋 Prerequisites

- Git repository with GitHub remote
- GitHub Actions enabled for the repository
- Push access to the repository

## 🏷️ Creating a Release

### 1. Update Version

Update the version in your code if needed, then commit your changes:

```bash
git add .
git commit -m "feat: prepare for v1.0.0 release"
git push origin main
```

### 2. Create and Push Tag

Use the Makefile to create a release:

```bash
# Create release v1.0.0
make release VERSION=1.0.0

# Or manually:
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 3. GitHub Actions

Once you push the tag, GitHub Actions will automatically:

1. ✅ Run tests
2. 🔨 Build binaries for all platforms:
   - `cogmit-linux-amd64`
   - `cogmit-linux-arm64`
   - `cogmit-darwin-amd64`
   - `cogmit-darwin-arm64`
   - `cogmit-windows-amd64.exe`
   - `cogmit-windows-arm64.exe`
3. 📦 Create a GitHub release with all binaries
4. 🔐 Generate checksums for verification

### 4. Verify Release

Check the [Releases page](https://github.com/nicoaudy/cogmit/releases) to see your new release with all the compiled binaries.

## 🎯 Supported Platforms

- **Linux**: amd64, arm64
- **macOS**: amd64, arm64 (Apple Silicon)
- **Windows**: amd64, arm64

## 📝 Release Notes

GitHub Actions will automatically generate release notes based on your commits since the last release. You can edit these in the GitHub UI after the release is created.

## 🔧 Manual Release (if needed)

If you need to create a release manually:

```bash
# Build all platforms
make build-all

# Create checksums
cd build
sha256sum * > checksums.txt

# Upload to GitHub releases manually
```

## 🐛 Troubleshooting

**Release not created?**
- Check GitHub Actions tab for any failed workflows
- Ensure the tag was pushed: `git ls-remote --tags origin`
- Verify GitHub Actions is enabled in repository settings

**Missing binaries?**
- Check the Actions log for build errors
- Ensure all platforms built successfully
- Verify the workflow file is in `.github/workflows/`

**Version not showing correctly?**
- Check that the version is set correctly in the tag
- Verify the LDFLAGS in the Makefile
- Test locally with `make build && ./build/cogmit --version`
