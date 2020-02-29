package graphics

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Text is a graphics object used to render a text on a rendering target.
type Text struct {
	font *ttf.Font

	// The text to be displayed.
	Text string
}

// NewText creates a Text graphics object.
func NewText(text string, f *ttf.Font) *Text {
	return &Text{f, text}
}

// Draw renders the text on a rendering target at the given position.
func (t *Text) Draw(r *Renderer, pos sdl.Point) error {
	s, err := t.font.RenderUTF8Solid(t.Text, sdl.Color{R: 255, G: 255, B: 0})
	if err != nil {
		return err
	}
	defer s.Free()

	tex, err := r.Renderer.CreateTextureFromSurface(s)
	if err != nil {
		return err
	}
	defer tex.Destroy() //nolint

	rc := sdl.Rect{X: pos.X, Y: pos.Y, W: s.W, H: s.H}
	if err := r.Renderer.Copy(tex, nil, &rc); err != nil {
		return err
	}

	return nil
}
