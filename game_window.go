package main

import "github.com/nsf/termbox-go"

type GameWindow struct {
	WM     *WindowManager
	GameID string

	ShowLog, ShowCommands    bool
	GameScrollX, GameScrollY int

	CommandInput *InputField
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
	w.CommandInput = NewInputField(wm)
}

func (w *GameWindow) Title() string {
	return w.GameID
}

func (w *GameWindow) Render() {
	w.WM.PrintHLine('═', w.WM.SizeY-3,
		termbox.ColorWhite, termbox.ColorDefault)
	w.WM.CursorX = 0
	w.WM.CursorY = w.WM.SizeY - 2
	w.WM.Print(
		"Command: ",
		termbox.ColorWhite|termbox.AttrBold,
		termbox.ColorDefault,
	)
	w.CommandInput.X = w.WM.CursorX
	w.CommandInput.Y = w.WM.CursorY
	w.CommandInput.Width = w.WM.SizeX - w.WM.CursorX
	w.CommandInput.HasFocus = true
	w.CommandInput.Render()

	if w.ShowLog {
		x := w.WM.SizeX * 2 / 3
		w.WM.PrintLine('║', x, 1, x, w.WM.SizeY-4,
			termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, w.WM.SizeY-3, '╩',
			termbox.ColorWhite, termbox.ColorDefault)
	}
}

func (w *GameWindow) Event(e termbox.Event) {
	switch {
	case e.Type == termbox.EventKey && e.Key == termbox.KeyEsc:
		w.WM.RemoveWindow(w)
	case e.Type == termbox.EventKey && e.Key == termbox.KeyF1:
		w.ShowLog = !w.ShowLog
		w.WM.Render()
	case e.Type == termbox.EventKey && e.Key == termbox.KeyEnter:
		w.CommandInput.Clear()
	default:
		w.CommandInput.Event(e)
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
