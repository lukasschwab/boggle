package counter

type Counter map[string]int

func (c Counter) Incr(token string) {
	if _, ok := c[token]; !ok {
		c[token] = 0
	}
	c[token]++
}

func (c Counter) LessThan(other Counter) bool {
	for letter, count := range c {
		otherCount, ok := other[letter]
		if !ok || count > otherCount {
			return false
		}
	}
	return true
}

type Count struct {
	Word  string
	Count int
}

func (c Counter) Counts() []Count {
	result := make([]Count, 0, len(c))
	for word, count := range c {
		result = append(result, Count{
			Word:  word,
			Count: count,
		})
	}
	return result
}
