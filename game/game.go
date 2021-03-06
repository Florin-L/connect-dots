package game

import (
	"connect-dots/config"
	"connect-dots/graphics"
	"connect-dots/ui"
	"fmt"
	"strconv"
	"strings"

	"log"
	"math"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"go.uber.org/zap"
)

type drawAction int

const (
	none drawAction = iota
	drawLine
	eraseLine
	completePath
)

const maxBoardSize = 10

type editPathState struct {
	// true if a path is being edited
	editingPath bool
	// the current path
	path *Path

	// the source dot (current selected dot) if any
	srcDot *Dot
	// the destination dot (if any)
	dstDot *Dot

	// the color of the current selected dot
	color graphics.Color
	// the square of the board that the mouse in hovering over
	square Coordinate
}

func (s *editPathState) reset() {
	s.editingPath = false
	s.srcDot = nil
	s.dstDot = nil
	s.color = graphics.NoColor
	s.square = Coordinate{-1, -1}
	s.path = nil
}

// Game implements the game.
type Game struct {
	// True if the level was completed:
	// all the dots are connected and all the squares are visited
	Completed bool

	// The game configuration.
	config *config.Config

	// The assets storage.
	assets *graphics.AssetsStorage

	// The window this game is rendered in.
	window *sdl.Window

	//
	movesText    *graphics.Text
	coverageText *graphics.Text

	// The current level.
	level *Level

	// The name of the file where the level is loaded from.
	file string

	// The  board.
	board *Board

	// The bounds of the dots graphics objects.
	dotBounds map[Dot]sdl.Rect

	// The bounds of the line graphics objects.
	lineBounds map[Line]sdl.Rect

	// The current state during a mouse move action.
	state *editPathState

	// The logger
	log *zap.Logger

	// The overall attempts to successfully connect the dots.
	Moves int32

	// The board coverage (the number of the squares which are covered
	// with dots or lines).
	coverage int32
}

func newState() *editPathState {
	return &editPathState{
		editingPath: false,
		srcDot:      nil,
		dstDot:      nil,
		color:       graphics.NoColor,
		square:      Coordinate{-1, -1},
		path:        nil,
	}
}

type option func(*Game)

