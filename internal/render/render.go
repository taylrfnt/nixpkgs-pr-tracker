// Package render handles formatting and output of PR status results
// in table and JSON formats with optional ANSI colors.
package render

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/core"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorPurple = "\033[35m"
	colorGray   = "\033[90m"

	iconPresent    = "✓"
	iconNotPresent = "✗"
	iconUnknown    = "?"
	iconDraft      = "\uf4dd"
	iconOpen       = "\uf407"
	iconMerged     = "\uf407"
	iconClosed     = "\uf4dc"
)

// Renderer outputs PR status in various formats.
type Renderer struct {
	useColor      bool
	useHyperlinks bool
	writer        io.Writer
}

// NewRenderer creates a new Renderer with the given output settings.
func NewRenderer(writer io.Writer, useColor bool, useHyperlinks bool) *Renderer {
	return &Renderer{
		useColor:      useColor,
		useHyperlinks: useHyperlinks,
		writer:        writer,
	}
}

// RenderTable outputs the PR status as a formatted ASCII table.
func (r *Renderer) RenderTable(status *core.PRStatus) error {
	r.renderPRStatusLine(status)
	fmt.Fprintln(r.writer)

	maxNameLen := len("CHANNEL")
	for _, ch := range status.Channels {
		if len(ch.Name) > maxNameLen {
			maxNameLen = len(ch.Name)
		}
	}

	headerFmt := fmt.Sprintf("%%-%ds  STATUS\n", maxNameLen)
	fmt.Fprintf(r.writer, headerFmt, "CHANNEL")

	dividerLen := maxNameLen + 2 + 6
	fmt.Fprintln(r.writer, strings.Repeat("-", dividerLen))

	rowFmt := fmt.Sprintf("%%-%ds  %%s\n", maxNameLen)
	for _, ch := range status.Channels {
		icon := r.formatChannelStatus(ch.Status)
		fmt.Fprintf(r.writer, rowFmt, ch.Name, icon)
	}

	return nil
}

func (r *Renderer) renderPRStatusLine(status *core.PRStatus) {
	icon := r.formatPRState(status.State)
	text := fmt.Sprintf("PR #%d", status.Number)
	url := fmt.Sprintf("https://github.com/NixOS/nixpkgs/pull/%d", status.Number)

	displayText := text
	if r.useColor {
		displayText = colorBold + text + colorReset
	}

	if r.useHyperlinks {
		displayText = fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, displayText)
	}

	fmt.Fprintf(r.writer, "%s %s", icon, displayText)
}

func (r *Renderer) formatPRState(state core.PRState) string {
	switch state {
	case core.PRStateDraft:
		if r.useColor {
			return colorGray + iconDraft + colorReset
		}
		return iconDraft
	case core.PRStateOpen:
		if r.useColor {
			return colorGreen + iconOpen + colorReset
		}
		return iconOpen
	case core.PRStateMerged:
		if r.useColor {
			return colorPurple + iconMerged + colorReset
		}
		return iconMerged
	case core.PRStateClosed:
		if r.useColor {
			return colorRed + iconClosed + colorReset
		}
		return iconClosed
	default:
		return iconOpen
	}
}

func (r *Renderer) formatChannelStatus(status core.ChannelStatus) string {
	switch status {
	case core.StatusPresent:
		if r.useColor {
			return colorGreen + iconPresent + colorReset
		}
		return iconPresent
	case core.StatusNotPresent:
		if r.useColor {
			return colorRed + iconNotPresent + colorReset
		}
		return iconNotPresent
	default:
		if r.useColor {
			return colorYellow + iconUnknown + colorReset
		}
		return iconUnknown
	}
}

// RenderJSON outputs the PR status as pretty-printed JSON.
func (r *Renderer) RenderJSON(status *core.PRStatus) error {
	output := struct {
		PR          int                  `json:"pr"`
		State       core.PRState         `json:"state"`
		MergeCommit string               `json:"merge_commit,omitempty"`
		Channels    []core.ChannelResult `json:"channels"`
	}{
		PR:          status.Number,
		State:       status.State,
		MergeCommit: status.MergeCommit,
		Channels:    status.Channels,
	}

	encoder := json.NewEncoder(r.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
