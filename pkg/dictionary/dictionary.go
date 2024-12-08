package dictionary

import (
	"bufio"
	"fmt"
	"strings"
)

type Interface interface {
	Add(word string)
	Contains(candidate string) bool
	CanBePrefix(pre string) bool
}

func Load(source Source, dict Interface) error {
	reader, err := source()
	if err != nil {
		return fmt.Errorf("error opening dictionary: %w", err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		normalized := strings.ToLower(scanner.Text())
		dict.Add(normalized)
	}
	return nil
}
