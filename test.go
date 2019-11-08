package main

import "strings"
import "os/exec"
import "fmt"

func main() {
	test := "sed s/STATION_ID/STATIONAD/g /env/command.env"
	r := strings.Split(test, " ")

	//	test := []string{"nmcli", "con", "show", "active"}

	out, err := exec.Command(r[0], r[1:]...).CombinedOutput()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println(string(out))
}
