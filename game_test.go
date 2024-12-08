package main

import (
	"testing"

	"github.com/lukasschwab/boggle/pkg/dictionary"
	"github.com/peterldowns/testy/assert"
)

func TestTraverse(t *testing.T) {
	// c  u  f  i
	// u  e  h  l
	// r  t  t  a
	// d  i  b  i
	b, err := deserialize("YyB1IGYgaSB1IGUgaCBsIHIgdCB0IGEgZCBpIGIgaQ==")
	assert.Nil(t, err)

	t.Logf("Board:\n%v\n", b.pretty())

	dict := dictionary.EmptyTrie()
	assert.Nil(t, dictionary.LoadFromFile(dict))
	assert.True(t, dict.Contains("cute"))

	words := b.words(dict, index{0, 0})
	assert.GreaterThan(t, len(words), 0)

	wordSet := map[string]bool{}
	for _, word := range words {
		assert.Equal(t, 'c', word[0])
		assert.LessThanOrEqual(t, len(word), 16)
		wordSet[word] = true
	}

	_, containsCute := wordSet["cute"]
	assert.True(t, containsCute)
}
