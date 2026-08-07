package openspec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeOpenSpecRunner struct {
	stdout []byte
	stderr []byte
	err    error
	dir    string
	args   []string
}

func (r *fakeOpenSpecRunner) Run(dir string, args ...string) ([]byte, []byte, error) {
	r.dir = dir
	r.args = append([]string(nil), args...)
	return r.stdout, r.stderr, r.err
}

func TestCLIListChanges(t *testing.T) {
	runner := &fakeOpenSpecRunner{stdout: []byte(`{
  "changes": [{"name":"fix-cache","completedTasks":2,"totalTasks":3,"lastModified":"2026-08-07T00:00:00Z","status":"in-progress"}],
  "root": {"path":"/project","source":"nearest"}
}`)}
	client := NewCLI(runner)

	changes, err := client.ListChanges("/project")
	if err != nil {
		t.Fatal(err)
	}
	if runner.dir != "/project" || strings.Join(runner.args, " ") != "list --json" {
		t.Fatalf("unexpected command: dir=%q args=%v", runner.dir, runner.args)
	}
	if len(changes) != 1 || changes[0].Name != "fix-cache" || changes[0].CompletedTasks != 2 {
		t.Fatalf("unexpected changes: %#v", changes)
	}
}

func TestCLISchemas(t *testing.T) {
	runner := &fakeOpenSpecRunner{stdout: []byte(`[
  {"name":"feature","description":"Feature workflow","artifacts":["proposal","specs","design","tasks"],"source":"project"},
  {"name":"spike","description":"Research workflow","artifacts":["proposal","research","decision"],"source":"user"}
]`)}
	client := NewCLI(runner)

	schemas, err := client.Schemas("/project")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(runner.args, " ") != "schemas --json" {
		t.Fatalf("unexpected args: %v", runner.args)
	}
	if len(schemas) != 2 || schemas[1].Name != "spike" || schemas[1].Artifacts[2] != "decision" {
		t.Fatalf("unexpected schemas: %#v", schemas)
	}
}

func TestCLIStatus(t *testing.T) {
	runner := &fakeOpenSpecRunner{stdout: []byte(`{
  "changeName":"fix-cache",
  "schemaName":"bugfix",
  "changeRoot":"/project/openspec/changes/fix-cache",
  "artifactPaths": {
    "proposal":{"outputPath":"proposal.md","resolvedOutputPath":"/project/openspec/changes/fix-cache/proposal.md","existingOutputPaths":["/project/openspec/changes/fix-cache/proposal.md"]},
    "diagnosis":{"outputPath":"diagnosis.md","resolvedOutputPath":"/project/openspec/changes/fix-cache/diagnosis.md","existingOutputPaths":[]}
  },
  "isComplete":false,
  "applyRequires":["tasks"],
  "artifacts":[
    {"id":"proposal","outputPath":"proposal.md","status":"done","requires":[]},
    {"id":"diagnosis","outputPath":"diagnosis.md","status":"ready","requires":["proposal","reproduction"]}
  ]
}`)}
	client := NewCLI(runner)

	status, err := client.Status("/project", "fix-cache")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(runner.args, " ") != "status --change fix-cache --json" {
		t.Fatalf("unexpected args: %v", runner.args)
	}
	if status.SchemaName != "bugfix" || len(status.Artifacts) != 2 {
		t.Fatalf("unexpected status: %#v", status)
	}
	if status.Artifacts[1].Requires[1] != "reproduction" {
		t.Fatalf("unexpected artifact requirements: %#v", status.Artifacts[1])
	}
}

