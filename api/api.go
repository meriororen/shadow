package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"shadow/cmd"
	"shadow/env"
	"shadow/status"
	"shadow/watcher"

	pmg "github.com/eclipse/paho.mqtt.golang"
)

type API struct {
	wt *watcher.Watcher
}

type MqttHandler func(client pmg.Client, msg pmg.Message)

func NewAPI(wt *watcher.Watcher) API {
	return API{
		wt: wt,
	}
}

func CommandHandler(client pmg.Client, msg pmg.Message) {
	wt := env.Default.Wtch
	topic := strings.TrimPrefix(msg.Topic(), env.Default.Topicprefix)
	log.Println(topic)

	var command string
	if i := strings.Index(topic, "/command"); i == 0 {
		command = strings.TrimPrefix(topic, "/command/")
	} else {
		log.Fatal("CommandHandler: Invalid topic, ", i)
	}

	for _, w := range wt.WatchConfigList {
		if w.ImageName != "System" {
			log.Println("Sending command to goroutine", w.ImageName)
			w.CommandChan <- cmd.Command{Type: command}
			log.Println("Sent command to goroutine", w.ImageName)
		}
	}
}

func ResponseHandler(client pmg.Client, msg pmg.Message) {
}

func (api API) MqttMonitor() {
	for {
		select {
		case stat := <-api.wt.StatusChan:
			//log.Println("Sending Status")

			var thestatus string
			switch v := stat.Payload.(type) {
			case status.System:
				thestatus = fmt.Sprintf("time: %s, total: %d, free: %d, cached: %d, used: %d",
					stat.LocalTime.Format(time.RFC3339),
					v.Memory.Total,
					v.Memory.Free,
					v.Memory.Cached,
					v.Memory.Used)
			}
			//log.Println("publishing to ", env.Default.Topicprefix+"/status")
			if token := env.Default.Mqtt.Publish(env.Default.Topicprefix+"/status", 0, false, thestatus); token.Wait() && token.Error() != nil {
				log.Println("Cannot publish status")
			}
		}
	}
}

func MqttSubscribe() {
	topics := map[string]pmg.MessageHandler{
		"/command/+": CommandHandler,
		"/response":  ResponseHandler,
	}

	for topic, topicHandler := range topics {
		env.Default.Mqtt.Subscribe(env.Default.Topicprefix+topic, 0, topicHandler)
	}
}

func (api API) Serve() {
	MqttSubscribe()

	go func() {
		api.MqttMonitor()
	}()

	api.wt.WatchAll()
}
