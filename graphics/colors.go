package graphics

import "github.com/veandco/go-sdl2/sdl"

type Color int

const (
	Red Color = iota
	Green
	Blue
	Yellow
	Magenta
	Pink
	Orange
	Brown
	White
	Black
	NoColor Color = -1
)

var Colors = []sdl.Color{
	sdl.Color{R: 255, G: 0, B: 0, A: 255},
	sdl.Color{R: 0, G: 255, B: 0, A: 255},
	sdl.Color{R: 0, G: 0, B: 255, A: 255},
	sdl.Color{R: 255, G: 255, B: 0, A: 255},
	sdl.Color{R: 255, G: 0, B: 255, A: 255},
	sdl.Color{R: 255, G: 20, B: 147, A: 255},
	sdl.Color{R: 255, G: 69, B: 0, A: 255},
	sdl.Color{R: 139, G: 69, B: 19, A: 255},
	sdl.Color{R: 255, G: 255, B: 255, A: 255},
	sdl.Color{R: 0, G: 0, B: 0, A: 255},
}
