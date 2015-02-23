package main

import (
	"flag"
	"fmt"
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

func HTTPBase() string {
	proto := "http"
	if SSL {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s", proto, Server)
}

func HTTPUrl(path string) string {
	return fmt.Sprintf("%s%s", HTTPBase(), path)
}

func WSBase() string {
	proto := "ws"
	if SSL {
		proto = "wss"
	}
	return fmt.Sprintf("%s://%s", proto, Server)
}

func WSUrl(path string) string {
	return fmt.Sprintf("%s%s", WSBase(), path)
}

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
