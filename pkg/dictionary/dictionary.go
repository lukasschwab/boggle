package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	SystemDictionary = "/usr/share/dict/words"
)

type Interface interface {
	Add(word string)
	Contains(candidate string) bool
	CanBePrefix(pre string) bool
}

func LoadFromFile(dict Interface) error {
	file, err := os.Open(SystemDictionary)
	if err != nil {
		return fmt.Errorf("error opening dictionary: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		normalized := strings.ToLower(scanner.Text())
		dict.Add(normalized)
	}
	return nil
}
