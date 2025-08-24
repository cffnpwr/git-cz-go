package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cffnpwr/git-cz-go/config"
	"github.com/cffnpwr/git-cz-go/internal/interface/repo"
	"github.com/cffnpwr/git-cz-go/pkg/component/confirm"
	"github.com/cffnpwr/git-cz-go/pkg/component/selector"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultScopePrompt           = "Enter scope (optional)"
	defaultSubjectPrompt         = "Enter commit subject"
	defaultBodyPrompt            = "Enter commit body (optional)"
	defaultConfirmPrompt         = "Commit this message?"
	defaultTypeSelectDisplaySize = 5
	defaultPromptSeparator       = ": "
	defaultIconCharQuestion      = "?"
	defaultIconCharEntered       = "✓"
)

var (
	quitKey = key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("Ctrl+C", "quit"),
	)
	enterKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "confirm"),
	)
	submitKey = key.NewBinding(
		key.WithKeys("alt+enter", "ctrl+enter"),
		key.WithHelp("⌘+Enter/Ctrl+Enter", "submit changes"),
	)

	fgColor = lipgloss.Color("#bb9af7")

	defaultPromptStyle               = lipgloss.NewStyle().Bold(true)
	defaultIconStyle                 = lipgloss.NewStyle().Foreground(fgColor).Bold(true).MarginRight(1)
	defaultCommitMessagePreviewStyle = lipgloss.NewStyle().Bold(true).Padding(0, 2).Margin(0, 1).Border(lipgloss.RoundedBorder())
)

// Stage represents the current stage of the commit message creation process
type Stage string

const (
	StageTypeSelect   Stage = "type_select"
	StageScope        Stage = "scope"
	StageTicketNumber Stage = "ticket_number"
	StageSubject      Stage = "subject"
	StageBody         Stage = "body"
	StageBreaking     Stage = "breaking"
	StageFooter       Stage = "footer"
	StageConfirm      Stage = "confirm"
)

// CommitData holds all the data collected from the user for generating commit message
type CommitData struct {
	Type            string // Selected commit type (feat, fix, etc.)
	Scope           string // Optional scope (api, ui, etc.)
	TicketNumber    string // Ticket number with prefix
	Subject         string // Commit message subject
	Body            string // Commit message body (multi-line)
	BreakingChanges string // Breaking changes description
	Footer          string // Footer information (validated format)
	IsBreaking      bool   // Whether there are breaking changes
}

// GenerateCommitMessage generates a conventional commit message from the collected data
func (cd CommitData) GenerateCommitMessage() string {
	var parts []string

	// Header: <type>(<scope>): <ticket_number> <subject>
	header := cd.Type
	if cd.Scope != "" {
		header += "(" + cd.Scope + ")"
	}
	if cd.IsBreaking {
		header += "!"
	}
	header += ":"

	if cd.TicketNumber != "" {
		header += " " + cd.TicketNumber
	}
	header += " " + cd.Subject

	parts = append(parts, header)

	// Body with blank line
	if cd.Body != "" {
		parts = append(parts, "", cd.Body)
	}

	// Footers with blank line
	hasFooter := false

	// BREAKING CHANGE exception
	if cd.BreakingChanges != "" {
		if !hasFooter {
			parts = append(parts, "")
		}
		parts = append(parts, "BREAKING CHANGE: "+cd.BreakingChanges)
		hasFooter = true
	}

	if cd.Footer != "" {
		if !hasFooter {
			parts = append(parts, "")
		}
		parts = append(parts, cd.Footer)
	}

	return strings.Join(parts, "\n")
}

var _ tea.Model = Model{}

// Model represents the main model for the git cz application
type Model struct {
	config       *config.Config
	gitRepo      repo.GitRepository
	currentStage Stage

	// Individual models
	typeSelect   selector.Model
	scopeInput   textinput.Model
	ticketNumber TicketNumberModel
	subjectInput textinput.Model
	bodyInput    textarea.Model
	breaking     BreakingChangesModel
	footer       FooterModel
	confirm      confirm.Model

	// Data collection
	commitData CommitData
}

