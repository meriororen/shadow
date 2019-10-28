package env

import (
	"fmt"
	"os"
	"shadow/watcher"
	"time"

	"shadow/mqtt"

	pmg "github.com/eclipse/paho.mqtt.golang"
)

type Env struct {
	Mqtt        pmg.Client
	Wtch        *watcher.Watcher
	Topicprefix string
}

var Default Env

func Init() {
	Default.Mqtt = mqtt.NewClient(
		os.Getenv("MQTT_BROKER_URL"),
		os.Getenv("MQTT_BROKER_USER"),
		os.Getenv("MQTT_BROKER_PASS"),
		"shadow-"+fmt.Sprint(time.Now().Local().Format(time.RFC3339)),
		"",
	)

	Default.Topicprefix = "/sensornetwork/" + os.Getenv("VENDOR_ID") + "/" + os.Getenv("TERMINAL_ID")

	Default.Wtch = watcher.NewWatcher()
}
