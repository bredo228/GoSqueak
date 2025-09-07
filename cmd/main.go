package main

import (
	"log"

	"github.com/bredo228/GoSqueak/cmd/config"
)

var (
	conf config.Config
)

func main() {
	log.Println("Starting GoSqueak v2")

	config.LoadConfig("./config.json", &conf)

}
