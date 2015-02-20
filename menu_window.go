package main

import (
	"code.google.com/p/goncurses"
)

type MenuWindow struct {
	WM       *WindowManager
	Selected int
}

func (w *MenuWindow) Init(wm *WindowManager) {
	w.WM = wm
}

func (w *MenuWindow) Title() string {
	return "Menu"
}

func (w *MenuWindow) Render(ncw *goncurses.Window) {
}

func (w *MenuWindow) GotChar(k1, k2 goncurses.Key) {
	switch k1 {
	case '\n':
		logger.Print("GAAAAME")
		w.WM.AddWindow(&GameWindow{})
	}
}

func (w *MenuWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{}
}
