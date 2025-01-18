package main

import (
	"fmt"
	"log"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

// NOTE: simulated annealing performs much better here. See some prior art:
// https://www.danvk.org/wp/2009-02-19/sky-high-boggle-scores-with-simulated-annealing/index.html
func main() {
	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.Load(dictionary.CSW19, dict); err != nil {
		log.Fatal(err.Error())
	}

	for _ = range 10_000 {
		b := boggle.Shake(boggle.ClassicDice)
		boardWordsDict := b.AllWords(dict)
		fmt.Println(boardWordsDict.Size())
	}
}
