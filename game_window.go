package main

import "github.com/nsf/termbox-go"

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

func (w *GameWindow) Render() {
}

func (w *GameWindow) Event(e termbox.Event) {
	if e.Type == termbox.EventKey && e.Key == termbox.KeyEsc {
		w.WM.RemoveWindow(w)
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
