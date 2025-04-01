package main

import (
	"log"

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

	var config Config

	config.Port = 9025

	oscClient := osc.NewClient("127.0.0.1", int(config.Port))

	sendOscMessage("testing", "/squeaknp/test", oscClient)

}
