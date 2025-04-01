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
	Port           int
	UpdateRate     float32
	LastFmEnabled  bool
	LastFmUsername string
	LastFmKey      string
}
