package main

import (
	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/counter"
)

var (
	DieCounts = dieCounts()
)

// histogramPlayable returns true if word is playable by letter counts, i.e. each
// letter in word occurs at least as many times on the Boggle die. Disregards
// inter-letter constraints: "quj" is histogramPlayable.
func histogramPlayable(word string) bool {
	return toCounts(word).LessThan(DieCounts)
}

func toCounts(word string) counter.Counter {
	lc := counter.Counter{}

	for len(word) > 0 {
		var head string
		head, word = boggle.HeadTail(word)
		lc.Incr(head)
	}

	return lc
}

func dieCounts() counter.Counter {
	counts := counter.Counter{}
	for _, d := range boggle.ClassicDice {
		for _, side := range d.Sides() {
			counts.Incr(side)
		}
	}
	return counts
}
