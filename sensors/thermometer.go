package sensors

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	// the base directory where we expect sensor files to be
	sensorBaseDirectory = "/sys/bus/w1/devices/"
)

var (
	sensorPaths = []string{}
)

type TemperatureReading struct {
	Celsius    float64
	Fahrenheit float64
}

func init() {
	// add the following modules if they aren't already
	command := exec.Command("modprobe", "w1-gpio")
	command.Run()
	command = exec.Command("modprobe", "w1-therm")
	command.Run()

	//get a list of sensor candidates
	dirs, err := ioutil.ReadDir(sensorBaseDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	//let's see what sensors might exist
	for _, dir := range dirs {
		//we're looking for the presence of w1_slave file inside the sensor folders
		sensorPath := sensorBaseDirectory + dir.Name() + "/w1_slave"

		//if there are multiple sensors, we can end up with a book keeping folder, so ignore it
		if _, err := os.Stat(sensorPath); err == nil {
			sensorPaths = append(sensorPaths, sensorPath)
		}
	}
}

func getThermometerReading(sensor string) TemperatureReading {

	stringContents := ""

	//we'll persistently try to read from the file until we get a good reading or we lose the sensor
	for goodReading := false; ; {
		contents, err := ioutil.ReadFile(sensor)
		if err != nil {
			panic(err)
		}

		//check for "YES" which will signify a valid reading
		stringContents = string(contents[:])
		goodReading = strings.Contains(stringContents, "YES")

		if goodReading {
			break
		}
		//we're polling the file, so wait a short while before retrying
		time.Sleep(time.Millisecond * 200)
	}

	//look for an actual temperature value
	tempAvailable := strings.Contains(stringContents, "t=")

	if !tempAvailable {
		panic("Could not find 't=' inside sensor file.")
	}

	//get rid of any new lines or space characters that will break parsing a float
	stringContents = strings.TrimSpace(stringContents)

	reading := strings.SplitAfter(stringContents, "t=")

	//try to parse the temperature reading
	raw, err := strconv.ParseFloat(reading[1], 64)
	if err != nil {
		panic(err)
	}

	celsius := raw / 1000.0

	return TemperatureReading{
		Celsius:    celsius,
		Fahrenheit: celsius*9.0/5.0 + 32.0,
	}
}

func GetThermometerReadings() []TemperatureReading {
	temperatures := []TemperatureReading{}

	for _, sensor := range sensorPaths {
		reading := getThermometerReading(sensor)
		temperatures = append(temperatures, reading)
	}

	return temperatures
}
