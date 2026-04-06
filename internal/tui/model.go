package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fishi0x01/vsh/internal/completer"
)

const maxVisibleSuggestions = 8

// CompleterFunc returns completion suggestions for the current input.
type CompleterFunc func(string) []completer.Suggestion

// PrefixFunc returns the current prompt prefix (e.g. "vault /secret/kv> ").
type PrefixFunc func() string

// Model is the bubbletea model for one prompt interaction.
// It reads a single command line, then quits — the caller loops.
type Model struct {
	input     textinput.Model
	completer CompleterFunc
	prefix    PrefixFunc

	// completion state — shown as hints only; Tab inserts, Enter always executes
	suggestions []completer.Suggestion
	selectedIdx int // index of the highlighted suggestion (-1 = none)

	// history passed in from the caller and returned on quit
	history    []string
	historyIdx int // == len(history) means "current (unsaved) input"
	savedInput string

	// result fields — read by the caller after Run()
	lastCommand string
	quitting    bool // true = Ctrl-C/Ctrl-D, skip execution
	done        bool // true whenever tea.Quit has been sent (clears the view)
}

// NewModel constructs the TUI model for one prompt cycle.
// history is the command history accumulated so far (may be nil on first run).
func NewModel(comp CompleterFunc, prefix PrefixFunc, history []string) Model {
	ti := textinput.New()
	ti.Prompt = "" // prefix is rendered separately
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 0

	if history == nil {
		history = []string{}
	}

	return Model{
		input:       ti,
		completer:   comp,
		prefix:      prefix,
		selectedIdx: -1,
		history:     history,
		historyIdx:  len(history),
	}
}

// LastCommand returns the command the user submitted (empty if they just quit).
func (m Model) LastCommand() string { return m.lastCommand }

// Quitting reports whether the user asked to exit (Ctrl-C / Ctrl-D).
func (m Model) Quitting() bool { return m.quitting }

// History returns the updated history slice to pass to the next NewModel call.
func (m Model) History() []string { return m.history }

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyCtrlD:
			m.quitting = true
			m.done = true
			return m, tea.Quit

		case tea.KeyEnter:
			// Always execute — never auto-accept completion on Enter.
			m.lastCommand = strings.TrimSpace(m.input.Value())
			if m.lastCommand != "" {
				m.history = append(m.history, m.lastCommand)
			}
			m.done = true
			return m, tea.Quit

		case tea.KeyTab:
			if len(m.suggestions) == 0 {
				m.suggestions = m.completer(m.input.Value())
				if len(m.suggestions) > 0 {
					m.selectedIdx = 0
				}
			} else {
				m.selectedIdx = (m.selectedIdx + 1) % len(m.suggestions)
			}
			if m.selectedIdx >= 0 {
				m = m.applyCompletion()
			}
			return m, nil

		case tea.KeyShiftTab:
			if len(m.suggestions) > 0 {
				m.selectedIdx = (m.selectedIdx - 1 + len(m.suggestions)) % len(m.suggestions)
				m = m.applyCompletion()
			}
			return m, nil

		case tea.KeyUp:
			if len(m.suggestions) > 0 {
				if m.selectedIdx > 0 {
					m.selectedIdx--
					m = m.applyCompletion()
				}
				return m, nil
			}
			// history navigation
			if len(m.history) == 0 {
				return m, nil
			}
			if m.historyIdx == len(m.history) {
				m.savedInput = m.input.Value()
			}
			if m.historyIdx > 0 {
				m.historyIdx--
			}
			m.input.SetValue(m.history[m.historyIdx])
			m.input.CursorEnd()
			return m, nil

		case tea.KeyDown:
			if len(m.suggestions) > 0 {
				if m.selectedIdx < len(m.suggestions)-1 {
					m.selectedIdx++
					m = m.applyCompletion()
				}
				return m, nil
			}
			// history navigation
			if m.historyIdx >= len(m.history) {
				return m, nil
			}
			m.historyIdx++
			if m.historyIdx == len(m.history) {
				m.input.SetValue(m.savedInput)
			} else {
				m.input.SetValue(m.history[m.historyIdx])
			}
			m.input.CursorEnd()
			return m, nil

		case tea.KeyEscape:
			m.suggestions = nil
			m.selectedIdx = -1
			return m, nil

		default:
			// Typing: update input and refresh the completion hint list.
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			m.historyIdx = len(m.history)
			m.suggestions = m.completer(m.input.Value())
			// Reset selection — user hasn't chosen anything yet.
			m.selectedIdx = -1
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.input.Width = msg.Width - len(m.prefix()) - 2
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// applyCompletion writes the currently selected suggestion into the input field.
// It keeps the prefix words intact and replaces only the last word.
func (m Model) applyCompletion() Model {
	if m.selectedIdx < 0 || len(m.suggestions) == 0 {
		return m
	}
	text := m.suggestions[m.selectedIdx].Text
	current := m.input.Value()
	idx := strings.LastIndex(current, " ")
	var newVal string
	if idx >= 0 {
		newVal = current[:idx+1] + text
	} else {
		newVal = text
	}
	m.input.SetValue(newVal)
	m.input.CursorEnd()
	return m
}

func (m Model) View() string {
	if m.done {
		return ""
	}

	prefix := m.prefix()
	promptLine := prefixStyle.Render(prefix) + m.input.View()

	if len(m.suggestions) == 0 {
		return promptLine
	}

	// Determine visible window of suggestions.
	start := 0
	if m.selectedIdx >= maxVisibleSuggestions {
		start = m.selectedIdx - maxVisibleSuggestions + 1
	}
	end := start + maxVisibleSuggestions
	if end > len(m.suggestions) {
		end = len(m.suggestions)
	}
	visible := m.suggestions[start:end]

	var rows []string
	for i, s := range visible {
		absIdx := i + start
		selected := absIdx == m.selectedIdx

		var nameStyle, dStyle lipgloss.Style
		indicator := "  "
		if selected {
			nameStyle = completionSelectedStyle
			dStyle = completionSelectedDescStyle
			indicator = "▸ "
		} else {
			nameStyle = completionNormalStyle
			dStyle = completionDescStyle
		}

		// Width() pads/truncates correctly even when the string contains ANSI codes.
		name := nameStyle.Width(completionNameWidth).Render(s.Text)
		desc := dStyle.Render(s.Description)
		rows = append(rows, indicator+name+" "+desc)
	}

	dropdown := completionBorderStyle.Render(strings.Join(rows, "\n"))
	return promptLine + "\n" + dropdown
}