// New creates a game.
func New(cfg *config.Config, assets *graphics.AssetsStorage, opts ...option) *Game {
	g := &Game{
		Completed:    false,
		config:       cfg,
		assets:       assets,
		window:       nil,
		movesText:    nil,
		coverageText: nil,
		level:        nil,
		board:        NewBoard(cfg.Size),
		dotBounds:    make(map[Dot]sdl.Rect),
		lineBounds:   make(map[Line]sdl.Rect),
		state:        newState(),
		log:          zap.NewNop(),
		Moves:        0,
		coverage:     0,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

//
func WithWindow(window *sdl.Window) option {
	return func(g *Game) {
		g.window = window
	}
}

//
func WithMoveText(text *graphics.Text) option {
	return func(g *Game) {
		g.movesText = text
	}
}

//
func WithCoverageText(text *graphics.Text) option {
	return func(g *Game) {
		g.coverageText = text
	}
}

// WithLogger creates a game and sets the logger.
func WithLogger(log *zap.Logger) option { //nolint
	return func(g *Game) {
		g.log = log
	}
}

// WithLevel creates a game and sets the level.
func WithLevel(l *Level) option { //nolint
	return func(g *Game) {
		for _, dot := range l.Dots {
			d := Dot{
				Location: dot.Location,
				Color:    dot.Color,
			}
			g.board.InitPath(d)

			gdot := g.assets.Dots[dot.Color]
			rc := sdl.Rect{
				X: g.assets.Grid.Bounds().X + dot.Location.X*g.config.SquareSize + (g.config.SquareSize-gdot.Bounds().W)/2,
				Y: g.assets.Grid.Bounds().Y + dot.Location.Y*g.config.SquareSize + (g.config.SquareSize-gdot.Bounds().H)/2,
				W: gdot.Bounds().W,
				H: gdot.Bounds().H,
			}
			g.dotBounds[d] = rc
		}
		g.level = l
		g.coverage = int32(len(g.dotBounds))
	}
}

// WithFile creates a game and sets the name of the file
// where the level was loaded from.
func WithFile(file string) option { //nolint
	return func(g *Game) {
		g.file = file
	}
}

// Repeat repeats the current level.
func (g *Game) Repeat() {
	for k := range g.lineBounds {
		delete(g.lineBounds, k)
	}

	g.Completed = false
	g.Moves = 0
	g.coverage = int32(len(g.dotBounds))

	g.board.Clear()
	g.board.InitPaths(g.level.Dots)

	g.state.reset()
}

// Continue tries to move on to the next level.
// It also triggers the creation of a new grid graphics asset
// if the size of the board has changed.
func (g *Game) Continue(gr *graphics.Renderer) {
	s := strings.Split(g.file, ".")
	if len(s) != 2 {
		g.log.Fatal("Wrong file name for the game level", zap.String("file name", g.file))
	}

	next, err := strconv.Atoi(s[0])
	if err != nil {
		g.log.Fatal("Internal error", zap.Error(err))
	}
	next++

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get the current working directory", zap.Error(err))
	}

	nextSize := g.board.size

	var filePath string
	for {
		filePath = fmt.Sprintf("%s/data/%d/%d.json", dir, nextSize, next)
		g.log.Debug("Try to load the level from", zap.String("file path", filePath))

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			if nextSize >= maxBoardSize {
				g.log.Info("No more files to load the levels from")
				ui.GameOver(g.window) //nolint
				os.Exit(0)
			}
			nextSize++
			next = 0
			continue
		}
		break
	}
	g.file = fmt.Sprintf("%d.json", next)

	g.log.Debug("Level file", zap.String("path", filePath))

	l, err := LoadFromFile(filePath)
	if err != nil {
		g.log.Fatal("Failed to load the level from file",
			zap.String("file path", filePath),
			zap.Error(err),
		)
	}

	for line := range g.lineBounds {
		delete(g.lineBounds, line)
	}

	for dot := range g.dotBounds {
		delete(g.dotBounds, dot)
	}

	g.state.reset()
	g.level = nil

	g.board.Clear()
	g.board = nil
	g.board = NewBoard(l.Size)

	g.Completed = false
	g.Moves = 0
	g.coverage = int32(len(g.dotBounds))

	if g.config.Size != g.board.size {
		g.config.Size = g.board.size
		g.assets.Grid.Destroy()
		g.assets.Grid = nil
		g.assets.Grid = graphics.CreateGrid(gr, g.config)
	}

	WithLevel(l)(g)
}

// Draw renders all the graphics objects on a rendering target.
func (g *Game) Draw(r *graphics.Renderer) {
	if g.movesText != nil {
		g.movesText.Text = fmt.Sprintf("Moves %d", g.Moves)
		err := g.movesText.Draw(r, sdl.Point{X: 0, Y: 0})
		if err != nil {
			g.log.Fatal("Draw text (moves) failed", zap.Error(err))
		}
	}

	if g.coverageText != nil {
		c := g.coverage - int32(len(g.dotBounds))
		sz := (g.board.size * g.board.size) - int32(len(g.dotBounds))
		pc := int(float64(c) / float64(sz) * 100.0)
		g.coverageText.Text = fmt.Sprintf("Coverage: %d %%", pc)
		err := g.coverageText.Draw(r, sdl.Point{X: 0, Y: 40})
		if err != nil {
			g.log.Fatal("Draw text (coverge) failed", zap.Error(err))
		}
	}

	g.assets.Grid.Blit(r)

	for dot, rc := range g.dotBounds {
		gdot := g.assets.Dots[dot.Color]
		gdot.BlitTo(r, &rc)
	}

	for line, rc := range g.lineBounds {
		var l graphics.Renderable
		if line.From.X == line.To.X {
			l = g.assets.VertLines[line.Color]
		} else {
			l = g.assets.HorizLines[line.Color]
		}
		l.BlitTo(r, &rc)
	}
}

