# Contributing to LinkBeam

Thank you for your interest in contributing to LinkBeam! This document provides guidelines and instructions for development.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)
- golangci-lint (for code quality checks)
- pre-commit (optional, for automated checks)

### Getting Started

1. Fork and clone the repository:

```bash
git clone https://github.com/5mdt/linkbeam.git
cd linkbeam
```

2. Install dependencies:

```bash
go mod download
```

3. Build the project:

```bash
make build
# or: go build -o dist/linkbeam cmd/linkbeam/main.go
```

## Pre-commit Hooks

This project uses pre-commit hooks to ensure code quality. Setting them up is highly recommended.

### Installation

```bash
# Install pre-commit (if not already installed)
pip install pre-commit
# or: brew install pre-commit  (on macOS)

# Install golangci-lint (if not already installed)
# See: https://golangci-lint.run/usage/install/
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install the git hook scripts
pre-commit install
# or: make install-hooks
```

### Running Hooks

```bash
# Run against all files
pre-commit run --all-files
# or: make pre-commit

# Hooks will automatically run on git commit
git commit -m "your message"
```

The pre-commit hooks check:
- Go code formatting (`go fmt`)
- Import organization (`goimports`)
- Code quality (`go vet`, `golangci-lint`)
- Tests (`go test`)
- Module tidiness (`go mod tidy`)
- YAML syntax
- Trailing whitespace
- End of file fixes
- And more

**Note**: The pre-commit framework will automatically download required Go tools on first run if they're not already installed.

## Development Commands

### Building

```bash
# Build binary (outputs to dist/linkbeam)
make build
# or: go build -o dist/linkbeam cmd/linkbeam/main.go
```

### Running

```bash
# Run with example config
make run-example
# or: go run cmd/linkbeam/main.go --config internal/config/testdata/config.yaml

# Run with custom config (long flags)
./linkbeam --config config.yaml --template templates/base.html --output dist/index.html

# Run with custom config (short flags)
./linkbeam -c config.yaml -t templates/base.html -o dist/index.html
```

### Testing

```bash
# Run all tests
make test
# or: go test -v ./...

# Run tests for a specific package
go test ./internal/config
go test ./internal/generator

# Run a single test function
go test -v -run TestLoad ./internal/config

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

```bash
# Format code
go fmt ./...
# or: make fmt

# Run go vet
go vet ./...

# Run golangci-lint
golangci-lint run
# or: make vet

# Run all quality checks
make pre-commit
```

### Cleaning

```bash
# Clean build artifacts
make clean
```

## Project Structure

```
linkbeam/
├── cmd/linkbeam/          # Main application entry point
│   ├── main.go
│   └── main_test.go
├── internal/
│   ├── config/            # Configuration loading and validation
│   │   ├── config.go
│   │   ├── config_test.go
│   │   └── testdata/
│   └── generator/         # Site generation and rendering
│       ├── generator.go
│       └── generator_test.go
├── templates/             # HTML templates
│   └── base.html
├── themes/                # CSS theme files
│   ├── light.css
│   └── dark.css
├── dist/                  # Build output: binary and generated site (gitignored)
├── .gitea/                # CI/CD workflows
│   └── workflows/
├── Makefile              # Build automation
├── go.mod                # Go module definition
└── go.sum                # Go dependencies
```

## Architecture

### Two-Layer Structure

1. **Config Layer** (`internal/config/`):
   - Loads and validates YAML configuration files
   - Defines data structures: `Config`, `Link`, `Social`
   - Validates theme (must be "light" or "dark") and required fields
   - Config validation happens automatically during `config.Load()`

2. **Generator Layer** (`internal/generator/`):
   - **HTML Generation**: Orchestrates template rendering using Go's `html/template`
     - `GenerateSite(cfg, templatePath, outPath)` - Creates HTML from templates
   - **Asset Management**: Copies CSS themes and static files to output directory
     - `CopyAssets(distDir, themeDirs...)` - Copies theme files to dist/themes/
   - **Text Rendering**: Plain-text representation for debugging
     - `RenderUserPage(cfg)` - Returns text-based preview of the page

### Testing Strategy

- Tests use table-driven patterns and dependency injection
- `cmd/linkbeam/main_test.go`: Mocks the config loader and captures stdout
- Config tests validate YAML parsing and business rules
- Generator tests verify file creation and template execution
- All tests use the testdata in `internal/config/testdata/config.yaml`

## Making Changes

### Adding New Features

1. Create a new branch for your feature:
```bash
git checkout -b feature/your-feature-name
```

2. Write tests first (TDD approach recommended)
3. Implement the feature
4. Ensure all tests pass: `go test ./...`
5. Run pre-commit hooks: `pre-commit run --all-files`
6. Commit your changes with a descriptive message

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (handled by pre-commit)
- Write clear, descriptive variable and function names
- Add comments for exported functions and types
- Keep functions focused and concise

### Testing Guidelines

- Write tests for all new functionality
- Aim for high test coverage
- Use table-driven tests for multiple test cases
- Use testdata files for fixtures
- Test both success and error cases

### Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in present tense (e.g., "Add", "Fix", "Update")
- Reference issue numbers when applicable
- Keep the first line under 72 characters

Examples:
```
Add support for custom CSS themes
Fix avatar path resolution in config loader
Update template to include meta tags
```

## CI/CD Pipeline

The project uses Gitea Actions for continuous integration and deployment.

### Workflows

1. **CI Pipeline** (`.gitea/workflows/ci.yml`):
   - Runs on push to main/develop and pull requests
   - **Lint**: golangci-lint, go fmt, go vet
   - **Test**: Run tests with race detection and coverage
   - **Build**: Multi-platform builds (Linux, Windows, macOS) for multiple architectures
   - **Docker**: Build multi-arch Docker images

2. **Docker Publish** (`.gitea/workflows/docker-publish.yml`):
   - Triggered on version tags (v*)
   - Builds and pushes Docker images to container registry
   - Multi-architecture: linux/amd64, linux/arm64, linux/arm/v7

3. **Release** (`.gitea/workflows/release.yml`):
   - Triggered on semver tags (v*.*.*)
   - Creates releases with binaries for all platforms
   - Generates SHA256 checksums for verification

### Local Testing with Docker

```bash
# Test Docker build
docker build -t linkbeam:test .

# Run Docker container
docker run -v $(pwd)/config.yaml:/app/config/config.yaml \
           -v $(pwd)/dist:/app/dist \
           linkbeam:test

# Test with docker-compose
docker-compose up
```

## Submitting Changes

1. Ensure your code passes all tests and checks
2. Push your branch to your fork
3. Open a pull request with a clear description of your changes
4. Wait for CI checks to pass
5. Respond to any review feedback

## Questions or Issues?

If you have questions or run into issues:
- Check existing issues on the repository
- Open a new issue with a clear description
- Join our community discussions (if available)

## License

By contributing to LinkBeam, you agree that your contributions will be licensed under the GNU General Public License v3.0.
