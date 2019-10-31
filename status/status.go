package status

import "time"

type Status struct {
	LocalTime time.Time
	Payload   interface{}
}

type Memory struct {
	Total  uint64 `json:"total"`
	Free   uint64 `json:"free"`
	Used   uint64 `json:"used"`
	Cached uint64 `json:"cached"`
}

type System struct {
	CPU        string `json:"cpu,omitempty"`
	Uname      string `json:"uname,omitempty"`
	Memory     Memory `json:"memory,omitempty"`
	MacAddress string `json:"mac_address"`
}
