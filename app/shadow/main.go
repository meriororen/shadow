package main

import (
	"log"

	"shadow/api"
	"shadow/docker"
	"shadow/env"
	"shadow/mqtt"
	"shadow/watcher"

	"github.com/joho/godotenv"
)

var err error

func main() {
	_ = godotenv.Load("/env/production.env", "./.env")
	env.Init()

	log.Println("Started shadow version: ", env.Version)

	docker.Default, err = docker.Init()
	if err != nil {
		log.Fatal("Error starting docker client: ", err)
	}

	mqtt.Default = mqtt.Init()
	watcher.Default = watcher.Init()

	api.Serve()
}
