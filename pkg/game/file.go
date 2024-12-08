package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

// TODO: test with a temp path.
func WriteFile(path string, f Frontmatter, words []string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error opening '%v': %w", path, err)
	}
	defer file.Close()

	sink := bufio.NewWriter(file)
	if err := f.WriteTo(sink); err != nil {
		return err
	}
	for _, word := range words {
		sink.WriteString(fmt.Sprintf("%s\n", word))
	}

	if err := sink.Flush(); err != nil {
		return fmt.Errorf("failed flushing file sink: %w", err)
	}
	return nil
}

func LoadFile(path string) (f Frontmatter, words []string, err error) {
	if file, err := os.Open(path); err != nil {
		return f, words, fmt.Errorf("error opening '%v': %w", path, err)
	} else if rest, err := frontmatter.Parse(file, &f); err != nil {
		return f, words, fmt.Errorf("error parsing frontmatter%w", path, err)
	} else {
		trimmedRest := strings.TrimSpace(string(rest))
		return f, strings.Split(trimmedRest, "\n"), nil
	}
}

type Frontmatter struct {
	// Board, serialized; see [boggle.Deserialize].
	Board        string
	TimerSeconds int
}

func (f Frontmatter) WriteTo(sink *bufio.Writer) error {
	sink.WriteString("---\n")
	bytes, err := yaml.Marshal(f)
	if err != nil {
		return fmt.Errorf("error marshaling frontmatter %+v: %w", f, err)
	}
	sink.Write(bytes)
	sink.WriteString("---\n")
	return nil
}
