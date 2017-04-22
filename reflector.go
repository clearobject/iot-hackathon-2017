package main

import (
	"os"

	"encoding/json"
	"math"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"gobot.io/x/gobot/platforms/mqtt"
)

const (
	lightTolerance = 45
	mqttHost       = "tcp://104.154.46.156:1883"
)

// Event structure; basis for JSON / MQTT messages
type Event struct {
	Name        string
	Time        int64
	Source      string
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
	host, _ := os.Hostname()
	event := Event{
		"sensorTrigger",
		time.Now().Unix(),
		host,
		ts.Temperature(),
	}
	output, _ := json.Marshal(event)
	return []byte(output)
}

func main() {
	host, _ := os.Hostname()
	e := edison.NewAdaptor()
	lightSensor := aio.NewGroveLightSensorDriver(e, "0", 250)
	temperatureSensor := aio.NewGroveTemperatureSensorDriver(e, "3")
	screen := i2c.NewGroveLcdDriver(e)
	mqttAdaptor := mqtt.NewAdaptorWithAuth(mqttHost, host, "test", "testpass")

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

	robot := gobot.NewRobot("reflectorBot",
		[]gobot.Connection{e, mqttAdaptor},
		[]gobot.Device{screen, lightSensor, temperatureSensor},
		work,
	)

	robot.Start()
}
