package ui

import (
	"path/filepath"
	"time"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"github.com/aemonge/dossier/internal/git"
	"github.com/aemonge/dossier/internal/openspec"
	"github.com/aemonge/dossier/internal/settings"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeIndex
	ModeViewingArchive
	ModeViewingSpec
	ModeViewingConfig
)

type errClearMsg struct{}
type indexStatusClearMsg string
type editorReturnMsg struct {
	indexIdentity    string
	artifactIdentity artifactOutputIdentity
}

// specRenderedMsg carries async glamour output for ModeViewingSpec.
type specRenderedMsg struct {
	content  string
	jumpLine int // line offset to scroll to after render; 0 = start of document
}

// renderedConfigMsg carries async glamour output for ModeViewingConfig.
type renderedConfigMsg struct {
	content string
}

type tickMsg time.Time

type indexActionMode int

const (
	indexActionIdle indexActionMode = iota
	indexActionSchema
	indexActionName
	indexActionConfirmArchive
	indexActionConfirmReactivate
	indexActionConfirmUndo
	indexActionPending
)

type indexOperation int

const (
	indexOperationNone indexOperation = iota
	indexOperationCreate
	indexOperationValidate
	indexOperationArchive
	indexOperationReactivate
	indexOperationUndo
)

type indexActionState struct {
	Mode           indexActionMode
	Operation      indexOperation
	Schema         string
	SchemaCursor   int
	SchemaFilter   string
	Name           string
	TargetIdentity string
	Status         string
	StatusError    bool
}

type indexActionResultMsg struct {
	Kind           indexOperation
	Name           string
	TargetIdentity string
	Validation     openspec.ValidationResult
	Archive        openspec.ArchiveResult
	Undo           *openspec.LifecycleUndo
	Err            error
}

type indexState struct {
	Items             []indexItem
	Cursor            int
	ExpandedSpecs     map[int]bool
	ExpandedChanges   map[string]bool
	ExpandedArtifacts map[string]bool
	CollapsedSections [3]bool
	SortBySuffix      bool
	Order             []int
	ArchiveChanges    []openspec.Change
	ArchiveCursor     int

	FilterText     string
	FilterActive   bool
	FilterIndices  []int
	PrevFilterText string
	Action         indexActionState
}

type specViewerState struct {
	Cursor     int
	JumpTarget string
	FocusMode  bool
	ReqCursor  int
}

type taskState struct {
	Items  []openspec.TaskItem
	Cursor int
}

type gitState struct {
	Files       []git.FileStatus
	Cursor      int
	ShowingDiff bool
	DiffLines   []DiffLine
	DiffFile    string
	ScrollX     int
	ErrMsg      string
}

type indexItemKind int

const (
	indexKindActive indexItemKind = iota
	indexKindArchived
	indexKindArtifact
	indexKindArtifactOutput
	indexKindSpec
	indexKindRequirement
	indexKindSection
)

const (
	sectionActive   = 0
	sectionSpecs    = 1
	sectionArchived = 2
)

var sectionNames = []string{
	sectionActive:   "Active Work",
	sectionSpecs:    "Canonical Specs",
	sectionArchived: "History",
}

type indexItem struct {
	kind        indexItemKind
	idx         int // change, archive, project spec, or section index
	artifactIdx int
	outputIdx   int
	reqIdx      int
	archived    bool
	identity    string
	depth       int
}

type Model struct {
	root   string
	loader *openspec.Loader

	project           *openspec.Project
	changeIdx         int
	artifactSelection artifactSelection
	viewingCode       bool

	openSpec              openSpecClient
	lifecycle             *openspec.LifecycleService
	lifecycleUndo         *openspec.LifecycleUndo
	schemaCatalog         []openspec.SchemaInfo
	schemaCatalogErr      string
	discoveryFingerprints map[string]string
	enrichedFingerprints  map[string]string
	pendingEnrichments    map[string]string
	enrichmentRetryAfter  map[string]time.Time

	vp      viewport.Model
	vpReady bool

	tasks taskState

	errMsg     string
	loading    bool
	singlePath bool
	readOnly   bool

	isGitRepo bool
	gitRoot   string
	gitState  gitState

	width, height int

	artifactRenderCache map[artifactOutputIdentity]string
	glamourRenderer     *glamour.TermRenderer
	lastRenderWidth     int

	mode                Mode
	prevMode            Mode
	returnIndexIdentity string
	index               indexState
	projectSpecs        []openspec.ProjectSpec
	specViewer          specViewerState
	projectConfig       openspec.ProjectConfig
	theme               Theme
	keyMap              settings.KeyConfig
}

func New(project *openspec.Project, cfg openspec.ProjectConfig, root string, loader *openspec.Loader, theme Theme, keyMap settings.KeyConfig, readOnly bool) Model {
	return NewWithStartView(project, cfg, root, loader, theme, keyMap, readOnly, "index")
}

func NewWithStartView(project *openspec.Project, cfg openspec.ProjectConfig, root string, loader *openspec.Loader, theme Theme, keyMap settings.KeyConfig, readOnly bool, startView string) Model {
	return newModel(project, cfg, root, loader, theme, keyMap, readOnly, startView, false)
}

