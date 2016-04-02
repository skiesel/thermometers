package thermometers

import (
	"fmt"
	"testing"
)

func TestBasic(*testing.T) {
	readings := GetThermometerReadings()
	for i, reading := range readings {
		fmt.Printf("%d) %g°C %g°F\n", i, reading.Celsius, reading.Fahrenheit)
	}
}
