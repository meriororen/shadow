package watcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Watcher struct {
	WatchList       map[string]*Actor
	WatchConfigList []string
	cli             *client.Client
	heartBeatPeriod time.Duration
}

type Actor struct {
	ImageNames   []string
	ContainerIDs []string
}

type Container interface {
	Exec(container string, cmd []string)
	Stop(container string)
	Remove(container string)
}

type containerRunConfig struct {
}

type Image interface {
	Run(image string, conf containerRunConfig)
}

func NewWatcher() *Watcher {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Cannot initialize client: ", err)
	}

	return &Watcher{
		WatchList:       make(map[string]*Actor),
		WatchConfigList: []string{},
		cli:             cli,
	}
}

func (w *Watcher) WatchImages() {
	for _, actor := range w.WatchList {
		for image := range actor.ImageNames {
			image = image
		}
	}
}

func (w *Watcher) WatchContainer() {
	for _, actor := range w.WatchList {
		for cont := range actor.ContainerIDs {
			cont = cont
		}
	}
}

func (w *Watcher) WatchAll() {
	for {
		w.WatchImages()
		w.WatchContainer()
	}
}

func (w *Watcher) isInWatchConfigList(imageName string) int {
	for i, im := range w.WatchConfigList {
		if imageName == im {
			return i
		}
	}
	return -1
}

func (w *Watcher) AddImageToWatch(imageName string) {
	if w.isInWatchConfigList(imageName) == -1 {
		w.WatchConfigList = append(w.WatchConfigList, imageName)
	}

	w.addImagesToWatchList()
	w.addRunningContainersToWatchList()
}

func (w *Watcher) RemoveImageFromWatchList(imageName string) {
	if index := w.isInWatchConfigList(imageName); index != -1 {
		l := len(w.WatchConfigList)
		w.WatchConfigList[index] = w.WatchConfigList[l-1]
		w.WatchConfigList[l-1] = ""
		w.WatchConfigList = w.WatchConfigList[:l-1]
		// remove from watchlist
		delete(w.WatchList, imageName)
	}
}

func (w *Watcher) addImagesToWatchList() {
	images, err := w.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal("Cannot add image in local to watchlist ", err)
	}

	for _, image := range images {
		for _, t := range image.RepoTags {
			for _, watch := range w.WatchConfigList {
				if watch == t {
					w.WatchList[watch] = &Actor{ImageNames: image.RepoTags}
				}
			}
		}
	}
}

func (w *Watcher) addRunningContainersToWatchList() {
	containers, err := w.cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal("Cannot get container list : ", err)
	}

	for _, actor := range w.WatchList {
		for _, in := range actor.ImageNames {
			for _, container := range containers {
				if container.Image == in {
					(*actor).ContainerIDs = append(actor.ContainerIDs, container.ID)
				}
			}
		}
	}
}
