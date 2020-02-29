package main

import (
	"connect-dots/config"
	"connect-dots/game"
	"connect-dots/graphics"
	"connect-dots/ui"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	var (
		size int
	)

	flag.IntVar(&size, "size", 5, "the board size")
	flag.Parse()

	log, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to create a zap logger", zap.Error(err))
	}
	defer log.Sync() //nolint

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Fatal("Failed to initialize SDL", zap.Error(err))
	}
	defer sdl.Quit()

	withSize := func(sz int32) config.Option {
		return func(c *config.Config) { c.Size = sz }
	}
	config := config.New(withSize(int32(size)))

	window, err := sdl.CreateWindow("dots connected",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		config.WindowWidth, config.WindowHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatal("Failed to create window", zap.Error(err))
	}
	defer window.Destroy() //nolint

	if err := ttf.Init(); err != nil {
		log.Fatal("Failed to init TTF API", zap.Error(err))
	}
	defer ttf.Quit()

	font, err := ttf.OpenFont("data/fonts/test.ttf", 32)
	if err != nil {
		log.Fatal("Failed to open font", zap.Error(err))
	}
	defer font.Close()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatal("Failed to create the SDL renderer", zap.Error(err))
	}

	if err := renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); err != nil {
		log.Fatal("Failed to set the blend mode for the renderer", zap.Error(err))
	}

	gr := graphics.NewRenderer(renderer, log)
	defer gr.Destroy()

	storage := graphics.NewAssetsStorage()
	if err := storage.Init(gr, config); err != nil {
		log.Fatal("Failed to load the graphics assets", zap.Error(err))
	}
	defer storage.Destroy()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get the current working directory", zap.Error(err))
	}

	fileName := "0.json"
	path := fmt.Sprintf("%s/data/%d/%s", dir, size, fileName)
	l, err := game.LoadFromFile(path)
	if err != nil {
		log.Fatal("Failed to load the level", zap.Error(err))
	}

	game := game.New(config, storage,
		game.WithWindow(window),
		game.WithMoveText(graphics.NewText("Moves: 0", font)),
		game.WithCoverageText(graphics.NewText("Coverage: 0%", font)),
		game.WithLogger(log),
		game.WithLevel(l),
		game.WithFile(fileName),
	)

	running := true
	for running {
		gr.SetDrawColor(0, 0, 0, 0)
		gr.Clear()

		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch t := ev.(type) {
			case *sdl.QuitEvent:
				running = false

			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONDOWN {
					game.MouseButtonDown(t)
					continue
				}

				if t.Type == sdl.MOUSEBUTTONUP {
					game.MouseButtonUp(t)
					continue
				}

			case *sdl.MouseMotionEvent:
				game.MouseMove(t)
			}
		}

		game.Draw(gr)
		gr.Present()
		sdl.Delay(5)

		if game.Completed {
			action, err := ui.LevelCompletedBox(game.Moves, window)
			if err != nil {
				log.Fatal("Internal error", zap.Error(err))
			}

			switch action {
			case ui.Continue:
				game.Continue(gr)
			case ui.Repeat:
				game.Repeat()
			case ui.Quit:
				os.Exit(0)
			}
		}
	}
	os.Exit(0)
}
