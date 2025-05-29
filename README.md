# VimCoplit

VimCoplit is a Neovim implementation of [Cline](https://github.com/cline/cline), an AI coding assistant that integrates directly into your editor. This project aims to bring the powerful capabilities of Cline to Neovim users, with a Go backend for performance and reliability.

## Features

- **AI-Powered Coding Assistant**: Leverage Claude 3.7 Sonnet's agentic coding capabilities directly in Neovim
- **File Operations**: Create, edit, and manage files with AI assistance
- **Terminal Integration**: Execute commands and monitor output
- **Browser Automation**: Test web applications with headless browser capabilities
- **Custom Tools**: Extend functionality through the Model Context Protocol (MCP)
- **Context Management**: Add URLs, problems, files, and folders to provide context to the AI

## Architecture

VimCoplit consists of two main components:

1. **Neovim Plugin**: Written in Lua, providing the user interface and integration with Neovim
2. **Go Backend**: Handling the heavy lifting of AI interactions, file operations, and system commands

## Requirements

- Neovim 0.9.0 or higher
- Go 1.21 or higher
- [Claude API key](https://console.anthropic.com/) or other supported AI provider

## Installation

```bash
# Using your preferred Neovim plugin manager
# For example, with lazy.nvim:
{
  "yourusername/vimcoplit",
  dependencies = {
    "nvim-lua/plenary.nvim",
    "nvim-telescope/telescope.nvim",
  },
  build = "go build -o bin/vimcoplit ./cmd/vimcoplit",
}
```

## Configuration

```lua
-- In your Neovim config
require('vimcoplit').setup({
  api_key = "your-claude-api-key",
  model = "claude-3-sonnet-20240229",
  -- Additional configuration options
})
```

## Usage

- `:VimCoplit` - Open the VimCoplit interface
- `:VimCoplitTask <task>` - Start a new task
- `:VimCoplitAddContext` - Add context to the current task

## Development

### Project Structure

```
vimcoplit/
├── cmd/                # Go backend entry points
├── internal/           # Go backend implementation
├── lua/                # Neovim plugin code
│   └── vimcoplit/      # Plugin modules
├── scripts/            # Build and utility scripts
└── test/               # Test files
```

### Building

```bash
# Build the Go backend
go build -o bin/vimcoplit ./cmd/vimcoplit

# Build the Neovim plugin
nvim --headless -c "luafile scripts/build.lua" -c "quit"
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache 2.0 © 2025 VimCoplit Contributors 