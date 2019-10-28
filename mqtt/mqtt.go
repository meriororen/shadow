package mqtt

import (
	"fmt"
	"log"
	"os"
	"time"

	pmg "github.com/eclipse/paho.mqtt.golang"
)

var Default pmg.Client

func Init() pmg.Client {
	return NewClient(
		os.Getenv("MQTT_BROKER_URL"),
		os.Getenv("MQTT_BROKER_USER"),
		os.Getenv("MQTT_BROKER_PASS"),
		"shadow-"+fmt.Sprint(time.Now().Local().Format(time.RFC3339)),
		"",
	)
}

func connLostHandler(c pmg.Client, err error) {
	fmt.Printf("Connection lost, reason: %v\n", err)
}

//NewClient create new client instance for mqtt
func NewClient(broker string, user string, pass string, clientid string, cert string) pmg.Client {
	var c pmg.Client

	if broker != "" {
		opts := pmg.NewClientOptions().AddBroker(broker).SetClientID(clientid)
		opts.SetKeepAlive(60 * time.Minute)
		opts.SetPingTimeout(2 * time.Second)
		opts.SetConnectionLostHandler(connLostHandler)

		c = pmg.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		log.Println("Started with Mqtt client_id ", clientid)
		return c
	}
	log.Println("Broker is None")

	return nil
}
