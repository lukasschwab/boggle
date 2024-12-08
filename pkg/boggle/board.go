package boggle

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type board struct {
	fields [4][4]string
}

const (
	serializedBoardSeparator = ' '
)

func indexToLinear(row, col int) int {
	return row*4 + col
}

func (b board) Pretty() string {
	builder := new(strings.Builder)
	for _, row := range b.fields {
		for _, cell := range row {
			builder.WriteString(cell)
			switch len(cell) {
			case 2:
				builder.WriteRune(' ')
			default:
				builder.WriteString("  ")
			}
		}
		builder.WriteRune('\n')
	}

	return builder.String()
}

func (b board) Serialize() string {
	builder := new(strings.Builder)
	for _, row := range b.fields {
		for _, cell := range row {
			builder.WriteString(cell)
		}
	}
	built := builder.String()
	return base64.StdEncoding.EncodeToString([]byte(built))
}

func Deserialize(serialized string) (board, error) {
	decoded, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return board{}, fmt.Errorf("invalid serialized board: %w", err)
	}

	decodedStr := string(decoded)
	result := board{}
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			switch decodedStr[0] {
			case 'q':
				result.fields[r][c], decodedStr = decodedStr[:2], decodedStr[2:]
			default:
				result.fields[r][c], decodedStr = decodedStr[:1], decodedStr[1:]
			}
		}
	}
	return result, nil
}
