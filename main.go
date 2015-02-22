package main

import (
	"log"
	"os"

	"github.com/nsf/termbox-go"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	if err := termbox.Init(); err != nil {
		log.Panic(err)
	}
	defer termbox.Close()
	NewWindowManager().Run()
}
