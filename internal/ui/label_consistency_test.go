package ui

import (
	"regexp"
	"strings"
	"testing"

	"github.com/aemonge/dossier/internal/openspec"
)

var labelStripANSI = regexp.MustCompile(`\x1b\[[0-9;:]*m`)

func plainIndexLines(t *testing.T, m *Model) []string {
	t.Helper()
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		lines = append(lines, labelStripANSI.ReplaceAllString(line, ""))
	}
	return lines
}

func lineContaining(lines []string, needle string) string {
	for _, line := range lines {
		if strings.Contains(line, needle) {
			return line
		}
	}
	return ""
}

// The label column starts after the 4-space cursor mark on non-cursor rows.
const artifactLabelColumn = 4

func TestRenderArtifactRowsKeepStableIDs(t *testing.T) {
	change := openspec.Change{
		Name:   "honest-debug",
		Schema: "spec",
		Artifacts: []openspec.ChangeArtifact{
			{
				ID: "idea", Status: openspec.ArtifactStatusDone, Source: openspec.ArtifactSourceStatus,
				Outputs: []openspec.ArtifactOutput{{RelativePath: "1-idea.md", DisplayName: "1-idea", Present: true}},
			},
			{ID: "requirements", Status: openspec.ArtifactStatusReady, Source: openspec.ArtifactSourceStatus},
			{
				ID: "plan", Status: openspec.ArtifactStatusDone, Source: openspec.ArtifactSourceStatus,
				Outputs: []openspec.ArtifactOutput{{RelativePath: "5-plan.md", DisplayName: "5-plan", Present: true}},
			},
			{ID: "tasks", Status: openspec.ArtifactStatusBlocked, Source: openspec.ArtifactSourceStatus},
		},
	}
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{change}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:honest-debug": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		theme: Theme{Styles: BuildStyles(DarkColors)},
		width: 100,
	}
	lines := plainIndexLines(t, m)

	idColumnEnd := artifactLabelColumn + len("requirements") + 1

	ideaRow := lineContaining(lines, "1-idea")
	if ideaRow == "" {
		t.Fatalf("expected a row exposing the 1-idea output:\n%s", strings.Join(lines, "\n"))
	}
	if idx := strings.Index(ideaRow, "idea"); idx != artifactLabelColumn {
		t.Errorf("idea row: primary label at column %d, want %d (schema ID, not filename stem): %q", idx, artifactLabelColumn, ideaRow)
	}
	if idx := strings.Index(ideaRow, "1-idea"); idx < idColumnEnd {
		t.Errorf("idea row: filename stem %q should render as secondary text after the label column (>= %d), got column %d: %q", "1-idea", idColumnEnd, idx, ideaRow)
	}

	planRow := lineContaining(lines, "5-plan")
	if planRow == "" {
		t.Fatalf("expected a row exposing the 5-plan output:\n%s", strings.Join(lines, "\n"))
	}
	if idx := strings.Index(planRow, "plan"); idx != artifactLabelColumn {
		t.Errorf("plan row: primary label at column %d, want %d (schema ID, not filename stem): %q", idx, artifactLabelColumn, planRow)
	}
	if idx := strings.Index(planRow, "5-plan"); idx < idColumnEnd {
		t.Errorf("plan row: filename stem %q should render as secondary text after the label column (>= %d), got column %d: %q", "5-plan", idColumnEnd, idx, planRow)
	}

	if row := lineContaining(lines, "requirements"); row == "" || strings.Index(row, "requirements") != artifactLabelColumn {
		t.Errorf("unwritten step should show its bare artifact ID at column %d: %q", artifactLabelColumn, row)
	}
	if row := lineContaining(lines, "tasks"); row == "" || strings.Index(row, "tasks") != artifactLabelColumn {
		t.Errorf("blocked step should show its bare artifact ID at column %d: %q", artifactLabelColumn, row)
	}
}

func TestRenderArtifactRowsSuppressRedundantSecondary(t *testing.T) {
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("add-themes", "spec-driven")}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:add-themes": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		theme: Theme{Styles: BuildStyles(DarkColors)},
		width: 100,
	}
	lines := plainIndexLines(t, m)
	proposalRow := lineContaining(lines, "proposal")
	if proposalRow == "" {
		t.Fatalf("expected a proposal artifact row:\n%s", strings.Join(lines, "\n"))
	}
	if got := strings.Count(proposalRow, "proposal"); got != 1 {
		t.Errorf("classic schema where output stem equals artifact ID should not repeat the name, got %d occurrences: %q", got, proposalRow)
	}
}
