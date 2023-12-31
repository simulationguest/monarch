package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
)

type Site struct {
	Hostname string
	Port     string
}

func NewSite(u *url.URL) Site {
	return Site{Hostname: u.Hostname(), Port: u.Port()}
}

type Config = map[string]Action

// type Path = string
// type Site = map[Path]Action
type Action struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data"`
}

const (
	ActionReverseProxy = "reverse_proxy"
	ActionRespond      = "respond"
)

func readConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(f, config)
	return config, err
}

func getConfig() Config {
	if len(os.Args) < 2 {
		log.Fatalln("No configuration file provided")
	}

	configPath := os.Args[1]
	config, err := readConfig(configPath)
	if err != nil {
		log.Fatalln("Error reading configuration file:", err)
	}
	return *config
}
