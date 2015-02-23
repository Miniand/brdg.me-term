package main

import (
	"fmt"
	"io"
	"net/http"
)

var Client = http.Client{}

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

func NewAuthRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, HTTPUrl(path), body)
	if err != nil {
		return req, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", Token))
	return req, nil
}
