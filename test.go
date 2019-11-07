package main

import "fmt"
import "os/exec"

func main() {
	test := []string{"nmcli", "con", "show", "active"}

	out, err := exec.Command(test[0], test[1:]...).CombinedOutput()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println(string(out))
}
