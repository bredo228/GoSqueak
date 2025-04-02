package main

import (
	"encoding/json"
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

func getCurrentTrack(obj dbus.BusObject) Track {

	var track Track

	data, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")

	if err != nil {
		log.Panicf("Failed getting property!")
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
		log.Fatalf("Failed getting progress")
	}

	if pos.Value() != nil {
		pos_i64 := pos.Value().(int64) / 1000000
		track.Position = int(pos_i64)
	}

	return track
}

func sendTrack(track Track, client *osc.Client) {

	sendOscMessage(track.Title, "/squeaknp/track_title", client)
	sendOscMessage(track.Album, "/squeaknp/track_album", client)
	sendOscMessage(track.Artist, "/squeaknp/track_artist", client)
	sendOscMessage(track.Artwork, "/squeaknp/lastfm_album_art", client)

	// TODO: these are floats not ints
	sendOscMessage(float32(track.Duration), "/squeaknp/timeline_end_time", client)
	sendOscMessage(float32(track.Position), "/squeaknp/timeline_position", client)

}

func main() {

	log.Println("Starting GoSqueak")

	log.Printf("Running on platform %s\n", runtime.GOOS)

	// init config

	config_file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Error when opening config file: %d\n", err)
	}

	var config Config

	err = json.Unmarshal(config_file, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %d\n", err)
	}

	// init osc client
	oscClient := osc.NewClient("127.0.0.1", int(config.Port))

	conn, err := dbus.ConnectSessionBus()

	if err != nil {
		log.Fatalf("Failed to connect to session bus: %d\n", err)
	}
	defer conn.Close()

	// find media player
	var mediaPlayer string

	for {
		var names []string

		err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)

		if err != nil {
			log.Fatalf("Failed to get list of owned names: %d\n", err)
		}

		for _, name := range names {
			if strings.HasPrefix(name, "org.mpris.MediaPlayer2") {
				log.Printf("Found media player with name %s", name)
				mediaPlayer = name
				break
			}
		}

		if mediaPlayer != "" { // we've found a media player, we can continue now
			break
		}

		log.Println("Didn't find a media player, retrying in 10 seconds.")

		time.Sleep(10 * time.Second)

	}

	obj := conn.Object(mediaPlayer, "/org/mpris/MediaPlayer2")

	var previousTrack Track
	var currentTrack Track

	// Main update loop
	for {

		previousTrack = currentTrack

		currentTrack = getCurrentTrack(obj)

		if previousTrack.Title != currentTrack.Title { // probably a better way to see if the track has changed
			log.Println("Track title has changed!")
			log.Printf("Found track %s by %s in %s\n", currentTrack.Title, currentTrack.Artist, currentTrack.Album)

			// send the track before we do any updates
			sendTrack(currentTrack, oscClient)

			currentTrack.Artwork = lastfm.GetTrackArtwork(currentTrack, config)

			log.Println("Updated artwork to " + currentTrack.Artwork)

		} else {
			if previousTrack.Artwork != "" {
				currentTrack.Artwork = previousTrack.Artwork
			}
		}

		sendTrack(currentTrack, oscClient)

		time.Sleep(time.Millisecond * time.Duration(config.UpdateRate))

	}

}
