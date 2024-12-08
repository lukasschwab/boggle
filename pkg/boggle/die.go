package boggle

import (
	"math/rand"
)

const (
	dieSides   = 6
	quLigature = "qu"
)

type die [dieSides]string

func (d die) roll() string {
	return d[rand.Intn(dieSides)]
}

// https://boardgames.stackexchange.com/questions/29264/boggle-what-is-the-dice-configuration-for-boggle-in-various-languages
var (
	die0  = die{"r", "i", "f", "o", "b", "x"}
	die1  = die{"i", "f", "e", "h", "e", "y"}
	die2  = die{"d", "e", "n", "o", "w", "s"}
	die3  = die{"u", "t", "o", "k", "n", "d"}
	die4  = die{"h", "m", "s", "r", "a", "o"}
	die5  = die{"l", "u", "p", "e", "t", "s"}
	die6  = die{"a", "c", "i", "t", "o", "a"}
	die7  = die{"y", "l", "g", "k", "u", "e"}
	die8  = die{quLigature, "b", "m", "j", "o", "a"}
	die9  = die{"e", "h", "i", "s", "p", "n"}
	die10 = die{"v", "e", "t", "i", "g", "n"}
	die11 = die{"b", "a", "l", "i", "y", "t"}
	die12 = die{"e", "z", "a", "v", "n", "d"}
	die13 = die{"r", "a", "l", "e", "s", "c"}
	die14 = die{"u", "w", "i", "l", "r", "g"}
	die15 = die{"p", "a", "c", "e", "m", "d"}

	allDice = []die{die0, die1, die2, die3, die4, die5, die6, die7, die8, die9, die10, die11, die12, die13, die14, die15}
)
