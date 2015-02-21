package main

import "code.google.com/p/goncurses"

type GameWindow struct {
	WM     *WindowManager
	GameID string
}

func (w *GameWindow) Init(wm *WindowManager) {
	w.WM = wm
}

func (w *GameWindow) Title() string {
	return "Zombie Dice"
}

func (w *GameWindow) Render(ncw *goncurses.Window) {
	goncurses.Cursor(0)
	ncw.Move(0, 0)
	ncw.Print("BOAH!")
}

func (w *GameWindow) GotChar(k1, k2 goncurses.Key) {
	switch k1 {
	case 27:
		if k2 == 0 {
			// Esc
			logger.Print("QUIIIT")
			w.WM.RemoveWindow(w)
		}
	}
}

func (w *GameWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{
		{"Arrows", "Scroll"},
		{"Tab", "Next game"},
		{"F1", "Log"},
		{"F2", "Cmd list"},
		{"Esc", "Menu"},
	}
}
