package main

import (
	"fmt"
	"strings"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

func main() {
	unfiltered := dictionary.Map{}
	dictionary.Load(dictionary.CSW19, unfiltered)

	sizePlayable, _ := filter(unfiltered, boggle.Playable)
	// fmt.Printf("Excluded for size: %v\n", sizeExcluded)
	classicPlayable := evaluate(sizePlayable, boggle.ClassicDice)
	newPlayable := evaluate(sizePlayable, boggle.NewDice)

	onlyClassic := filterExclusive(classicPlayable, newPlayable)
	// fmt.Printf("Only classic-playable:\n %+v\n", onlyClassic)

	fucklike, _ := filterSlice(onlyClassic, func(word string) bool {
		return strings.ContainsRune(word, 'f') && strings.ContainsRune(word, 'k')
	})
	fmt.Printf("%d/%d Classic-exclusive words contain 'f' and 'k'", fucklike.Size(), len(onlyClassic))

	// onlyNew := filterExclusive(newPlayable, classicPlayable)
	// fmt.Printf("Only new-playable:\n %+v\n", onlyNew)
}

func evaluate(sizePlayable dictionary.Map, dice [16]boggle.Die) dictionary.Map {
	histogramPlayable, _ := filter(sizePlayable, histogramPlayableFilter(dice))
	// fmt.Printf("Additionally excluded by letter count: %v\n", histogramExcluded)

	dicePlayable, _ := filter(histogramPlayable, dicePlayableFilter(dice))
	// fmt.Printf("Additionally excluded by letter count: %v\n", diceExcluded)

	dicts := []dictionary.Map{sizePlayable, histogramPlayable, dicePlayable}
	sizes := make([]int, len(dicts))
	for i := range dicts {
		sizes[i] = dicts[i].Size()
	}

	fmt.Printf("sizes: %+v\n", sizes)

	return dicePlayable
}

func filter(dict dictionary.Map, filter func(string) bool) (dictionary.Map, []string) {
	return filterSlice(dict.Members(), filter)
}

func filterSlice(dict []string, filter func(string) bool) (dictionary.Map, []string) {
	rejects := []string{}
	result := dictionary.Map{}
	for _, word := range dict {
		if filter(word) {
			result.Add(word)
		} else {
			rejects = append(rejects, word)
		}
	}
	return result, rejects
}

func filterExclusive(primary, comp dictionary.Map) []string {
	_, exclusives := filter(primary, func(word string) bool {
		return comp.Contains(word)
	})
	return exclusives
}
