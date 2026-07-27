# GitForge

<div align="center">

![GitForge Logo](https://raw.githubusercontent.com/phantom1785322-hub/gitforge/main/assets/logo.png)

**A beautiful, powerful Git client for your terminal and browser.**

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/phantom1785322-hub/gitforge/actions/workflows/ci.yml/badge.svg)](https://github.com/phantom1785322-hub/gitforge/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/phantom1785322-hub/gitforge)](https://github.com/phantom1785322-hub/gitforge/releases)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows%20%7C%20termux%20%7C%20bsd-lightgrey)](https://github.com/phantom1785322-hub/gitforge/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/phantom1785322-hub/gitforge)](https://goreportcard.com/report/github.com/phantom1785322-hub/gitforge)
[![Go Reference](https://pkg.go.dev/badge/github.com/phantom1785322-hub/gitforge.svg)](https://pkg.go.dev/github.com/phantom1785322-hub/gitforge)
[![Discord](https://img.shields.io/discord/123456789?color=5865F2&logo=discord&logoColor=white)](https://discord.gg/gitforge)

[Installation](#installation) • [Features](#features) • [Usage](#usage) • [Configuration](#configuration) • [Plugins](#plugins) • [Contributing](#contributing) • [Security](#security) • [Code of Conduct](#code-of-conduct)

</div>

---

## Why GitForge?

You know that feeling when you're in the flow, coding away, and then you need to git something? You switch to the terminal, type `git status`, `git add -p`, `git commit -m "..."`, `git push`... and suddenly you've lost your momentum.

GitForge brings the Git experience into the 2020s. It's a **terminal UI that feels like a modern app**, a **web UI that works offline**, and an **AI assistant that actually helps** — all in one tool that runs everywhere.

### The Philosophy

- **Terminal-first, not terminal-only** — Beautiful TUI with mouse support, but also a full web UI
- **Local-first, cloud-optional** — Works completely offline, syncs when you want
- **AI that respects your code** — Local models via WASM, nothing leaves your machine
- **Extensible by default** — Everything is a plugin, even core features
- **Runs everywhere** — Linux, macOS, Windows, Termux (Android), FreeBSD, OpenBSD

---

## Features

### Terminal UI (TUI)
```
┌─────────────────────────────────────────────────────────────────────────────┐
│ 📁 my-project          🌿 main ↑2 ↓1    🔗 origin                            │
│ ┌─────────────┐ ┌────────────────────────────────────┐ ┌─────────────────┐ │
│ │ 📁 Repo     │ │ * a1b2c3d Add amazing feature      │ │ 📝 Diff         │ │
│ │   my-proj   │ │ | Author: You <you@example.com>     │ │ ┌─────────────┐ │ │
│ │             │ │ | Date:   Mon Jul 27 10:30:00 2025  │ │ │ + import "fmt"  │ │
│ │ 🌿 Branches │ │ |                                    │ │ │               │ │
│ │ * main      │ │ * e4f5g6h Fix critical bug          │ │ │ func main() { │ │
│ │   feature/  │ │ | Author: You <you@example.com>     │ │ │ -  println    │ │
│ │   bugfix/   │ │ | Date:   Mon Jul 26 15:45:00 2025  │ │ │ +  fmt.Println│ │
│ │             │ │ |                                    │ │ │ }             │ │
│ │ 🔗 Remotes  │ │ * h7i8j9k Initial commit            │ │ └─────────────┘ │ │
│ │   origin    │ │                                       │                  │ │
│ │             │ │ Press Tab to switch panels           │ │                 │ │
│ │ ⚡ Actions  │ └────────────────────────────────────┘ └─────────────────┘ │
│ │ Stage All   │                                                             │
│ │ Commit      │ Tab/Shift+Tab: Switch panels  ↑/↓: Navigate  Enter: Select │
│ │ Push        │ r: Refresh  s: Stage all  c: Commit  p: Push  P: Pull      │
│ │ Pull        │ q: Quit  ?: Help                                              │
│ │ Fetch       │                                                             │
│ │ Stash       │                                                             │
│ └─────────────┘                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

- **Visual commit graph** — See branches, merges, and history at a glance
- **Interactive staging** — Stage/unstage hunks, not just files
- **Smart conflict resolution** — Side-by-side diff with inline actions
- **Keyboard-driven** — Vim-like bindings, fully customizable
- **Mouse support** — Click to select, scroll, resize panels
- **Live updates** — Watch status change as you edit files

### Web UI
- **Same features, in your browser** — Runs locally, no server needed
- **PWA support** — Install as app, works offline
- **Real-time collaboration** — Share a session with team (opt-in)
- **Mobile responsive** — Works on phone/tablet

### AI Assistant (Local-First)
```bash
gitforge ai commit       # Generate commit message from staged changes
gitforge ai explain      # Explain what a commit does
gitforge ai review       # Review diff for issues
gitforge ai suggest      # Suggest branch names, commit messages
```

- **Runs locally via WASM** — Your code never leaves your machine
- **Multiple models** — Phi-3, CodeLlama, DeepSeek Coder (auto-downloaded)
- **Context-aware** — Knows your repo, conventions, and history

### GitHub/GitLab Integration
```bash
gitforge pr create       # Create PR from current branch
gitforge pr list         # List open PRs
gitforge pr checkout 42  # Checkout PR #42 locally
gitforge issue create    # Create issue from TUI
```

### Plugin System
```bash
gitforge plugin install gitforge/jira      # Jira integration
gitforge plugin install gitforge/linear    # Linear integration
gitforge plugin install gitforge/slack     # Slack notifications
gitforge plugin install gitforge/ai-review # AI code review
```

- **WASM plugins** — Safe, sandboxed, language-agnostic
- **Go plugins** — Native performance for core extensions
- **Marketplace** — Discover, install, update with one command

---

## Installation

### Quick Install (Linux/macOS/WSL)
```bash
curl -fsSL https://gitforge.dev/install.sh | bash
```

### Package Managers
```bash
# Homebrew (macOS/Linux)
brew tap phantom1785322-hub/gitforge && brew install gitforge

# Scoop (Windows)
scoop bucket add gitforge https://github.com/phantom1785322-hub/scoop-gitforge
scoop install gitforge

# Chocolatey (Windows)
choco install gitforge

# Arch Linux (AUR)
yay -S gitforge

# Nix
nix profile install github:phantom1785322-hub/gitforge

# Termux (Android)
pkg install gitforge
# or
curl -fsSL https://gitforge.dev/install-termux.sh | bash
```

### From Source
```bash
git clone https://github.com/phantom1785322-hub/gitforge
cd gitforge
make build
# Binary at ./bin/gitforge
```

### Docker
```bash
docker run -it --rm -v $(pwd):/repo ghcr.io/phantom1785322-hub/gitforge:latest tui
```

### Verify Installation
```bash
gitforge version
gitforge doctor    # Run system diagnostics
```

---

## Usage

### Getting Started
```bash
# In any git repo
gitforge tui        # Launch terminal UI
gitforge web        # Launch web UI (opens browser)

# Outside a repo
gitforge init my-project
gitforge clone https://github.com/user/repo
```

### Core Commands
```bash
# Repository
gitforge init [path]           # Initialize new repo
gitforge clone <url> [path]    # Clone repository
gitforge status                # Show status (staged/unstaged/untracked)

# History
gitforge log                   # Beautiful commit log
gitforge log --graph           # ASCII graph
gitforge log --oneline         # Compact view
gitforge show <commit>         # Show commit details

# Branches
gitforge branch                # List branches
gitforge branch <name>         # Create branch
gitforge branch -d <name>      # Delete branch
gitforge checkout <branch>     # Switch branch
gitforge checkout -b <name>    # Create and switch

# Changes
gitforge add <files...>        # Stage files
gitforge add -p                # Interactive staging (hunks)
gitforge commit -m "msg"       # Commit
gitforge commit --amend        # Amend last commit
gitforge diff                  # Show unstaged changes
gitforge diff --cached         # Show staged changes

# Remotes
gitforge remote                # List remotes
gitforge remote add <name> <url>
gitforge fetch [remote]        # Fetch from remote
gitforge pull [remote] [branch]
gitforge push [remote] [branch]

# Stash
gitforge stash                 # Stash changes
gitforge stash list            # List stashes
gitforge stash pop             # Pop latest stash

# Tags
gitforge tag                   # List tags
gitforge tag <name> -m "msg"   # Create annotated tag

# Configuration
gitforge config                # Show config
gitforge config user.name "You"
gitforge config --global core.editor "nvim"

# AI Assistant
gitforge ai commit             # Generate commit message
gitforge ai explain <commit>   # Explain commit
gitforge ai review             # Review staged changes
gitforge ai suggest-branch     # Suggest branch name

# Plugins
gitforge plugin list           # List installed plugins
gitforge plugin install <name> # Install plugin
gitforge plugin remove <name>  # Remove plugin
gitforge plugin update         # Update all plugins

# Diagnostics
gitforge doctor                # System health check
gitforge version               # Version info
```

### TUI Keybindings
| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Switch panels |
| `↑` `↓` / `j` `k` | Navigate |
| `Enter` | Select / Open |
| `Space` | Toggle stage (in status) |
| `s` | Stage all |
| `c` | Commit |
| `p` | Push |
| `P` | Pull |
| `r` | Refresh |
| `f` | Fetch |
| `d` | Show diff |
| `l` | Show log |
| `b` | Branch menu |
| `?` | Help |
| `q` / `Ctrl+C` | Quit |

---

## Configuration

GitForge uses TOML config files with environment variable overrides.

### Config Locations (in order of precedence)
1. `./gitforge.toml` (repo-specific)
2. `~/.config/gitforge/config.toml` (user)
3. `/etc/gitforge/config.toml` (system)

### Example Config
```toml
# ~/.config/gitforge/config.toml

[core]
editor = "nvim"
pager = "delta"
auto_fetch = true
fetch_interval = 300  # seconds

[ui]
theme = "catppuccin-mocha"  # or "dracula", "nord", "tokyo-night", "auto"
show_line_numbers = true
mouse_enabled = true
compact_mode = false

[ai]
enabled = true
model = "phi-3-mini-4k-instruct-q4_k_m.gguf"
auto_download_models = true
max_tokens = 2048
temperature = 0.3

[github]
token = "${GITHUB_TOKEN}"  # From environment
default_reviewers = ["team-lead"]

[gitlab]
token = "${GITLAB_TOKEN}"
url = "https://gitlab.com"

[plugins]
auto_update = true
allow_unsigned = false
trusted_sources = ["gitforge", "my-org"]

[aliases]
co = "checkout"
br = "branch"
st = "status"
lg = "log --graph --oneline"
cm = "commit -m"
```

### Environment Variables
```bash
export GITFORGE_CONFIG="/path/to/config.toml"
export GITFORGE_THEME="dracula"
export GITFORGE_AI_MODEL="codellama-7b-instruct-q4_k_m.gguf"
export GITHUB_TOKEN="ghp_xxx"
export GITLAB_TOKEN="glpat_xxx"
```

---

## Themes

GitForge comes with beautiful built-in themes:

| Theme | Preview |
|-------|---------|
| `catppuccin-mocha` | ![Catppuccin](https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/palette/mocha.png) |
| `dracula` | ![Dracula](https://draculatheme.com/img/dracula.png) |
| `nord` | ![Nord](https://www.nordtheme.com/img/nord.png) |
| `tokyo-night` | ![Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme/raw/main/images/screenshot.png) |
| `gruvbox` | ![Gruvbox](https://github.com/morhetz/gruvbox/raw/master/gruvbox.png) |
| `everforest` | ![Everforest](https://github.com/sainnhe/everforest/raw/master/screenshots/everforest-dark.png) |
| `rose-pine` | ![Rose Pine](https://rosepinetheme.com/palette.png) |
| `auto` | Follows system/terminal theme |

Create custom themes in `~/.config/gitforge/themes/my-theme.toml`.

---

## AI Models

GitForge downloads models on first use (~2-4GB each).

| Model | Size | Best For | RAM Required |
|-------|------|----------|--------------|
| `phi-3-mini-4k-instruct-q4_k_m` | 2.3 GB | General, fast | 4 GB |
| `phi-3-medium-4k-instruct-q4_k_m` | 3.8 GB | Better quality | 6 GB |
| `codellama-7b-instruct-q4_k_m` | 3.9 GB | Code tasks | 6 GB |
| `codellama-13b-instruct-q4_k_m` | 7.2 GB | Complex code | 10 GB |
| `deepseek-coder-6.7b-instruct-q4_k_m` | 3.8 GB | Code generation | 6 GB |
| `deepseek-coder-33b-instruct-q4_k_m` | 18 GB | Best code quality | 24 GB |

```bash
# List available models
gitforge ai models

# Download specific model
gitforge ai download codellama-7b-instruct-q4_k_m

# Set default model
gitforge config ai.model codellama-7b-instruct-q4_k_m
```

---

## Plugins

### Official Plugins
| Plugin | Description |
|--------|-------------|
| `gitforge/github` | GitHub PRs, Issues, Actions |
| `gitforge/gitlab` | GitLab MRs, Issues, CI |
| `gitforge/jira` | Jira issue linking, transitions |
| `gitforge/linear` | Linear issue integration |
| `gitforge/slack` | Slack notifications |
| `gitforge/discord` | Discord webhook notifications |
| `gitforge/ai-review` | AI code review on commit/push |
| `gitforge/conventional` | Conventional commits enforcement |
| `gitforge/sign` | GPG/SSH commit signing |
| `gitforge/lfs` | Git LFS management |

### Installing Plugins
```bash
# From marketplace
gitforge plugin install gitforge/github

# From GitHub
gitforge plugin install github.com/myorg/gitforge-jira

# From local path
gitforge plugin install ./my-plugin

# List installed
gitforge plugin list

# Update all
gitforge plugin update
```

### Creating Plugins
See [Plugin Development Guide](docs/plugins/development.md).

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        GitForge App                              │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐              │
│  │    TUI      │  │    Web      │  │    CLI      │              │
│  │  (Bubble    │  │  (React +   │  │  (Cobra)    │              │
│  │   Tea)      │  │   Vite)     │  │             │              │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘              │
│         │                │                │                      │
│         └────────────────┼────────────────┘                      │
│                          ▼                                       │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                    Core Engine (Go)                      │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │   │
│  │  │  Git     │  │  Config  │  │  Plugin  │  │   AI       │  │   │
│  │  │  Ops     │  │  Manager │  │  Manager │  │  Engine    │  │   │
│  │  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                          │                                       │
│         ┌────────────────┼────────────────┐                     │
│         ▼                ▼                ▼                     │
│  ┌──────────┐      ┌──────────┐      ┌──────────┐              │
│  │  libgit2 │      │  SQLite  │      │  WASM    │              │
│  │  (CLI    │      │  (config,│      │  Runtime │              │
│  │   fallback)      │   cache) │      │  (plugins,│              │
│  └──────────┘      └──────────┘      │   AI)    │              │
│                                       └──────────┘              │
└─────────────────────────────────────────────────────────────────┘
```

### Tech Stack
- **Core**: Go 1.23+, single binary, zero runtime dependencies
- **TUI**: Bubble Tea + Lip Gloss (Charmbracelet)
- **Web**: React 18 + TypeScript + Vite + Tailwind CSS
- **Git**: Git CLI (primary) + libgit2 (optional, for performance)
- **AI**: llama.cpp via WASM (wasmer/wazero)
- **Plugins**: WASM (wasmer) + Go plugins
- **Config**: TOML + CUE validation
- **Database**: SQLite (embedded) + etcd (distributed sync)

---

## Cross-Platform Support

| Platform | TUI | Web UI | AI | Plugins | Status |
|----------|-----|--------|-----|---------|--------|
| Linux (x86_64) | ✅ | ✅ | ✅ | ✅ | ✅ Stable |
| Linux (ARM64) | ✅ | ✅ | ✅ | ✅ | ✅ Stable |
| macOS (Intel) | ✅ | ✅ | ✅ | ✅ | ✅ Stable |
| macOS (Apple Silicon) | ✅ | ✅ | ✅ | ✅ | ✅ Stable |
| Windows (x86_64) | ✅ | ✅ | ✅ | ✅ | ✅ Stable |
| Windows (ARM64) | ✅ | ✅ | ✅ | ✅ | ✅ Beta |
| Termux (Android ARM64) | ✅ | ⚠️ | ⚠️ | ⚠️ | 🧪 Experimental |
| Termux (Android ARMv7) | ✅ | ❌ | ❌ | ❌ | 🧪 Experimental |
| FreeBSD | ✅ | ✅ | ✅ | ✅ | ✅ Stable |
| OpenBSD | ✅ | ✅ | ⚠️ | ⚠️ | 🧪 Beta |
| WSL2 | ✅ | ✅ | ✅ | ✅ | ✅ Stable |

**Legend**: ✅ Full support • ⚠️ Partial • ❌ Not supported • 🧪 Experimental

---

## Performance

GitForge is built for speed:

| Operation | Time (typical repo) |
|-----------|---------------------|
| Startup (cold) | ~50ms |
| Startup (warm) | ~15ms |
| `gitforge status` | ~30ms |
| `gitforge log -20` | ~40ms |
| `gitforge diff` | ~50ms |
| TUI first frame | ~100ms |
| AI commit message | ~500ms-2s (local) |

**Optimizations:**
- Lazy loading — only loads what you view
- Incremental diff — reuses previous computations
- SIMD-accelerated parsing (AVX2/NEON/SVE)
- Object pooling — zero-allocation hot paths
- Smart caching — SQLite with intelligent invalidation

---

## Contributing

We ❤️ contributors! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Quick Start
```bash
git clone https://github.com/phantom1785322-hub/gitforge
cd gitforge
make install-tools
make dev          # Run TUI in dev mode
make dev-web      # Run web UI in dev mode
make test         # Run tests
make check        # Run all checks (fmt, vet, lint, test)
```

### Development Workflow
1. Fork the repo
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes with tests
4. Run `make check`
5. Submit PR

### Code Style
- Standard Go formatting (`gofmt`, `goimports`)
- Effective Go guidelines
- Conventional commits (`feat:`, `fix:`, `docs:`, etc.)
- 80%+ test coverage for new code

---

## Security

See [SECURITY.md](SECURITY.md) for our vulnerability disclosure policy.

- **No cloud required** — Local-first architecture
- **No telemetry** — Zero data collection
- **Supply chain security** — Verified dependencies, SBOM generation
- **Regular audits** — govulncheck, gosec, CodeQL

---

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). We follow the [Contributor Covenant](https://www.contributor-covenant.org/).

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Acknowledgments

GitForge stands on the shoulders of giants:

- **[Charmbracelet](https://charmbracelet.com/)** — Bubble Tea, Lip Gloss, Bubbles
- **[libgit2](https://libgit2.org/)** — Git library
- **[llama.cpp](https://github.com/ggerganov/llama.cpp)** — Local LLM inference
- **[wasmer](https://wasmer.io/)** — WASM runtime
- **[Tree-sitter](https://tree-sitter.github.io/)** — Incremental parsing
- **[Cobra](https://cobra.dev/)** & **[Kong](https://github.com/alecthomas/kong)** — CLI frameworks
- **[Viper](https://github.com/spf13/viper)** — Config management
- **[aio-stack](https://github.com/phantom1785322-hub/aio-stack)** — Shared foundation

---

## Support

- 📖 [Documentation](https://gitforge.dev/docs)
- 💬 [Discord Community](https://discord.gg/gitforge)
- 🐛 [Issue Tracker](https://github.com/phantom1785322-hub/gitforge/issues)
- 💡 [Feature Requests](https://github.com/phantom1785322-hub/gitforge/discussions/categories/ideas)
- 📧 [Email](mailto:hello@gitforge.dev)

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.

---

<div align="center">

**Made with ❤️ by the GitForge community**

[Website](https://gitforge.dev) • [Twitter](https://twitter.com/gitforge) • [Mastodon](https://fosstodon.org/@gitforge)

</div>