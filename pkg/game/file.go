package game

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v3"
)

const (
	aboutThisFile = "https://github.com/lukasschwab/boggle"
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

	sort.Strings(words)
	for _, word := range words {
		sink.WriteString(fmt.Sprintf("%s\n", word)) //nolint:errcheck
	}

	if err := sink.Flush(); err != nil {
		return fmt.Errorf("failed flushing file sink: %w", err)
	}
	return nil
}

func LoadFile(path string) (f Frontmatter, words []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return f, words, fmt.Errorf("error opening '%v': %w", path, err)
	}
	defer file.Close()

	if rest, err := frontmatter.Parse(file, &f); err != nil {
		return f, words, fmt.Errorf("error parsing frontmatter from '%v': %w", path, err)
	} else {
		trimmedRest := strings.TrimSpace(string(rest))
		return f, strings.Split(trimmedRest, "\n"), nil
	}
}

type Frontmatter struct {
	// Board, serialized; see [boggle.Deserialize].
	Board         string
	TimerSeconds  int
	AboutThisFile string
}

//nolint:errcheck
func (f Frontmatter) WriteTo(sink *bufio.Writer) error {
	f.AboutThisFile = aboutThisFile

	sink.WriteString("---\n")
	bytes, err := yaml.Marshal(f)
	if err != nil {
		return fmt.Errorf("error marshaling frontmatter %+v: %w", f, err)
	}
	sink.Write(bytes)
	sink.WriteString("---\n")
	return nil
}
