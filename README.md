# 🤖 cogmit

**AI-powered Git commit message generator using local Ollama models**

cogmit is a beautiful CLI tool that analyzes your Git changes and generates intelligent commit message suggestions using local AI models via Ollama. No data leaves your machine - everything runs locally!

## ✨ Features

- 🧠 **Local AI**: Uses Ollama for completely local AI processing
- 🎨 **Beautiful UI**: Interactive selection with arrow key navigation using Bubble Tea
- ⚙️ **Configurable**: Customize model, host, and behavior settings
- 🔒 **Privacy-First**: All processing happens locally, no external API calls
- 🚀 **Fast**: Lightweight models like `llama3.2:1b` for quick responses
- 📝 **Conventional Commits**: Generates properly formatted commit messages

## 🚀 Quick Start

### Prerequisites

1. **Ollama**: Install and run Ollama locally
   ```bash
   # Install Ollama (if not already installed)
   curl -fsSL https://ollama.ai/install.sh | sh

   # Pull a lightweight model
   ollama pull llama3.2:1b

   # Start Ollama (if not running)
   ollama serve
   ```

2. **Go**: Go 1.19 or later

### Installation

**Option 1: One-liner install (Recommended)**
```bash
curl -fsSL https://raw.githubusercontent.com/nicoaudy/cogmit/main/install.sh | bash
```

**Option 2: Manual installation**
```bash
# Download the binary for your platform
# Linux (amd64)
wget https://github.com/nicoaudy/cogmit/releases/latest/download/cogmit-linux-amd64
chmod +x cogmit-linux-amd64
sudo mv cogmit-linux-amd64 /usr/local/bin/cogmit

# macOS (Intel)
wget https://github.com/nicoaudy/cogmit/releases/latest/download/cogmit-darwin-amd64
chmod +x cogmit-darwin-amd64
sudo mv cogmit-darwin-amd64 /usr/local/bin/cogmit

# macOS (Apple Silicon)
wget https://github.com/nicoaudy/cogmit/releases/latest/download/cogmit-darwin-arm64
chmod +x cogmit-darwin-arm64
sudo mv cogmit-darwin-arm64 /usr/local/bin/cogmit

# Windows (amd64)
# Download cogmit-windows-amd64.exe and add to PATH
```

**Option 3: Build from source**
```bash
# Clone the repository
git clone https://github.com/nicoaudy/cogmit.git
cd cogmit

# Build the binary
go build -o cogmit .

# Install globally (optional)
sudo mv cogmit /usr/local/bin/
```

### Setup

Configure cogmit with your preferences:

```bash
cogmit setup
```

This will prompt you for:
- Ollama host (default: `http://localhost:11434`)
- Model name (default: `llama3.2:1b`)
- Number of suggestions (default: `3`)
- Auto-commit behavior (default: `false`)

## 📖 Usage

### Basic Usage

```bash
# Generate commit messages for staged changes
cogmit

# Or explicitly use the generate command
cogmit generate
```

### Workflow Example

```bash
# 1. Make some changes to your code
echo "console.log('Hello, World!');" > hello.js

# 2. Stage your changes
git add hello.js

# 3. Generate commit messages
cogmit
```

You'll see an interactive interface like this:

```
🔍 Analyzing changes...
📝 Found staged changes
🤖 Generating commit messages using llama3.2:1b...
✨ Generated 3 commit message suggestions

🤖 Choose a commit message:
> feat: add hello world console log
  fix: add missing console.log statement
  chore: add hello.js file

↑/↓ or k/j: navigate • enter: select • e: edit • q: quit
```

### Interactive Controls

- **↑/↓ or k/j**: Navigate between options
- **Enter**: Select the highlighted option
- **e**: Edit the highlighted option (future feature)
- **q**: Quit without selecting

## ⚙️ Configuration

Configuration is stored in `~/.config/cogmit/config.json`:

```json
{
  "model": "llama3.2:1b",
  "ollama_host": "http://localhost:11434",
  "num_suggestions": 3,
  "auto_commit": false
}
```

### Recommended Models

For the best balance of speed and quality:

- **`llama3.2:1b`** - Fastest, good for simple changes
- **`llama3.2:3b`** - Better quality, still very fast
- **`llama3.2:7b`** - Best quality, slower but still reasonable

## 🛠️ Development

### Project Structure

```
cogmit/
├── cmd/
│   ├── root.go         # Main command entrypoint
│   ├── setup.go        # Setup command
│   └── generate.go     # Generate command
├── internal/
│   ├── config/         # Configuration management
│   ├── git/           # Git operations
│   ├── ai/            # Ollama API client
│   └── ui/            # Bubble Tea UI components
└── main.go            # Application entrypoint
```

### Building from Source

```bash
# Clone and build
git clone https://github.com/nicoaudy/cogmit.git
cd cogmit
go build -o cogmit .

# Run tests
go test ./...

# Run with debug logging
DEBUG=1 ./cogmit generate
```

## 🐛 Troubleshooting

### Common Issues

**"Ollama API returned status 404"**
- Make sure Ollama is running: `ollama serve`
- Check if the model exists: `ollama list`
- Pull the model: `ollama pull llama3.2:1b`

**"Not in a Git repository"**
- Make sure you're in a Git repository
- Initialize one: `git init`

**"No changes found to commit"**
- Stage some changes: `git add .`
- Or make some changes to your files

**"Failed to connect to Ollama"**
- Check if Ollama is running on the configured host
- Verify the host URL in your config

### Debug Mode

Enable debug logging:

```bash
DEBUG=1 cogmit generate
```

This will create a `debug.log` file with detailed information.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Ollama](https://ollama.ai/) for providing local AI capabilities
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the beautiful TUI framework
- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for terminal styling

## 🔮 Future Ideas

- [ ] Support for multiple AI providers (OpenAI, Gemini, etc.)
- [ ] `--dry-run` mode to preview without committing
- [ ] `cogmit history` to review past generated commits
- [ ] `cogmit config show/edit` commands
- [ ] Conventional commit mode enforcement
- [ ] Custom prompt templates
- [ ] Integration with popular Git hooks

---

**Made with ❤️ for developers who love clean, meaningful commit messages**
