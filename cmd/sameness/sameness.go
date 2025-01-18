package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

const (
	GamesPerSession = 10
)

func main() {
	classicUniqueCounts := []string{}
	newUniqueCounts := []string{}
	for _ = range 1_000 {
		classicUnique, _ := evaluate(boggle.ClassicDice)
		classicUniqueCounts = append(classicUniqueCounts, strconv.Itoa(classicUnique))
		newUnique, _ := evaluate(boggle.NewDice)
		newUniqueCounts = append(newUniqueCounts, strconv.Itoa(newUnique))
	}
	fmt.Println("CLASSIC")
	fmt.Println(strings.Join(classicUniqueCounts, "\n"))
	fmt.Println("NEW")
	fmt.Println(strings.Join(newUniqueCounts, "\n"))
}

func evaluate(dice [16]boggle.Die) (unique, common int) {
	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.Load(dictionary.CSW19, dict); err != nil {
		log.Fatal(err.Error())
	}

	wordCounts := map[string]int{}

	for _ = range GamesPerSession {
		b := boggle.Shake(boggle.ClassicDice)
		boardWordsDict := b.AllWords(dict)

		for _, word := range boardWordsDict.Members() {
			if _, ok := wordCounts[word]; ok {
				wordCounts[word]++
			} else {
				wordCounts[word] = 1
			}
		}
	}

	// fmt.Printf("COUNT UNIQUE WORDS: %v\n", len(wordCounts))

	commonWords := []string{}
	for word, count := range wordCounts {
		if count >= GamesPerSession/20 {
			commonWords = append(commonWords, word)
		}
	}
	// fmt.Printf("COUNT COMMON WORDS: %v\n", len(commonWords))

	return len(wordCounts), len(commonWords)
}
