package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/config"
	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/core"
	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/github"
	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/render"
)

var version = "dev"

const usage = `Usage: nprt [options] <PR number | PR URL>

Track which nixpkgs channels contain a given pull request.

Arguments:
  PR number    A pull request number (e.g., 476497)
  PR URL       A full GitHub PR URL (e.g., https://github.com/NixOS/nixpkgs/pull/476497)

Options:
  --channels   Comma-separated list of channels to check (default: master,staging-next,nixpkgs-unstable,nixos-unstable-small,nixos-unstable)
  --color      Color output mode: auto, always, never (default: auto)
  --json       Output results as JSON
  --verbose    Show detailed progress and debug information
  --version    Print version and exit
  -h, --help   Show this help message

Environment:
  GITHUB_TOKEN  GitHub personal access token for higher rate limits
`

func main() {
	os.Exit(run())
}

func run() int {
	var (
		channelsFlag string
		colorMode    string
		jsonOutput   bool
		verbose      bool
		showVersion  bool
	)

	flag.StringVar(&channelsFlag, "channels", "", "Comma-separated list of channels to check")
	flag.StringVar(&colorMode, "color", "auto", "Color output: auto, always, never")
	flag.BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	flag.BoolVar(&verbose, "verbose", false, "Show detailed progress and debug information")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("nprt version %s\n", version)
		return 0
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprint(os.Stderr, usage)
		return 2
	}

	prNumber, err := config.ParsePRInput(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 2
	}

	channels, err := config.ParseChannels(channelsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 2
	}

	token := config.GetGitHubToken()
	useColor := config.ShouldUseColor(colorMode)
	useHyperlinks := config.IsTerminal()

	if verbose {
		fmt.Fprintf(os.Stderr, "Fetching PR #%d from NixOS/nixpkgs...\n", prNumber)
	}

	client := github.NewClient(token, verbose)
	checker := core.NewChecker(client, verbose)

	// Set up context with signal handling for clean cancellation
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	status, err := checker.CheckPR(ctx, prNumber, channels)
	if err != nil {
		var apiErr *github.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == 403 {
			fmt.Fprintf(os.Stderr, "Error: %s\n", apiErr.Message)
			return 3
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	renderer := render.NewRenderer(os.Stdout, useColor, useHyperlinks)

	if jsonOutput {
		if err := renderer.RenderJSON(status); err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering output: %s\n", err)
			return 1
		}
	} else {
		if err := renderer.RenderTable(status); err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering output: %s\n", err)
			return 1
		}
	}

	return 0
}
