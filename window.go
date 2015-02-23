package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/nsf/termbox-go"
)

type KeyInfo struct {
	Name, Info string
}

type Window interface {
	Init(wm *WindowManager)
	Title() string
	Render()
	Event(e termbox.Event)
	KeyInfo() []KeyInfo
}

type WindowManager struct {
	WindowStack      []Window
	ColourStack      []string
	BoldStack        int
	CursorX, CursorY int
	SizeX, SizeY     int

	renderMut sync.Mutex
}

func NewWindowManager() *WindowManager {
	x, y := termbox.Size()
	wm := &WindowManager{
		WindowStack: []Window{},
		ColourStack: []string{},
		SizeX:       x,
		SizeY:       y,
	}
	if Token != "" {
		wm.AddWindow(NewMenuWindow())
	} else {
		wm.LogOut()
	}
	return wm
}

func (wm *WindowManager) LogOut() {
	Token = ""
	if ConfigDir != "" {
		if err := SaveConfig(); err != nil {
			log.Printf("Could not save config, %v", err)
		}
	}
	wm.ClearWindows()
	wm.AddWindow(NewAuthWindow())
}

func (wm *WindowManager) ClearWindows() {
	wm.WindowStack = []Window{}
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
		signal.Notify(c, syscall.SIGWINCH)
		for running {
			<-c
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			x, y := termbox.Size()
			wm.SizeX = x
			wm.SizeY = y
			wm.Render()
		}
	}()
	wm.Render()
	for running {
		e := termbox.PollEvent()
		if e.Type == termbox.EventKey && e.Key == termbox.KeyF12 {
			running = false
		} else {
			if c := wm.CurrentWindow(); c != nil {
				c.Event(e)
			}
		}
	}
}

func (wm *WindowManager) Render() {
	wm.renderMut.Lock()
	defer wm.renderMut.Unlock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	wm.RenderHeader()
	wm.RenderFooter()
	if c := wm.CurrentWindow(); c != nil {
		c.Render()
	}
	termbox.Flush()
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
	fullTitle := fmt.Sprintf("brdg.me%s%s", titlePrefix, title)

	wm.PrintHLine(' ', 0, termbox.ColorDefault, termbox.ColorBlack)
	wm.CursorX = (wm.SizeX - len(fullTitle)) / 2
	wm.CursorY = 0
	wm.Print(fullTitle, termbox.ColorGreen|termbox.AttrBold, termbox.ColorBlack)
}

func (wm *WindowManager) RenderFooter() {
	wm.PrintHLine(' ', wm.SizeY-1, termbox.ColorDefault, termbox.ColorBlack)

	wm.CursorX = 0
	wm.CursorY = wm.SizeY - 1

	keyInfo := []KeyInfo{}
	if cur := wm.CurrentWindow(); cur != nil {
		keyInfo = append(keyInfo, cur.KeyInfo()...)
	}
	keyInfo = append(keyInfo, wm.KeyInfo()...)

	for _, ki := range keyInfo {
		wm.Print(fmt.Sprintf("%s ", ki.Name),
			termbox.ColorYellow|termbox.AttrBold, termbox.ColorBlack)
		wm.Print(fmt.Sprintf("%s  ", ki.Info),
			termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack)
	}
}

var actionRegexp = regexp.MustCompile(`^\{\{([_a-z]+)(\s+['"]?(.+?)['"]?)?\}\}`)

func (wm *WindowManager) Print(input string, fg, bg termbox.Attribute) {
	runes := bytes.Runes([]byte(input))
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '{' {
			// Read ahead and see if it's for formatting.
			if match := actionRegexp.FindStringSubmatch(string(runes[i:])); match != nil {
				// Skip
				i += len(bytes.Runes([]byte(match[0]))) - 1
				switch match[1] {
				case "b":
					wm.BoldStack++
				case "_b":
					wm.BoldStack--
				case "c":
					wm.ColourStack = append(wm.ColourStack, match[3])
				case "_c":
					if l := len(wm.ColourStack); l > 0 {
						wm.ColourStack = wm.ColourStack[:l-1]
					}
				}
				continue
			}
		}
		f := fg
		if l := len(wm.ColourStack); l > 0 {
			if c, ok := Colours[wm.ColourStack[l-1]]; ok {
				f = c
			}
		}
		if wm.BoldStack > 0 {
			f |= termbox.AttrBold
		}
		wm.PrintRune(r, f, bg)
	}
}

func (wm *WindowManager) PrintRune(r rune, fg, bg termbox.Attribute) {
	termbox.SetCell(wm.CursorX, wm.CursorY, r, fg, bg)
	wm.CursorX++
}

func (wm *WindowManager) PrintLine(r rune, x1, y1, x2, y2 int, fg, bg termbox.Attribute) {
	steps := max(abs(x2-x1), abs(y2-y1))
	xPer := float64(x2-x1) / float64(steps)
	yPer := float64(y2-y1) / float64(steps)
	for i := 0; i <= steps; i++ {
		termbox.SetCell(
			x1+int(xPer*float64(i)+0.5),
			y1+int(yPer*float64(i)+0.5),
			r, fg, bg)
	}
}

func (wm *WindowManager) PrintVLine(r rune, x int, fg, bg termbox.Attribute) {
	wm.PrintLine(r, x, 0, x, wm.SizeY-1, fg, bg)
}

func (wm *WindowManager) PrintHLine(r rune, y int, fg, bg termbox.Attribute) {
	wm.PrintLine(r, 0, y, wm.SizeX-1, y, fg, bg)
}

func (wm *WindowManager) KeyInfo() []KeyInfo {
	return []KeyInfo{
		{"F12", "Quit"},
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func unit(i int) int {
	if i == 0 {
		return 0
	}
	return abs(i) / i
}
