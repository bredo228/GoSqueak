package lastfm

import (
	"encoding/json"
	"fmt"
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
	Url   string      `json:"url"`
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

func GetTrackArtwork(track LastFmTrack) string {

	if len(track.Album.Images) < 1 {
		return ""
	}

	image := track.Album.Images[len(track.Album.Images)-1]

	return image.Text
}

func GetTrackInfo(track Track, config Config) LastFmTrack {

	var queryString string
	var lastFmTrack LastFmTrack

	httpClient := &http.Client{}

	if track.Album != "" {
		queryString = fmt.Sprintf("track=%s&album=%s&artist=%s", url.QueryEscape(track.Title), url.QueryEscape(track.Album), url.QueryEscape(track.Artist))
	} else {
		queryString = fmt.Sprintf("track=%s&artist=%s", url.QueryEscape(track.Title), url.QueryEscape(track.Artist))
	}

	requestUrl := fmt.Sprintf("https://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=%s&%s&format=json", config.LastFmApiKey, queryString)

	req, err := http.NewRequest("GET", requestUrl, nil)

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
	if err != nil {
		log.Printf("Error reading body with error %s\n", err.Error())

	}
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

	var trackResponse LastFmTrackResponse

	err = json.Unmarshal(body, &trackResponse)
	if err != nil {
		log.Printf("Failed unmarshaling json with error %d\n", err)
		return lastFmTrack
	}

	lastFmTrack = trackResponse.Track

	return lastFmTrack

}
