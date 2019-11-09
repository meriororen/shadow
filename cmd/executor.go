package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"shadow/docker"
	"shadow/rsp"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/joho/godotenv"
)

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
	rs := rsp.Response{}
	rs.Type = cmd.Type

	switch cmd.Type {
	case "envfile":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}

		envfilecmd := CmdEnvFile{}
		if err := json.Unmarshal(cmd.Payload, &envfilecmd); err != nil {
			log.Println("Unprocessable pull payload", err)
		}

		log.Println(envfilecmd)

		if envfilecmd.Path == "" {
			return nil, fmt.Errorf("Path is null")
		}

		var currentEnv map[string]string
		rsp := rsp.RspEnvFile{Status: "Error"}
		if envfilecmd.SetGet == "get" {
			currentEnv, err = godotenv.Read(envfilecmd.Path)
			if err != nil {
				return nil, fmt.Errorf("env path read: ", err)
			}
			rsp.Status = "OK"
			rsp.Env = currentEnv
		} else if envfilecmd.SetGet == "set" {
			err = godotenv.Write(envfilecmd.Env, envfilecmd.Path)
			if err != nil {
				return nil, fmt.Errorf("env path write: ", err)
			}
			rsp.Status = "OK"
		}

		rs.Payload = rsp
		return rs, nil
	case "composefile":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}
	case "shell":
		if cmd.Payload == nil {
			return "", fmt.Errorf("Error executing command, payload must not be empty!")
		}

		shellcmd := CmdShell{}
		if err := json.Unmarshal(cmd.Payload, &shellcmd); err != nil {
			log.Println("Unprocessable pull payload", err)
		}

		cmds := strings.Split(shellcmd.Cmd, " ")

		out, err := exec.Command(cmds[0], cmds[1:]...).CombinedOutput()
		if err != nil {
			rs.Error = err.Error()
		}
		fmt.Println(string(out))
		rs.Payload = rsp.RspShell{Output: string(out)}
		return rs, nil
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

		if res, err = docker.Default.RegistryLogin(logincmd.URL, logincmd.Username,
			logincmd.Password, logincmd.Token); err != nil {
			return "Failed to log in", err
		}
	case "listimages":
		log.Println("listimage")
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
			ExposedPorts: nat.PortSet(exposedPorts),
		}

		hcfg := container.HostConfig{
			Privileged:   runcmd.HostPrivileged,
			PortBindings: nat.PortMap(bindings),
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
