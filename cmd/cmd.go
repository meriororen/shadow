package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"shadow/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	//	"github.com/docker/docker/api/types"
	//"github.com/docker/docker/client"
)

type Command struct {
	Type         string
	ProgressChan chan []byte
	Payload      []byte
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

func parseMountPoints(volumes []string) (res []mount.Mount) {
	for _, vol := range volumes {
		mnt := strings.Split(vol, ":")

		var readonly bool
		if len(mnt) > 2 && mnt[2] == "ro" {
			readonly = true
		}
		res = append(res, mount.Mount{
			Type:     mount.TypeBind,
			Source:   mnt[0],
			Target:   mnt[1],
			ReadOnly: readonly,
		})
	}

	return res
}

func Exec(cmd Command, args ...interface{}) (res interface{}, err error) {
	switch cmd.Type {
	case "pull":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}

		pullcmd := CmdPull{}
		if err := json.Unmarshal(cmd.Payload, &pullcmd); err != nil {
			log.Println("Unprocessable pull payload", err)
		}

		log.Println("pulling image: ", pullcmd.ImageName)
		if res, err = docker.Default.ImagePull(pullcmd.ImageName, cmd.ProgressChan); err != nil {
			return "Error Pulling Image!", err
		}
	case "login":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}

		logincmd := CmdLogin{}
		if err := json.Unmarshal(cmd.Payload, &logincmd); err != nil {
			log.Println("Unprocessable login payload", err)
		}

		if res, err = docker.Default.RegistryLogin(logincmd.URL, logincmd.Username, logincmd.Password, logincmd.Token); err != nil {
			return "Failed to log in", err
		}
	case "listimages":
		var pl string
		if cmd.Payload == nil {
			pl = ""
		} else {
			imagelistcmd := CmdImageList{}
			if err := json.Unmarshal(cmd.Payload, &imagelistcmd); err != nil {
				log.Println("Unprocessable Image List payload", err)
			}
			pl = imagelistcmd.ImageName
		}

		if res, err = docker.Default.ImageList(pl); err != nil {
			return "Error listing out images!", err
		}
	case "listcontainers":
		var pl string
		if cmd.Payload == nil {
			pl = ""
		} else {
			containerlistcmd := CmdContainerList{}
			if err := json.Unmarshal(cmd.Payload, &containerlistcmd); err != nil {
				log.Println("Unprocessable Container List payload", err)
			}
			pl = containerlistcmd.ImageID
		}

		if res, err = docker.Default.ContainerList(pl); err != nil {
			return "Error listing out images!", err
		}
	case "run":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}

		runcmd := CmdContainerRun{}
		if err := json.Unmarshal(cmd.Payload, &runcmd); err != nil {
			log.Println("Unprocessable pull payload", err)
		}

		exposedPorts, bindings, err := nat.ParsePortSpecs(runcmd.HostPortBindings)
		log.Println("bindings: ", bindings)
		log.Println("exposed: ", exposedPorts)
		if err != nil {
			return "Error Port bindings format", err
		}

		mountBindings := parseMountPoints(runcmd.Volumes)

		ccfg := container.Config{
			Image:        runcmd.ImageName,
			Env:          runcmd.Env,
			Cmd:          runcmd.Cmd,
			ExposedPorts: exposedPorts,
		}

		hcfg := container.HostConfig{
			Privileged:   runcmd.HostPrivileged,
			PortBindings: bindings,
			Mounts:       mountBindings,
		}

		ncfg := network.NetworkingConfig{}

		if res, err = docker.Default.ContainerRun(&ccfg, &hcfg, &ncfg, runcmd.Name); err != nil {
			return "Error Running!", err
		}

	case "stop":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}

		stopcmd := CmdContainerStop{}
		if err := json.Unmarshal(cmd.Payload, &stopcmd); err != nil {
			log.Println("Unprocessable pull payload", err)
		}

		if res, err = docker.Default.ContainerStop(stopcmd.Id); err != nil {
			return "Error Stopping!", err
		}
	}

	return res, nil
}
