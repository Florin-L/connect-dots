package game

import (
	"connect-dots/graphics"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// Level is a struct which stores the configuration of game level:
// - the board size
// - the dots (colors and board coordinations)
type Level struct {
	// Size is the size of the board (5,6,7,8,9 or 10).
	Size int32

	// The dots loaded from the file level.
	Dots []Dot
}

// LoadFromFile loads the level data from a file.
// The path is relative to the directory where the game process runs in
// and has the following structure:
// data/<n>x<n>/level<d>/<m>.json
// where:
// - data is a directory relative to the working directory
// - <n> is the board dimension (5,6,7,8,9,10)
// - <d> is the difficulty (0: easy, 1: intermediate, 2: advanced)
// - <m> is the m-th file in the directory where we
func LoadFromFile(path string) (*Level, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Load(data)
}

// Load decodes a Json Blob and instantiate a Level struct.
func Load(data []byte) (*Level, error) {
	var level struct {
		Size int32 `json:"size"`
		Dots []struct {
			X     int32  `json:"x"`
			Y     int32  `json:"y"`
			Color string `json:"color"`
		}
	}

	err := json.Unmarshal(data, &level)
	if err != nil {
		return nil, err
	}

	if level.Size <= 0 {
		return nil, fmt.Errorf("Invalid value for size: %d", level.Size)
	}

	if len(level.Dots) == 0 {
		return nil, errors.New("No dots found in the level file")
	}

	l := &Level{}

	l.Size = level.Size
	l.Dots = []Dot{}
	for _, dot := range level.Dots {
		var c graphics.Color
		switch dot.Color {
		case "red":
			c = graphics.Red
		case "green":
			c = graphics.Green
		case "blue":
			c = graphics.Blue
		case "yellow":
			c = graphics.Yellow
		case "pink":
			c = graphics.Pink
		case "orange":
			c = graphics.Orange
		case "brown":
			c = graphics.Brown
		case "white":
			c = graphics.White
		case "black":
			c = graphics.Black
		}

		l.Dots = append(l.Dots, Dot{
			Location: Coordinate{dot.X, dot.Y},
			Color:    c,
		})
	}

	return l, nil
}
