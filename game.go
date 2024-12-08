package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lukasschwab/boggle/pkg/dictionary"
)

func main() {
	b := randomBoard()
	fmt.Printf("%+v\n", b.fields)
	fmt.Println(b.pretty())

	serialized := b.serialize()
	fmt.Println(serialized)

	after, err := deserialize(serialized)
	if err != nil {
		panic(err)
	}
	// fmt.Println(after.pretty())

	fmt.Printf("Match? %v\n", after.pretty() == b.pretty())

	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     dictionary.Playable,
	}
	if err := dictionary.LoadFromFile(dict); err != nil {
		panic(err)
	}
	fmt.Printf("Contains 'test'? %v\n", dict.Contains("test"))

	boardWordsDict := b.AllWords(dict)

	// Print some stats.
	boardWordsSlice := boardWordsDict.Members()
	sort.Slice(boardWordsSlice, func(i, j int) bool {
		// Sort by length, decreasing.
		return len(boardWordsSlice[i]) > len(boardWordsSlice[j])
	})
	fmt.Printf("Word count in board: %v\n", len(boardWordsSlice))
	if len(boardWordsSlice) > 10 {
		fmt.Printf("%+v\n", boardWordsSlice[:10])
	} else {
		fmt.Printf("%+v\n", boardWordsSlice)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		candidate := strings.ToLower(text[:len(text)-1]) // Trailing newline.

		if boardWordsDict.Contains(candidate) {
			fmt.Printf("%v is playable\n", candidate)
		} else if dict.Contains(candidate) {
			fmt.Printf("%v is unplayable\n", candidate)
		} else {
			fmt.Printf("%v is not a word\n", candidate)
		}
	}
}
