package main

import (
	"log"
	"os"

	"code.google.com/p/goncurses"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	window, err := goncurses.Init()
	if err != nil {
		log.Fatal("Error initialising ncurses:", err)
	}
	defer goncurses.End()
	initColors()
	goncurses.CBreak(true)
	goncurses.Echo(false)
	window.SetBackground(goncurses.ColorPair(ColorPairBackground))
	window.Keypad(true)
	NewWindowManager(window).Run()
}
