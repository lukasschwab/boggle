package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/counter"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

func main() {
	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.Load(dictionary.CSW19G, dict); err != nil {
		log.Fatal(err.Error())
	}

	wordCounter := counter.Counter{}

	for range 10_000 {
		b := boggle.Shake()
		boardWordsDict := b.AllWords(dict)

		for _, word := range boardWordsDict.Members() {
			wordCounter.Incr(word)
		}
	}

	wordCounts := wordCounter.Counts()
	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].Count > wordCounts[j].Count
	})

	for i := range 1000 {
		fmt.Printf("%s\t%d\n", wordCounts[i].Word, wordCounts[i].Count)
	}
}
