package main

import (
	"log"
	"time"

	"shadow/api"
	"shadow/docker"
	"shadow/env"
	"shadow/mqtt"
	"shadow/watcher"

	"github.com/subosito/gotenv"
)

var err error

func main() {
	gotenv.Load()
	env.Init()

	docker.Default, err = docker.NewDocker()
	if err != nil {
		log.Fatal("Error starting docker client: ", err)
	}

	mqtt.Default = mqtt.Init()
	watcher.Default = watcher.NewWatcher()

	backend := "registry.gitlab.com/sangkuriang-dev/transmissor-be/backend:devel"
	golang := "golang:1.13-alpine"

	watcher.Default.AddImageToWatch(watcher.WatchConfig{
		ImageName: golang,
		AutoPull:  false,
		HBPeriod:  3 * time.Second,
	})

	watcher.Default.AddImageToWatch(watcher.WatchConfig{
		ImageName: backend,
		AutoPull:  true,
		HBPeriod:  2 * time.Second,
	})

	for _, w := range watcher.Default.WatchList {
		log.Println((*w).ImageNames, " => ", (*w).ContainerIDs)
	}

	watcher.Default.RemoveImageFromWatchList("golang:1.13-alpine")

	for _, w := range watcher.Default.WatchList {
		log.Println((*w).ImageNames, " => ", (*w).ContainerIDs)
	}

	api.Serve()
}
