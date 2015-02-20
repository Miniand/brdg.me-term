package main

import (
	"code.google.com/p/goncurses"
)

type MenuWindow struct {
	Selected int
}

func (w *MenuWindow) Title() string {
	return "Menu"
}

func (w *MenuWindow) Render(ncw *goncurses.Window) {
}

func (w *MenuWindow) GotChar(k goncurses.Key) {
	logger.Printf("Got key: %v", k)
}

func (w *MenuWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{}
}
