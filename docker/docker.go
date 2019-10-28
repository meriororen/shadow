package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Docker struct {
	Client *client.Client
}

var Default Docker

func NewDocker() (Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Cannot initialize client: ", err)
	}

	return Docker{Client: cli}, nil
}

func (d Docker) RegistryLogin(url string, user string, pass string, token string) (interface{}, error) {
	auth := types.AuthConfig{
		ServerAddress: url,
		Username:      user,
		Password:      pass,
		IdentityToken: token,
	}

	return d.Client.RegistryLogin(context.Background(), auth)
}
