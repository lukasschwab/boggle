package main

import (
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/log"
	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/lukasschwab/boggle/pkg/dictionary"
	"github.com/lukasschwab/boggle/pkg/game"
)

const usage = `boggle - play an unsanctioned game of boggle.

This game uses the Collins Scrabble Words 2019 dictionary, a three-minute timer,
and requires words be at least four letters long.

Each game outputs a .boggle file describing the board and your performance. You
can replay a board by passing a .boggle file to this program:

    $ boggle -file past-game.boggle

Or by providing the "serialized" short-form description of the board, included
in the YAML frontmatter of each .boggle file:

    $ boggle -board Y3VkbnF1dG5kZHVybHllYXg=

The following options are available:

`

func main() {
	var filenameFlag = flag.String("file", "", ".boggle file to configure game")
	var boardFlag = flag.String("board", "", "serialized board string")
	var boardUrlFlag = flag.String("url", "", "web URL of a public .boggle file to configure game")
	var solveFlag = flag.Bool("solve", false, "print all possible words on board after game")
	var skipFlag = flag.Bool("skip", false, "skip interactive game")
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	var err error
	b := boggle.Shake()
	duration := 3 * time.Minute

	if boardUrlFlag != nil && len(*boardUrlFlag) != 0 {
		path, err := game.DownloadFile(*boardUrlFlag)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("Downloaded %v", path)
		filenameFlag = &path
	}

	if filenameFlag != nil && len(*filenameFlag) != 0 {
		frontmatter, _, err := game.LoadFile(*filenameFlag)
		if err != nil {
			log.Fatal(err.Error())
		}
		duration = time.Duration(frontmatter.TimerSeconds) * time.Second
		if b, err = boggle.Deserialize(frontmatter.Board); err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("Loaded %v", *filenameFlag)
	} else if boardFlag != nil && len(*boardFlag) != 0 {
		if b, err = boggle.Deserialize(*boardFlag); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Loading a trie makes the b.AllWords operation faster: prefix-based
	// exclusion.
	dict := dictionary.Filtered{
		Underlying: dictionary.EmptyTrie(),
		Filter:     boggle.Playable,
	}
	if err := dictionary.Load(dictionary.CSW19G, dict); err != nil {
		log.Fatal(err.Error())
	}

	boardWordsDict := b.AllWords(dict)

	// Game loop.
	if !*skipFlag {
		words, err := game.Run(boardWordsDict, b, duration)
		if err != nil {
			log.Fatal(err.Error())
		}

		filename := fmt.Sprintf("./%s.boggle", time.Now().UTC().Format(time.RFC3339))
		if err := game.WriteFile(filename, game.Frontmatter{
			Board:        b.Serialize(),
			TimerSeconds: 180,
		}, words); err != nil {
			log.Fatal(err.Error())
		}
		log.Infof("Wrote %v", filename)
	}

	if *solveFlag {
		allAvailableWords := boardWordsDict.Members()
		sort.Strings(allAvailableWords)
		for _, word := range allAvailableWords {
			fmt.Println(word)
		}
	}
}
