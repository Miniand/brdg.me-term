package main

import (
	"log"
	"time"

	"code.google.com/p/goncurses"
)

func main() {
	window, err := goncurses.Init()
	if err != nil {
		log.Fatal("Error initialising ncurses:", err)
	}
	goncurses.Echo(false)
	window.MovePrint(0, 0, "brdg.me")
	window.Refresh()
	defer goncurses.End()
	time.Sleep(3 * time.Second)
}