// MouseButtonDown handles the mouse button down events.
func (g *Game) MouseButtonDown(ev *sdl.MouseButtonEvent) {
	if ev.Button != sdl.BUTTON_LEFT {
		return
	}

	var ok bool
	grid, ok := g.assets.Grid.(*graphics.Grid)
	if !ok {
		g.log.Fatal("Wrong concrete type for the grid graphics asset")
	}

	cx, cy, inside := grid.ScreenToGrid(ev.X, ev.Y, g.config.SquareSize)
	if !inside {
		return
	}

	clr := *g.board.ColorAt(cx, cy)
	if clr == graphics.NoColor {
		return
	}

	c := Coordinate{cx, cy}
	dot := Dot{
		Location: c,
		Color:    clr,
	}

	var path *Path
	path, ok = g.board.Paths[dot]
	if !ok {
		return
	}

	if path.EndDot != nil && len(path.Lines) == 0 {
		path, ok = g.board.Paths[*path.EndDot]
		if !ok {
			return
		}

		if len(path.Lines) == 0 {
			g.log.Fatal("Internal error: inconsistent paths")
		}
	}

	if path.EndDot != nil {
		for _, line := range path.Lines {
			if line.From != path.StartDot.Location {
				*(g.board.ColorAt(line.From.X, line.From.Y)) = graphics.NoColor
			}

			l := Line{
				From:  line.From,
				To:    line.To,
				Color: clr,
			}
			delete(g.lineBounds, l)
		}

		endDot := *path.EndDot
		if p, ok := g.board.Paths[*path.StartDot]; ok {
			p.EndDot = nil
			p.Lines = nil
		}

		if p, ok := g.board.Paths[endDot]; ok {
			p.EndDot = nil
			p.Lines = nil
		}
	}

	g.state.srcDot = &dot
	g.state.path = path
	g.state.dstDot = nil
	g.state.color = clr
	g.state.square = c
	g.state.editingPath = true
}

// MouseButtonUp handles the mouse button up events.
func (g *Game) MouseButtonUp(ev *sdl.MouseButtonEvent) {
	if !g.state.editingPath || g.state.srcDot == nil {
		return
	}

	if g.state.dstDot != nil {
		path, ok := g.board.Paths[*g.state.dstDot]
		if !ok {
			g.log.Fatal("Destination dot not found in paths",
				zap.Int32("x", g.state.dstDot.Location.X),
				zap.Int32("y", g.state.dstDot.Location.Y))
		}
		path.StartDot = g.state.dstDot
		path.EndDot = g.state.srcDot

		g.Moves++
		if g.coverage == g.board.size*g.board.size {
			g.Completed = true
		}
	} else {
		path, ok := g.board.Paths[*g.state.srcDot]
		if ok && len(path.Lines) > 0 {
			for _, line := range path.Lines {
				*(g.board.ColorAt(line.To.X, line.To.Y)) = graphics.NoColor

				l := Line{
					From:  line.From,
					To:    line.To,
					Color: g.state.color,
				}
				delete(g.lineBounds, l)
			}
			path.Lines = nil
		}
	}

	g.coverage = g.board.Coverage()
	g.state.reset()
}

