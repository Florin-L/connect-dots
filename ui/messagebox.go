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
func LevelCompletedBox(moves int32) (int32, error) {
	buttons := []sdl.MessageBoxButtonData{
		{sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT, Continue, "Continue"},
		{0, Repeat, "Repeat"},
		{sdl.MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT, Quit, "Quit"},
	}

	text := fmt.Sprintf("Completed the level in %d moves", int(moves))
	mbdata := sdl.MessageBoxData{
		sdl.MESSAGEBOX_INFORMATION,
		nil,
		"Level completed",
		text,
		buttons,
		nil,
	}

	id, err := sdl.ShowMessageBox(&mbdata)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func GameOver() error {
	buttons := []sdl.MessageBoxButtonData{
		{sdl.MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT, Ok, "Ok"},
	}

	mbdata := sdl.MessageBoxData{
		sdl.MESSAGEBOX_INFORMATION,
		nil,
		"Game over",
		"No more levels to be played. The game will be terminated.",
		buttons,
		nil,
	}

	_, err := sdl.ShowMessageBox(&mbdata)
	if err != nil {
		return err
	}

	return nil
}
