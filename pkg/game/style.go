package game

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultStyle = NewStyle(lipgloss.DefaultRenderer())
)

type Style struct {
	blurredStyle      lipgloss.Style
	boardStyle        lipgloss.Style
	promptStyle       lipgloss.Style
	historyStyle      lipgloss.Style
	enteredWordStyles map[entryOutcome]lipgloss.Style
	timerStyle        func(t timer.Model) lipgloss.Style
}

func NewStyle(r *lipgloss.Renderer) Style {
	return Style{
		blurredStyle: r.NewStyle().
			Foreground(lipgloss.Color("240")),
		boardStyle: r.NewStyle().
			PaddingLeft(2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			MarginRight(2),
		promptStyle: r.NewStyle().
			Foreground(lipgloss.Color("63")),
		historyStyle: r.NewStyle().
			PaddingTop(1).PaddingBottom(1),
		enteredWordStyles: map[entryOutcome]lipgloss.Style{
			EntryRepeat: r.NewStyle().
				Foreground(lipgloss.Color("240")),
			EntryScored: r.NewStyle().
				Foreground(lipgloss.ANSIColor(2)),
			EntryInvalid: r.NewStyle().
				Foreground(lipgloss.ANSIColor(1)),
		},
		timerStyle: func(t timer.Model) lipgloss.Style {
			if t.Timeout <= 15*time.Second {
				return r.NewStyle().Foreground(lipgloss.ANSIColor(11))
			}
			return r.NewStyle().
				Foreground(lipgloss.Color("240"))
		},
	}
}
