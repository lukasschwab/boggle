package dictionary

type Trie struct {
	root *trieNode
}

func (d Trie) Add(word string) {
	d.root.addSuffix(word)
}

func (d Trie) Contains(candidate string) bool {
	return d.root.traverse(candidate, IsWord)
}

func (d Trie) CanBePrefix(pre string) bool {
	return d.root.traverse(pre, Exists)
}

func EmptyTrie() Trie {
	return Trie{newTrieNode()}
}

type trieNode struct {
	isWord   bool
	children map[rune]*trieNode
}

func newTrieNode() *trieNode {
	return &trieNode{
		isWord:   false,
		children: make(map[rune]*trieNode),
	}
}

func (tn *trieNode) addSuffix(suf string) {
	if suf == "" {
		tn.isWord = true
		return
	}

	head, tail := rune(suf[0]), suf[1:]
	if _, ok := tn.children[head]; !ok {
		tn.children[head] = newTrieNode()
	}
	tn.children[head].addSuffix(tail)
}

// traverse returns false if the suffix is unfound in the trie. If it does find
// a node, it evaluates eval on that node.
func (tn *trieNode) traverse(suf string, eval func(*trieNode) bool) bool {
	if suf == "" {
		return eval(tn)
	}
	head, tail := rune(suf[0]), suf[1:]
	if next, ok := tn.children[head]; !ok {
		return false
	} else {
		return next.traverse(tail, eval)
	}
}

// IsWord returns whether tn represents a word. If the dictionary only contains
// "test", IsWord will return false for the 'e' node for the search string "te".
func IsWord(tn *trieNode) bool {
	return tn.isWord
}

// Exists returns whether tn exists at all. If the dictionary only contains
// "test", Exists will return true for the 'e' node for the search string "te".
func Exists(tn *trieNode) bool {
	return true
}
