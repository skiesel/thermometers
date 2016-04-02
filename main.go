package main

import (
	"fmt"
	"github.com/skiesel/thermometers/sensors"
	"time"
)

func main() {
	for {
		readings := sensors.GetThermometerReadings()
		for i, reading := range readings {
			fmt.Printf("%d) %g°C %g°F\n", i, reading.Celsius, reading.Fahrenheit)
		}
		time.Sleep(time.Second)
		fmt.Println("")

	}
}