// MouseMove handles the mouse move events.
func (g *Game) MouseMove(ev *sdl.MouseMotionEvent) {
	if !g.state.editingPath {
		return
	}

	grid, ok := g.assets.Grid.(*graphics.Grid)
	if !ok {
		g.log.Fatal("Wrong concrete type for the grid graphics asset")
	}

	cx, cy, inside := grid.ScreenToGrid(ev.X, ev.Y, g.config.SquareSize)
	if !inside {
		return
	}

	c := NewCoord(cx, cy)
	if g.state.srcDot != nil && c != g.state.square {
		path, ok := g.board.Paths[*g.state.srcDot]
		if !ok {
			return
		}

		action := g.nextAction(g.state.square, c,
			g.state.color, *g.board.ColorAt(c.X, c.Y), path)
		switch action {
		case drawLine:
			g.addLine(g.state.square, c, g.state.color, path)
		case eraseLine:
			g.removeLine(g.state.square, c, g.state.color, path)
		case completePath:
			g.state.dstDot = &Dot{c, g.state.color}
			g.addLine(g.state.square, c, g.state.color, path)
		}
		g.state.square = c

		g.coverage = g.board.Coverage()
	}
}

func (g *Game) addLine(from, to Coordinate, clr graphics.Color, path *Path) {
	r := sdl.Rect{
		X: 0,
		Y: 0,
		W: g.config.SquareSize,
		H: g.config.SquareSize,
	}

	if from.X == to.X {
		r.X = g.assets.Grid.Bounds().X + from.X*g.config.SquareSize

		var dir int32
		if from.Y < to.Y {
			dir = 1
		} else {
			dir = -1
		}
		r.Y = g.assets.Grid.Bounds().Y + from.Y*g.config.SquareSize + dir*g.config.SquareSize/2

	}

	if from.Y == to.Y {
		var dir int32
		if from.X < to.X {
			dir = 1
		} else {
			dir = -1
		}
		r.X = g.assets.Grid.Bounds().X + from.X*g.config.SquareSize + dir*g.config.SquareSize/2
		r.Y = g.assets.Grid.Bounds().Y + from.Y*g.config.SquareSize
	}

	l := Line{
		From:  from,
		To:    to,
		Color: clr,
	}
	g.lineBounds[l] = r

	path.AddLine(from, to)
	if g.state.dstDot != nil {
		path.EndDot = g.state.dstDot
	}
	*(g.board.ColorAt(to.X, to.Y)) = clr
}

func (g *Game) removeLine(from, to Coordinate, clr graphics.Color, path *Path) {
	l := Line{
		From:  to,
		To:    from,
		Color: clr,
	}
	delete(g.lineBounds, l)
	path.RemoveLine(l.From, l.To)

	if g.state.dstDot == nil || (g.state.dstDot != nil && g.state.dstDot.Location != l.To) {
		*(g.board.ColorAt(from.X, from.Y)) = graphics.NoColor
	}

	if g.state.dstDot != nil {
		path.EndDot = nil
		g.state.dstDot = nil
	}
}

func (g *Game) nextAction(from, to Coordinate,
	clrSrc graphics.Color, clrDst graphics.Color, path *Path) drawAction {

	if len(path.Lines) > 0 {
		if from != path.Lines[len(path.Lines)-1].To {
			// we did not get here from the current path
			return none
		}
	}

	if clrDst == clrSrc {
		dot := Dot{
			Location: NewCoord(to.X, to.Y),
			Color:    clrSrc,
		}

		_, ok := g.board.Paths[dot]
		if ok && dot != *g.state.srcDot && g.state.dstDot == nil {
			return completePath
		}

		if path.ContainsLine(to, from) {
			// get back
			return eraseLine
		}
	}

	if clrDst == graphics.NoColor {
		if g.state.dstDot != nil {
			return none
		}

		var lastVisited Coordinate
		if len(path.Lines) > 0 {
			lastVisited = path.Lines[len(path.Lines)-1].To
		} else {
			lastVisited = path.StartDot.Location
		}

		if !((to.X == lastVisited.X && math.Abs(float64(to.Y-lastVisited.Y)) == 1) ||
			(to.Y == lastVisited.Y && math.Abs(float64(to.X-lastVisited.X)) == 1)) {
			// we can only draw horizontal and vertical lines
			return none
		}

		return drawLine
	}

	return none
}
