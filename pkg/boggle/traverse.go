package boggle

import (
	"github.com/lukasschwab/boggle/pkg/dictionary"
)

// There are basically two strategies for handling gameplay:
//
// 1. Precalculate all the words in the board, then determine membership when
//    the user inputs a candidate (simple lookup).
// 2. Determine "reachability" on the fly (can this word be found on this board)
//    and do lookup in the total dictionary.

func (b Board) AllWords(dict dictionary.Interface) dictionary.Map {
	result := dictionary.Map{}
	for _, idx := range allIndices() {
		for _, word := range b.words(dict, idx) {
			result.Add(word)
		}
	}
	return result
}

// ContainsString if candidate is reachable on b, regardless of any relationship
// to a dictionary.
func (b Board) ContainsString(candidate string) bool {
	for _, idx := range allIndices() {
		if b.hasStringAtDepthFirst(idx, candidate, visited{}) {
			return true
		}
	}
	return false
}

// NOTE: closely resembles wordsDepthFirst, but with a different visitation
// pattern and early escape. Don't think it's worth generalizing.
func (b Board) hasStringAtDepthFirst(idx index, candidate string, vis visited) bool {
	vis[idx] = true
	defer delete(vis, idx)

	head, tail := HeadTail(candidate)
	if b.get(idx) != string(head) {
		return false
	}
	if tail == "" {
		return true
	}

	for _, n := range idx.neighbors() {
		if _, visited := vis[n]; !visited && b.hasStringAtDepthFirst(n, tail, vis) {
			return true
		}
	}
	return false
}

// allIndices for a 4x4 game.
func allIndices() []index {
	results := make([]index, 0, 16)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			results = append(results, index{i, j})
		}
	}
	return results
}

type index [2]int

func (i index) neighbors() []index {
	results := make([]index, 0)

	for drow := -1; drow <= 1; drow++ {
		for dcol := -1; dcol <= 1; dcol++ {
			if drow == 0 && dcol == 0 {
				continue
			}
			candidate := index{i[0] + drow, i[1] + dcol}
			if candidate.valid() {
				results = append(results, candidate)
			}
		}
	}

	return results
}

func (i index) valid() bool {
	return i[0] >= 0 && i[0] < 4 &&
		i[1] >= 0 && i[1] < 4
}

type visited map[index]bool

// words rooted at a given index.
func (b Board) words(dict dictionary.Interface, start index) []string {
	return b.wordsDepthFirst(dict, start, visited{}, "")
}

func (b Board) wordsDepthFirst(
	dict dictionary.Interface,
	idx index,
	vis visited,
	soFar string,
) []string {
	vis[idx] = true
	defer delete(vis, idx)
	soFar = soFar + b.get(idx)

	result := []string{}

	if !dict.CanBePrefix(soFar) {
		return result
	} else if dict.Contains(soFar) {
		result = append(result, soFar)
	}

	for _, n := range idx.neighbors() {
		if _, visited := vis[n]; !visited {
			result = append(result, b.wordsDepthFirst(dict, n, vis, soFar)...)
		}
	}
	return result
}

func (b Board) get(idx index) string {
	return b.fields[idx[0]][idx[1]]
}

// HeadTail of s, treating `qu` as a single character.
func HeadTail(s string) (string, string) {
	if len(s) >= 2 && s[:2] == quLigature {
		return s[:2], s[2:]
	}
	return string(s[0]), s[1:]
}
