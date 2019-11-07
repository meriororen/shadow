package status

import "time"

type Status struct {
	LocalTime time.Time
	Payload   interface{}
}

type SystemMemory struct {
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Used   uint64 `json:"used"`
	Cached uint64 `json:"cached"`
}

type System struct {
	CPUTemp       string              `json:"cputemp,omitempty"`
	Memory        SystemMemory        `json:"memory,omitempty"`
	WANIP         string              `json:"wanip,omitempty"`
	ServiceUptime []map[string]string `json:"uptime,omitempty"`
	SystemUptime  string              `json:"sysuptime,omitempty"`
}
