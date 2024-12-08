package dictionary

type Filtered struct {
	Underlying Interface
	Filter     func(string) bool
}

func (d Filtered) Add(word string) {
	if d.Filter(word) {
		d.Underlying.Add(word)
	}
}

func (d Filtered) Contains(candidate string) bool {
	return d.Filter(candidate) && d.Underlying.Contains(candidate)
}

func (d Filtered) CanBePrefix(pre string) bool {
	return d.Underlying.CanBePrefix(pre)
}
