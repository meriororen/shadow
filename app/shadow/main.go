package main

import (
	"log"
	"time"

	"shadow/api"
	"shadow/env"
	"shadow/watcher"

	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	env.Init()

	wt := env.Default.Wtch

	backend := "registry.gitlab.com/sangkuriang-dev/transmissor-be/backend:devel"
	golang := "golang:1.13-alpine"

	wt.AddImageToWatch(watcher.WatchConfig{
		ImageName: golang,
		AutoPull:  false,
		HBPeriod:  3 * time.Second,
	})

	wt.AddImageToWatch(watcher.WatchConfig{
		ImageName: backend,
		AutoPull:  true,
		HBPeriod:  2 * time.Second,
	})

	for _, w := range wt.WatchList {
		log.Println((*w).ImageNames, " => ", (*w).ContainerIDs)
	}

	wt.RemoveImageFromWatchList("golang:1.13-alpine")

	for _, w := range wt.WatchList {
		log.Println((*w).ImageNames, " => ", (*w).ContainerIDs)
	}

	api := api.NewAPI(wt)

	api.Serve()
}
