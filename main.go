package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"

	"code.google.com/p/goncurses"
)

var logger = log.New(os.Stderr, "", 0)

func mainOld() {
	if os.Getenv("ESCDELAY") == "" {
		os.Setenv("ESCDELAY", "25")
	}
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

func layout(g *gocui.Gui) error {
	if v, err := g.SetView("center", 0, 0, 50, 50); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Wrap = true
		v.FgColor = gocui.ColorBlue | gocui.AttrBold
		fmt.Fprintln(v, "This is a test, does this mean that text wraps? This is a real question.")
		v.FgColor = gocui.ColorRed | gocui.AttrBold
		fmt.Fprintln(v, "This is a test, does this mean that text wraps? This is a real question.")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func main() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panic(err)
	}
	defer g.Close()

	g.FgColor = gocui.ColorDefault
	g.BgColor = gocui.ColorDefault

	g.SetLayout(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panic(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.Quit {
		log.Panic(err)
	}
}
