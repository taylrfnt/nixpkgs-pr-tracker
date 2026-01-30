# Release Notes: v0.5.0

**Summary:** Issue warning with related PRs

## Overview

When users provide an issue number instead of a PR number, nprt now displays
a helpful warning with the issue details and discovers related pull requests
from the issue's timeline.

## Features

### Issue Warning with Related PRs

- **Issue Detection**: When the input is an issue (not a PR), displays a clear
  warning with issue title and state
- **Related PR Discovery**: Fetches the issue's timeline to find cross-referenced
  PRs in the same repository
- **Related PRs Table**: Displays related PRs with state icons, number, and title
- **Clickable Links**: Issue and related PR lines are hyperlinked in supported
  terminals

### New CLI Flag

- `--timeline-pages`: Controls how many pages of timeline events to fetch when
  looking for related PRs (default: 3)

### Visual Enhancements

- **Issue Icons**: New Nerd Font icons for issues:
  - `\uf41b` (nf-oct-issue_opened) for open issues
  - `\uf41d` (nf-oct-issue_closed) for closed issues
  - `\uf4e7` (nf-oct-issue_draft) for draft issues
- **Purple Color**: Closed issues display in purple (palette 13) to distinguish
  from closed PRs (which display in red)

## API Usage

- Uses GitHub's Timeline API (`GET /repos/{owner}/{repo}/issues/{number}/timeline`)
  to discover related PRs via `cross-referenced` events
- Timeline pagination is controlled by `--timeline-pages` flag
