// Package render handles formatting and output of PR status, issue warnings,
// and related PR lists in table and JSON formats with optional ANSI colors.
package render

import (
	"encoding/json"
	"io"
	"os"

	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/core"
)

// Renderer outputs PR status in various formats.
type Renderer struct {
	useColor      bool
	useHyperlinks bool
	useNerdFonts  bool
	writer        io.Writer
}

// NewRenderer creates a new Renderer with the given output settings.
// Nerd Font icons are enabled by default; set NO_NERD_FONTS=1 to disable.
func NewRenderer(writer io.Writer, useColor bool, useHyperlinks bool) *Renderer {
	return &Renderer{
		useColor:      useColor,
		useHyperlinks: useHyperlinks,
		useNerdFonts:  os.Getenv("NO_NERD_FONTS") == "",
		writer:        writer,
	}
}

// RenderJSON outputs the PR status as pretty-printed JSON.
func (r *Renderer) RenderJSON(status *core.PRStatus) error {
	encoder := json.NewEncoder(r.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(status)
}
