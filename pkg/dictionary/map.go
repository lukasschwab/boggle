package dictionary

type Map map[string]bool

func (d Map) Add(word string) {
	d[word] = true
}

func (d Map) Contains(candidate string) bool {
	_, ok := d[candidate]
	return ok
}

func (d Map) CanBePrefix(pre string) bool {
	return true
}

func (d Map) Members() []string {
	result := make([]string, 0, len(d))
	for word, _ := range d {
		result = append(result, word)
	}
	return result
}

func (d Map) Size() int {
	return len(d)
}
