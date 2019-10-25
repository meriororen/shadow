package main

import (
	"log"

	"shadow/watcher"

	pmg "github.com/eclipse/paho.mqtt.golang"
	"github.com/subosito/gotenv"
)

var (
	mq *pmg.Client
)

func main() {
	gotenv.Load()

	wt := watcher.NewWatcher()
	wt.AddImageToWatch("golang:1.13-alpine")
	wt.AddImageToWatch("registry.gitlab.com/sangkuriang-dev/transmissor-be/backend:devel")

	//	wt.WatchAll()

	/*
		mq = &mqtt.NewClient(
			os.Getenv("MQTT_BROKER_URL"),
			os.Getenv("MQTT_BROKER_USER"),
			os.Getenv("MQTT_BROKER_PASS"),
			"shadow-"+fmt.Sprint(time.Now().Local().Format(time.RFC3339)),
			"",
		)
	*/

	for _, w := range wt.WatchList {
		log.Println((*w).ImageNames, " => ", (*w).ContainerIDs)
	}
}
