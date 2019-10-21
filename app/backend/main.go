package main

import (
	"log"
	"time"

	"transmissor-be/api"
	"transmissor-be/gps"
	"transmissor-be/sensormanager"
	"transmissor-be/test"

	"github.com/subosito/gotenv"
)

var sm *sensormanager.SensorManager
var err error

func main() {
	gotenv.Load()

	if sm, err = sensormanager.NewSensorManager(&test.DummySensor{}); err != nil {
		log.Fatal(err)
	}
	sm.ActivateDataGetter(4 * time.Second)

	theapi, _ := api.NewAPIServer(sm, &gps.GPS{})

	log.Println("Started the server at: ", time.Now().Local())
	theapi.Serve()

	/*
		  gps, _ := gps.NewGPS("/dev/ttyACM0", 9600)
		  gps.StartListening()

			water, err := wqc24.NewWQC24("/dev/ttyUSB0", 9600)
			if err != nil {
				log.Fatal(err)
			}

			if sm, err = sensormanager.NewSensorManager(&water); err != nil {
				log.Fatal(err)
			}

			var pollperiod int
			if pollperiod, err = strconv.Atoi(os.Getenv("SENSOR_POLL_PERIOD_S")); err != nil {
				pollperiod = 10
			}
			sm.ActivateDataGetter(time.Duration(pollperiod) * time.Second)

			theapi, _ := api.NewAPIServer(sm, gps)

			log.Println("Started the server at: ", time.Now().Local())
			theapi.Serve()

		  //var windData wind.Data
		  //wind, _ := wind.NewWind("/dev/ttyUSB1", 9600)
		  //err = wind.GetData(&windData, 10*time.Second)
		  //if err != nil {
		  //  log.Println(err)
		  //}
	*/
}
