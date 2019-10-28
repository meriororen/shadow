package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"shadow/docker"
	//	"github.com/docker/docker/api/types"
	//"github.com/docker/docker/client"
)

type Command struct {
	Type    string
	Payload []byte
}

type CmdPull struct {
	ImageName string `json:"image"`
}

type CmdLogin struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Token    string `json:"token,omitempty"`
	Password string `json:"password,omitempty"`
}

func Exec(cmd Command, args ...interface{}) (res interface{}, err error) {
	if cmd.Payload == nil {
		return "", fmt.Errorf("Error executing command, payload must not be empty!")
	}

	switch cmd.Type {
	case "pull":
		pullcmd := CmdPull{}
		if err := json.Unmarshal(cmd.Payload, &pullcmd); err != nil {
			log.Println("Unprocessable pull payload", err)
		}

		if res, err = PullImage(pullcmd); err != nil {
			return "Error Pulling!", err
		}
	case "login":
		logincmd := CmdLogin{}
		if err := json.Unmarshal(cmd.Payload, &logincmd); err != nil {
			log.Println("Unprocessable login payload", err)
		}

		if res, err = Login(logincmd); err != nil {
			return "Error Logging in!", err
		}
	}

	return res, nil
}

func PullImage(pull CmdPull) (resp interface{}, err error) {
	log.Println("pulling image: ", pull.ImageName)

	if resp, err = docker.Default.ImagePull(pull.ImageName); err != nil {
		log.Fatal("Failed to pull image: ", err)
	}

	return resp, nil
}

func Login(login CmdLogin) (resp interface{}, err error) {
	log.Println("Loggin in..")
	log.Println("Registry: ", login.URL)
	log.Println("Username: ", login.Username)
	log.Println("Password: ", login.Password)
	log.Println("Token: ", login.Token)

	if resp, err = docker.Default.RegistryLogin(login.URL, login.Username, login.Password, login.Token); err != nil {
		log.Fatal("Failed to log in: ", err)
	}

	return resp, nil
}
