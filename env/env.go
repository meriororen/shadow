package env

import "os"

type Env struct {
	Topicprefix string
}

var Default Env

func Init() {
	Default.Topicprefix = "/sensornetwork/" + os.Getenv("VENDOR_ID") + "/" + os.Getenv("TERMINAL_ID")
}
