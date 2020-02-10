package config

import "github.com/veandco/go-sdl2/sdl"

type Difficulty int

const (
	Easy Difficulty = iota
	Intermediate
	Advanced
)

// Config stores the game configuration.
type Config struct {
	// The width of the application window.
	WindowWidth int32
	// The height of the application window.
	WindowHeight int32

	// The size (the number of rows and columns) of the board.
	Size int32

	// The size of a board square.
	SquareSize int32
	// The radius of the dot.
	DotRadius int32

	// The difficulty level.
	Difficulty Difficulty

	// The board color.
	Color *sdl.Color
}

type Option func(*Config)

// New returns a new game configuration.
func New(opts ...Option) *Config {
	cfg := &Config{
		WindowWidth:  800,
		WindowHeight: 600,
		Size:         10,
		SquareSize:   48,
		DotRadius:    16,
		Difficulty:   Easy,
		Color:        &sdl.Color{R: 168, G: 168, B: 168, A: 255},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
