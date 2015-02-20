package main

import (
	"code.google.com/p/goncurses"
)

const (
	ColorPairBackground = iota
	ColorPairBrdgmeTitle
	ColorPairKeyInfoName
	ColorPairKeyInfoInfo
)

func initColors() {
	goncurses.StartColor()
	goncurses.UseDefaultColors()
	goncurses.InitPair(ColorPairBackground, goncurses.C_WHITE, goncurses.C_BLACK)
	goncurses.InitPair(ColorPairBrdgmeTitle, goncurses.C_GREEN, goncurses.C_BLACK)
	goncurses.InitPair(ColorPairKeyInfoName, goncurses.C_YELLOW, goncurses.C_BLACK)
	goncurses.InitPair(ColorPairKeyInfoInfo, goncurses.C_WHITE, goncurses.C_BLACK)
}
