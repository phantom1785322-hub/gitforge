# Security Policy

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | ✅ Yes (current)   |
| < 0.1   | ❌ No              |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please report it responsibly.

### How to Report

**Preferred: GitHub Security Advisories (Private)**
1. Go to the [Security tab](https://github.com/phantom1785322-hub/gitforge/security)
2. Click "Report a vulnerability"
3. Fill out the form with details

**Alternative: Email**
- Email: security@gitforge.dev
- Subject: [SECURITY] Vulnerability Report
- Include: Description, reproduction steps, impact, suggested fix

### What to Include
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)
- Affected versions
- Your contact information (for follow-up)

### Response Timeline
- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: Within 30 days
  - High: Within 60 days
  - Medium: Within 90 days
  - Low: Next release cycle

### What We'll Do
1. Acknowledge receipt
2. Validate the vulnerability
3. Develop and test a fix
4. Coordinate disclosure timeline
5. Release patch with credit (if desired)

## Security Best Practices for Users

### Keep Updated
```bash
# Check for updates
gitforge version

# Update via package manager
brew upgrade gitforge        # macOS/Linux
scoop update gitforge        # Windows
# Or re-run installer
curl -fsSL https://gitforge.dev/install.sh | bash
```

### Secure Configuration
- Store tokens in environment variables, not config files
- Use `GITHUB_TOKEN`, `GITLAB_TOKEN` env vars
- Enable GPG signing for commits
- Review plugin permissions before installing

### Plugin Security
- Only install plugins from trusted sources
- Review plugin permissions before installing
- Prefer official plugins (`gitforge/`)
- Check `gitforge plugin list` regularly

## Security Features

### Built-in Protections
- **Local-first AI** — Your code never leaves your machine
- **WASM sandboxing** — Plugins run in isolated WebAssembly runtime
- **Capability-based security** — Plugins declare required permissions
- **Signed releases** — All binaries signed with cosign/Sigstore
- **Dependency scanning** — Automated via govulncheck and gosec

### Supply Chain Security
- **Go modules** with checksums (`go.sum`)
- **Reproducible builds** — Go 1.23+ reproducible builds
- **SBOM generation** — CycloneDX format on every release
- **Sigstore signing** — Cosign keyless signing

## Vulnerability Disclosure History

| Date | Version | CVE | Severity | Description |
|------|---------|-----|----------|-------------|
| — | — | — | — | No vulnerabilities reported yet |

## Security Contacts

- **Primary**: security@gitforge.dev
- **Maintainer**: @phantom1785322-hub
- **PGP Key**: Available on [GitHub profile](https://github.com/phantom1785322-hub.gpg)

## Bug Bounty

We don't have a formal bug bounty program yet, but we recognize security researchers in our:
- Release notes
- Hall of Fame (README)
- Annual contributor report

## Responsible Disclosure Guidelines

1. **Do not** publicly disclose before fix
2. **Do not** access/modify data not yours
3. **Do not** disrupt services
4. **Do** report promptly
5. **Do** allow reasonable time for fix

---

**Last Updated**: 2025-07-27
**Next Review**: 2025-10-27