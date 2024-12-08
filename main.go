package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
	"github.com/lukasschwab/boggle/pkg/game"
)

// TODO: introduce a time limit.
// TODO: parameterize a dictionary file (optional). Make a good one avail.
// TODO: test in CI.

func main() {
	var filenameFlag = flag.String("f", "", ".boggle file to configure game")
	var boardFlag = flag.String("b", "", "serialized board string")
	var boardUrlFlag = flag.String("url", "", "web URL of a public .boggle file to configure game")
	flag.Parse()

	var err error
	b := boggle.Shake()
	duration := 3 * time.Minute

	if boardUrlFlag != nil && len(*boardUrlFlag) != 0 {
		path, err := game.DownloadFile(*boardUrlFlag)
		if err != nil {
			panic(err)
		}
		log.Printf("Downloaded %v", path)
		filenameFlag = &path
	}

	if filenameFlag != nil && len(*filenameFlag) != 0 {
		frontmatter, _, err := game.LoadFile(*filenameFlag)
		if err != nil {
			panic(err)
		}
		duration = time.Duration(frontmatter.TimerSeconds) * time.Second
		if b, err = boggle.Deserialize(frontmatter.Board); err != nil {
			panic(err)
		}
		log.Printf("Loaded %v", *filenameFlag)
	} else if boardFlag != nil && len(*boardFlag) != 0 {
		if b, err = boggle.Deserialize(*boardFlag); err != nil {
			panic(err)
		}
	}
	// else if boardUrlFlag != nil && len(*boardUrlFlag) != 0 {
	// 	if f, _, err := game.DownloadFile(*boardUrlFlag); err != nil {
	// 		panic(err)
	// 	} else if b, err = boggle.Deserialize(f.Board); err != nil {
	// 		panic(err)
	// 	}
	// }

	// Loading a trie makes the b.AllWords operation faster: prefix-based
	// exclusion.
	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.LoadFromFile(dict); err != nil {
		panic(err)
	}

	boardWordsDict := b.AllWords(dict)
	words, err := game.Run(boardWordsDict, b, duration)
	if err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("./%s.boggle", b.Serialize())
	game.WriteFile(filename, game.Frontmatter{
		Board:        b.Serialize(),
		TimerSeconds: 180,
	}, words)
	log.Printf("Wrote %v", filename)
}
