package main

import (
	"context"
	"fmt"
	/*
		"os"
		"time"

		"shadow/mqtt"
	*/

	"shadow/watcher"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	pmg "github.com/eclipse/paho.mqtt.golang"
	"github.com/subosito/gotenv"
)

var (
	mq              *pmg.Client
	watchConfigList []string
	watchList       []Actor
)

func main() {
	gotenv.Load()

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Cannot initialize client: ", err)
	}

	/*
		mq = &mqtt.NewClient(
			os.Getenv("MQTT_BROKER_URL"),
			os.Getenv("MQTT_BROKER_USER"),
			os.Getenv("MQTT_BROKER_PASS"),
			"shadow-"+fmt.Sprint(time.Now().Local().Format(time.RFC3339)),
			"",
		)
	*/

	watchConfigList = []string{"registry.gitlab.com/sangkuriang-dev/transmissor-be/backend:devel", "golang/1.13-alpine"}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println("Cannot get container list : ", err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})

	for _, image := range images {
		for _, t := range image.RepoTags {
			for _, w := range watchConfigList {
				if w == t {
					//fmt.Printf("%v %v is being watched\n", image.Containers, image.RepoTags)
					watchList = append(watchList, Actor{ImageNames: image.RepoTags})
				}
			}
		}
	}

	for i, w := range watchList { // among the watched list
		for _, in := range w.ImageNames { // which has these images alias
			for _, container := range containers { // find containers
				fmt.Printf("<%v> (%v)", container.ID, container.Image)
				if container.Image == in {
					//fmt.Printf("container id %s, with image %s is being watched", container.ID[:10], in)
					watchList[i].ContainerIDs = append(watchList[i].ContainerIDs, container.ID)
				}
				fmt.Printf("\n")
			}
		}
	}

	fmt.Println(watchList)
}
