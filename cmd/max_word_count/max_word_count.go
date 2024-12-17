package main

import (
	"fmt"
	"log"

	"github.com/lukasschwab/boggle/pkg/boggle"
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

	var bestBoard boggle.Board
	var bestBoardCount = 0

	for _ = range 100_000 {
		b := boggle.Shake(boggle.ClassicDice)
		boardWordsDict := b.AllWords(dict)
		// fmt.Println(boardWordsDict.Size())

		if boardWordsDict.Size() > bestBoardCount {
			bestBoardCount = boardWordsDict.Size()
			bestBoard = b
		}
	}

	fmt.Printf("%d words\n", bestBoardCount)
	fmt.Println(bestBoard.Pretty())
}
