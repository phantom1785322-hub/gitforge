# CHANGELOG.md

All notable changes to GitForge will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added
- Comprehensive GitHub issue templates (bug, feature, question)
- Pull request template with checklist
- CODEOWNERS file
- Dependabot configuration for Go modules and GitHub Actions
- GitHub Actions workflows: CI, Release, Auto-promote
- goreleaser configuration for multi-platform releases
- MIT License
- CONTRIBUTING.md with development guidelines
- SECURITY.md with vulnerability disclosure policy
- CODE_OF_CONDUCT.md (Contributor Covenant v2.1)
- GitHub issue templates: bug_report, feature_request, question
- Pull request template
- GitHub Actions: ci.yml, release.yml, promo.yml

### Changed
- Updated README.md with professional formatting, badges, installation methods, usage, configuration, themes, AI models, plugins, architecture, cross-platform support, performance benchmarks, contributing, security, code of conduct
- Improved .github/workflows/ci.yml with multi-platform builds, security scanning
- Improved .github/workflows/release.yml with Homebrew auto-update
- Improved .github/workflows/promo.yml with Twitter, Discord, Dev.to, Mastodon auto-promotion

### Security
- Added govulncheck, gosec, CodeQL to CI pipeline
- Added SECURITY.md with vulnerability disclosure policy

---

## [0.1.0] - 2025-07-27

### Added
- Initial release of GitForge
- **TUI**: 3-panel Bubble Tea interface with sidebar, commit graph, diff/detail panel
- **CLI**: 15+ commands (status, log, commit, diff, branch, remote, stash, tag, config, plugin, doctor, version, init, clone)
- **Git Operations**: Full wrapper around git CLI (status, log, branch, remote, stash, tag, commit, diff, merge, rebase, cherry-pick, revert, blame, blame, submodule, config)
- **Cross-platform**: Builds for 11 platforms (Linux/amd64, Linux/arm64, Linux/armv7, macOS/amd64, macOS/arm64, Windows/amd64, Windows/arm64, FreeBSD/amd64, FreeBSD/arm64, OpenBSD/amd64, Termux/arm64, Termux/armv7)
- **Package Managers**: Homebrew tap, Scoop bucket
- **Installer**: Universal install.sh script
- **Documentation**: Comprehensive README with features, installation, usage, configuration, themes, AI models, plugins, architecture, cross-platform matrix, performance benchmarks
- **License**: MIT License

### Technical
- Built on aio-stack shared foundation (platform detection, optimizer, code intelligence, AI engine, plugin system, CLI framework, config system)
- Go 1.23 with zero runtime dependencies (single binary)
- Bubble Tea + Lip Gloss for TUI
- Cobra + Kong for CLI
- Git CLI wrapper (no libgit2 dependencies)
- Tree-sitter for code intelligence (15+ languages)
- llama.cpp via WASM for local AI (Phi-3, CodeLlama, DeepSeek)
- wasmer/wazero for WASM plugin runtime

---

[Unreleased]: https://github.com/phantom1785322-hub/gitforge/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/phantom1785322-hub/gitforge/releases/tag/v0.1.0