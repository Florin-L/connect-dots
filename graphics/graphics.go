package graphics

import "github.com/veandco/go-sdl2/sdl"

// Bounded is an interface that the objects which have boundaries
// should implement.
type Bounded interface {
	// Bounds returns the bounds of the graphics object.
	Bounds() *sdl.Rect
}

// Destroyer is an interface that the types which owns graphics resources
// should implement.
type Destroyer interface {
	// Destroy releases the graphics object.
	Destroy()
}

// Blitter is an interface that all the objects (textures) which are to be
// blit/copied to a rendering target should implement.
type Blitter interface {
	// Blit copies a graphics object (texture) to a rendering target.
	Blit(r *Renderer)

	// BlitTo copies a graphics object (texture) to a rendering target
	// into a given rectangle.
	BlitTo(r *Renderer, dst *sdl.Rect)
}

// Renderable
type Renderable interface {
	Bounded
	Destroyer
	Blitter
}
