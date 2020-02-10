package main

import (
	"connect-dots/config"
	"connect-dots/game"
	"connect-dots/graphics"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	var (
		size  int
		level int
	)

	flag.IntVar(&size, "size", 5, "the board size")
	flag.IntVar(&level, "level", int(config.Easy), "the dificulty level")
	flag.Parse()

	if level < int(config.Easy) {
		level = int(config.Easy)
	}

	if level > int(config.Advanced) {
		level = int(config.Advanced)
	}

	log, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to create a zap logger", zap.Error(err))
	}
	defer log.Sync()

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Fatal("Failed to initialize SDL", zap.Error(err))
	}
	defer sdl.Quit()

	withSize := func(sz int32) config.Option {
		return func(c *config.Config) { c.Size = sz }
	}
	withLevel := func(l config.Difficulty) config.Option {
		return func(c *config.Config) { c.Difficulty = l }
	}
	config := config.New(
		withSize(int32(size)),
		withLevel(config.Difficulty(level)))

	window, err := sdl.CreateWindow("dots connected",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		config.WindowWidth, config.WindowHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatal("Failed to create window", zap.Error(err))
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
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

	fileName := fmt.Sprintf("%d.json", level)
	path := fmt.Sprintf("%s/data/%dx%d/level%d/%s", dir, size, size, level, fileName)
	log.Debug("Level file", zap.String("file", path))

	l, err := game.LoadFromFile(path)
	if err != nil {
		log.Fatal("Failed to load the level", zap.Error(err))
	}

	game := game.New(config, storage,
		game.WithLogger(log), game.WithLevel(l), game.WithFile(fileName))

	running := true
	for running {
		gr.SetDrawColor(0, 0, 0, 0)
		gr.Clear()

		for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
			switch t := ev.(type) {
			case *sdl.QuitEvent:
				running = false
				break

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
		sdl.Delay(10)
	}
}
