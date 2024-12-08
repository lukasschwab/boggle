package boggle

import (
	"math/rand"
	"time"
)

func Playable(s string) bool {
	return len(s) >= 4
}

func Shake() board {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allDice), func(i, j int) {
		allDice[i], allDice[j] = allDice[j], allDice[i]
	})

	result := board{}
	for rowIdx, row := range result.fields {
		for colIdx := range row {
			result.fields[rowIdx][colIdx] = allDice[indexToLinear(rowIdx, colIdx)].roll()
		}
	}
	return result
}
