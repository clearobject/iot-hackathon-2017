package main

import (
	"os"

	"fmt"
	"math"
	"time"
	"encoding/json"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"gobot.io/x/gobot/platforms/mqtt"
)

const (
	lightTolerance = 50
	mqttHost = "tcp://104.154.233.174:1883"
)

type Event struct {
	Name string
	Time int64
	Source string
	Temperature float64
}

func average(slc []int) int {
	total := 0
	for _, v := range slc {
		total += v
	}
	return total / len(slc)
}

func calibrateLighting(lightSensor *aio.GroveLightSensorDriver) int {
	completeTime := time.Now().Add(3 * time.Second)
	readings := []int{}
	for time.Now().Unix() < completeTime.Unix() {
		reading, _ := lightSensor.Read()
		readings = append(readings, reading)
	}
	return average(readings)
}

func createEventJSON(ts *aio.GroveTemperatureSensorDriver) []byte {
	event := Event{
		"edgeTrigger",
		time.Now().Unix(),
		os.Getenv("HOSTNAME"),
		ts.Temperature(),
	}
	output, _ := json.Marshal(event)
	return []byte(output)
}

func main() {
	e := edison.NewAdaptor()
	lightSensor := aio.NewGroveLightSensorDriver(e, "0", 250)
	temperatureSensor := aio.NewGroveTemperatureSensorDriver(e, "3")
	screen := i2c.NewGroveLcdDriver(e)
	mqttAdaptor := mqtt.NewAdaptorWithAuth(mqttHost, os.Getenv("HOSTNAME"), "test", "testpass")

	fmt.Println(calibrateLighting(lightSensor))
	averageLight := calibrateLighting(lightSensor)

	work := func() {
		for {
			lightScore, _ := lightSensor.Read()
			delta := averageLight - lightScore
			if math.Abs(float64(delta)) > lightTolerance {
				screen.SetRGB(255, 0, 0)
				mqttAdaptor.Publish("reflectors", createEventJSON(temperatureSensor))
				time.Sleep(1 * time.Second)
			} else {
				screen.SetRGB(0, 255, 0)
			}

		}
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{e, mqttAdaptor},
		[]gobot.Device{screen, lightSensor, temperatureSensor},
		work,
	)

	robot.Start()
}
