# nprt - NixPkgs PR Tracker

A CLI tool to track which [nixpkgs](https://github.com/NixOS/nixpkgs) channels contain a given pull request.

## Features

- Check if a PR has been merged into various nixpkgs channels
- Support for PR numbers or full GitHub URLs
- Colored terminal output with Nerd Font icons
- Clickable hyperlinks to PRs (in supported terminals)
- JSON output for scripting
- Parallel channel checking for fast results

## Installation

### From Source

Requires Go 1.23 or later.

```bash
git clone https://github.com/taylrfnt/nixpkgs-pr-tracker.git
cd nixpkgs-pr-tracker
make build
```

The binary will be at `bin/nprt`.

## Usage

```bash
# Check by PR number
nprt 475593

# Check by PR URL
nprt https://github.com/NixOS/nixpkgs/pull/475593

# Check specific channels only
nprt --channels=master,nixos-unstable 475593

# JSON output for scripting
nprt --json 475593

# Force colors (useful for piping)
nprt --color=always 475593

# Verbose output for debugging
nprt --verbose 475593
```

### Example Output

```
 PR #475593

CHANNEL               STATUS
----------------------------
master                ✓
staging-next          ✓
nixpkgs-unstable      ✓
nixos-unstable-small  ✓
nixos-unstable        ✗
```

### Options

| Option | Description |
|--------|-------------|
| `--channels` | Comma-separated list of channels to check |
| `--color` | Color mode: `auto`, `always`, `never` (default: `auto`) |
| `--json` | Output results as JSON |
| `--verbose` | Show detailed progress and debug information |
| `--version` | Print version and exit |
| `-h, --help` | Show help message |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `GITHUB_TOKEN` | GitHub personal access token for higher API rate limits |
| `NO_COLOR` | Disable colors when set (respects [NO_COLOR](https://no-color.org/) standard) |

## Channels

By default, the following channels are checked:

- `master` - Main development branch
- `staging-next` - Staging integration branch
- `nixpkgs-unstable` - Unstable channel for non-NixOS users
- `nixos-unstable-small` - Fast-moving unstable channel with fewer packages
- `nixos-unstable` - Main unstable channel for NixOS

## PR Status Icons

The PR status line shows the current state of the pull request:

| Icon | Color | State |
|------|-------|-------|
|  | Gray | Draft |
|  | Green | Open |
|  | Purple | Merged |
|  | Red | Closed |

## Development

### Prerequisites

- Go 1.23+
- [gofumpt](https://github.com/mvdan/gofumpt) for formatting

### Commands

```bash
# Build the binary
make build

# Run tests
make test

# Check formatting
make format-check

# Format code
make format

# Bump version (patch/minor/major)
make version-bump TYPE=patch
make version-bump TYPE=minor
make version-bump TYPE=major

# Clean build artifacts
make clean
```

### Project Structure

```
cmd/
  nprt/
    main.go              # CLI entry point

internal/
  config/
    config.go            # Configuration and input parsing
  github/
    client.go            # GitHub API client
  core/
    core.go              # Domain logic for PR checking
  render/
    render.go            # Terminal and JSON output

tests/
    *_test.go            # Unit tests
```

### Architecture

The tool uses the GitHub REST API to:

1. Fetch PR metadata to get the merge commit SHA
2. Compare the merge commit against each channel branch
3. Determine if the commit is present based on the `behind_by` count

Channel checks run in parallel for faster results.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | PR/repo issue (not found, not merged, etc.) |
| 2 | CLI usage error (bad arguments) |
| 3 | Network/API error (rate limit, auth failure) |

## License

MIT
