package watcher

import (
	"fmt"
)

type Watcher struct {
	WatchList       []Actor
	WatchConfigList []string
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

type Image interface {
	Run(image string, conf containerRunConfig)
}

func NewWatcher(period time.Duration) *Watcher {
}

func (w *Watcher) Watch() {
	go func() {

	}()
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
	if isInWatchConfigList(imageName) == -1 {
		w.WatchConfigList = append(w.WatchConfigList, imageName)
	}
}

func (w *Watcher) RemoveImageFromWatchList(imageName string) {
	if index := isInWatchConfigList(imageName); index != -1 {
		l := len(w.WatchConfigList)
		w.WatchConfigList[index] = w.WatchConfigList[l-1]
		w.WatchConfigList[l-1] = ""
		w.WatchConfigList = w.WatchConfigList[:l-1]
	}
}
