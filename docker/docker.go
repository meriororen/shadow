package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

type Docker struct {
	Client           *client.Client
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

	log.Println(auth)

	var strauth string
	if encauth, err := json.Marshal(auth); err != nil {
		log.Println("Cannot marshal auth")
		return nil, err
	} else {
		log.Println("encauth :", encauth)
		strauth = base64.URLEncoding.EncodeToString(encauth)
	}

	log.Println("Trying to pull with registry auth: ", strauth)

	out, err := d.Client.ImagePull(context.Background(), imageName, types.ImagePullOptions{RegistryAuth: strauth})
	if err != nil {
		log.Fatal("Cannot image pull")
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	return "Success Pulled", nil
}

func (d *Docker) ContainerList(imageName string) (interface{}, error) {

	return "Listed Container", nil
}

func (d *Docker) ContainerStart(imageName string) (interface{}, error) {

	return "Started Container", nil
}

func (d *Docker) ContainerStop(id string) (interface{}, error) {

	return "Stopped Container", nil
}
