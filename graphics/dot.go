package graphics

import "github.com/veandco/go-sdl2/sdl"

// Dot is the graphic object used to draw a dot.
type Dot struct {
	bounds  *sdl.Rect
	texture *sdl.Texture
}

func NewDot(r *sdl.Rect, t *sdl.Texture) *Dot {
	return &Dot{bounds: r, texture: t}
}

func (d *Dot) Blit(r *Renderer) {
	r.Copy(d.texture, nil, d.bounds)
}

func (d *Dot) BlitTo(r *Renderer, dst *sdl.Rect) {
	r.Copy(d.texture, nil, dst)
}

func (d *Dot) Bounds() *sdl.Rect {
	return d.bounds
}

func (d *Dot) Destroy() {
	d.texture.Destroy() //nolint
}