// NewModel creates a new main model for git cz
func NewModel(cfg *config.Config, gitRepo repo.GitRepository) (Model, error) {
	// Initialize type select model
	size := min(len(cfg.Types), defaultTypeSelectDisplaySize)
	selectItems := make([]selector.SelectItem, len(cfg.Types))
	for i, t := range cfg.Types {
		selectItems[i] = t
	}
	typeSelect, err := selector.New(selectItems, size)
	if err != nil {
		return Model{}, err
	}
	typeSelect = typeSelect.SetCyclic(true).SetShowSelectedItem(true)
	if cfg.Messages.Type != "" {
		typeSelect.Prompt = cfg.Messages.Type
	}

	// Initialize scope input
	scopeInput := textinput.New()
	scopePrompt := defaultScopePrompt
	if cfg.Messages.Scope != "" {
		scopePrompt = cfg.Messages.Scope
	}
	scopeInput.Prompt = scopePrompt + defaultPromptSeparator
	scopeInput.PromptStyle = defaultPromptStyle

	// Initialize ticket number model
	ticketNumber := NewTicketNumberModel(cfg.Messages.TicketNumber, cfg.TicketNumber, gitRepo)

	// Initialize subject input
	subjectInput := textinput.New()
	subjectPrompt := defaultSubjectPrompt
	if cfg.Messages.Subject != "" {
		subjectPrompt = cfg.Messages.Subject
	}
	subjectInput.Prompt = subjectPrompt + defaultPromptSeparator
	subjectInput.PromptStyle = defaultPromptStyle

	// Initialize body input
	bodyInput := textarea.New()
	bodyPrompt := defaultBodyPrompt
	if cfg.Messages.Body != "" {
		bodyPrompt = cfg.Messages.Body
	}
	bodyInput.Prompt = bodyPrompt + defaultPromptSeparator
	bodyInput.FocusedStyle = textarea.Style{
		Prompt: defaultPromptStyle,
	}

	// Initialize breaking changes model
	breaking := NewBreakingChangesModel(cfg.Messages.BreakingConfirm, cfg.Messages.BreakingMessage)

	// Initialize footer model
	footerModel := NewFooterModel(cfg.Messages.Footer)

	// Initialize confirm model
	confirmModel := confirm.New()
	confirmPrompt := defaultConfirmPrompt
	if cfg.Messages.ConfirmCommit != "" {
		confirmPrompt = cfg.Messages.ConfirmCommit
	}
	confirmModel.Prompt = confirmPrompt

	model := Model{
		config:       cfg,
		gitRepo:      gitRepo,
		currentStage: StageTypeSelect,
		typeSelect:   typeSelect,
		scopeInput:   scopeInput,
		ticketNumber: ticketNumber,
		subjectInput: subjectInput,
		bodyInput:    bodyInput,
		breaking:     breaking,
		footer:       footerModel,
		confirm:      confirmModel,
	}

	return model, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.typeSelect.Init(),
		m.ticketNumber.Init(),
		m.breaking.Init(),
		m.footer.Init(),
		m.confirm.Init(),
		textarea.Blink,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, quitKey) {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	isFinished := false
	switch m.currentStage {
	case StageTypeSelect:
		m.typeSelect, cmd = m.typeSelect.Update(msg)
		isFinished = m.typeSelect.IsSelected()
	case StageScope:
		isFinished, m.scopeInput, cmd = handleTextInput(m.scopeInput, msg)
	case StageTicketNumber:
		m.ticketNumber, cmd = m.ticketNumber.Update(msg)
		isFinished = m.ticketNumber.IsFinished()
	case StageSubject:
		isFinished, m.subjectInput, cmd = handleTextInput(m.subjectInput, msg)
	case StageBody:
		if msg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(msg, submitKey) {
				isFinished = true
				break
			}
		}
		m.bodyInput, cmd = m.bodyInput.Update(msg)
	case StageBreaking:
		m.breaking, cmd = m.breaking.Update(msg)
		isFinished = m.breaking.IsFinished()
	case StageFooter:
		m.footer, cmd = m.footer.Update(msg)
		isFinished = m.footer.IsFinished()
	case StageConfirm:
		m.confirm, cmd = m.confirm.Update(msg)
		isFinished = m.confirm.IsConfirmed()
	}
	if isFinished {
		// Store data before moving to next stage
		switch m.currentStage {
		case StageTypeSelect:
			m.commitData.Type = m.typeSelect.GetSelectedItem().(config.TypeValue).Value
		case StageScope:
			m.commitData.Scope = m.scopeInput.Value()
		case StageTicketNumber:
			m.commitData.TicketNumber = m.ticketNumber.GetValue()
		case StageSubject:
			m.commitData.Subject = m.subjectInput.Value()
		case StageBody:
			m.commitData.Body = m.bodyInput.Value()
		case StageBreaking:
			m.commitData.BreakingChanges = m.breaking.GetValue()
			m.commitData.IsBreaking = m.breaking.HasBreakingChanges()
		case StageFooter:
			m.commitData.Footer = m.footer.GetValue()
		case StageConfirm:
			if m.confirm.GetValue() {
				// User confirmed - generate and commit the message
				commitMsg := m.commitData.GenerateCommitMessage()
				if err := m.gitRepo.Commit(commitMsg); err != nil {
					fmt.Printf("Error committing changes: %v\n", err)
					return m, tea.Quit
				}
				return m, tea.Quit
			} else {
				// User declined
				return m, tea.Quit
			}
		}

		m.currentStage = m.nextStage(m.currentStage)
		switch m.currentStage {
		case StageScope:
			m.scopeInput.Focus()
		case StageTicketNumber:
			m.ticketNumber.Focus()
		case StageSubject:
			m.subjectInput.Focus()
		case StageBody:
			m.bodyInput.Focus()
		case StageBreaking:
			m.breaking.Focus()
		case StageFooter:
			m.footer.Focus()
		}
	}

	return m, cmd
}

