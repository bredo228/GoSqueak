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
	Title  string
	Album  string
	Artist string
}

func sendOscMessage(message string, path string, client *osc.Client) {
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

	log.Println(artist.Value())

	track.Title = title.String()
	track.Album = album.String()
	track.Artist = artist.String()

	log.Printf("Found track %s by %s in %s\n", track.Title, track.Artist, track.Album)

	return track
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

	log.Println(t)

}
