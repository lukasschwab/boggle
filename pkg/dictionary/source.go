package dictionary

import (
	"embed"
	"io"
	"os"
)

//go:embed csw19.txt
var embedFilesystem embed.FS

const (
	// CSW19 file from https://boardgames.stackexchange.com/a/38386
	csw19Path            = "csw19.txt"
	systemDictionaryPath = "/usr/share/dict/words"
)

type Source func() (io.ReadCloser, error)

func CSW19() (io.ReadCloser, error) {
	return embedFilesystem.Open(csw19Path)
}

func SystemDictionary() (io.ReadCloser, error) {
	return os.Open(systemDictionaryPath)
}
