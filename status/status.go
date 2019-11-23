package status

type Status struct {
	LocalTime int64       `json:"time"`
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
	CpuTemp       float32             `json:"cputemp,omitempty"`
	Memory        SystemMemory        `json:"memory,omitempty"`
	WanIp         string              `json:"wanip,omitempty"`
	ServiceUpTime []map[string]string `json:"uptime,omitempty"`
	SystemUpTime  string              `json:"sysuptime,omitempty"`
}
