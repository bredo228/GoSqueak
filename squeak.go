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

func sendOscMessage(message string, path string, client *osc.Client) {
	msg := osc.NewMessage(path)
	msg.Append(message)
	client.Send(msg)
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

	// var status string

	obj := conn.Object(mediaPlayer, "/org/mpris/MediaPlayer2")

	log.Println(obj.Destination())

	log.Println(obj.Path())

	test, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")

	if err != nil {
		log.Fatalf("Metadata: %d, %d", err, test)
	}

	log.Println(test)

}
