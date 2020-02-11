package ui

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	// Move on to the next level
	Continue int32 = 0

	// Repeat the completed level
	Repeat int32 = 1

	// Quit the game
	Quit int32 = 2

	// Ok confirms the action
	Ok int32 = 3
)

// LevelCompletedBox informs the user that the level gets completed.
// The user may choose to repeat the current level or to move on
// to the next level or to quit the game.
func LevelCompletedBox(moves int32, window *sdl.Window) (int32, error) {
	buttons := []sdl.MessageBoxButtonData{
		{Flags: sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT, ButtonID: Continue, Text: "Continue"},
		{Flags: 0, ButtonID: Repeat, Text: "Repeat"},
		{Flags: sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT, ButtonID: Quit, Text: "Quit"},
	}

	text := fmt.Sprintf("Completed the level in %d moves", int(moves))
	mbdata := sdl.MessageBoxData{
		Flags:       sdl.MESSAGEBOX_INFORMATION,
		Window:      window,
		Title:       "Level completed",
		Message:     text,
		Buttons:     buttons,
		ColorScheme: nil,
	}

	id, err := sdl.ShowMessageBox(&mbdata)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func GameOver(window *sdl.Window) error {
	buttons := []sdl.MessageBoxButtonData{
		{Flags: sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT, ButtonID: Ok, Text: "Ok"},
	}

	mbdata := sdl.MessageBoxData{
		Flags:       sdl.MESSAGEBOX_INFORMATION,
		Window:      window,
		Title:       "Game over",
		Message:     "No more levels to be played. The game will be terminated.",
		Buttons:     buttons,
		ColorScheme: nil,
	}

	_, err := sdl.ShowMessageBox(&mbdata)
	if err != nil {
		return err
	}

	return nil
}
