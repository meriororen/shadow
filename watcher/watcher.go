package watcher

import (
	"context"
	"log"
	//	"syscall"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"shadow/cmd"
	"shadow/docker"
	"shadow/status"

	"github.com/docker/docker/api/types"
	//	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

var wg sync.WaitGroup
var Default *Watcher

type WatchConfig struct {
	ImageName string
	AutoPull  bool
	HBPeriod  time.Duration
}

type Watcher struct {
	WatchList       map[string]*Actor
	WatchConfigList []WatchConfig
	StatusChan      chan status.Status
	CommandChan     chan cmd.Command

	heartBeatPeriod time.Duration
}

type Actor struct {
	ImageNames   []string
	ContainerIDs []string
	watchConfig  *WatchConfig
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

func Init() *Watcher {
	wt := &Watcher{
		WatchList:       make(map[string]*Actor),
		WatchConfigList: []WatchConfig{},
		StatusChan:      make(chan status.Status),
		CommandChan:     make(chan cmd.Command),
	}

	wt.AddImageToWatch(WatchConfig{
		ImageName: "System",
		HBPeriod:  3 * time.Second,
	})

	return wt
}

func (w *Watcher) WatchImages() {
	log.Println("Waiting..")
	wg.Wait()
}

func getSystemStatus() (status.System, error) {
	ms, err := memory.Get()
	if err != nil {
		log.Fatal("SystemStat: Cannot get memory status")
	}

	wlan := os.Getenv("WIFI_INTERFACE_NAME")
	ma, err := ioutil.ReadFile("/sys/class/net/" + wlan + "/address")
	if err != nil {
		log.Println("Cannot get mac address")
		ma = []byte("Error")
	}

	return status.System{
		Memory: status.Memory{
			Total:  ms.Total,
			Free:   ms.Free,
			Used:   ms.Used,
			Cached: ms.Cached,
		},
		MacAddress: string(ma),
	}, nil
}

func (w *Watcher) WatchContainer() {
	for _, actor := range w.WatchList {
		for cont := range actor.ContainerIDs {
			cont = cont
		}
	}
}

func (w *Watcher) WatchAll() {
	w.WatchImages()
}

func (w *Watcher) isInWatchConfigList(imageName string) int {
	for i, im := range w.WatchConfigList {
		if imageName == im.ImageName {
			return i
		}
	}
	return -1
}

func (w *Watcher) AddImageToWatch(config WatchConfig) {
	if w.isInWatchConfigList(config.ImageName) == -1 {
		w.WatchConfigList = append(w.WatchConfigList, config)
	}

	// dispatch a goroutine for each watched item
	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(config.HBPeriod):
				//log.Println("HB for", config.ImageName)

				status := status.Status{
					LocalTime: time.Now().Local(),
				}
				if config.ImageName == "System" {
					if ss, err := getSystemStatus(); err != nil {
						log.Println("Error getting system status:", err)
						break
					} else {
						status.Payload = ss
					}
				} else {
					// TODO: container status
				}

				w.StatusChan <- status
			}
		}
	}()

	w.addImagesToWatchList()
	w.addRunningContainersToWatchList()
}

func (w *Watcher) RemoveImageFromWatchList(imageName string) {
	if index := w.isInWatchConfigList(imageName); index != -1 {
		l := len(w.WatchConfigList)
		log.Println("Removing ", imageName)
		wg.Done()
		w.WatchConfigList[index] = w.WatchConfigList[l-1]
		w.WatchConfigList = w.WatchConfigList[:l-1]
		// remove from watchlist
		delete(w.WatchList, imageName)
	}
}

func (w *Watcher) addImagesToWatchList() {
	images, err := docker.Default.Client.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal("Cannot add image in local to watchlist ", err)
	}

	for _, image := range images {
		for _, t := range image.RepoTags {
			for k, watch := range w.WatchConfigList {
				if watch.ImageName == t {
					w.WatchList[watch.ImageName] = &Actor{
						ImageNames:  image.RepoTags,
						watchConfig: &w.WatchConfigList[k],
					}
				}
			}
		}
	}
}

func (w *Watcher) addRunningContainersToWatchList() {
	containers, err := docker.Default.Client.ContainerList(context.Background(), types.ContainerListOptions{})
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
