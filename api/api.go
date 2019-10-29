package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"shadow/cmd"
	"shadow/docker"
	"shadow/env"
	"shadow/mqtt"
	//	"shadow/rsp"
	"shadow/status"
	"shadow/watcher"

	pmg "github.com/eclipse/paho.mqtt.golang"
)

type MqttHandler func(client pmg.Client, msg pmg.Message)

func commandExecutor(command cmd.Command) {
	rsc := make(chan interface{})
	erc := make(chan error)

	go func() {
		if res, err := cmd.Exec(command); err != nil {
			erc <- err
		} else {
			rsc <- res
		}
		close(rsc)
		close(erc)
	}()

	select {
	case result := <-rsc:
		log.Println("Result of cmd ", command, " -> ", result)
	case err := <-erc:
		log.Println("Error running cmd ", command, " -> ", err)
	}
}

func MqttMonitor() {
	var err error
	for {
		select {
		case resp := <-docker.Default.ResponseChan:
			log.Println("Got response for ", resp.Type, "->", resp.Payload)
			var theresp []byte
			if theresp, err = json.Marshal(resp.Payload); err != nil {
				log.Println("MQTTMON: Cannot marshal response struct")
			}
			if token := mqtt.Default.Publish(env.Default.Topicprefix+"/response", 0, false, theresp); token.Wait() && token.Error() != nil {
				log.Println("MQTTMON: Cannot publish response")
			}
		case stat := <-watcher.Default.StatusChan:
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
			if token := mqtt.Default.Publish(env.Default.Topicprefix+"/status", 0, false, thestatus); token.Wait() && token.Error() != nil {
				log.Println("MQTTMON: Cannot publish status")
			}
		case command := <-watcher.Default.CommandChan:
			commandExecutor(command)
		}
	}
}

func CommandHandler(client pmg.Client, msg pmg.Message) {
	//wt := watcher.Default
	topic := strings.TrimPrefix(msg.Topic(), env.Default.Topicprefix)
	log.Println(topic)

	var command string
	if i := strings.Index(topic, "/command"); i == 0 {
		command = strings.TrimPrefix(topic, "/command/")
	} else {
		log.Fatal("CommandHandler: Invalid topic, ", i)
	}

	watcher.Default.CommandChan <- cmd.Command{Type: command, Payload: msg.Payload()}

	/*
		for _, w := range wt.WatchConfigList {
			if w.ImageName != "System" {
				w.CommandChan <- cmd.Command{Type: command, Payload: msg.Payload()}
			}
		}
	*/
}

func MqttSubscribe() {
	topics := map[string]pmg.MessageHandler{
		"/command/+": CommandHandler,
	}

	for topic, topicHandler := range topics {
		mqtt.Default.Subscribe(env.Default.Topicprefix+topic, 0, topicHandler)
	}
}

func Serve() {
	MqttSubscribe()

	go func() {
		MqttMonitor()
	}()

	watcher.Default.WatchAll()
}
