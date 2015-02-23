package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

var Summary GameSummary

type Log struct {
	Time time.Time `json:"time"`
	Text string    `json:"text"`
}

type Game struct {
	Commands   string    `json:"commands"`
	FinishedAt time.Time `json:"finishedAt"`
	Game       string    `json:"game"`
	ID         string    `json:"id"`
	Identifier string    `json:"identifier"`
	IsFinished bool      `json:"isFinished"`
	Log        []Log     `json:"log"`
	Name       string    `json:"name"`
	PlayerList []string  `json:"playerList"`
	WhoseTurn  []string  `json:"whoseTurn"`
	Winners    []string  `json:"winners"`
}

type GameSummary struct {
	CurrentTurn      []Game `json:"currentTurn"`
	OtherActive      []Game `json:"otherActive"`
	RecentlyFinished []Game `json:"recentlyFinished"`
}

type MenuWindow struct {
	WM       *WindowManager
	Selected int
}

func NewMenuWindow() *MenuWindow {
	return &MenuWindow{}
}

func (w *MenuWindow) Init(wm *WindowManager) {
	w.WM = wm
	w.UpdateGames()
}

func (w *MenuWindow) UpdateGames() {
	go func() {
		req, err := NewAuthRequest("GET", "/game/summary?renderer=raw", nil)
		if err != nil {
			log.Printf("Unable to create game summary request, %v", err)
		}
		resp, err := Client.Do(req)
		if err != nil {
			log.Printf("Unable to request game summary, %v", err)
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&Summary); err != nil {
			log.Printf("Unable to parse summary, %v", err)
		}
		w.WM.Render()
	}()
}

func (w *MenuWindow) Title() string {
	return "Menu"
}

func (w *MenuWindow) Render() {
	termbox.HideCursor()
	y := 2
	for _, l := range [][]Game{
		Summary.CurrentTurn,
		Summary.OtherActive,
		Summary.RecentlyFinished,
	} {
		for _, g := range l {
			if y >= w.WM.SizeY-1 {
				return
			}
			w.WM.CursorX = 0
			w.WM.CursorY = y
			w.WM.Print(
				fmt.Sprintf("%s with %s", g.Name, strings.Join(g.PlayerList, ", ")),
				termbox.ColorDefault,
				termbox.ColorDefault,
			)
			y++
		}
		y++
	}
}

func (w *MenuWindow) Event(e termbox.Event) {
	if e.Type == termbox.EventKey && e.Key == termbox.KeyEnter {
		w.WM.AddWindow(NewGameWindow("Roll Through the Ages"))
	}
}

func (w *MenuWindow) KeyInfo() []KeyInfo {
	return []KeyInfo{}
}
