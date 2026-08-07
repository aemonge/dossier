package ui

import (
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/fselich/dossier/internal/git"
)

func (m Model) updateViewer(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.tab == TabGit && m.gitState.ErrMsg != "" {
		m.gitState.ErrMsg = ""
	}

	keys := m.effectiveKeyMap().Viewer
	switch {

	case matchesKey(msg, keys.Quit):
		return m, tea.Quit

	case matchesKey(msg, keys.Stage):
		if m.readOnly || m.tab != TabGit || m.gitState.ShowingDiff || len(m.gitState.Files) == 0 {
			return m, nil
		}
		f := m.gitState.Files[m.gitState.Cursor]

		var err error
		if f.Y != ' ' || f.X == '?' && f.Y == '?' {
			paths := []string{f.Path}
			if f.OldPath != "" {
				paths = []string{f.OldPath, f.Path}
			}
			err = git.Stage(m.gitRoot, paths...)
		} else {
			paths := []string{f.Path}
			if f.OldPath != "" {
				paths = []string{f.OldPath, f.Path}
			}
			err = git.Unstage(m.gitRoot, paths...)
		}
		if err != nil {
			m.gitState.ErrMsg = strings.TrimSpace(err.Error())
		}
		files, statusErr := git.Status(m.gitRoot)
		if statusErr == nil {
			m.gitState.Files = files
		}
		m.restoreGitCursor(f.Path)
		m.refreshGitViewport()
		return m, nil

	case matchesKey(msg, keys.Info):
		m.prevMode = m.mode
		m.mode = ModeViewingConfig
		return m.commitStateChange()

	case matchesKey(msg, keys.Back):
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.toggleGitDiff()
			return m, nil
		}
		m.enterIndex()
		return m, nil

	case matchesKey(msg, keys.Previous):
		if len(m.project.Changes) > 0 {
			m.changeIdx = (m.changeIdx - 1 + len(m.project.Changes)) % len(m.project.Changes)
			m.renderCache = make(map[Tab]string)
			m.artifactRenderCache = make(map[artifactOutputIdentity]string)
			m.loadTaskItems()
			m.tab = m.defaultTab()
			m.selectDefaultArtifact()
			m.specIdx = 0
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.Next):
		if len(m.project.Changes) > 0 {
			m.changeIdx = (m.changeIdx + 1) % len(m.project.Changes)
			m.renderCache = make(map[Tab]string)
			m.artifactRenderCache = make(map[artifactOutputIdentity]string)
			m.loadTaskItems()
			m.tab = m.defaultTab()
			m.selectDefaultArtifact()
			m.specIdx = 0
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.Right):
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.gitState.ScrollX += 10
			m.refreshGitViewport()
		}

	case matchesKey(msg, keys.Left):
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.gitState.ScrollX -= 10
			if m.gitState.ScrollX < 0 {
				m.gitState.ScrollX = 0
			}
			m.refreshGitViewport()
		}

	case matchesKey(msg, keys.NextDiff):
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.moveGitDiffCursorDown()
			m.loadDiffForFile(m.gitState.Cursor)
			return m, nil
		}

	case matchesKey(msg, keys.PreviousDiff):
		if m.tab == TabGit && m.gitState.ShowingDiff {
			m.moveGitDiffCursorUp()
			m.loadDiffForFile(m.gitState.Cursor)
			return m, nil
		}

	case matchesKey(msg, keys.ProposalTab):
		if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
			if m.selectArtifactPosition(0) {
				return m.commitStateChange()
			}
		} else if m.tabAvailable(TabProposal) {
			m.tab = TabProposal
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.DesignTab):
		if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
			if m.selectArtifactPosition(1) {
				return m.commitStateChange()
			}
		} else if m.tabAvailable(TabDesign) {
			m.tab = TabDesign
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.SpecsTab):
		if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
			if m.selectArtifactPosition(2) {
				return m.commitStateChange()
			}
		} else if m.tabAvailable(TabSpecs) {
			if m.tab == TabSpecs {
				ch := m.current()
				if ch != nil && len(ch.SpecFiles) > 1 {
					m.specIdx = (m.specIdx + 1) % len(ch.SpecFiles)
					delete(m.renderCache, TabSpecs)
				}
			} else {
				m.tab = TabSpecs
				m.specIdx = 0
			}
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.TasksTab):
		if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
			if m.selectArtifactPosition(3) {
				return m.commitStateChange()
			}
		} else if m.tabAvailable(TabTasks) {
			m.tab = TabTasks
			return m.commitStateChange()
		}
	case matchesKey(msg, keys.GitTab):
		if m.tabAvailable(TabGit) {
			m.viewingCode = true
			m.tab = TabGit
			return m.commitStateChange()
		}

	case matchesKey(msg, keys.NextTab):
		if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
			if m.moveViewerDestination(1) {
				return m.commitStateChange()
			}
		} else {
			nxt := m.nextAvailableTab(m.tab, 1)
			if nxt != m.tab {
				m.tab = nxt
				return m.commitStateChange()
			}
		}
	case matchesKey(msg, keys.PreviousTab):
		if ch := m.current(); ch != nil && len(ch.Artifacts) > 0 {
			if m.moveViewerDestination(-1) {
				return m.commitStateChange()
			}
		} else {
			prv := m.nextAvailableTab(m.tab, -1)
			if prv != m.tab {
				m.tab = prv
				return m.commitStateChange()
			}
		}

	case matchesKey(msg, keys.ViewDiff):
		if m.tab == TabGit {
			m.toggleGitDiff()
		}

	case matchesKey(msg, keys.Down):
		switch m.tab {
		case TabTasks:
			if m.mode == ModeViewingArchive {
				m.vp.ScrollDown(1)
			} else {
				m.moveCursorDown()
				m.refreshTasksViewport()
			}
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.ScrollDown(1)
			} else {
				m.moveGitCursorDown()
				m.refreshGitViewport()
			}
		default:
			m.vp.ScrollDown(1)
		}

	case matchesKey(msg, keys.PageDown):
		switch m.tab {
		case TabTasks:
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.PageDown()
			}
		default:
			m.vp.PageDown()
		}

	case matchesKey(msg, keys.Up):
		switch m.tab {
		case TabTasks:
			if m.mode == ModeViewingArchive {
				m.vp.ScrollUp(1)
			} else {
				m.moveCursorUp()
				m.refreshTasksViewport()
			}
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.ScrollUp(1)
			} else {
				m.moveGitCursorUp()
				m.refreshGitViewport()
			}
		default:
			m.vp.ScrollUp(1)
		}

	case matchesKey(msg, keys.PageUp):
		switch m.tab {
		case TabTasks:
		case TabGit:
			if m.gitState.ShowingDiff {
				m.vp.PageUp()
			}
		default:
			m.vp.PageUp()
		}

	case matchesKey(msg, keys.ToggleTask):
		if m.mode == ModeViewingArchive {
			return m, nil
		}
		if m.tab == TabTasks {
			return m, m.doToggle()
		}

	case matchesKey(msg, keys.Open):
		if m.tab == TabGit {
			m.toggleGitDiff()
			return m, nil
		}
		if m.readOnly || m.mode == ModeViewingArchive {
			return m, nil
		}
		if m.tabAvailable(m.tab) {
			path := m.artifactPath()
			if path != "" {
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vi"
				}
				cmd := exec.Command(editor, path)
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return editorReturnMsg{}
				})
			}
		}
	}
	return m, nil
}
