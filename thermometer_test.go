package main

import (
	"fmt"
	"github.com/skiesel/thermometers/sensors"
	"testing"
)

func TestBasic(*testing.T) {
	readings := sensors.GetThermometerReadings()
	for i, reading := range readings {
		fmt.Printf("%d) %g°C %g°F\n", i, reading.Celsius, reading.Fahrenheit)
	}
}
