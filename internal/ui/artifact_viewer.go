package ui

import (
	"fmt"

	"github.com/fselich/dossier/internal/openspec"
)

type artifactSelection struct {
	ArtifactID string
	OutputPath string
}

type artifactOutputIdentity struct {
	ChangeName string
	Archived   bool
	ArtifactID string
	OutputPath string
}

type artifactRenderedMsg struct {
	key     artifactOutputIdentity
	content string
}

func (m *Model) selectDefaultArtifact() bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		m.artifactSelection = artifactSelection{}
		return false
	}
	return m.selectArtifactOutput(ch.Artifacts[0].ID, firstArtifactOutputPath(ch.Artifacts[0]))
}

func firstArtifactOutputPath(artifact openspec.ChangeArtifact) string {
	for _, output := range artifact.Outputs {
		if output.Present {
			return output.RelativePath
		}
	}
	return ""
}

func (m *Model) selectArtifactOutput(artifactID, outputPath string) bool {
	ch := m.current()
	if ch == nil {
		return false
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != artifactID {
			continue
		}
		if outputPath != "" {
			found := false
			for _, output := range artifact.Outputs {
				if output.RelativePath == outputPath && output.Present {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		} else {
			outputPath = firstArtifactOutputPath(artifact)
		}
		m.artifactSelection = artifactSelection{ArtifactID: artifactID, OutputPath: outputPath}
		m.viewingCode = false
		m.syncCompatibilityTab(artifactID)
		return true
	}
	return false
}

func (m *Model) syncCompatibilityTab(artifactID string) {
	switch artifactID {
	case "proposal":
		m.tab = TabProposal
	case "design":
		m.tab = TabDesign
	case "specs":
		m.tab = TabSpecs
	case "tasks":
		m.tab = TabTasks
	default:
		m.tab = TabProposal
	}
}

func (m *Model) selectedArtifactOutput() (openspec.ChangeArtifact, *openspec.ArtifactOutput, bool) {
	ch := m.current()
	if ch == nil {
		return openspec.ChangeArtifact{}, nil, false
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != m.artifactSelection.ArtifactID {
			continue
		}
		if m.artifactSelection.OutputPath == "" {
			return artifact, nil, true
		}
		for outputIndex := range artifact.Outputs {
			if artifact.Outputs[outputIndex].RelativePath == m.artifactSelection.OutputPath && artifact.Outputs[outputIndex].Present {
				return artifact, &artifact.Outputs[outputIndex], true
			}
		}
		return artifact, nil, true
	}
	return openspec.ChangeArtifact{}, nil, false
}

func (m *Model) selectedArtifactKey() (artifactOutputIdentity, bool) {
	artifact, output, ok := m.selectedArtifactOutput()
	if !ok {
		return artifactOutputIdentity{}, false
	}
	ch := m.current()
	key := artifactOutputIdentity{
		ChangeName: ch.Name,
		Archived:   m.mode == ModeViewingArchive,
		ArtifactID: artifact.ID,
	}
	if output != nil {
		key.OutputPath = output.RelativePath
	}
	return key, true
}

func (m *Model) invalidateArtifactRenderCache(changeName string) {
	for key := range m.artifactRenderCache {
		if key.ChangeName == changeName {
			delete(m.artifactRenderCache, key)
		}
	}
}

func (m *Model) reconcileArtifactSelection() bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		m.artifactSelection = artifactSelection{}
		return false
	}
	if m.artifactSelection.ArtifactID == "" {
		return m.selectDefaultArtifact()
	}
	for _, artifact := range ch.Artifacts {
		if artifact.ID != m.artifactSelection.ArtifactID {
			continue
		}
		if m.artifactSelection.OutputPath == "" {
			return true
		}
		for _, output := range artifact.Outputs {
			if output.RelativePath == m.artifactSelection.OutputPath && output.Present {
				return true
			}
		}
		return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
	}
	return m.selectDefaultArtifact()
}

func (m *Model) moveArtifactSelection(delta int) bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		return false
	}
	current := 0
	for index := range ch.Artifacts {
		if ch.Artifacts[index].ID == m.artifactSelection.ArtifactID {
			current = index
			break
		}
	}
	next := (current + delta%len(ch.Artifacts) + len(ch.Artifacts)) % len(ch.Artifacts)
	artifact := ch.Artifacts[next]
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func (m *Model) moveViewerDestination(delta int) bool {
	ch := m.current()
	if ch == nil || len(ch.Artifacts) == 0 {
		return false
	}
	includeCode := m.isGitRepo && m.mode == ModeNormal && len(m.gitState.Files) > 0
	count := len(ch.Artifacts)
	if includeCode {
		count++
	}
	current := 0
	if m.viewingCode {
		current = len(ch.Artifacts)
	} else {
		for index := range ch.Artifacts {
			if ch.Artifacts[index].ID == m.artifactSelection.ArtifactID {
				current = index
				break
			}
		}
	}
	next := (current + delta%count + count) % count
	if includeCode && next == len(ch.Artifacts) {
		m.viewingCode = true
		m.tab = TabGit
		return true
	}
	artifact := ch.Artifacts[next]
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func (m *Model) selectArtifactPosition(position int) bool {
	ch := m.current()
	if ch == nil || position < 0 || position >= len(ch.Artifacts) {
		return false
	}
	artifact := ch.Artifacts[position]
	if artifact.ID == m.artifactSelection.ArtifactID && len(artifact.Outputs) > 1 {
		current := 0
		for index := range artifact.Outputs {
			if artifact.Outputs[index].RelativePath == m.artifactSelection.OutputPath {
				current = index
				break
			}
		}
		next := (current + 1) % len(artifact.Outputs)
		return m.selectArtifactOutput(artifact.ID, artifact.Outputs[next].RelativePath)
	}
	return m.selectArtifactOutput(artifact.ID, firstArtifactOutputPath(artifact))
}

func artifactStatusView(artifact openspec.ChangeArtifact) string {
	status := "unknown"
	switch artifact.Status {
	case openspec.ArtifactStatusDone:
		status = "authored"
	case openspec.ArtifactStatusReady:
		status = "ready to author"
	case openspec.ArtifactStatusBlocked:
		status = "blocked"
	default:
		if artifact.Source == openspec.ArtifactSourceDiscovered {
			status = "discovered"
		}
	}
	text := fmt.Sprintf("# %s\n\nStatus: %s", artifact.ID, status)
	if len(artifact.Requires) > 0 {
		text += "\n\nRequires: "
		for index, requirement := range artifact.Requires {
			if index > 0 {
				text += ", "
			}
			text += requirement
		}
	}
	return text + "\n"
}
