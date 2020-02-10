package graphics

import "github.com/veandco/go-sdl2/sdl"

// Line is the graphic object used to draw a line connecting squares.
type Line struct {
	bounds  *sdl.Rect
	texture *sdl.Texture
}

func NewLine(r *sdl.Rect, t *sdl.Texture) *Line {
	return &Line{bounds: r, texture: t}
}

func (l *Line) Blit(r *Renderer) {
	r.Copy(l.texture, nil, l.bounds)
}

func (l *Line) BlitTo(r *Renderer, dst *sdl.Rect) {
	r.Copy(l.texture, nil, dst)
}

func (l *Line) Bounds() *sdl.Rect {
	return l.bounds
}

func (l *Line) Destroy() {
	l.texture.Destroy()
}
