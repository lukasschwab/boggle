package boggle

import (
	"math/rand"
)

func Playable(s string) bool {
	return len(s) >= 4 && len(s) <= 16
}

func Shake(dice [16]Die) Board {
	rand.Shuffle(len(dice), func(i, j int) {
		dice[i], dice[j] = dice[j], dice[i]
	})
	result := Board{}
	for rowIdx, row := range result.fields {
		for colIdx := range row {
			result.fields[rowIdx][colIdx] = dice[indexToLinear(rowIdx, colIdx)].roll()
		}
	}
	return result
}
