package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"time"
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

func randomBoard() board {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(allDice), func(i, j int) {
		allDice[i], allDice[j] = allDice[j], allDice[i]
	})

	fmt.Println(allDice)

	result := board{}
	for rowIdx, row := range result.fields {
		for colIdx := range row {
			result.fields[rowIdx][colIdx] = allDice[indexToLinear(rowIdx, colIdx)].roll()
		}
	}
	return result
}

func (b board) pretty() string {
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

func (b board) serialize() string {
	builder := new(strings.Builder)
	for _, row := range b.fields {
		for _, cell := range row {
			builder.WriteString(cell)
			builder.WriteRune(serializedBoardSeparator)
		}
	}

	// BODGE: drop the trailing separator.
	built := builder.String()
	bytes := []byte(built[:len(built)-1])
	return base64.StdEncoding.EncodeToString(bytes)
}

// TODO: test round-trip.
func deserialize(serialized string) (board, error) {
	decoded, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return board{}, fmt.Errorf("invalid serialized board: %w", err)
	}

	parts := strings.Split(string(decoded), string(serializedBoardSeparator))
	if len(parts) != 16 {
		return board{}, fmt.Errorf("unexpected number of fields %v in %v; expected 16", len(parts), parts)
	}

	result := board{}
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			result.fields[r][c] = parts[indexToLinear(r, c)]
		}
	}
	return result, nil
}
