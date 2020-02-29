package graphics

import (
	"connect-dots/config"

	"github.com/veandco/go-sdl2/sdl"
)

// AssetsStorage keeps the graphics assests/objects
// which are to be rendered on a rendering target:
// - the board (grid)
// - the dots for each color
// - the vertical and horizontal lines for each color
// The assests may be loaded from files where generated programatically.
type AssetsStorage struct {
	Grid       Renderable
	Dots       []Renderable
	VertLines  []Renderable
	HorizLines []Renderable
}

// NewAssetsStorage creates a new graphics assests storage.
func NewAssetsStorage() *AssetsStorage {
	return &AssetsStorage{}
}

// Init creates the graphics assests.
func (s *AssetsStorage) Init(renderer *Renderer, config *config.Config) error {
	grid := CreateGrid(renderer, config)
	s.Grid = grid

	s.Dots = make([]Renderable, len(Colors))
	s.VertLines = make([]Renderable, len(Colors))
	s.HorizLines = make([]Renderable, len(Colors))
	for i, c := range Colors {
		s.Dots[i] = createDot(c, renderer, config)
		s.VertLines[i] = createVLine(c, renderer, config)
		s.HorizLines[i] = createHLine(c, renderer, config)
	}

	return nil
}

// Destroy releases all the graphics assests.
func (s *AssetsStorage) Destroy() {
	s.Grid.Destroy()
	for _, dot := range s.Dots {
		dot.Destroy()
	}
	for _, vl := range s.VertLines {
		vl.Destroy()
	}
	for _, hl := range s.HorizLines {
		hl.Destroy()
	}
}

// CreateGrid creates a graphics object which is used to render
// a grid of a given size.
func CreateGrid(renderer *Renderer, cfg *config.Config) *Grid {
	crtTarget := renderer.GetRenderTarget()
	defer renderer.SetRenderTarget(crtTarget)

	w := cfg.Size * cfg.SquareSize
	h := w
	r := &sdl.Rect{
		X: (cfg.WindowWidth - w) / 2,
		Y: (cfg.WindowHeight - h) / 2,
		W: w,
		H: h,
	}

	t := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		cfg.SquareSize,
		cfg.SquareSize,
	)
	defer t.Destroy() //nolint

	renderer.SetRenderTarget(t)
	renderer.SetDrawColor(cfg.Color.R, cfg.Color.G, cfg.Color.B, cfg.Color.A)

	rc := sdl.Rect{X: 0, Y: 0, W: cfg.SquareSize, H: cfg.SquareSize}
	renderer.FillRect(&rc)
	renderer.SetDrawColor(255, 255, 0, 255)
	renderer.DrawRect(&rc)

	gt := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		w,
		h,
	)

	renderer.SetRenderTarget(gt)
	for x := 0; x < int(cfg.Size); x++ {
		for y := 0; y < int(cfg.Size); y++ {
			renderer.Copy(t,
				nil,
				&sdl.Rect{
					X: int32(x) * cfg.SquareSize,
					Y: int32(y) * cfg.SquareSize,
					W: cfg.SquareSize,
					H: cfg.SquareSize,
				},
			)
		}
	}

	return NewGrid(r, gt)
}

func createDot(color sdl.Color,
	renderer *Renderer, cfg *config.Config) *Dot {

	crtTarget := renderer.GetRenderTarget()
	defer renderer.SetRenderTarget(crtTarget)

	r := &sdl.Rect{
		X: 0,
		Y: 0,
		W: 2 * cfg.DotRadius,
		H: 2 * cfg.DotRadius,
	}

	t := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		r.W,
		r.H,
	)
	renderer.SetRenderTarget(t)

	renderer.SetDrawColor(cfg.Color.R, cfg.Color.G, cfg.Color.B, cfg.Color.A)
	renderer.Clear()
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.FillCircle(r.W/2, r.H/2, cfg.DotRadius-1, color)

	return NewDot(r, t)
}

func createVLine(color sdl.Color,
	renderer *Renderer, cfg *config.Config) *Line {

	crtTarget := renderer.GetRenderTarget()
	defer renderer.SetRenderTarget(crtTarget)

	r := &sdl.Rect{
		X: 0,
		Y: 0,
		W: cfg.SquareSize,
		H: cfg.SquareSize,
	}

	t := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		r.W,
		r.H,
	)
	t.SetBlendMode(sdl.BLENDMODE_BLEND) //nolint
	t.SetAlphaMod(50)                   //nolint
	renderer.SetRenderTarget(t)

	renderer.SetDrawColor(cfg.Color.R, cfg.Color.G, cfg.Color.B, cfg.Color.A)
	renderer.Clear()
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.DrawVLine(r.W/2, 0, r.H, 4)

	return NewLine(r, t)
}

func createHLine(color sdl.Color,
	renderer *Renderer, cfg *config.Config) *Line {

	crtTarget := renderer.GetRenderTarget()
	defer renderer.SetRenderTarget(crtTarget)

	r := &sdl.Rect{
		X: 0,
		Y: 0,
		W: cfg.SquareSize,
		H: cfg.SquareSize,
	}

	t := renderer.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		r.W,
		r.H,
	)
	t.SetBlendMode(sdl.BLENDMODE_BLEND) //nolint
	t.SetAlphaMod(50)                   //nolint
	renderer.SetRenderTarget(t)

	renderer.SetDrawColor(cfg.Color.R, cfg.Color.G, cfg.Color.B, cfg.Color.A)
	renderer.Clear()
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	renderer.DrawHLine(0, r.W, r.H/2, 4)

	return NewLine(r, t)
}
