package watcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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
	out, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	} else {
		uptime := strings.Split(string(out), " ")
		uptimef, _ := strconv.ParseFloat(uptime[0], 32)
		uptimei := int(uptimef)
		hours := uptimei / (60 * 60)
		minutes := (uptimei % (60 * 60)) / 60
		seconds := (uptimei % 60)
		up := fmt.Sprint(hours, ".", minutes, ".", seconds)
		return up, nil
	}
}
