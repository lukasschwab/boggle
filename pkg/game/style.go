package game

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/charmbracelet/lipgloss"
)

var (
	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	boardStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			MarginRight(2)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63"))

	historyStyle = lipgloss.NewStyle().
			PaddingTop(1).PaddingBottom(1)

	enteredWordStyles = map[entryOutcome]lipgloss.Style{
		EntryRepeat:  blurredStyle,
		EntryScored:  lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(2)),
		EntryInvalid: lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1)),
	}
)

func timerStyle(t timer.Model) lipgloss.Style {
	if t.Timeout <= 15*time.Second {
		return lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(11))
	}
	return blurredStyle
}
