package main

import (
	"testing"

	"github.com/peterldowns/testy/assert"
)

func TestDicePlayable(t *testing.T) {
	// Quasi-tests.
	assert.False(t, dicePlayable("ququ"))
	assert.True(t, dicePlayable("cat"))
	assert.True(t, dicePlayable("xysdosaeanntdcgd"))
	assert.False(t, dicePlayable("xysdosaeanntdcgda"))
	assert.True(t, dicePlayable("qu"))
	assert.False(t, dicePlayable("qu"+"j"))
}
