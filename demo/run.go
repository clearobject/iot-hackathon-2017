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
	sensitivity = 25
)

func average(slc[]int)int {
	total := 0
	for _, v := range slc {
		total +=v
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
	screen := i2c.NewGroveLcdDriver(e)

	fmt.Println(calibrateLighting(lightSensor))
	averageLight := calibrateLighting(lightSensor)

	work := func() {
		for {
			lightScore, _ := lightSensor.Read()
			delta := averageLight - lightScore
			if math.Abs(float64(delta)) > sensitivity {
				screen.SetRGB(255, 0, 0)
				time.Sleep(1 * time.Second)
			}
			screen.SetRGB(0, 255, 0)
		}
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{e},
		[]gobot.Device{screen, lightSensor},
		work,
	)

	robot.Start()
}