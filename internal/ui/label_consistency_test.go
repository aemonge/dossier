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

func indexOfLineContaining(lines []string, needle string) int {
	for i, line := range lines {
		if strings.Contains(line, needle) {
			return i
		}
	}
	return -1
}

func TestRenderStrayArtifactsDistinctAndUncounted(t *testing.T) {
	specStatus := func(id string, status openspec.ArtifactStatus, outputs ...openspec.ArtifactOutput) openspec.ChangeArtifact {
		return openspec.ChangeArtifact{ID: id, Status: status, Source: openspec.ArtifactSourceStatus, Outputs: outputs}
	}
	stray := func(id, file string) openspec.ChangeArtifact {
		return openspec.ChangeArtifact{
			ID: id, Status: openspec.ArtifactStatusDone, Source: openspec.ArtifactSourceDiscovered,
			Outputs: []openspec.ArtifactOutput{{RelativePath: file, DisplayName: strings.TrimSuffix(file, ".md"), Present: true}},
		}
	}
	change := openspec.Change{
		Name: "fix-startup-race", Schema: "spec", StatusEnriched: true,
		Artifacts: []openspec.ChangeArtifact{
			specStatus("idea", openspec.ArtifactStatusDone, openspec.ArtifactOutput{RelativePath: "1-idea.md", DisplayName: "1-idea", Present: true}),
			specStatus("requirements", openspec.ArtifactStatusDone, openspec.ArtifactOutput{RelativePath: "2-drafted-requirements.md", DisplayName: "2-drafted-requirements", Present: true}),
			specStatus("research", openspec.ArtifactStatusDone, openspec.ArtifactOutput{RelativePath: "3-research.md", DisplayName: "3-research", Present: true}),
			specStatus("refined-requirements", openspec.ArtifactStatusReady),
			specStatus("plan", openspec.ArtifactStatusBlocked),
			specStatus("tasks", openspec.ArtifactStatusBlocked),
			stray("2-proposition", "2-proposition.md"),
			stray("2-requirements-a", "2-requirements-a.md"),
			stray("3-proposal", "3-proposal.md"),
		},
	}
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{change}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-startup-race": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		theme: Theme{Styles: BuildStyles(DarkColors)},
		width: 100,
	}
	lines := plainIndexLines(t, m)

	changeLine := lineContaining(lines, "fix-startup-race")
	if !strings.Contains(changeLine, "planning 3/6 · 3 strays") {
		t.Errorf("planning fraction must count status-source steps only, with the stray suffix: %q", changeLine)
	}
	lastStep := indexOfLineContaining(lines, "refined-requirements")
	firstStray := indexOfLineContaining(lines, "2-proposition")
	if lastStep < 0 || firstStray < 0 || firstStray < lastStep {
		t.Errorf("strays must render after the Spec's steps: lastStep=%d firstStray=%d", lastStep, firstStray)
	}
	for _, stem := range []string{"2-proposition", "2-requirements-a", "3-proposal"} {
		row := lineContaining(lines, stem)
		if row == "" {
			t.Fatalf("expected stray row %q:\n%s", stem, strings.Join(lines, "\n"))
		}
		if !strings.Contains(row, "[stray]") {
			t.Errorf("stray row %q must badge [stray]: %q", stem, row)
		}
		if strings.Contains(row, "[authored]") {
			t.Errorf("stray row %q must not badge [authored]: %q", stem, row)
		}
	}
}
