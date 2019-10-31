package main

import (
	"log"
	//	"time"

	"shadow/api"
	"shadow/docker"
	"shadow/env"
	"shadow/mqtt"
	"shadow/watcher"

	"github.com/subosito/gotenv"
)

var err error

func main() {
	gotenv.Load("/env/shadow.env", "/env/production.env")
	env.Init()

	docker.Default, err = docker.Init()
	if err != nil {
		log.Fatal("Error starting docker client: ", err)
	}

	mqtt.Default = mqtt.Init()
	watcher.Default = watcher.Init()

	api.Serve()
}
