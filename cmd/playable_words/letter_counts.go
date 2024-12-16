package main

import "github.com/lukasschwab/boggle/pkg/boggle"

var (
	DieCounts = dieCounts()
)

// histogramPlayable returns true if word is playable by letter counts, i.e. each
// letter in word occurs at least as many times on the Boggle die. Disregards
// inter-letter constraints: "quj" is histogramPlayable.
func histogramPlayable(word string) bool {
	return toCounts(word).LessThan(DieCounts)
}

type letterCounts map[string]int

func (lc letterCounts) Incr(head string) {
	if _, ok := lc[head]; !ok {
		lc[head] = 0
	}
	lc[head]++
}

func (lc letterCounts) LessThan(other letterCounts) bool {
	for letter, count := range lc {
		otherCount, ok := other[letter]
		if !ok || count > otherCount {
			return false
		}
	}
	return true
}

func toCounts(word string) letterCounts {
	lc := letterCounts{}

	for len(word) > 0 {
		var head string
		head, word = boggle.HeadTail(word)
		lc.Incr(head)
	}

	return lc
}

func dieCounts() letterCounts {
	counts := letterCounts{}
	for _, d := range boggle.ClassicDice {
		for _, side := range d.Sides() {
			counts.Incr(side)
		}
	}
	return counts
}
