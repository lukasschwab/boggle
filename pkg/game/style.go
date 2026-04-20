package game

import (
	"time"

	"charm.land/bubbles/v2/timer"
	"charm.land/lipgloss/v2"
)

var DefaultStyle = NewStyle()

type Style struct {
	blurredStyle      lipgloss.Style
	boardStyle        lipgloss.Style
	promptStyle       lipgloss.Style
	historyStyle      lipgloss.Style
	enteredWordStyles map[entryOutcome]lipgloss.Style
	timerStyle        func(t timer.Model) lipgloss.Style
}

func NewStyle() Style {
	return Style{
		blurredStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		boardStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			MarginRight(2),
		promptStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")),
		historyStyle: lipgloss.NewStyle().
			PaddingTop(1).PaddingBottom(1),
		enteredWordStyles: map[entryOutcome]lipgloss.Style{
			EntryRepeat: lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")),
			EntryScored: lipgloss.NewStyle().
				Foreground(lipgloss.ANSIColor(2)),
			EntryInvalid: lipgloss.NewStyle().
				Foreground(lipgloss.ANSIColor(1)),
		},
		timerStyle: func(t timer.Model) lipgloss.Style {
			if t.Timeout <= 15*time.Second {
				return lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(11))
			}
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))
		},
	}
}
