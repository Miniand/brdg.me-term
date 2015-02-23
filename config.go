package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

const (
	ConfigPrefsFile = "prefs.json"
)

type ConfigJSON struct {
	Email, Token string
}

func CreateConfigDirIfNotExists() error {
	_, err := os.Stat(ConfigDir)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("could not stat config directory, %v", err)
	}
	return os.MkdirAll(ConfigDir, 0700)
}

func LoadConfig() error {
	f, err := os.Open(path.Join(ConfigDir, ConfigPrefsFile))
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return fmt.Errorf("could not open prefs file, %v", err)
	}
	r := json.NewDecoder(f)
	conf := ConfigJSON{}
	if err := r.Decode(&conf); err != nil {
		return fmt.Errorf("could not read prefs file, %v", err)
	}
	Email = conf.Email
	Token = conf.Token
	return nil
}

func SaveConfig() error {
	conf := ConfigJSON{
		Email: Email,
		Token: Token,
	}
	f, err := os.OpenFile(
		path.Join(ConfigDir, ConfigPrefsFile),
		os.O_CREATE|os.O_WRONLY,
		0600,
	)
	if err != nil {
		return fmt.Errorf("could not open prefs file, %v", err)
	}
	w := json.NewEncoder(f)
	if err := w.Encode(conf); err != nil {
		return fmt.Errorf("could not write prefs file, %v", err)
	}
	return nil
}
