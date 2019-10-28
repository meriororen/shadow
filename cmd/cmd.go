package cmd

import (
	"log"
	/*
		"github.com/docker/docker/api/types"
		"github.com/docker/docker/client"
	*/)

type Command struct {
	Type    string
	Payload interface{}
}

func Exec(cmd string, args ...interface{}) (res string, err error) {
	switch cmd {
	case "pull":
		if res, err = PullImage(args[0].(string)); err != nil {
			return "Error Pulling!", err
		}
	}

	return res, nil
}

func PullImage(name string) (string, error) {
	log.Println("pulling image: ", name)

	return "Pulled!", nil
}
