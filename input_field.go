package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type InputField struct {
	WM                  *WindowManager
	X, Y, Width, Cursor int
	FgColor, BgColor    termbox.Attribute
	Value               string
	HasFocus, IsEnabled bool
}

func NewInputField(wm *WindowManager) *InputField {
	return &InputField{
		WM:        wm,
		Width:     50,
		FgColor:   termbox.ColorDefault | termbox.AttrUnderline,
		BgColor:   termbox.ColorDefault,
		IsEnabled: true,
	}
}

func (f *InputField) Event(e termbox.Event) {
	if !f.IsEnabled {
		return
	}
	if e.Type == termbox.EventKey {
		switch {
		case e.Key == termbox.KeyHome:
			f.Home()
		case e.Key == termbox.KeyEnd:
			f.End()
		case e.Key == termbox.KeySpace:
			f.InsertRune(' ')
		case e.Key == termbox.KeyBackspace || e.Key == termbox.KeyBackspace2:
			f.Backspace()
		case e.Key == termbox.KeyDelete:
			f.Delete()
		case e.Key == termbox.KeyArrowLeft && f.Cursor > 0:
			f.Cursor--
		case e.Key == termbox.KeyArrowRight && f.Cursor < len(f.Value):
			f.Cursor++
		case e.Ch != 0:
			f.InsertRune(e.Ch)
		default:
			return
		}
		f.WM.Render()
	}
}

func (f *InputField) InsertRune(r rune) {
	f.Value = f.Value[:f.Cursor] + string(r) + f.Value[f.Cursor:]
	f.Cursor++
}

func (f *InputField) Home() {
	f.Cursor = 0
}

func (f *InputField) End() {
	f.Cursor = len(f.Value)
}

func (f *InputField) Render() {
	value := f.Value
	l := len(value)
	if l >= f.Width {
		value = fmt.Sprintf("...%s", value[l+4-f.Width:l])
	}
	value = fmt.Sprintf(fmt.Sprintf("%%-%ds", f.Width), value)
	f.WM.CursorX = f.X
	f.WM.CursorY = f.Y
	f.WM.Print(value, f.FgColor, f.BgColor)
	if f.IsEnabled && f.HasFocus {
		termbox.SetCursor(f.X+f.Cursor, f.Y)
	}
}

func (f *InputField) Backspace() {
	if f.Cursor > 0 {
		f.Value = f.Value[:f.Cursor-1] + f.Value[f.Cursor:]
		f.Cursor--
	}
}

func (f *InputField) Delete() {
	if f.Cursor < len(f.Value)-1 {
		f.Value = f.Value[:f.Cursor] + f.Value[f.Cursor+1:]
	}
}

func (f *InputField) Clear() {
	f.Value = ""
	f.Cursor = 0
	f.WM.Render()
}
