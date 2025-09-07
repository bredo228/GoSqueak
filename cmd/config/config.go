package config

type Config struct {
	OscPort             int    `json:"oscPort"`
	WebsocketListenPort int    `json:"websocketListenPort"`
	UpdateRate          int    `json:"updateRate"`
	LastFmEnabled       bool   `json:"lastFmEnabled"`
	LastFmUsername      string `json:"lastFmUsername"`
	LastFmApiKey        string `json:"lastFmApiKey"`
}
