package watcher

import (
	"context"
	"log"
	"os"
	"strconv"
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
	WanIp           string

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

	syshbp := 20 // default system heartbeat
	if hbpstr := os.Getenv("SYSTEM_HB_PERIOD_S"); hbpstr != "" {
		syshbp, _ = strconv.Atoi(hbpstr)
	}

	wanIpPollPeriod := 300 // default wanip poll

	wt.AddImageToWatch(WatchConfig{
		ImageName: "System",
		HBPeriod:  time.Duration(syshbp) * time.Second,
	})

	wt.AddImageToWatch(WatchConfig{
		ImageName: "WanIp",
		HBPeriod:  time.Duration(wanIpPollPeriod) * time.Second,
	})

	return wt
}

func (w *Watcher) WatchImages() {
	log.Println("Waiting..")
	wg.Wait()
}

func (w *Watcher) getSystemStatus() (status.System, error) {
	ms, err := memory.Get()
	if err != nil {
		log.Println("SystemStat: Cannot get memory status")
	}

	cputemp, err := CheckCpuTemp()
	if err != nil {
		log.Println("SystemStat: Cannot get Cpu temperature", err)
	}

	sysuptime, err := CheckUpTime()
	if err != nil {
		log.Println("Cannot check System Uptime")
	}

	return status.System{
		WanIp:        w.WanIp,
		CpuTemp:      cputemp,
		SystemUpTime: sysuptime,
		Memory: status.SystemMemory{
			Total: ms.Total,
			Free:  ms.Free,
			//			Used:   ms.Used,
			//			Cached: ms.Cached,
		},
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

func (w *Watcher) statusPoll(imageName string) (status.Status, error) {
	stat := status.Status{
		LocalTime: time.Now().Local(),
	}
	switch imageName {
	case "System":
		if ss, err := w.getSystemStatus(); err != nil {
			log.Println("Error getting system status:", err)
			return status.NilStatus, err
		} else {
			stat.Type = "System"
			stat.Payload = ss
		}
	case "WanIp": /* doesn't need to return status */
		if ip, err := FetchWanIp(); err != nil {
			log.Println(err)
			return status.NilStatus, err
		} else {
			w.WanIp = ip
			return status.NilStatus, nil
		}
	default:
		// container status
	}
	return stat, nil
}

func (w *Watcher) AddImageToWatch(config WatchConfig) {
	if w.isInWatchConfigList(config.ImageName) == -1 {
		w.WatchConfigList = append(w.WatchConfigList, config)
	}

	// dispatch a goroutine for each watched item
	wg.Add(1)
	go func() {
		/* always send first time */
		if status, err := w.statusPoll(config.ImageName); err == nil {
			w.StatusChan <- status
		}

		for {
			select {
			case <-time.After(config.HBPeriod):
				//log.Println("HB for", config.ImageName)
				if s, err := w.statusPoll(config.ImageName); err == nil {
					w.StatusChan <- s
				} else {
					break
				}
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
		log.Println("Cannot add image in local to watchlist ", err)
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
		log.Println("Cannot get container list : ", err)
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
