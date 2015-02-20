package main

import (
	"code.google.com/p/goncurses"
)

const (
	ColorPairBackground = iota
	ColorPairBrdgmeTitle
)

func initColors() {
	goncurses.StartColor()
	goncurses.UseDefaultColors()
	goncurses.InitPair(ColorPairBackground, goncurses.C_WHITE, goncurses.C_BLACK)
	goncurses.InitPair(ColorPairBrdgmeTitle, goncurses.C_GREEN, goncurses.C_BLACK)
}
