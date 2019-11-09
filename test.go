package main

/*
import "strings"
import "os/exec"
*/
import "log"
import "github.com/ghodss/yaml"
import "io/ioutil"

func main() {
	b, _ := ioutil.ReadFile("/srv/docker-compose.yml")
	y, err := yaml.YAMLToJSON(b)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(y))

	/*
		test := "sed s/STATION_ID/STATIONAD/g /env/command.env"
		r := strings.Split(test, " ")

		//	test := []string{"nmcli", "con", "show", "active"}

		out, err := exec.Command(r[0], r[1:]...).CombinedOutput()
		if err != nil {
			fmt.Println("ERROR", err)
		}
		fmt.Println(string(out))
	*/
}
