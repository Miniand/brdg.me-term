package main

import "github.com/nsf/termbox-go"

type GameWindow struct {
	WM     *WindowManager
	GameID string

	ShowLog, ShowCommands    bool
	GameScrollX, GameScrollY int
}

func NewGameWindow(gameID string) *GameWindow {
	return &GameWindow{
		GameID:       gameID,
		ShowLog:      true,
		ShowCommands: true,
	}
}

func (w *GameWindow) Init(wm *WindowManager) {
	w.WM = wm
}

func (w *GameWindow) Title() string {
	return w.GameID
}

func (w *GameWindow) Render() {
	if w.ShowLog {
		x := w.WM.SizeX * 2 / 3
		w.WM.PrintLine('║', x, 1, x, w.WM.SizeY-2,
			termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (w *GameWindow) Event(e termbox.Event) {
	switch {
	case e.Type == termbox.EventKey && e.Key == termbox.KeyEsc:
		w.WM.RemoveWindow(w)
	case e.Type == termbox.EventKey && e.Key == termbox.KeyF1:
		w.ShowLog = !w.ShowLog
		w.WM.Render()
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
