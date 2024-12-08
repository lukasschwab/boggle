package dictionary_test

import (
	"testing"

	"github.com/lukasschwab/boggle/pkg/dictionary"
	"github.com/peterldowns/testy/assert"
)

func TestMapDict(t *testing.T) {
	d := dictionary.Map{}

	assert.False(t, d.Contains("test"))

	d.Add("test")
	assert.True(t, d.Contains("test"))
	assert.False(t, d.Contains("tes"))
	assert.False(t, d.Contains("other"))
}

func TestTrieDict(t *testing.T) {
	d := dictionary.EmptyTrie()

	assert.False(t, d.Contains("test"))

	d.Add("test")
	assert.True(t, d.Contains("test"))
	assert.False(t, d.Contains("tes"))
	assert.False(t, d.Contains("other"))
}
