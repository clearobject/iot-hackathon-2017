package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/aio"
	"time"
	"fmt"
	"math"
)

const (
	lightTolerance = 50
)

func average(slc[]int)int {
	total := 0
	for _, v := range slc {
		total += v
	}
	return total/len(slc)
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

func main() {
	e := edison.NewAdaptor()
	lightSensor := aio.NewGroveLightSensorDriver(e, "0", 250)
	temperatureSensor := aio.NewGroveTemperatureSensorDriver(e, "3")
	screen := i2c.NewGroveLcdDriver(e)

	fmt.Println(calibrateLighting(lightSensor))
	averageLight := calibrateLighting(lightSensor)

	work := func() {
		for {
			lightScore, _ := lightSensor.Read()
			temperature := temperatureSensor.Temperature()
			fmt.Println(temperature)
			delta := averageLight - lightScore
			if math.Abs(float64(delta)) > lightTolerance {
				screen.SetRGB(255, 0, 0)
			} else {
				screen.SetRGB(0, 255, 0)
			}

		}
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{e},
		[]gobot.Device{screen, lightSensor, temperatureSensor},
		work,
	)

	robot.Start()
}