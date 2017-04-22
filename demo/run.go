package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/aio"
)

func main() {
	e := edison.NewAdaptor()
	lightSensor := aio.NewGroveLightSensorDriver(e, "0")
	led := gpio.NewLedDriver(e, "8")
	screen := i2c.NewGroveLcdDriver(e)

	work := func() {
		for {
			lightScore, _ := lightSensor.Read()
			if lightScore < 350 {
				screen.SetRGB(255, 0, 0)
			} else {
				screen.SetRGB(0, 255, 0)
			}
		}
	}

	robot := gobot.NewRobot("buttonBot",
		[]gobot.Connection{e},
		[]gobot.Device{screen, lightSensor, led},
		work,
	)

	robot.Start()
}