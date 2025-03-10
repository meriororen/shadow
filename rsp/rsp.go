package rsp

import (
	"github.com/docker/docker/api/types"
)

type Response struct {
	Type    string      `json:"type,omitempty"`
	Status  string      `json:"status"`
	Error   string      `json:"error,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type ShadowVersion struct {
	Version string `json:"version"`
}

type RspEnvFile struct {
	Status string            `json:"status"`
	Env    map[string]string `json:"env,omitempty"`
}

type RspComposeFile struct {
	Status      string      `json:"status"`
	ComposeFile interface{} `json:"composefile,omitempty"`
}

type PullProgress struct {
	Percentage int `json:"percentage"`
}

type RspShell struct {
	Output string `json:"output"`
}

type RspPull struct {
	Status   string       `json:"status"`
	Progress PullProgress `json:"progress"`
}

type RspLogin struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"errormsg,omitempty"`
}

type ImageItem struct {
	Id      string   `json:"id"`
	Names   []string `json:"names"`
	Created int64    `json:"created"`
	Size    int64    `json:"size"`
}

type RspImageList struct {
	Status string      `json:"status"`
	Images []ImageItem `json:"images,omitempty"`
}

type ContainerItem struct {
	Id      string       `json:"id"`
	ImageId string       `json:"imageid"`
	Command string       `json:"command"`
	Created int64        `json:"created"`
	Status  string       `json:"status"`
	Ports   []types.Port `json:"ports"`
	Names   []string     `json:"names"`
}

type RspContainerList struct {
	Status     string `json:"status"`
	Containers []ContainerItem
}

type RspContainerStart struct {
	Status string `json:"status"`
}

type RspContainerStop struct {
	Status string `json:"status"`
}
