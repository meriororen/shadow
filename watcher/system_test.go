package watcher

import (
	"testing"
)

func TestCheckCpuTemp(t *testing.T) {
	if res, err := CheckCpuTemp(); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}

func TestCheckUpTime(t *testing.T) {
	if res, err := CheckUpTime(); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}
