package openspec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const openSpecTimeout = 10 * time.Second

type OpenSpecRunner interface {
	Run(dir string, args ...string) (stdout, stderr []byte, err error)
}

type OSOpenSpecRunner struct{}

func (OSOpenSpecRunner) Run(dir string, args ...string) ([]byte, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), openSpecTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "openspec", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		err = fmt.Errorf("openspec command timed out after %s: %w", openSpecTimeout, ctx.Err())
	}
	return stdout.Bytes(), stderr.Bytes(), err
}

type CLI struct {
	runner OpenSpecRunner
}

func NewCLI(runner OpenSpecRunner) *CLI {
	return &CLI{runner: runner}
}

func NewOSCLI() *CLI {
	return NewCLI(OSOpenSpecRunner{})
}

type ChangeListItem struct {
	Name           string `json:"name"`
	CompletedTasks int    `json:"completedTasks"`
	TotalTasks     int    `json:"totalTasks"`
	LastModified   string `json:"lastModified"`
	Status         string `json:"status"`
}

type SchemaInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Artifacts   []string `json:"artifacts"`
	Source      string   `json:"source"`
}

type ArtifactPathStatus struct {
	OutputPath          string   `json:"outputPath"`
	ResolvedOutputPath  string   `json:"resolvedOutputPath"`
	ExistingOutputPaths []string `json:"existingOutputPaths"`
}

type StatusArtifact struct {
	ID         string   `json:"id"`
	OutputPath string   `json:"outputPath"`
	Status     string   `json:"status"`
	Requires   []string `json:"requires"`
}

type ChangeStatus struct {
	ChangeName    string                        `json:"changeName"`
	SchemaName    string                        `json:"schemaName"`
	ChangeRoot    string                        `json:"changeRoot"`
	ArtifactPaths map[string]ArtifactPathStatus `json:"artifactPaths"`
	IsComplete    bool                          `json:"isComplete"`
	ApplyRequires []string                      `json:"applyRequires"`
	Artifacts     []StatusArtifact              `json:"artifacts"`
}

type listChangesResponse struct {
	Changes []ChangeListItem `json:"changes"`
}

