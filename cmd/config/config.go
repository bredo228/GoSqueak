package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	OscPort             int    `json:"oscPort"`
	WebsocketListenPort int    `json:"websocketListenPort"`
	UpdateRate          int    `json:"updateRate"`
	LastFmEnabled       bool   `json:"lastFmEnabled"`
	LastFmUsername      string `json:"lastFmUsername"`
	LastFmApiKey        string `json:"lastFmApiKey"`
}

func LoadConfig(filename string, config *Config) error {

	c, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(c, &config)
	return err

}
