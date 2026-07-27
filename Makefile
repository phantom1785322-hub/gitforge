# GitForge Makefile
# Build, test, and release automation

.DEFAULT_GOAL := build

# Variables
BINARY_NAME := gitforge
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -s -w"

# Go settings
GO := go
GOFLAGS := -trimpath
CGO_ENABLED := 0

# Directories
BUILD_DIR := ./bin
DIST_DIR := ./dist
CMD_DIR := ./cmd/gitforge

# Platforms for release
PLATFORMS := \
	linux_amd64 \
	linux_arm64 \
	darwin_amd64 \
	darwin_arm64 \
	windows_amd64 \
	windows_arm64 \
	freebsd_amd64 \
	freebsd_arm64 \
	openbsd_amd64 \
	linux_armv7

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build for current platform
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-all
build-all: ## Build for all platforms
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d_ -f1); \
		ARCH=$$(echo $$platform | cut -d_ -f2); \
		EXT=""; \
		if [ "$$OS" = "windows" ]; then EXT=".exe"; fi; \
		echo "Building for $$OS/$$ARCH..."; \
		GOOS=$$OS GOARCH=$$ARCH CGO_ENABLED=$(CGO_ENABLED) \
			$(GO) build $(GOFLAGS) $(LDFLAGS) \
			-o $(DIST_DIR)/$(BINARY_NAME)_$$OS_$$ARCH$$EXT $(CMD_DIR); \
	done

.PHONY: test
test: ## Run tests
	$(GO) test -v -race -count=1 ./...

.PHONY: test-short
test-short: ## Run tests (short)
	$(GO) test -short -count=1 ./...

.PHONY: lint
lint: ## Run linter
	@which golangci-lint >/dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	$(GO) fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	$(GO) mod tidy

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR) $(DIST_DIR)

.PHONY: install
install: build ## Install to GOPATH/bin
	$(GO) install $(GOFLAGS) $(LDFLAGS) $(CMD_DIR)

.PHONY: dev
dev: ## Run in development mode
	$(GO) run $(GOFLAGS) $(LDFLAGS) $(CMD_DIR) tui

.PHONY: dev-web
dev-web: ## Run web UI in development mode
	$(GO) run $(GOFLAGS) $(LDFLAGS) $(CMD_DIR) web

.PHONY: release
release: clean test lint build-all ## Create release artifacts
	@mkdir -p $(DIST_DIR)/checksums
	@cd $(DIST_DIR) && for f in $(BINARY_NAME)_*; do \
		if [ -f "$$f" ]; then \
			sha256sum "$$f" > checksums/$$f.sha256; \
		fi; \
	done
	@echo "Release artifacts in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

.PHONY: release-dry
release-dry: ## Dry-run release (no publish)
	@echo "Would release version $(VERSION)"
	@echo "Platforms: $(PLATFORMS)"

.PHONY: install-tools
install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/goreleaser/goreleaser@latest

.PHONY: check
check: fmt vet lint test-short ## Run all checks

.PHONY: benchmark
benchmark: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

.PHONY: coverage
coverage: ## Generate coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: docs
docs: ## Generate documentation
	@mkdir -p docs
	$(GO) run $(CMD_DIR) --help > docs/cli-reference.md
	@echo "Generated docs/cli-reference.md"

.PHONY: web-build
web-build: ## Build web UI
	@echo "Building web UI..."
	@cd web && npm ci && npm run build
	@echo "Web UI built to web/dist/"

.PHONY: web-dev
web-dev: ## Run web UI dev server
	@cd web && npm run dev

# Cross-compilation for Termux (Android)
.PHONY: build-termux
build-termux: ## Build for Termux (Android ARM64)
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
		$(GO) build $(GOFLAGS) $(LDFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)_termux_arm64 $(CMD_DIR)

.PHONY: build-termux-armv7
build-termux-armv7: ## Build for Termux (Android ARMv7)
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
		$(GO) build $(GOFLAGS) $(LDFLAGS) \
		-o $(DIST_DIR)/$(BINARY_NAME)_termux_armv7 $(CMD_DIR)

# Generate shell completions
.PHONY: gen-completions
gen-completions: build ## Generate shell completions
	@mkdir -p $(BUILD_DIR)/completions
	$(BUILD_DIR)/$(BINARY_NAME) completion bash > $(BUILD_DIR)/completions/gitforge.bash
	$(BUILD_DIR)/$(BINARY_NAME) completion zsh > $(BUILD_DIR)/completions/_gitforge
	$(BUILD_DIR)/$(BINARY_NAME) completion fish > $(BUILD_DIR)/completions/gitforge.fish
	$(BUILD_DIR)/$(BINARY_NAME) completion powershell > $(BUILD_DIR)/completions/gitforge.ps1
	@echo "Completions generated in $(BUILD_DIR)/completions/"

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t gitforge:$(VERSION) -t gitforge:latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -it --rm -v $(PWD):/repo gitforge:latest tui

# Dependency management
.PHONY: deps
deps: ## Download dependencies
	$(GO) mod download

.PHONY: deps-graph
deps-graph: ## Show dependency graph
	$(GO) mod graph | head -50

# Security
.PHONY: audit
audit: ## Run security audit
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Profiling
.PHONY: profile-cpu
profile-cpu: ## CPU profile
	$(GO) test -cpuprofile=cpu.prof -bench=. ./...
	$(GO) tool pprof cpu.prof

.PHONY: profile-mem
profile-mem: ## Memory profile
	$(GO) test -memprofile=mem.prof -bench=. ./...
	$(GO) tool pprof mem.prof

# Default target
all: check build