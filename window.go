package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"

	"code.google.com/p/goncurses"
)

type KeyInfo struct {
	Name, Info string
}

type Window interface {
	Init(wm *WindowManager)
	Title() string
	Render(ncw *goncurses.Window)
	GotChar(k1, k2 goncurses.Key)
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
	wm := &WindowManager{
		NCW:         ncw,
		ContentNCW:  ncw.Derived(y-2, x, 1, 0),
		WindowStack: []Window{},
	}
	wm.AddWindow(&MenuWindow{})
	return wm
}

func (wm *WindowManager) CurrentWindow() Window {
	l := len(wm.WindowStack)
	if l == 0 {
		return nil
	}
	return wm.WindowStack[l-1]
}

func (wm *WindowManager) AddWindow(w Window) {
	w.Init(wm)
	wm.WindowStack = append(wm.WindowStack, w)
	wm.Render()
}

func (wm *WindowManager) RemoveWindow(w Window) {
	for i, wsw := range wm.WindowStack {
		if wsw == w {
			wm.WindowStack = append(
				wm.WindowStack[:i],
				wm.WindowStack[i+1:]...,
			)
			wm.Render()
			return
		}
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
				if w, h, err := GetTerminalSize(); err == nil {
					wm.renderMut.Lock()
					wm.NCW.Resize(h, w)
					wm.ContentNCW.Resize(h-2, w)
					wm.renderMut.Unlock()
					wm.Render()
				}
			default:
				running = false
			}
		}
	}()
	for running {
		k1 := wm.NCW.GetChar()
		wm.NCW.Timeout(0)
		k2 := wm.NCW.GetChar()
		wm.NCW.Timeout(-1)
		logger.Printf("Got %d %d", k1, k2)
		switch k1 {
		case goncurses.KEY_F12:
			running = false
		default:
			if cur := wm.CurrentWindow(); cur != nil {
				cur.GotChar(k1, k2)
			}
		}
	}
}

func (wm *WindowManager) Render() {
	wm.renderMut.Lock()
	defer wm.renderMut.Unlock()
	wm.NCW.Clear()
	wm.RenderHeader()
	wm.RenderFooter()
	if cur := wm.CurrentWindow(); cur != nil {
		cur.Render(wm.ContentNCW)
	}
	wm.NCW.Refresh()
}

func (wm *WindowManager) RenderHeader() {
	titlePrefix := ""
	title := ""
	if cur := wm.CurrentWindow(); cur != nil {
		title = cur.Title()
	}
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
	y, x := wm.NCW.MaxYX()

	wm.NCW.ColorOn(ColorPairBrdgmeTitle)
	wm.NCW.HLine(y-1, 0, ' ', x)
	wm.NCW.ColorOff(ColorPairBrdgmeTitle)

	wm.NCW.Move(y-1, 0)

	keyInfo := []KeyInfo{}
	if cur := wm.CurrentWindow(); cur != nil {
		keyInfo = append(keyInfo, cur.KeyInfo()...)
	}
	keyInfo = append(keyInfo, wm.KeyInfo()...)

	wm.NCW.AttrOn(goncurses.A_BOLD)
	for _, ki := range keyInfo {
		wm.NCW.ColorOn(ColorPairKeyInfoName)
		wm.NCW.Printf("%s ", ki.Name)
		wm.NCW.ColorOff(ColorPairKeyInfoName)

		wm.NCW.ColorOn(ColorPairKeyInfoInfo)
		wm.NCW.Printf("%s  ", ki.Info)
		wm.NCW.ColorOff(ColorPairKeyInfoInfo)
	}
	wm.NCW.AttrOff(goncurses.A_BOLD)
}

func (wm *WindowManager) KeyInfo() []KeyInfo {
	return []KeyInfo{
		{"Esc", "Menu"},
		{"F12", "Quit"},
	}
}

func GetTerminalSize() (width, height int, err error) {
	var dimensions [4]uint16
	if _, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&dimensions)),
		0,
		0,
		0,
	); err != 0 {
		return -1, -1, err
	}
	return int(dimensions[1]), int(dimensions[0]), nil
}
