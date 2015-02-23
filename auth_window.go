package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/nsf/termbox-go"
)

type AuthWindow struct {
	WM *WindowManager

	ShowConfirm                            bool
	EmailInput, ConfirmInput, FocusedInput *InputField
}

func NewAuthWindow() *AuthWindow {
	return &AuthWindow{}
}

func (w *AuthWindow) Init(wm *WindowManager) {
	w.WM = wm
	w.EmailInput = NewInputField(wm)
	w.EmailInput.Width = 50
	w.ConfirmInput = NewInputField(wm)
	w.ConfirmInput.Width = 50
	w.FocusedInput = w.EmailInput
}

func (w *AuthWindow) Title() string {
	return "Log in"
}

func (w *AuthWindow) Render() {
	fields := []struct {
		Label string
		Field *InputField
	}{
		{"Email", w.EmailInput},
		{"Confirm code", w.ConfirmInput},
	}
	l := len(fields)
	height := l*2 - 1
	maxLabelWidth := 0
	maxFieldWidth := 0
	maxWidth := 0
	for _, f := range fields {
		ll := len(f.Label)
		fl := f.Field.Width
		if ll > maxLabelWidth {
			maxLabelWidth = ll
		}
		if fl > maxFieldWidth {
			maxFieldWidth = fl
		}
		if total := ll + fl + 2; total > maxWidth {
			maxWidth = total
		}
	}
	xStart := (w.WM.SizeX - maxWidth) / 2
	y := (w.WM.SizeY - height) / 2
	for _, f := range fields {
		ll := len(f.Label)
		w.WM.CursorX = xStart + maxLabelWidth - ll
		w.WM.CursorY = y
		w.WM.Print(
			fmt.Sprintf("%s: ", f.Label),
			termbox.AttrBold,
			termbox.ColorDefault,
		)
		f.Field.X = xStart + maxLabelWidth + 2
		f.Field.Y = y
		f.Field.Render()
		y += 2
	}
}

func (w *AuthWindow) Event(e termbox.Event) {
	switch {
	case e.Type == termbox.EventKey && e.Key == termbox.KeyF2:
		w.ToggleConfirmField()
	case e.Type == termbox.EventKey && e.Key == termbox.KeyEnter:
		w.LogIn()
	case e.Type == termbox.EventKey && e.Key == termbox.KeyTab:
		w.NextField()
	default:
		w.FocusedInput.Event(e)
	}
}

func (w *AuthWindow) NextField() {
	fields := w.Fields()
	l := len(fields)
	if l == 0 {
		return
	}
	cur := 0
	for i, f := range fields {
		if f.HasFocus {
			cur = i
			f.HasFocus = false
			break
		}
	}
	next := fields[(cur+1)%l]
	next.HasFocus = true
	w.FocusedInput = next
	w.WM.Render()
}

func (w *AuthWindow) Fields() []*InputField {
	return []*InputField{
		w.EmailInput,
		w.ConfirmInput,
	}
}

func (w *AuthWindow) LogIn() {
	if w.ConfirmInput.Value != "" {
		resp, err := http.PostForm(HTTPUrl("/auth/confirm"), url.Values{
			"email":        {w.EmailInput.Value},
			"confirmation": {w.ConfirmInput.Value},
		})
		if err != nil {
			log.Printf("Error confirming auth, %v", err)
			return
		}
		defer resp.Body.Close()
		r := json.NewDecoder(resp.Body)
		if err := r.Decode(&Token); err != nil {
			log.Printf("Error reading token response, %v", err)
			return
		}
		log.Printf("Logged in with token %s", Token)
		Email = w.EmailInput.Value
		if ConfigDir != "" {
			if err := SaveConfig(); err != nil {
				log.Printf("Could not save config, %v", err)
			}
		}
		w.WM.ClearWindows()
		w.WM.AddWindow(NewMenuWindow())
	} else if w.EmailInput.Value != "" {
		if _, err := http.PostForm(HTTPUrl("/auth/request"), url.Values{
			"email": {w.EmailInput.Value},
		}); err != nil {
			log.Printf("Error requesting auth, %v", err)
			return
		}
	}
}

func (w *AuthWindow) ToggleConfirmField() {
	w.ShowConfirm = !w.ShowConfirm
	w.ConfirmInput.IsEnabled = w.ShowConfirm
}

func (w *AuthWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{
		{"Tab", "Next field"},
		{"Enter", "Submit"},
		{"F2", "Show confirm field"},
	}
}
