package watcher

import (
	"testing"
)

func TestFetchWanIP(t *testing.T) {
	if ip, err := FetchWanIP(); err != nil {
		t.Error("Error: ", err)
	}
	t.Log("Fetched IP: ", ip)
}