func newModel(project *openspec.Project, cfg openspec.ProjectConfig, root string, loader *openspec.Loader, theme Theme, keyMap settings.KeyConfig, readOnly bool, startView string, singlePath bool) Model {
	m := Model{
		root:                  root,
		loader:                loader,
		project:               project,
		openSpec:              openspec.NewOSCLI(),
		lifecycle:             openspec.NewLifecycleService(root, openspec.NewOSCLI()),
		discoveryFingerprints: discoveryFingerprints(project.Changes),
		enrichedFingerprints:  make(map[string]string),
		pendingEnrichments:    make(map[string]string),
		enrichmentRetryAfter:  make(map[string]time.Time),
		artifactRenderCache:   make(map[artifactOutputIdentity]string),
		projectConfig:         cfg,
		theme:                 theme,
		keyMap:                keyMap,
		readOnly:              readOnly,
		singlePath:            singlePath,
		isGitRepo:             git.IsInsideWorkTree(root),
	}
	if m.isGitRepo {
		if r, err := git.WorkTreeRoot(root); err == nil {
			m.gitRoot = r
		} else {
			m.gitRoot = root
		}
	}
	m.pollGitStatus()
	if len(project.Changes) > 0 {
		m.selectDefaultArtifact()
		m.loadTaskItems()
	}
	if !singlePath && (startView == "index" || len(project.Changes) == 0) {
		var archiveErr error
		m.index.ArchiveChanges, archiveErr = loader.ListArchiveChangesFrom(root)
		if archiveErr != nil {
			m.errMsg = "error loading archive changes: " + archiveErr.Error()
		}
		var specErr error
		m.projectSpecs, specErr = loader.LoadProjectSpecsFrom(root)
		if specErr != nil {
			m.errMsg = "error loading project specs: " + specErr.Error()
		}
		m.index.ExpandedSpecs = make(map[int]bool)
		m.buildIndexItems()
		m.mode = ModeIndex
	}
	return m
}

func NewSinglePath(project *openspec.Project, cfg openspec.ProjectConfig, root string, loader *openspec.Loader, theme Theme, keyMap settings.KeyConfig, readOnly bool) Model {
	m := newModel(project, cfg, root, loader, theme, keyMap, readOnly, "change", true)
	if filepath.Base(filepath.Dir(filepath.Clean(root))) == "archive" && len(project.Changes) > 0 {
		m.index.ArchiveChanges = append([]openspec.Change(nil), project.Changes...)
		m.index.ArchiveCursor = 0
		m.mode = ModeViewingArchive
		m.selectDefaultArtifact()
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) }),
		func() tea.Msg { return startEnrichmentMsg{} },
	)
}

func (m Model) View() tea.View {
	if !m.vpReady {
		return tea.NewView("")
	}

	var content string
	if m.mode == ModeViewingConfig {
		content = m.viewContentWithChrome()
	} else if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		content = m.viewContentWithChrome()
	} else if len(m.project.Changes) == 0 && m.mode == ModeNormal {
		content = m.emptyViewContent()
	} else {
		content = m.mainViewContent()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	v.BackgroundColor = m.theme.ViewBg
	return v
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (m *Model) current() *openspec.Change {
	if m.mode == ModeViewingArchive {
		return m.currentArchive()
	}
	if m.mode == ModeIndex || m.mode == ModeViewingSpec {
		return nil
	}
	if len(m.project.Changes) == 0 {
		return nil
	}
	return &m.project.Changes[m.changeIdx]
}

func (m *Model) currentArchive() *openspec.Change {
	if m.index.ArchiveCursor < len(m.index.ArchiveChanges) {
		return &m.index.ArchiveChanges[m.index.ArchiveCursor]
	}
	return nil
}

const (
	chromeTop        = 1
	chromeHeader     = 1
	chromeTabBar     = 1
	chromeInnerSep   = 1
	chromeSpecSubnav = 1
	chromeHelpBar    = 1
	chromeBottom     = 1
)

func (m *Model) contentHeight() int {
	if m.mode == ModeIndex || m.mode == ModeViewingSpec || m.mode == ModeViewingConfig {
		h := m.height - (chromeTop + chromeHeader + chromeInnerSep + chromeInnerSep + chromeHelpBar + chromeBottom)
		if h < 1 {
			h = 1
		}
		return h
	}
	h := m.height - (chromeTop + chromeHeader + chromeTabBar + chromeInnerSep + chromeInnerSep + chromeHelpBar + chromeBottom)
	if m.hasSpecSubnav() {
		h -= chromeSpecSubnav
	}
	if h < 1 {
		h = 1
	}
	return h
}

// mergeReloadedChange updates the selected active change from dynamic artifact
// outputs while preserving the selected artifact/output identity when possible.
func (m *Model) mergeReloadedChange(fresh openspec.Change) (tasksChanged bool, viewportDirty bool) {
	ch := m.current()
	if ch == nil || m.mode != ModeNormal || m.changeIdx < 0 || m.changeIdx >= len(m.project.Changes) {
		return false, false
	}
	oldFingerprint := changeDiscoveryFingerprint(*ch)
	newFingerprint := changeDiscoveryFingerprint(fresh)
	if oldFingerprint == newFingerprint {
		return false, false
	}
	_, oldTasks, _ := artifactOutputByID(ch, "tasks")
	_, newTasks, _ := artifactOutputByID(&fresh, "tasks")
	oldTaskContent, newTaskContent := "", ""
	if oldTasks != nil {
		oldTaskContent = oldTasks.Content
	}
	if newTasks != nil {
		newTaskContent = newTasks.Content
	}
	tasksChanged = oldTaskContent != newTaskContent

	selectedBefore := m.artifactSelection
	oldArtifactIndex, oldOutputIndex := artifactSelectionPosition(ch, selectedBefore)
	m.project.Changes[m.changeIdx] = fresh
	m.invalidateArtifactRenderCache(fresh.Name)
	m.reconcileArtifactSelectionNear(oldArtifactIndex, oldOutputIndex)
	viewportDirty = !m.viewingCode && (selectedBefore != m.artifactSelection || oldFingerprint != newFingerprint)
	if tasksChanged {
		m.tasks.Items = openspec.ParseTasks(newTaskContent)
	}
	return tasksChanged, viewportDirty
}
