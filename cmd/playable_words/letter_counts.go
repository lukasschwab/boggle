package main

import (
	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/counter"
)

func histogramPlayableFilter(dice [16]boggle.Die) func(string) bool {
	counts := dieCounts(dice)
	return func(word string) bool {
		return toCounts(word).LessThan(counts)
	}
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

func dieCounts(dice [16]boggle.Die) counter.Counter {
	counts := counter.Counter{}
	for _, d := range dice {
		for _, side := range d.Sides() {
			counts.Incr(side)
		}
	}
	return counts
}
