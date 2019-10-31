package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
	//	"os"

	"shadow/rsp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

type Docker struct {
	Client           *client.Client
	savedCredentials types.AuthConfig
}

var Default *Docker

func Init() (*Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println("Cannot initialize client: ", err)
		return nil, err
	}

	return &Docker{
		Client: cli,
	}, nil
}

func attemptLogin(creds types.AuthConfig) (idtoken string, err error) {
	var res registry.AuthenticateOKBody
	if res, err = Default.Client.RegistryLogin(context.Background(), creds); err != nil {
		return "", err
	} else {
		log.Println("Logged in -> ", res)
	}

	return res.IdentityToken, nil
}

func (d *Docker) RegistryLogin(url string, user string, pass string, token string) (rsp.Response, error) {
	creds := types.AuthConfig{
		ServerAddress: url,
		Username:      user,
		Password:      pass,
		IdentityToken: token,
	}

	if idtoken, err := attemptLogin(creds); err != nil {
		return rsp.Response{}, err
	} else {
		d.savedCredentials = creds

		if idtoken != "" {
			d.savedCredentials.IdentityToken = idtoken
			log.Println("Login success, id token: ", idtoken)
		} else {
			log.Println("Id token is empty")
		}
	}

	return rsp.Response{Type: "login"}, nil
}

func (d *Docker) ImagePull(imageName string, prg chan []byte) (rsp.Response, error) {
	auth := types.AuthConfig{
		Username: d.savedCredentials.Username,
		Password: d.savedCredentials.Password,
	}

	var strauth string
	if encauth, err := json.Marshal(auth); err != nil {
		return rsp.Response{}, err
	} else {
		strauth = base64.URLEncoding.EncodeToString(encauth)
	}

	out, err := d.Client.ImagePull(context.Background(), imageName, types.ImagePullOptions{RegistryAuth: strauth})
	if err != nil {
		log.Println("Cannot image pull")
		return rsp.Response{}, err
	}
	defer out.Close()

	// progress update
	reader := bufio.NewReader(out)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err)
			}
		}
		prg <- []byte(line)
	}

	return rsp.Response{Type: "pull", Payload: "Pull Success"}, nil
}

func (d *Docker) ImageList(imageName string) (rsp.Response, error) {
	images, err := d.Client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Println("Cannot list out all images: ", err)
		return rsp.Response{}, err
	}

	log.Println("Listing image for: ", imageName)

	imgs := []rsp.ImageItem{}
	for _, img := range images {
		for _, t := range img.RepoTags {
			imgi := rsp.ImageItem{
				Id:      img.ID[7:19],
				Created: img.Created,
				Names:   img.RepoTags,
				Size:    img.VirtualSize,
			}

			if imageName == "" {
				imgs = append(imgs, imgi)
			} else {
				if imageName == t {
					imgs = append(imgs, imgi)
				}
			}
		}
	}

	return rsp.Response{Type: "listimages", Payload: imgs}, nil
}

func (d *Docker) ContainerList(imageID string) (rsp.Response, error) {
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Println("Cannot get container list : ", err)
		return rsp.Response{}, nil
	}

	log.Println("Listing container for: ", imageID)

	cnts := []rsp.ContainerItem{}
	for _, cnt := range containers {
		cnti := rsp.ContainerItem{
			Id:      cnt.ID[0:11],
			ImageId: cnt.ImageID[7:19],
			Command: cnt.Command,
			Created: cnt.Created,
			Status:  cnt.Status,
			Ports:   cnt.Ports,
			Names:   cnt.Names,
		}

		if imageID == "" {
			cnts = append(cnts, cnti)
		} else {
			if imageID == cnt.ImageID[7:19] {
				cnts = append(cnts, cnti)
			}
		}
	}

	return rsp.Response{Type: "listcontainers", Payload: cnts}, nil
}

func (d *Docker) ContainerRun(ccfg *container.Config, hcfg *container.HostConfig, ncfg *network.NetworkingConfig, name string) (rsp.Response, error) {
	log.Println("starting container for image: ", ccfg)
	resp, err := d.Client.ContainerCreate(context.Background(), ccfg, hcfg, ncfg, name)

	if err != nil {
		return rsp.Response{}, err
	} else {
		log.Println("running created container with ID: ", resp.ID)
		if err = d.Client.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
			return rsp.Response{}, err
		}
	}

	return rsp.Response{Type: "run", Payload: fmt.Sprint("Container ran with id: ", resp.ID)}, nil
}

func (d *Docker) ContainerStop(id string) (rsp.Response, error) {
	timeout := time.Second * 10
	if err := d.Client.ContainerStop(context.Background(), id, &timeout); err != nil {
		return rsp.Response{}, err
	}

	return rsp.Response{Type: "stop", Payload: fmt.Sprint("Container Stopped for id: ", id)}, nil
}
