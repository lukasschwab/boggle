package boggle

import (
	"math/rand"
)

type Die interface {
	roll() string
	Contains(string) bool
	Sides() [dieSides]string
}

const (
	dieSides   = 6
	quLigature = "qu"
)

type die [dieSides]string

func (d die) roll() string {
	return d[rand.Intn(dieSides)]
}

func (d die) Contains(head string) bool {
	for _, side := range d {
		if side == head {
			return true
		}
	}
	return false
}

func (d die) Sides() [dieSides]string {
	return d
}

// "Classic" scrabble dice. See
// - https://boardgames.stackexchange.com/questions/29264/boggle-what-is-the-dice-configuration-for-boggle-in-various-languages
// - https://www.bananagrammer.com/2013/10/the-boggle-cube-redesign-and-its-effect.html
var (
	// NilDie is a special die for experiments
	NilDie = die{}

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

	ClassicDice = [16]Die{die0, die1, die2, die3, die4, die5, die6, die7, die8, die9, die10, die11, die12, die13, die14, die15}
)

// "New" scrabble dice. See https://www.bananagrammer.com/2013/10/the-boggle-cube-redesign-and-its-effect.html
// Separate types to prevent mixing.
type newDie die

func (d newDie) roll() string {
	return die(d).roll()
}

func (d newDie) Contains(head string) bool {
	return die(d).Contains(head)
}

func (d newDie) Sides() [dieSides]string {
	return die(d).Sides()
}

// "New" scrabble dice. See https://www.bananagrammer.com/2013/10/the-boggle-cube-redesign-and-its-effect.html
var (
	newDie0  = newDie{"a", "a", "e", "e", "g", "n"}
	newDie1  = newDie{"a", "b", "b", "j", "o", "o"}
	newDie2  = newDie{"a", "c", "h", "o", "p", "s"}
	newDie3  = newDie{"a", "f", "f", "k", "p", "s"}
	newDie4  = newDie{"a", "o", "o", "t", "t", "w"}
	newDie5  = newDie{"c", "i", "m", "o", "t", "u"}
	newDie6  = newDie{"d", "e", "i", "l", "r", "x"}
	newDie7  = newDie{"d", "e", "l", "r", "v", "y"}
	newDie8  = newDie{"d", "i", "s", "t", "t", "y"}
	newDie9  = newDie{"e", "e", "g", "h", "n", "w"}
	newDie10 = newDie{"e", "e", "i", "n", "s", "v"}
	newDie11 = newDie{"e", "h", "r", "t", "v", "w"}
	newDie12 = newDie{"e", "i", "o", "s", "s", "t"}
	newDie13 = newDie{"e", "l", "r", "t", "t", "y"}
	newDie14 = newDie{"h", "i", "m", "n", quLigature, "u"}
	newDie15 = newDie{"h", "l", "n", "n", "r", "z"}

	NewDice = [16]Die{newDie0, newDie1, newDie2, newDie3, newDie4, newDie5, newDie6, newDie7, newDie8, newDie9, newDie10, newDie11, newDie12, newDie13, newDie14, newDie15}
)