func TestCLIErrors(t *testing.T) {
	t.Run("command error includes stderr", func(t *testing.T) {
		runner := &fakeOpenSpecRunner{stderr: []byte("unknown change\n"), err: errors.New("exit status 1")}
		_, err := NewCLI(runner).Status("/project", "missing")
		if err == nil || !strings.Contains(err.Error(), "unknown change") || !strings.Contains(err.Error(), "status --change missing --json") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing executable remains actionable", func(t *testing.T) {
		runner := &fakeOpenSpecRunner{err: errors.New("executable file not found in $PATH")}
		_, err := NewCLI(runner).Schemas("/project")
		if err == nil || !strings.Contains(err.Error(), "executable file not found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("malformed JSON identifies command", func(t *testing.T) {
		runner := &fakeOpenSpecRunner{stdout: []byte(`{"changes": [`)}
		_, err := NewCLI(runner).ListChanges("/project")
		if err == nil || !strings.Contains(err.Error(), "decode openspec list --json") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestEnrichRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "openspec", "changes", "unsafe-change")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(root, "outside.md")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "linked.md")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	loader := NewLoader(OSFS{})
	change := Change{Name: "unsafe-change", Path: dir}
	status := ChangeStatus{
		ChangeName: "unsafe-change",
		ChangeRoot: dir,
		ArtifactPaths: map[string]ArtifactPathStatus{
			"linked": {OutputPath: "linked.md", ExistingOutputPaths: []string{link}},
		},
		Artifacts: []StatusArtifact{{ID: "linked", OutputPath: "linked.md", Status: "done"}},
	}

	enriched := loader.EnrichChange(change, status)
	artifact, ok := enriched.ArtifactByID("linked")
	if !ok {
		t.Fatal("expected status artifact")
	}
	if len(artifact.Outputs) != 0 {
		t.Fatalf("symlink escape must not be read: %#v", artifact.Outputs)
	}
	if !strings.Contains(enriched.Diagnostic, "outside change root") {
		t.Fatalf("expected escape diagnostic, got %q", enriched.Diagnostic)
	}
}

func TestEnrichChangeFromStatus(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "openspec", "changes", "fix-cache")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		".openspec.yaml": "schema: bugfix\ncreated: 2026-08-07\n",
		"proposal.md":    "# Proposal",
		"diagnosis.md":   "# Diagnosis",
		"notes.md":       "# Fallback notes",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	loader := NewLoader(OSFS{})
	project, err := loader.LoadFrom(root)
	if err != nil {
		t.Fatal(err)
	}
	change := project.Changes[0]
	outside := filepath.Join(root, "outside.md")
	if err := os.WriteFile(outside, []byte("outside"), 0644); err != nil {
		t.Fatal(err)
	}
	status := ChangeStatus{
		ChangeName: "fix-cache",
		SchemaName: "bugfix",
		ChangeRoot: dir,
		ArtifactPaths: map[string]ArtifactPathStatus{
			"proposal": {
				OutputPath:          "proposal.md",
				ResolvedOutputPath:  filepath.Join(dir, "proposal.md"),
				ExistingOutputPaths: []string{filepath.Join(dir, "proposal.md")},
			},
			"diagnosis": {
				OutputPath:          "diagnosis.md",
				ResolvedOutputPath:  filepath.Join(dir, "diagnosis.md"),
				ExistingOutputPaths: []string{filepath.Join(dir, "diagnosis.md"), outside},
			},
			"tasks": {OutputPath: "tasks.md"},
		},
		Artifacts: []StatusArtifact{
			{ID: "proposal", OutputPath: "proposal.md", Status: "done"},
			{ID: "diagnosis", OutputPath: "diagnosis.md", Status: "ready", Requires: []string{"proposal", "reproduction"}},
			{ID: "tasks", OutputPath: "tasks.md", Status: "mystery", Requires: []string{"diagnosis"}},
		},
	}

	enriched := loader.EnrichChange(change, status)
	if enriched.Schema != "bugfix" {
		t.Fatalf("expected bugfix schema, got %q", enriched.Schema)
	}
	wantIDs := []string{"proposal", "diagnosis", "tasks", "notes"}
	if len(enriched.Artifacts) != len(wantIDs) {
		t.Fatalf("expected %d artifacts, got %#v", len(wantIDs), enriched.Artifacts)
	}
	for i, id := range wantIDs {
		if enriched.Artifacts[i].ID != id {
			t.Errorf("artifact %d: expected %q, got %q", i, id, enriched.Artifacts[i].ID)
		}
	}
	diagnosis, _ := enriched.ArtifactByID("diagnosis")
	if diagnosis.Status != ArtifactStatusReady || len(diagnosis.Requires) != 2 {
		t.Fatalf("unexpected diagnosis metadata: %#v", diagnosis)
	}
	if len(diagnosis.Outputs) != 1 || diagnosis.Outputs[0].Content != "# Diagnosis" {
		t.Fatalf("unsafe output was not rejected: %#v", diagnosis.Outputs)
	}
	tasks, _ := enriched.ArtifactByID("tasks")
	if tasks.Status != ArtifactStatusUnknown {
		t.Fatalf("unknown status must map to unknown, got %q", tasks.Status)
	}
	if !strings.Contains(enriched.Diagnostic, "outside change root") || !strings.Contains(enriched.Diagnostic, `unknown artifact status "mystery"`) {
		t.Fatalf("expected enrichment diagnostics, got %q", enriched.Diagnostic)
	}
	notes, _ := enriched.ArtifactByID("notes")
	if notes.Source != ArtifactSourceDiscovered {
		t.Fatalf("expected unmatched fallback artifact, got %#v", notes)
	}
}
