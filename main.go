package main

import (
	"contactUs/configuration"
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Config file path")
	config, err := configuration.ParseFile(*configPath)
	if err != nil {
		panic(err)
	}
	if len(config.Forms) < 1 {
		log.Panicln("Could not find any form configured", config.Forms)
	}
	log.Println("Found ", len(config.Forms), " forms")

	setupHttp()
}
