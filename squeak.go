package main

import (
	"log"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/hypebeast/go-osc/osc"
)

type Config struct {
	Port int
}

type Track struct {
	Title    string
	Album    string
	Artist   string
	Duration int
	Position int
}

func sendOscMessage(message string, path string, client *osc.Client) {
	msg := osc.NewMessage(path)
	msg.Append(message)
	client.Send(msg)
}

func sendOscMessageInt(message int, path string, client *osc.Client) {
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

	track.Title = title.Value().(string)
	track.Album = album.Value().(string)

	artists := artist.Value().([]string)
	track.Artist = artists[0]

	length := metadata["mpris:length"]
	length_i64 := length.Value().(int64) / 1000000
	duration := int(length_i64)

	track.Duration = duration

	// get progress
	pos, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")

	if err != nil {
		log.Fatalf("Failed getting progress")
	}

	pos_i64 := pos.Value().(int64) / 1000000
	track.Position = int(pos_i64)

	log.Printf("Found track %s by %s in %s with duration %d\n", track.Title, track.Artist, track.Album, track.Duration)

	return track
}

func sendTrack(track Track, client *osc.Client) {

	sendOscMessage(track.Title, "/squeaknp/track_title", client)
	sendOscMessage(track.Album, "/squeaknp/track_album", client)
	sendOscMessage(track.Artist, "/squeaknp/track_artist", client)

	// TODO: these are floats not ints
	sendOscMessageInt(track.Duration, "/squeaknp/timeline_end_time", client)
	sendOscMessageInt(track.Position, "/squeaknp/timeline_position", client)

}

func main() {

	log.Println("Starting SqueakLinux")

	// init config
	var config Config
	config.Port = 9025

	// init osc client
	oscClient := osc.NewClient("127.0.0.1", int(config.Port))

	sendOscMessage("testing", "/squeaknp/test", oscClient)

	conn, err := dbus.ConnectSessionBus()

	if err != nil {
		log.Fatalf("Failed to connect to session bus: %d\n", err)
	}
	defer conn.Close()

	// find media player
	var names []string
	var mediaPlayer string

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

	obj := conn.Object(mediaPlayer, "/org/mpris/MediaPlayer2")

	t := getCurrentTrack(obj)

	sendTrack(t, oscClient)

}
