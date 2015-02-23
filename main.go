package main

import (
	"flag"
	"log"
	"os/user"
	"path"

	"github.com/nsf/termbox-go"
)

var (
	Server string
	SSL    bool

	Email, Token string

	ConfigDir string
)

func main() {
	flag.StringVar(&Server, "server", "api.brdg.me",
		"Server to connect to")
	flag.BoolVar(&SSL, "ssl", true, "Whether to use SSL")
	defaultConfigDir := ""
	if usr, err := user.Current(); err == nil {
		defaultConfigDir = path.Join(usr.HomeDir, ".brdg.me")
	}
	flag.StringVar(&ConfigDir, "config-dir", defaultConfigDir,
		"Directory to save config")
	flag.Parse()
	if ConfigDir != "" {
		if err := CreateConfigDirIfNotExists(); err != nil {
			log.Printf("Could not make config dir, %v", err)
		} else if err := LoadConfig(); err != nil {
			log.Printf("Could not load config, %v", err)
		}
	}
	if err := termbox.Init(); err != nil {
		log.Panic(err)
	}
	defer termbox.Close()
	NewWindowManager().Run()
}
