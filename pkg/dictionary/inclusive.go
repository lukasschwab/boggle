package dictionary

// Inclusive of all strings.
type Inclusive struct{}

func (d Inclusive) Add(word string) { return }

func (d Inclusive) Contains(candidate string) bool { return true }

func (d Inclusive) CanBePrefix(pre string) bool { return true }
