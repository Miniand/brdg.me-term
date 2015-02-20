package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"code.google.com/p/goncurses"
)

type KeyInfo struct {
	Name, Info string
}

type Window interface {
	Title() string
	Render(ncw *goncurses.Window)
	GotChar(k goncurses.Key)
	KeyInfo() []KeyInfo
}

type WindowManager struct {
	NCW         *goncurses.Window
	ContentNCW  *goncurses.Window
	WindowStack []Window

	renderMut sync.Mutex
}

func NewWindowManager(ncw *goncurses.Window) *WindowManager {
	y, x := ncw.MaxYX()
	return &WindowManager{
		NCW:         ncw,
		ContentNCW:  ncw.Derived(y-2, x, 1, 0),
		WindowStack: []Window{&GameWindow{}},
	}
}

func (wm *WindowManager) Run() {
	running := true
	wm.Render()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGWINCH, os.Interrupt, os.Kill)
		for running {
			switch <-c {
			case syscall.SIGWINCH:
				wm.Render()
			default:
				running = false
			}
		}
	}()
	for running {
		c := wm.NCW.GetChar()
		switch c {
		case goncurses.KEY_F12:
			running = false
		default:
			wm.WindowStack[0].GotChar(c)
		}
	}
}

func (wm *WindowManager) Render() {
	wm.renderMut.Lock()
	defer wm.renderMut.Unlock()
	goncurses.Cursor(0)
	wm.NCW.Clear()
	wm.RenderHeader()
	wm.WindowStack[0].Render(wm.ContentNCW)
	wm.NCW.Refresh()
}

func (wm *WindowManager) RenderHeader() {
	titlePrefix := ""
	title := wm.WindowStack[0].Title()
	if title != "" {
		titlePrefix = " - "
	}
	fullTitle := fmt.Sprintf(" brdg.me%s%s ", titlePrefix, title)
	_, x := wm.NCW.MaxYX()
	wm.NCW.AttrOn(goncurses.A_BOLD)
	wm.NCW.ColorOn(ColorPairBrdgmeTitle)
	wm.NCW.HLine(0, 0, ' ', x)
	wm.NCW.MovePrintf(0, (x-len(fullTitle))/2, fullTitle)
	wm.NCW.ColorOff(ColorPairBrdgmeTitle)
	wm.NCW.AttrOff(goncurses.A_BOLD)
}

func (wm *WindowManager) RenderFooter() {
}
