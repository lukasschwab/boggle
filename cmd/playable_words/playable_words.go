package main

import (
	"fmt"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

func main() {

	unfiltered := dictionary.Map{}
	dictionary.Load(dictionary.CSW19G, unfiltered)

	sizePlayable, sizeExcluded := filter(unfiltered, boggle.Playable)
	fmt.Printf("Excluded for size: %v\n", sizeExcluded)

	histogramPlayable, histogramExcluded := filter(sizePlayable, histogramPlayable)
	fmt.Printf("Additionally excluded by letter count: %v\n", histogramExcluded)

	dicePlayable, diceExcluded := filter(histogramPlayable, dicePlayable)
	fmt.Printf("Additionally excluded by letter count: %v\n", diceExcluded)

	dicts := []dictionary.Map{unfiltered, sizePlayable, histogramPlayable, dicePlayable}
	sizes := make([]int, len(dicts))
	for i := range dicts {
		sizes[i] = dicts[i].Size()
	}

	fmt.Printf("sizes: %+v\n", sizes)

	// Print the final set of playable words.
	// for _, word := range dicePlayable.Members() {
	// 	fmt.Println(word)
	// }
}

func filter(dict dictionary.Map, filter func(string) bool) (dictionary.Map, []string) {
	rejects := []string{}
	result := dictionary.Map{}
	for _, word := range dict.Members() {
		if filter(word) {
			result.Add(word)
		} else {
			rejects = append(rejects, word)
		}
	}
	return result, rejects
}
