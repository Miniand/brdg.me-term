package main

import (
	"code.google.com/p/goncurses"
)

type GameWindow struct {
	GameID string
}

func (g *GameWindow) Title() string {
	return "Zombie Dice"
}

func (g *GameWindow) Render(ncw *goncurses.Window) {
	goncurses.Cursor(0)
	ncw.Box(0, 0)
}

func (g *GameWindow) GotChar(k goncurses.Key) {
	logger.Printf("Got key: %v", k)
}

func (g *GameWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{
		{"Arrows", "Scroll"},
		{"F1", "Log"},
		{"F2", "Cmd list"},
	}
}
