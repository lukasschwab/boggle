package main

import "github.com/lukasschwab/boggle/pkg/boggle"

func dicePlayableFilter(dice [16]boggle.Die) func(string) bool {
	return func(word string) bool {
		return dicePlayableTail(word, dice[:])
	}
}

func dicePlayableTail(subword string, diceLeft []boggle.Die) bool {
	if len(subword) == 0 {
		return true
	}

	head, tail := boggle.HeadTail(subword)

	for i := range diceLeft {
		if diceLeft[i].Contains(head) {
			// Copy to list without i
			removed := diceLeft[i]
			diceLeft[i] = boggle.NilDie

			if dicePlayableTail(tail, diceLeft) {
				diceLeft[i] = removed
				return true
			}

			diceLeft[i] = removed
		}
	}

	return false
}
