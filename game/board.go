package game

import (
	"connect-dots/graphics"
	"fmt"
)

// Coordinate stores a board coordinates (Y|row and X|column).
type Coordinate struct {
	// The column coordinate.
	X int32
	// The row coordinate.
	Y int32
}

// NewCoord returns a new coordinate.
func NewCoord(x int32, y int32) Coordinate {
	return Coordinate{x, y}
}

func (c *Coordinate) IsValid() bool {
	return c.X >= 0 && c.Y >= 0
}

// Line connects two adjacent squares on the board.
type Line struct {
	// The starting coordinate.
	From Coordinate
	// The end coordinate.
	To Coordinate
	// The color of the line.
	Color graphics.Color
}

// NewLine returns a new line.
func NewLine(from, to Coordinate, c graphics.Color) Line {
	return Line{from, to, c}
}

// Dot represents a dot from the board.
type Dot struct {
	// The coordinates of the dot.
	Location Coordinate
	// The color of the dot.
	Color graphics.Color
}

// Path is a path between two dots on the board.
type Path struct {
	// The start of the path.
	StartDot *Dot
	// The end of the path.
	EndDot *Dot
	// The lines which connect the start and end dots.
	Lines []*Line
}

// AddLine adds a line to the current path.
func (p *Path) AddLine(from, to Coordinate) {
	p.Lines = append(p.Lines,
		&Line{
			From:  from,
			To:    to,
			Color: p.StartDot.Color,
		})
}

// RemoveLine removes a line from the current path.
func (p *Path) RemoveLine(from, to Coordinate) bool {
	for i, l := range p.Lines {
		if l.From == from && l.To == to {
			if i < len(p.Lines)-1 {
				copy(p.Lines[i:], p.Lines[i+1:])
			}
			p.Lines[len(p.Lines)-1] = nil
			p.Lines = p.Lines[:len(p.Lines)-1]
			return true
		}
	}
	return false
}

// ContainsLine checks if the path contains a line.
func (p *Path) ContainsLine(from, to Coordinate) bool {
	for _, l := range p.Lines {
		if l.From == from && l.To == to {
			return true
		}
	}
	return false
}

// Board stores all the dots and the paths connecting the dots
// as well as the state of the squares (the colors).
type Board struct {
	// The paths which connect the pairs of dots having the same color.
	Paths map[Dot]*Path

	// The states (colors) of the board squares/cells.
	colors []graphics.Color

	// The size of the board.
	size int32
}

// NewBoard creates a board.
func NewBoard(size int32) *Board {
	cs := make([]graphics.Color, size*size)
	for i := range cs {
		cs[i] = graphics.NoColor
	}

	return &Board{
		Paths:  make(map[Dot]*Path),
		colors: cs,
		size:   size,
	}
}

// InitPath initilizes the paths.
func (b *Board) InitPath(dot Dot) {
	b.Paths[dot] = &Path{StartDot: &Dot{
		Location: Coordinate{dot.Location.X, dot.Location.Y},
		Color:    dot.Color,
	}}
	*(b.ColorAt(dot.Location.X, dot.Location.Y)) = dot.Color
}

// ColorAt returns a pointer to the graphics.Color object
// corresponding to the given board coordinates.
func (b *Board) ColorAt(x, y int32) *graphics.Color {
	return &(b.colors[x*b.size+y])
}

// InitPaths initializes the map of the paths:
// for each dot set the start dots and their colors
func (b *Board) InitPaths(dots []Dot) {
	for _, dot := range dots {
		b.Paths[dot] = &Path{StartDot: &Dot{
			Location: Coordinate{dot.Location.X, dot.Location.Y},
			Color:    dot.Color,
		}}
		*(b.ColorAt(dot.Location.X, dot.Location.Y)) = dot.Color
	}
}

// Clear removes the paths and resets the slice of colors.
func (b *Board) Clear() {
	for path := range b.Paths {
		delete(b.Paths, path)
	}

	for i := range b.colors {
		b.colors[i] = graphics.NoColor
	}
}

// Coverage returns the number of covered squares.
func (b *Board) Coverage() int32 {
	count := int32(0)
	for _, c := range b.colors {
		if c != graphics.NoColor {
			count++
		}
	}
	return count
}

func (b *Board) Dump() {
	fmt.Println("Board:")
	for y := 0; y < int(b.size); y++ {
		for x := 0; x < int(b.size); x++ {
			fmt.Printf("|% 2d| ", int(*b.ColorAt(int32(x), int32(y))))
		}
		fmt.Println()
	}
}
