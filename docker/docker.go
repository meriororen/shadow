package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"shadow/rsp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

type Docker struct {
	Client           *client.Client
	ResponseChan     chan rsp.Response
	savedCredentials types.AuthConfig
}

var Default Docker

func NewDocker() (Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Cannot initialize client: ", err)
	}

	return Docker{Client: cli}, nil
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

func (d *Docker) RegistryLogin(url string, user string, pass string, token string) (interface{}, error) {
	creds := types.AuthConfig{
		ServerAddress: url,
		Username:      user,
		Password:      pass,
		IdentityToken: token,
	}

	if idtoken, err := attemptLogin(creds); err != nil {
		return "", err
	} else {
		d.savedCredentials = creds

		if idtoken != "" {
			d.savedCredentials.IdentityToken = idtoken
			log.Println("Login success, id token: ", idtoken)
		} else {
			log.Println("Id token is empty")
		}
	}

	return "Login Successful", nil
}

func (d *Docker) ImagePull(imageName string) (interface{}, error) {
	auth := types.AuthConfig{
		Username: d.savedCredentials.Username,
		Password: d.savedCredentials.Password,
	}

	var strauth string
	if encauth, err := json.Marshal(auth); err != nil {
		return nil, err
	} else {
		strauth = base64.URLEncoding.EncodeToString(encauth)
	}

	out, err := d.Client.ImagePull(context.Background(), imageName, types.ImagePullOptions{RegistryAuth: strauth})
	if err != nil {
		log.Fatal("Cannot image pull")
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return "Success Pulled", nil
}

func (d *Docker) ImageList(imageName string) (interface{}, error) {
	images, err := d.Client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal("Cannot list out all images: ", err)
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

	log.Println(imgs)

	return "Success List Image", nil
}

func (d *Docker) ContainerList(imageID string) (interface{}, error) {
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal("Cannot get container list : ", err)
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

	log.Println(cnts)

	return "Listed Container", nil
}

func (d *Docker) ContainerStart(imageName string) (interface{}, error) {

	return "Started Container", nil
}

func (d *Docker) ContainerStop(id string) (interface{}, error) {

	return "Stopped Container", nil
}
