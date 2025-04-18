package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bredo228/GoSqueak/lastfm"
	"github.com/godbus/dbus/v5"
	"github.com/hypebeast/go-osc/osc"

	. "github.com/bredo228/GoSqueak/datatypes"
)

func sendOscMessage(message any, path string, client *osc.Client) {
	msg := osc.NewMessage(path)
	msg.Append(message)
	client.Send(msg)
}

func getCurrentTrack(obj dbus.BusObject) (Track, error) {

	var track Track

	data, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")

	if err != nil {
		log.Println("Failed getting property!")
		return track, err
	}

	metadata := data.Value().(map[string]dbus.Variant)

	title := metadata["xesam:title"]
	album := metadata["xesam:album"]
	artist := metadata["xesam:artist"]

	if title.Value() != nil {
		track.Title = title.Value().(string)
	}

	if album.Value() != nil {
		track.Album = album.Value().(string)
	}

	if artist.Value() != nil {
		artists := artist.Value().([]string)
		track.Artist = artists[0]
	}

	length := metadata["mpris:length"]

	if length.Value() != nil {
		length_i64 := length.Value().(int64) / 1000000
		duration := int(length_i64)

		track.Duration = duration
	}

	// get progress
	pos, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")

	if err != nil {
		log.Printf("Failed getting progress")
		// This is the current position in the song, I'm just going to have the position be 0 if this fails.
		track.Position = 0
		return track, nil
	}

	if pos.Value() != nil {
		pos_i64 := pos.Value().(int64) / 1000000
		track.Position = int(pos_i64)
	}

	return track, nil
}

// Calling this will lock the program up until it finds a media player.
func findMusicPlayer(conn *dbus.Conn) string {

	var player string

	for {

		var names []string

		err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)

		if err != nil {
			log.Fatalf("Failed to get list of owned names: %s\n", err.Error())
		}

		for _, name := range names {
			if strings.HasPrefix(name, "org.mpris.MediaPlayer2") {
				log.Printf("Found media player with name %s", name)
				player = name
				break
			}
		}

		if player != "" { // If we've found a player, return it
			return player
		}

		log.Println("Didn't find a media player, retrying in 10 seconds.")

		time.Sleep(10 * time.Second)

	}

}

func sendTrack(track Track, client *osc.Client) {

	sendOscMessage(track.Title, "/squeaknp/track_title", client)
	sendOscMessage(track.Album, "/squeaknp/track_album", client)
	sendOscMessage(track.Artist, "/squeaknp/track_artist", client)
	sendOscMessage(track.Artwork, "/squeaknp/lastfm_album_art", client)
	sendOscMessage(track.Url, "/squeaknp/lastfm_url", client)

	// TODO: these are floats not ints
	sendOscMessage(float32(track.Duration), "/squeaknp/timeline_end_time", client)
	sendOscMessage(float32(track.Position), "/squeaknp/timeline_position", client)

}

func main() {

	log.Println("Starting GoSqueak")

	log.Printf("Running on platform %s\n", runtime.GOOS)

	// init config

	var useFallbackConfig bool

	config_file, err := os.ReadFile("./config.json")

	if errors.Is(err, os.ErrNotExist) {
		log.Println("Config file does not exist, falling back.")
		useFallbackConfig = true
	} else if err != nil {
		log.Fatalf("Error when opening config file: %d\n", err)
	}

	var config Config

	if !useFallbackConfig {
		err = json.Unmarshal(config_file, &config)
		if err != nil {
			log.Fatalf("Error unmarshaling config: %d\n", err)
		}
	} else {
		config.Port = 9025
		config.UpdateRate = 500
		config.LastFmEnabled = false
	}

	// init osc client
	oscClient := osc.NewClient("127.0.0.1", int(config.Port))

	conn, err := dbus.ConnectSessionBus()

	if err != nil {
		log.Fatalf("Failed to connect to session bus: %d\n", err)
	}
	defer conn.Close()

	// find media player
	mediaPlayer := findMusicPlayer(conn)

	obj := conn.Object(mediaPlayer, "/org/mpris/MediaPlayer2")

	var previousTrack Track
	var currentTrack Track

	// Main update loop
	for {

		previousTrack = currentTrack

		currentTrack, err = getCurrentTrack(obj)

		if err != nil {

			// If we get a specific error, the music player has probably been closed so try and find it again.
			if err.Error() != "The name is not activatable" {
				log.Fatalf("Error in getCurrentTrack: %s\n", err.Error())
			}

			log.Println("Failed getting current track - music player has probably been closed.")
			mediaPlayer = findMusicPlayer(conn)
			obj = conn.Object(mediaPlayer, "/org/mpris/MediaPlayer2")
		}

		if previousTrack.Title != currentTrack.Title { // probably a better way to see if the track has changed
			log.Println("Track title has changed!")
			log.Printf("Found track %s by %s in %s\n", currentTrack.Title, currentTrack.Artist, currentTrack.Album)

			// send the track before we do any updates
			sendTrack(currentTrack, oscClient)

			if config.LastFmEnabled {
				lastFmTrack := lastfm.GetTrackInfo(currentTrack, config)
				currentTrack.Artwork = lastfm.GetTrackArtwork(lastFmTrack)
				currentTrack.Url = lastFmTrack.Url
			}

			log.Println("Updated artwork to " + currentTrack.Artwork)

		} else {
			if previousTrack.Artwork != "" {
				currentTrack.Artwork = previousTrack.Artwork
			}
			if previousTrack.Url != "" {
				currentTrack.Url = previousTrack.Url
			}
		}

		sendTrack(currentTrack, oscClient)

		time.Sleep(time.Millisecond * time.Duration(config.UpdateRate))

	}

}
