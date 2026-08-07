package ui

import (
	"strings"
	"testing"

	"github.com/fselich/dossier/internal/openspec"
)

func hierarchyChange(name, schema string) openspec.Change {
	return openspec.Change{
		Name:   name,
		Schema: schema,
		Artifacts: []openspec.ChangeArtifact{
			{
				ID:     "proposal",
				Status: openspec.ArtifactStatusDone,
				Source: openspec.ArtifactSourceStatus,
				Outputs: []openspec.ArtifactOutput{{
					RelativePath: "proposal.md", DisplayName: "proposal", Present: true,
				}},
			},
			{
				ID:       "specs",
				Status:   openspec.ArtifactStatusDone,
				Requires: []string{"proposal"},
				Source:   openspec.ArtifactSourceStatus,
				Outputs: []openspec.ArtifactOutput{
					{RelativePath: "specs/auth/spec.md", DisplayName: "auth", Present: true},
					{RelativePath: "specs/payments/spec.md", DisplayName: "payments", Present: true},
				},
			},
			{
				ID:       "tasks",
				Status:   openspec.ArtifactStatusReady,
				Requires: []string{"specs", "design"},
				Source:   openspec.ArtifactSourceStatus,
			},
		},
	}
}

func TestBuildSchemaAwareHierarchy(t *testing.T) {
	active := hierarchyChange("fix-cache", "bugfix")
	archived := openspec.Change{
		Name: "old-spike", Schema: "spike", DisplayDate: "01/05/2026",
		Artifacts: []openspec.ChangeArtifact{{
			ID: "research", Status: openspec.ArtifactStatusDone, Source: openspec.ArtifactSourceDiscovered,
			Outputs: []openspec.ArtifactOutput{{RelativePath: "research.md", DisplayName: "research", Present: true}},
		}},
	}
	m := &Model{
		project:      &openspec.Project{Changes: []openspec.Change{active}},
		projectSpecs: []openspec.ProjectSpec{{Name: "canonical-auth", RequirementNames: []string{"Login"}}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true, "archive:old-spike": true},
			ExpandedArtifacts: map[string]bool{"active:fix-cache:artifact:specs": true},
			ArchiveChanges:    []openspec.Change{archived},
		},
		width: 100,
	}

	m.buildIndexItems()
	want := []struct {
		kind     indexItemKind
		identity string
	}{
		{indexKindSection, "section:active-work"},
		{indexKindActive, "active:fix-cache"},
		{indexKindArtifact, "active:fix-cache:artifact:proposal"},
		{indexKindArtifact, "active:fix-cache:artifact:specs"},
		{indexKindArtifactOutput, "active:fix-cache:artifact:specs:output:specs/auth/spec.md"},
		{indexKindArtifactOutput, "active:fix-cache:artifact:specs:output:specs/payments/spec.md"},
		{indexKindArtifact, "active:fix-cache:artifact:tasks"},
		{indexKindSection, "section:canonical-specs"},
		{indexKindSpec, "canonical-spec:canonical-auth"},
		{indexKindSection, "section:history"},
		{indexKindArchived, "archive:old-spike"},
		{indexKindArtifact, "archive:old-spike:artifact:research"},
	}
	if len(m.index.Items) != len(want) {
		t.Fatalf("expected %d hierarchy rows, got %d: %#v", len(want), len(m.index.Items), m.index.Items)
	}
	for index, expected := range want {
		item := m.index.Items[index]
		if item.kind != expected.kind || item.identity != expected.identity {
			t.Errorf("row %d: expected kind=%d identity=%q, got kind=%d identity=%q", index, expected.kind, expected.identity, item.kind, item.identity)
		}
	}
}

func TestHierarchyIdentitySurvivesSiblingInsertion(t *testing.T) {
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: make(map[string]bool),
		},
	}
	m.buildIndexItems()
	identity := "active:fix-cache:artifact:specs"
	before := indexItemByIdentity(m.index.Items, identity)
	if before < 0 {
		t.Fatal("expected specs artifact before insertion")
	}

	m.project.Changes = append([]openspec.Change{hierarchyChange("add-auth", "feature")}, m.project.Changes...)
	m.buildIndexItems()
	after := indexItemByIdentity(m.index.Items, identity)
	if after < 0 {
		t.Fatal("expected specs artifact after insertion")
	}
	if m.index.Items[after].identity != m.index.Items[before+1].identity {
		// Numeric position changes; semantic identity must not.
		if m.index.Items[after].identity != identity {
			t.Fatalf("expected stable identity %q, got %q", identity, m.index.Items[after].identity)
		}
	}
}

func TestRebuildHierarchyPreservesCursorIdentity(t *testing.T) {
	m := &Model{
		project: &openspec.Project{Changes: []openspec.Change{hierarchyChange("fix-cache", "bugfix")}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: make(map[string]bool),
		},
	}
	m.buildIndexItems()
	identity := "active:fix-cache:artifact:specs"
	m.index.Cursor = indexItemByIdentity(m.index.Items, identity)
	if m.index.Cursor < 0 {
		t.Fatal("expected initial cursor identity")
	}

	m.project.Changes = append([]openspec.Change{hierarchyChange("add-auth", "feature")}, m.project.Changes...)
	m.rebuildIndexPreservingCursor()
	if m.index.Cursor < 0 || m.index.Cursor >= len(m.index.Items) {
		t.Fatalf("invalid restored cursor %d", m.index.Cursor)
	}
	if got := m.index.Items[m.index.Cursor].identity; got != identity {
		t.Fatalf("expected cursor identity %q, got %q", identity, got)
	}
}

func TestRenderSchemaAwareHierarchy(t *testing.T) {
	change := hierarchyChange("fix-cache", "bugfix")
	change.Diagnostic = "degraded status"
	m := &Model{
		project:      &openspec.Project{Changes: []openspec.Change{change}},
		projectSpecs: []openspec.ProjectSpec{{Name: "auth", RequirementCount: 2}},
		index: indexState{
			ExpandedSpecs:     make(map[int]bool),
			ExpandedChanges:   map[string]bool{"active:fix-cache": true},
			ExpandedArtifacts: make(map[string]bool),
		},
		width: 100,
	}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	for _, expected := range []string{"Active Work", "Canonical Specs", "History", "fix-cache", "bugfix", "2/3 artifacts", "proposal", "tasks", "ready", "requires specs, design", "degraded status"} {
		if !strings.Contains(content, expected) {
			t.Errorf("expected rendered hierarchy to contain %q:\n%s", expected, content)
		}
	}
}

func TestRenderSchemaAwareHierarchyEmptyStates(t *testing.T) {
	m := &Model{project: &openspec.Project{}, width: 80}
	m.buildIndexItems()
	content, _ := m.renderIndexContent()
	for _, expected := range []string{"No active work", "No canonical specifications available", "No history"} {
		if !strings.Contains(content, expected) {
			t.Errorf("expected empty hierarchy to contain %q:\n%s", expected, content)
		}
	}
}
