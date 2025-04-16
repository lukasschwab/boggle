package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

// This file provides a charmbracelet/bubbletea UI for boggle. Like all my
// bubbletea UIs, it's a rat's nest.

type model struct {
	Dict           dictionary.Interface
	Board          boggle.Board
	totalWordCount int

	// Controlled by the user
	userInput   textinput.Model
	scoredWords map[string]bool
	history     history

	// Misc. UI
	quitting bool
	keymap   keymap
	help     help.Model
	timer    timer.Model
	style    Style
}

type keymap struct {
	submit key.Binding
	quit   key.Binding
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.timer.Init(),
		textinput.Blink,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		return m.quit()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m.quit()
		case key.Matches(msg, m.keymap.submit):
			return m.handleSubmission(), nil
		}
	}

	var cmd tea.Cmd
	m.userInput, cmd = m.userInput.Update(msg)
	return m, cmd
}

func (m model) quit() (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

func (m model) handleSubmission() model {
	input := strings.TrimSpace(m.userInput.Value())
	ew := entry{
		word:    input,
		outcome: EntryScored,
	}

	if m.Dict.Contains(input) {
		if _, seen := m.scoredWords[input]; seen {
			ew.outcome = EntryRepeat
		}
		m.scoredWords[input] = true
	} else {
		ew.outcome = EntryInvalid
	}

	m.userInput.Reset()
	m.history = m.history.append(ew)
	return m
}

func (m model) View() string {
	builder := new(strings.Builder)

	scoreView := fmt.Sprintf("%d", len(m.scoredWords)) + m.style.blurredStyle.Render(fmt.Sprintf("/%d", m.totalWordCount))
	timerView := m.style.timerStyle(m.timer).Render(m.timer.View())
	serializedBoardView := m.style.blurredStyle.Render(m.Board.Serialize())
	builder.WriteString(strings.Join([]string{scoreView, timerView, serializedBoardView}, m.style.blurredStyle.Render(" • ")))
	builder.WriteRune('\n')

	builder.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		m.style.boardStyle.Render(strings.TrimSpace(m.Board.Pretty())),
		m.style.historyStyle.Render(m.history.view(m.style)),
	) + "\n")

	if !m.quitting {
		builder.WriteString(m.userInput.View())
		builder.WriteString(m.helpView())
	}

	return builder.String()
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.submit,
		m.keymap.quit,
	})
}

// BODGE: should really take dictionary.Interface, but dictionary.Map exposes a
// word count.
func Run(
	dict dictionary.Map,
	board boggle.Board,
	duration time.Duration,
) ([]string, error) {
	final, err := tea.NewProgram(Model(dict, board, duration, DefaultStyle)).Run()

	if err != nil {
		return []string{}, fmt.Errorf("error running game: %w", err)
	}

	finalModel := final.(model)
	scoredWords := make([]string, 0, len(finalModel.scoredWords))
	for word := range finalModel.scoredWords {
		scoredWords = append(scoredWords, word)
	}
	return scoredWords, nil
}

func Model(
	dict dictionary.Map,
	board boggle.Board,
	duration time.Duration,
	style Style,
) tea.Model {
	ti := textinput.New()
	ti.Focus()
	ti.PromptStyle = style.promptStyle

	return model{
		Dict:           dict,
		Board:          board,
		totalWordCount: len(dict.Members()),
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "esc"),
				key.WithHelp("[ctrl+c]", "quit"),
			),
			submit: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("[enter]", "submit word"),
			),
		},
		scoredWords: make(map[string]bool),
		timer:       timer.New(duration),
		userInput:   ti,
		help:        help.New(),
		style:       style,
	}
}
