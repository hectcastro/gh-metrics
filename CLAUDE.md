# CLAUDE.md

This document provides context about the `gh-metrics` project for Claude AI assistant.

## Project Overview

`gh-metrics` is a GitHub CLI (`gh`) extension written in Go that provides summary pull request metrics. It helps developers and teams analyze PR performance by calculating various metrics including time to first review, feature lead time, review durations, and more.

**Key Features:**
- Calculates PR metrics: commits, additions, deletions, changed files, time to first review, comments, participants, feature lead time, etc.
- Supports flexible date ranges and query filters
- Outputs data in table or CSV format
- Uses GitHub GraphQL API for efficient data retrieval

## Project Structure

```
gh-metrics/
├── main.go              # Entry point - delegates to cmd package
├── cmd/                 # Command implementation and business logic
│   ├── root.go          # Main CLI command implementation
│   ├── graphql.go       # GitHub GraphQL API interactions
│   ├── ghrepo.go        # Repository handling utilities
│   ├── ui.go            # Output formatting (table/CSV)
│   └── *_test.go        # Test files
├── scripts/             # Build and validation scripts
│   └── gofmtcheck       # Go formatting validation
├── Makefile             # Build automation
├── .goreleaser.yml      # Release configuration
├── go.mod               # Go module dependencies
└── README.md            # User-facing documentation
```

## Key Files

### main.go
Simple entry point that delegates to `cmd.Execute()`.

### cmd/root.go
Contains the main command logic, including:
- CLI flag parsing (--repo, --start, --end, --query, --csv)
- PR data fetching orchestration
- Metric calculations
- Output rendering

### cmd/graphql.go
Handles GitHub GraphQL API interactions:
- PR queries with filters
- Commit history retrieval
- Review timeline data
- Rate limit handling

### cmd/ui.go
Output formatting:
- Table rendering using go-pretty
- CSV generation
- Duration formatting

### cmd/ghrepo.go
Repository parsing and validation utilities.

## Development Workflow

### Prerequisites
- Go 1.25+ installed
- GitHub CLI (`gh`) installed and authenticated
- Access to GitHub repositories for testing

### Testing
```bash
# Run all tests with format checking
make test

# Run tests with coverage
make cover

# Format Go code
make fmt

# Check formatting
make fmtcheck
```

### Building
```bash
# Build for local testing
go build -o gh-metrics

# Install as gh extension (from source)
gh extension install .
```

### Running Locally
```bash
# After building
./gh-metrics --repo owner/repo

# Or as installed extension
gh metrics --repo owner/repo
```

## Key Dependencies

- **github.com/cli/go-gh**: Official GitHub CLI library for API access
- **github.com/cli/shurcooL-graphql**: GraphQL client
- **github.com/jedib0t/go-pretty/v6**: Table formatting
- **github.com/spf13/cobra**: CLI framework
- **github.com/spf13/pflag**: Flag parsing
- **github.com/rickar/cal/v2**: Calendar/business day calculations

## Metric Definitions

Understanding these metrics is crucial when modifying calculations:

- **Time to first review**: Duration from PR creation (or marked "Ready for review") to first review completion
- **Feature lead time**: Duration from first commit creation to PR merge
- **First review to last review**: Duration between first non-author review and last approving non-author review
- **First approval to merge**: Duration from first approval review to PR merge

## Coding Conventions

1. **Formatting**: Use `gofmt` for all Go code (enforced by `make fmtcheck`)
2. **Testing**: Write tests for new functionality in `*_test.go` files
3. **Error Handling**: Return errors explicitly; handle them at appropriate levels
4. **GraphQL**: Use the shurcooL GraphQL client for type-safe queries
5. **CLI**: Follow cobra/pflag patterns for command structure

## Common Tasks

### Adding a New Metric
1. Update GraphQL query in `cmd/graphql.go` to fetch required data
2. Add calculation logic in `cmd/root.go`
3. Update output formatters in `cmd/ui.go` (both table and CSV)
4. Add tests in `cmd/root_test.go`
5. Update README.md with metric definition

### Modifying GraphQL Queries
- Queries are in `cmd/graphql.go`
- Use GitHub GraphQL Explorer to test: https://docs.github.com/en/graphql/overview/explorer
- Be mindful of rate limits and pagination

### Changing Output Format
- Table formatting: `cmd/ui.go` using go-pretty library
- CSV formatting: Also in `cmd/ui.go`
- Ensure both formats stay in sync

### Updating Dependencies
```bash
go get -u <package>
go mod tidy
```

## Release Process

Releases are managed through GoReleaser:
1. Tag a version: `git tag vX.Y.Z`
2. Push tag: `git push origin vX.Y.Z`
3. GoReleaser builds and publishes release artifacts

Configuration: `.goreleaser.yml`

## Troubleshooting

### Common Issues
- **Authentication errors**: Ensure `gh auth status` shows valid authentication
- **GraphQL rate limits**: Use smaller date ranges or add delays between requests
- **Test failures**: Run `make fmtcheck` to ensure code is properly formatted

## Additional Context

- Inspired by: [jmartin82/mkpis](https://github.com/jmartin82/mkpis)
- GitHub CLI extensions: https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions
- GraphQL API docs: https://docs.github.com/en/graphql
