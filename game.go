package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

// TODO: introduce a time limit.
// TODO: parameterize a dictionary file (optional). Make a good one avail.

func main() {
	b := boggle.Shake()
	fmt.Println(b.Pretty())

	serialized := b.Serialize()
	fmt.Println(serialized)

	after, err := boggle.Deserialize(serialized)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Match? %v\n", after.Pretty() == b.Pretty())

	// Loading a trie makes the b.AllWords operation faster: prefix-based
	// exclusion.
	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.LoadFromFile(dict); err != nil {
		panic(err)
	}
	fmt.Printf("Contains 'test'? %v\n", dict.Contains("test"))

	boardWordsDict := b.AllWords(dict)

	// Print some stats. Sort by length, decreasing.
	boardWordsSlice := boardWordsDict.Members()
	sort.Slice(boardWordsSlice, func(i, j int) bool {
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
