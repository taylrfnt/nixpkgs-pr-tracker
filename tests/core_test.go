package tests

import (
	"testing"

	"github.com/taylrfnt/nixpkgs-pr-tracker/internal/core"
)

func TestSortChannelResults_PresentFirst(t *testing.T) {
	results := []core.ChannelResult{
		{Name: "alpha", Status: core.StatusNotPresent},
		{Name: "beta", Status: core.StatusPresent},
		{Name: "gamma", Status: core.StatusNotPresent},
		{Name: "delta", Status: core.StatusPresent},
	}

	sorted := core.SortChannelResults(results)

	if sorted[0].Status != core.StatusPresent || sorted[1].Status != core.StatusPresent {
		t.Error("Present channels should come first")
	}

	if sorted[0].Name != "beta" || sorted[1].Name != "delta" {
		t.Errorf("Present channels should maintain relative order: got %s, %s", sorted[0].Name, sorted[1].Name)
	}
}

func TestSortChannelResults_NotPresentAlphabetical(t *testing.T) {
	results := []core.ChannelResult{
		{Name: "zebra", Status: core.StatusNotPresent},
		{Name: "alpha", Status: core.StatusNotPresent},
		{Name: "master", Status: core.StatusPresent},
		{Name: "beta", Status: core.StatusNotPresent},
	}

	sorted := core.SortChannelResults(results)

	if sorted[0].Name != "master" {
		t.Errorf("First should be master (present), got %s", sorted[0].Name)
	}
	if sorted[1].Name != "alpha" {
		t.Errorf("Second should be alpha, got %s", sorted[1].Name)
	}
	if sorted[2].Name != "beta" {
		t.Errorf("Third should be beta, got %s", sorted[2].Name)
	}
	if sorted[3].Name != "zebra" {
		t.Errorf("Fourth should be zebra, got %s", sorted[3].Name)
	}
}

func TestSortChannelResults_UnknownWithNotPresent(t *testing.T) {
	results := []core.ChannelResult{
		{Name: "unknown-channel", Status: core.StatusUnknown},
		{Name: "master", Status: core.StatusPresent},
		{Name: "alpha", Status: core.StatusNotPresent},
	}

	sorted := core.SortChannelResults(results)

	if sorted[0].Name != "master" {
		t.Errorf("First should be master (present), got %s", sorted[0].Name)
	}
	if sorted[1].Name != "alpha" {
		t.Errorf("Second should be alpha (alphabetically first among non-present), got %s", sorted[1].Name)
	}
	if sorted[2].Name != "unknown-channel" {
		t.Errorf("Third should be unknown-channel, got %s", sorted[2].Name)
	}
}

func TestSortChannelResults_Empty(t *testing.T) {
	results := []core.ChannelResult{}
	sorted := core.SortChannelResults(results)
	if len(sorted) != 0 {
		t.Error("Empty input should return empty output")
	}
}

func TestSortChannelResults_AllPresent(t *testing.T) {
	results := []core.ChannelResult{
		{Name: "gamma", Status: core.StatusPresent},
		{Name: "alpha", Status: core.StatusPresent},
		{Name: "beta", Status: core.StatusPresent},
	}

	sorted := core.SortChannelResults(results)

	if sorted[0].Name != "gamma" || sorted[1].Name != "alpha" || sorted[2].Name != "beta" {
		t.Error("All present channels should maintain original order")
	}
}
