package datatypes

type Track struct {
	Title    string
	Album    string
	Artist   string
	Artwork  string
	Duration int
	Position int
}

type Config struct {
	Port           int     `json:"port"`
	UpdateRate     float32 `json:"updateRate"`
	LastFmEnabled  bool    `json:"lastFmEnabled"`
	LastFmUsername string  `json:"lastFmUsername"`
	LastFmApiKey   string  `json:"lastFmApiKey"`
}
