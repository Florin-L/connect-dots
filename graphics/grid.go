package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Grid is the graphic object used to display the board.
type Grid struct {
	bounds  *sdl.Rect
	texture *sdl.Texture
}

// NewGrid creates a grid graphic object.
func NewGrid(r *sdl.Rect, t *sdl.Texture) *Grid {
	return &Grid{
		bounds:  r,
		texture: t,
	}
}

// Blit renders the grid on a surface.
func (g *Grid) Blit(r *Renderer) {
	r.Copy(g.texture, nil, g.bounds)
}

//
func (g *Grid) BlitTo(r *Renderer, dst *sdl.Rect) {
	r.Copy(g.texture, nil, dst)
}

func (g *Grid) Bounds() *sdl.Rect {
	return g.bounds
}

// Destroy releases the graphics resources used by the grid.
func (g *Grid) Destroy() {
	g.texture.Destroy()
}

// ClickedInside checks if the given screen coordinate is within
// the bounds of the grid.
func (g *Grid) ClickedInside(x, y int32) bool {
	p := &sdl.Point{X: x, Y: y}
	return p.InRect(g.bounds)
}

// ScreenToGrid returns the grid coordinates (column and row) corresponding
// to a mouse event screen coordinates if the latter is moving/clicked
// over/on the grid.
func (g *Grid) ScreenToGrid(x, y, squareSize int32) (int32, int32, bool) {
	if !g.ClickedInside(x, y) {
		return -1, -1, false
	}

	cx := (x - g.bounds.X) / squareSize
	cy := (y - g.bounds.Y) / squareSize

	return cx, cy, true
}
