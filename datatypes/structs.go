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
	UpdateRate     int
	LastFmEnabled  bool
	LastFmUsername string
	LastFmKey      string
}
