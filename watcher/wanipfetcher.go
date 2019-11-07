package watcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func FetchWanIP() (string, error) {
	che := make(chan error)
	chr := make(chan string)

	go func() {
		resp, err := http.Get("ifconfig.me")
		if err != nil {
			che <- err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			che <- err
		}

		chr <- body
	}()

	select {
	case err := <-che:
		return "", err
	case res := <-chr:
		return res, nil
	case <-time.After(5 * time.Second):
		return "", fmt.Error("Timeout waiting for wan IP")
	}
}