func (m Model) View() string {
	var sections []string
	// Progress display
	sections = append(sections, m.buildProgressView())
	// Current stage view
	sections = append(sections, defaultIconStyle.Render(defaultIconCharQuestion)+m.getStageView(m.currentStage))

	return strings.Join(sections, "\n") + "\n"
}

func (m Model) buildProgressView() string {
	var sections []string
	s := StageTypeSelect
	for {
		if s == m.currentStage {
			break
		}

		sections = append(sections, defaultIconStyle.Render(defaultIconCharEntered)+m.getStageView(s))
		s = m.nextStage(s)
	}
	return strings.Join(sections, "\n")
}

func (m Model) getStageView(stage Stage) string {
	switch stage {
	case StageTypeSelect:
		return m.typeSelect.View()
	case StageScope:
		return m.scopeInput.View()
	case StageTicketNumber:
		return m.ticketNumber.View()
	case StageSubject:
		return m.subjectInput.View()
	case StageBody:
		return m.bodyInput.View()
	case StageBreaking:
		return m.breaking.View()
	case StageFooter:
		return m.footer.View()
	case StageConfirm:
		commitMessagePreview := m.commitData.GenerateCommitMessage()
		m.confirm.Prompt += "\n" + defaultCommitMessagePreviewStyle.Render(commitMessagePreview)
		return m.confirm.View()
	default:
		return ""
	}
}

func (m Model) nextStage(stage Stage) Stage {
	switch stage {
	case StageTypeSelect:
		if !slices.Contains(m.config.SkipQuestions, string(StageScope)) {
			return StageScope
		}
		fallthrough
	case StageScope:
		if m.config.TicketNumber.Enable {
			return StageTicketNumber
		}
		fallthrough
	case StageTicketNumber:
		return StageSubject
	case StageSubject:
		if !slices.Contains(m.config.SkipQuestions, string(StageBody)) {
			return StageBody
		}
		fallthrough
	case StageBody:
		if !slices.Contains(m.config.SkipQuestions, string(StageBreaking)) {
			return StageBreaking
		}
		fallthrough
	case StageBreaking:
		if !slices.Contains(m.config.SkipQuestions, string(StageFooter)) {
			return StageFooter
		}
	}
	return StageConfirm
}

func handleTextInput(m textinput.Model, msg tea.Msg) (bool, textinput.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, enterKey) {
			return true, m, nil
		}
	}

	var cmd tea.Cmd
	m, cmd = m.Update(msg)
	return false, m, cmd
}
