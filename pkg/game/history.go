package game

import "strings"

type history []entry

func (h history) append(ew entry) history {
	if len(h) == 4 {
		h = h[1:]
	}
	return append(h, ew)
}

func (h history) view() string {
	styled := make([]string, len(h))
	for i := range h {
		styled[i] = enteredWordStyles[h[i].outcome].Render(h[i].word)
	}
	return strings.Join(styled, "\n")
}

type entry struct {
	word    string
	outcome entryOutcome
}

type entryOutcome int

const (
	EntryInvalid entryOutcome = iota
	EntryRepeat  entryOutcome = iota
	EntryScored  entryOutcome = iota
)
