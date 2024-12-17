package boggle

import (
	"testing"

	"github.com/peterldowns/testy/assert"
)

func TestNewSerialize(t *testing.T) {
	board := Shake(ClassicDice)
	// Cover the complex case: 'qu' ligature.
	board.fields[0][0] = quLigature

	serialized := board.Serialize()
	deserialized, err := Deserialize(serialized)

	assert.Nil(t, err)
	assert.Equal(t, board.Pretty(), deserialized.Pretty())
}
