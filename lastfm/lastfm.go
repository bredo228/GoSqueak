package lastfm

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	. "github.com/bredo228/GoSqueak/datatypes"
)

type LastFmTrackResponse struct {
	Track LastFmTrack `json:"track"`
}

type LastFmTrack struct {
	Name  string      `json:"name"`
	Album LastFmAlbum `json:"album"`
}

type LastFmAlbum struct {
	Artist string        `json:"artist"`
	Title  string        `json:"title"`
	Images []LastFmImage `json:"image"`
}

type LastFmImage struct {
	Text string `json:"#text"`
	Size string `json:"size"`
}

func GetTrackArtwork(track Track, config Config) string {

	lastFmTrack := getTrackInfo(track, config)

	log.Println(lastFmTrack.Album.Images)

	return "some album art"
}

func getTrackInfo(track Track, config Config) LastFmTrack {

	var lastFmTrack LastFmTrack

	httpClient := &http.Client{}

	req, err := http.NewRequest("GET", "https://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key="+config.LastFmApiKey+"&artist="+url.QueryEscape(track.Artist)+"&track="+url.QueryEscape(track.Title)+"&format=json", nil)

	log.Println(req.URL)

	if err != nil {
		log.Printf("Error in http request: %d\n", err)
		return lastFmTrack
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Error in doing http request: %d\n", err)
		return lastFmTrack
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode > 299 {
		log.Printf("Response failed with status code %d\n", resp.StatusCode)
		return lastFmTrack
	}

	isValidJson := json.Valid(body)

	if !isValidJson {
		log.Println("didn't get valid json in body!")
		return lastFmTrack
	}

	log.Println(string(body))

	var trackResponse LastFmTrackResponse

	err = json.Unmarshal(body, &trackResponse)
	if err != nil {
		log.Printf("Failed unmarshaling json with error %d\n", err)
		return lastFmTrack
	}

	lastFmTrack = trackResponse.Track

	return lastFmTrack

}
