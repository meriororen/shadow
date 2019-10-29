package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"shadow/cmd"
	"shadow/env"
	"shadow/mqtt"
	"shadow/rsp"
	"shadow/status"
	"shadow/watcher"

	pmg "github.com/eclipse/paho.mqtt.golang"
)

type MqttHandler func(client pmg.Client, msg pmg.Message)

func commandExecutor(command cmd.Command) {
	rsc := make(chan interface{})
	erc := make(chan error)
	prg := make(chan []byte) // progress channel

	go func() {
		if command.Type == "pull" {
			command.ProgressChan = prg

			// dispatch another goroutine to monitor progress
			go func() {
				for {
					select {
					case progress, ok := <-prg:
						if !ok {
							return
						}
						if token := mqtt.Default.Publish(env.Default.Topicprefix+"/progress", 0, false, progress); token.Wait() && token.Error() != nil {
							log.Println("MQTTMON: Cannot publish progress")
						}
					}
				}
			}()
		}

		if res, err := cmd.Exec(command); err != nil {
			erc <- err
		} else {
			rsc <- res
		}
		close(prg)
		close(rsc)
		close(erc)
	}()

	var (
		err     error
		theresp []byte
	)
	select {
	case result := <-rsc:
		if result == nil {
			break
		}
		resp := result.(rsp.Response)
		log.Println("Result of cmd ", resp.Type, " -> ", resp.Payload)

		if theresp, err = json.Marshal(resp.Payload); err != nil {
			log.Println("MQTTMON: Cannot marshal response struct")
		}
	case er := <-erc:
		log.Println("Error running cmd ", command, " -> ", er)

		errsp := rsp.Response{
			Type:  command.Type,
			Error: er.Error(),
		}

		if theresp, err = json.Marshal(errsp); err != nil {
			log.Println("MQTTMON: Cannot marshal error")
		}
	}

	if token := mqtt.Default.Publish(env.Default.Topicprefix+"/response", 0, false, theresp); token.Wait() && token.Error() != nil {
		log.Println("MQTTMON: Cannot publish response")
	}
}

func MqttMonitor() {
	// monitor Status
	go func() {
		for {
			select {
			case stat, ok := <-watcher.Default.StatusChan:
				if !ok {
					return
				}
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
			}
		}
	}()

	// monitor command execution
	go func() {
		for {
			select {
			case command, ok := <-watcher.Default.CommandChan:
				if !ok {
					return
				}
				commandExecutor(command)
			}
		}
	}()
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

	MqttMonitor()

	watcher.Default.WatchAll()
}
