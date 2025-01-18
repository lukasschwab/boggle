package main

import (
	"testing"

	"github.com/lukasschwab/boggle/pkg/boggle"
	"github.com/peterldowns/testy/assert"
)

func TestDicePlayable(t *testing.T) {
	dicePlayable := dicePlayableFilter(boggle.ClassicDice)

	// Quasi-tests.
	assert.False(t, dicePlayable("ququ"))
	assert.True(t, dicePlayable("cat"))
	assert.True(t, dicePlayable("xysdosaeanntdcgd"))
	assert.False(t, dicePlayable("xysdosaeanntdcgda"))
	assert.True(t, dicePlayable("qu"))
	assert.False(t, dicePlayable("qu"+"j"))
}