func (c *CLI) ListChanges(root string) ([]ChangeListItem, error) {
	args := []string{"list", "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return nil, err
	}
	var response listChangesResponse
	if err := json.Unmarshal(stdout, &response); err != nil {
		return nil, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return response.Changes, nil
}

func (c *CLI) Schemas(root string) ([]SchemaInfo, error) {
	args := []string{"schemas", "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return nil, err
	}
	var schemas []SchemaInfo
	if err := json.Unmarshal(stdout, &schemas); err != nil {
		return nil, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return schemas, nil
}

func (c *CLI) Status(root, changeName string) (ChangeStatus, error) {
	args := []string{"status", "--change", changeName, "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return ChangeStatus{}, err
	}
	var status ChangeStatus
	if err := json.Unmarshal(stdout, &status); err != nil {
		return ChangeStatus{}, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return status, nil
}

func (c *CLI) run(root string, args []string) ([]byte, error) {
	stdout, stderr, err := c.runner.Run(root, args...)
	if err == nil {
		return stdout, nil
	}
	message := strings.TrimSpace(string(stderr))
	if message == "" {
		message = err.Error()
	}
	return nil, fmt.Errorf("openspec %s: %s: %w", strings.Join(args, " "), message, err)
}

func (l *Loader) EnrichChange(ch Change, status ChangeStatus) Change {
	fallback := append([]ChangeArtifact(nil), ch.Artifacts...)
	fallbackByID := make(map[string]ChangeArtifact, len(fallback))
	for _, artifact := range fallback {
		fallbackByID[artifact.ID] = artifact
	}

	if status.ChangeName != "" && status.ChangeName != ch.Name {
		appendChangeDiagnostic(&ch, fmt.Sprintf("status change name %q does not match %q", status.ChangeName, ch.Name))
	}
	if status.ChangeRoot != "" && !sameCleanPath(status.ChangeRoot, ch.Path) {
		appendChangeDiagnostic(&ch, fmt.Sprintf("status change root %q does not match %q", status.ChangeRoot, ch.Path))
	}
	if status.SchemaName != "" {
		ch.Schema = status.SchemaName
	}

	ch.Artifacts = make([]ChangeArtifact, 0, len(status.Artifacts)+len(fallback))
	seen := make(map[string]bool, len(status.Artifacts))
	for _, raw := range status.Artifacts {
		seen[raw.ID] = true
		pathStatus := status.ArtifactPaths[raw.ID]
		outputPattern := raw.OutputPath
		if outputPattern == "" {
			outputPattern = pathStatus.OutputPath
		}
		artifact := ChangeArtifact{
			ID:            raw.ID,
			OutputPattern: outputPattern,
			Status:        parseArtifactStatus(raw.Status, &ch),
			Requires:      append([]string(nil), raw.Requires...),
			Source:        ArtifactSourceStatus,
		}
		for _, outputPath := range pathStatus.ExistingOutputPaths {
			relativePath, err := safeRelativeOutputPath(ch.Path, outputPath)
			if err != nil {
				appendChangeDiagnostic(&ch, fmt.Sprintf("artifact %s output %q: %v", raw.ID, outputPath, err))
				continue
			}
			content, err := l.fs.ReadFile(filepath.Join(ch.Path, filepath.FromSlash(relativePath)))
			if err != nil {
				if fallbackOutput, ok := findFallbackOutput(fallbackByID[raw.ID], relativePath); ok {
					artifact.Outputs = append(artifact.Outputs, fallbackOutput)
					continue
				}
				appendChangeDiagnostic(&ch, fmt.Sprintf("read status output %s: %v", relativePath, err))
				continue
			}
			artifact.Outputs = append(artifact.Outputs, ArtifactOutput{
				RelativePath: relativePath,
				DisplayName:  discoveredOutputName(relativePath),
				Content:      string(content),
				Present:      true,
			})
		}
		ch.Artifacts = append(ch.Artifacts, artifact)
	}
	for _, artifact := range fallback {
		if !seen[artifact.ID] {
			ch.Artifacts = append(ch.Artifacts, artifact)
		}
	}
	deriveConventionalFields(&ch)
	return ch
}

func parseArtifactStatus(status string, ch *Change) ArtifactStatus {
	switch ArtifactStatus(status) {
	case ArtifactStatusBlocked, ArtifactStatusReady, ArtifactStatusDone:
		return ArtifactStatus(status)
	case ArtifactStatusUnknown:
		return ArtifactStatusUnknown
	default:
		appendChangeDiagnostic(ch, fmt.Sprintf("unknown artifact status %q", status))
		return ArtifactStatusUnknown
	}
}

func safeRelativeOutputPath(root, path string) (string, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve change root: %w", err)
	}
	absolutePath := path
	if !filepath.IsAbs(absolutePath) {
		absolutePath = filepath.Join(absoluteRoot, absolutePath)
	}
	absolutePath, err = filepath.Abs(absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve output: %w", err)
	}
	relativePath, err := filepath.Rel(absoluteRoot, absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve output relative path: %w", err)
	}
	if !isContainedRelativePath(relativePath) {
		return "", fmt.Errorf("outside change root")
	}

	evaluatedRoot, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return "", fmt.Errorf("resolve change root symlinks: %w", err)
	}
	evaluatedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve output symlinks: %w", err)
	}
	evaluatedRelative, err := filepath.Rel(evaluatedRoot, evaluatedPath)
	if err != nil {
		return "", fmt.Errorf("resolve evaluated output relative path: %w", err)
	}
	if !isContainedRelativePath(evaluatedRelative) {
		return "", fmt.Errorf("outside change root after resolving symlinks")
	}
	return filepath.ToSlash(relativePath), nil
}

func isContainedRelativePath(path string) bool {
	return path != ".." && !strings.HasPrefix(path, ".."+string(filepath.Separator)) && !filepath.IsAbs(path)
}

func sameCleanPath(a, b string) bool {
	absoluteA, errA := filepath.Abs(a)
	absoluteB, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return filepath.Clean(absoluteA) == filepath.Clean(absoluteB)
}

func findFallbackOutput(artifact ChangeArtifact, relativePath string) (ArtifactOutput, bool) {
	for _, output := range artifact.Outputs {
		if output.RelativePath == relativePath {
			return output, true
		}
	}
	return ArtifactOutput{}, false
}
