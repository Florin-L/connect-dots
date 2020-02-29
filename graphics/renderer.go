package graphics

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"go.uber.org/zap"
)

// Renderer wrapps a sdl.Renderer (a rendering target).
// It "overrides" some functions from sdl.Renderer, checks the return value and
// terminates the applcation in case of an error (these kind of errors are un-recoverables).
type Renderer struct {
	Renderer *sdl.Renderer
	log      *zap.Logger
}

// NewRenderer creates a new renderer.
func NewRenderer(r *sdl.Renderer, l *zap.Logger) *Renderer {
	return &Renderer{r, l}
}

// Destroy releases the resources allocated to the renderer.
func (r *Renderer) Destroy() {
	r.Renderer.Destroy() //nolint
}

// Clear clears the current rendering target with the drawing color.
func (r *Renderer) Clear() {
	if err := r.Renderer.Clear(); err != nil {
		r.log.Fatal("Clear failed", zap.Error(err))
	}
}

func (r *Renderer) Present() {
	r.Renderer.Present()
}

// CreateTexture returns a new texture for a rendering context.
func (r *Renderer) CreateTexture(format uint32, access int, w, h int32) *sdl.Texture {
	t, err := r.Renderer.CreateTexture(format, access, w, h)
	if err != nil {
		r.log.Fatal("CreateTexture failed", zap.Error(err))
	}
	return t
}

// Copy copies a portion of the texture to the current rendering target.
func (r *Renderer) Copy(t *sdl.Texture, src *sdl.Rect, dst *sdl.Rect) {
	if err := r.Renderer.Copy(t, src, dst); err != nil {
		r.log.Fatal("Copy failed", zap.Error(err))
	}
}

func (r *Renderer) SetDrawColor(rc, gc, bc, a uint8) {
	if err := r.Renderer.SetDrawColor(rc, gc, bc, a); err != nil {
		r.log.Fatal("SetDrawColor failed", zap.Error(err))
	}
}

func (r *Renderer) SetRenderTarget(t *sdl.Texture) {
	if err := r.Renderer.SetRenderTarget(t); err != nil {
		r.log.Fatal("SetRenderTarget failed", zap.Error(err))
	}
}

func (r *Renderer) GetRenderTarget() *sdl.Texture {
	return r.Renderer.GetRenderTarget()
}

func (r *Renderer) DrawRect(rc *sdl.Rect) {
	if err := r.Renderer.DrawRect(rc); err != nil {
		r.log.Fatal("DrawRect failed", zap.Error(err))
	}
}

func (r *Renderer) FillRect(rc *sdl.Rect) {
	if err := r.Renderer.FillRect(rc); err != nil {
		r.log.Fatal("FillRect failed", zap.Error(err))
	}
}

func (r *Renderer) FillCircle(x, y, radius int32, c sdl.Color) {
	ok := gfx.FilledCircleColor(r.Renderer, x, y, radius, c)
	if !ok {
		r.log.Fatal("FilledCircleColor failed")
	}
}

func (r *Renderer) DrawVLine(x, y1, y2 int32, thick int32) {
	rc := sdl.Rect{
		X: x - thick/2,
		Y: y1,
		W: thick,
		H: y2 - y1,
	}
	r.FillRect(&rc)
}

func (r *Renderer) DrawHLine(x1, x2, y int32, thick int32) {
	rc := sdl.Rect{
		X: x1,
		Y: y - thick/2,
		W: x2 - x1,
		H: thick,
	}
	r.FillRect(&rc)
}
