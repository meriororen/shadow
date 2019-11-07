package watcher

import (
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
)

func CheckCpuTemp() (float32, error) {
	res, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		log.Println("Cannot read from cpu thermal sysfs")
	}

	if resint, err := strconv.Atoi(string(res[:len(res)-1])); err != nil {
		return float32(-1.0), err
	} else {
		return float32(resint) / 1000, nil
	}
}

func CheckUpTime() (string, error) {
	cmd := exec.Command("/usr/bin/uptime", "-p")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	} else {
		return string(out[:len(out)-1]), nil
	}
}
