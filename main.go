package main

import (
	"github.com/FG-GIS/boot-go-gator/internal/config"
	"log"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading the config file: %v", err)
	}
	err = conf.SetUser("Void")
	if err != nil {
		log.Fatalf("Error writing the config file: %v", err)
	}
	conf, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading the config file: %v", err)
	}
	log.Print(conf)
}
