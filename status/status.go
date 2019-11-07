package status

import "time"

type Status struct {
	LocalTime time.Time   `json:"time"`
	Type      string      `json:"type"`
	Payload   interface{} `json:"status"`
}

var NilStatus Status

type SystemMemory struct {
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Used   uint64 `json:"used,omitempty"`
	Cached uint64 `json:"cached,omitempty"`
}

type System struct {
	CPUTemp       string              `json:"cputemp,omitempty"`
	Memory        SystemMemory        `json:"memory,omitempty"`
	WanIp         string              `json:"wanip,omitempty"`
	ServiceUptime []map[string]string `json:"uptime,omitempty"`
	SystemUptime  string              `json:"sysuptime,omitempty"`
}
