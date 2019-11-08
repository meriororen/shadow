package cmd

import "github.com/docker/docker/api/types/strslice"

type Command struct {
	Type         string
	ProgressChan chan []byte
	Payload      []byte
}

type CmdGetEnvFile struct {
	Path string
}

type CmdSetEnvFile struct {
	Path    string
	EnvFile interface{}
}

type CmdGetComposeYaml struct {
	Path string
}

type CmdSetComposeYaml struct {
	Path        string
	ComposeYaml interface{}
}

type CmdShell struct {
	Cmd string `json:"cmd"`
}

type CmdPull struct {
	ImageName string `json:"image_name"`
}

type CmdLogin struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

type CmdImageList struct {
	ImageName string `json:"image_name,omitempty"`
}

type CmdContainerList struct {
	ImageID string `json:"image_id,omitempty"`
}

type CmdContainerRun struct {
	ImageName        string            `json:"image_name"`
	Name             string            `json:"name,omitempty"`
	Volumes          []string          `json:"volumes,omitempty"`
	Networks         []string          `json:"networks,omitempty"`
	Env              []string          `json:"env,omitempty"`
	WorkingDir       string            `json:"workdir,omitempty"`
	EntryPoint       string            `json:"entrypoint,omitempty"`
	Cmd              strslice.StrSlice `json:"cmd,omitempty"`
	HostPrivileged   bool              `json:"privileged,omitempty"`
	HostPortBindings []string          `json:"ports,omitempty"`
}

type CmdContainerStop struct {
	Id string `json:"id"`
}
