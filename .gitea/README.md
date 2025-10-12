# Gitea CI/CD Configuration

This directory contains Gitea Actions workflows for automated testing, building, and deployment.

## Workflows

### `ci.yml` - Main CI/CD Pipeline

Runs on every push and pull request:

1. **Lint Stage**
   - golangci-lint with all enabled linters
   - go fmt verification
   - go vet static analysis

2. **Test Stage**
   - Run all tests with `-race` flag
   - Generate coverage reports
   - Upload coverage artifacts

3. **Build Stage**
   - Matrix build for multiple platforms:
     - Linux: amd64, 386, arm64, armv7
     - Windows: amd64, 386, arm64
     - macOS: amd64, arm64
   - Build with version and commit info in ldflags
   - Upload build artifacts

4. **Docker Stage**
   - Multi-architecture Docker build
   - Cache optimization with GitHub Actions cache
   - Build for: linux/amd64, linux/arm64, linux/arm/v7

### `docker-publish.yml` - Docker Registry Publishing

Triggered on version tags (`v*`):
- Builds multi-arch Docker images
- Pushes to container registry
- Tags with version numbers and `latest`
- Uses registry caching for faster builds

### `release.yml` - Release Creation

Triggered on semver tags (`v*.*.*`):
- Creates GitHub/Gitea release
- Uploads platform-specific binaries (tar.gz for Unix, zip for Windows)
- Generates SHA256 checksums
- Includes release notes

## Configuration Required

### Secrets

Set these in your Gitea repository settings:

- `GITEA_TOKEN` - Gitea API token with package write permissions
- `GITHUB_TOKEN` - Automatically provided by Gitea Actions

### Registry Configuration

Update the registry URL in `docker-publish.yml`:
```yaml
env:
  REGISTRY: gitea.example.com  # Change to your Gitea instance
  IMAGE_NAME: ${{ github.repository }}
```

## Local Testing

### Test builds locally:

```bash
# Test Go build
make build

# Test multi-platform build
GOOS=linux GOARCH=arm64 go build -o linkbeam-linux-arm64 cmd/linkbeam/main.go

# Test Docker build (requires Docker)
docker build -t linkbeam:test .

# Test Docker with compose
docker-compose up
```

## Triggering Releases

To create a new release:

```bash
# Tag with semantic version
git tag v1.0.0
git push origin v1.0.0

# This will trigger:
# 1. Full CI pipeline
# 2. Docker image build and push
# 3. Release creation with binaries
```

## Build Artifacts

Artifacts are uploaded for each workflow run:
- Coverage reports (test stage)
- Platform binaries (build stage)
- Docker images (docker stage)
- Release assets (release stage)

Artifacts are available for 90 days after the workflow run.
