package main

import "github.com/nsf/termbox-go"

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

func (w *MenuWindow) Render() {
}

func (w *MenuWindow) Event(e termbox.Event) {
	if e.Type == termbox.EventKey && e.Key == termbox.KeyEnter {
		w.WM.AddWindow(&GameWindow{
			GameID: "Roll Through the Ages",
		})
	}
}

func (w *MenuWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{}
}
