package env

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Env struct {
	Topicprefix string
}

var Default Env

func Init() {
	var terminalid string

	wlan := os.Getenv("WIFI_INTERFACE_NAME")
	mac, err := ioutil.ReadFile("/sys/class/net/" + wlan + "/address")

	if err != nil {
		log.Println("Cannot get mac address")
		mac = []byte(fmt.Sprintf("Unknown_Mac_%d", time.Now().Unix()))
	}

	if terminalid = os.Getenv("TERMINAL_ID"); terminalid == "UNPROVISIONED_TERMINAL" {
		terminalid = strings.TrimSpace(string(mac))
	}

	Default.Topicprefix = "sensornetwork/" + os.Getenv("VENDOR_ID") + "/" + terminalid
}
