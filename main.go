package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"strconv"
)

type City struct {
	name string;
	temp float32;
}

type CityState struct {
	name string;
	max float32;
	min float32;
	mean float32;
}

func Split(r rune) bool {
	return r == ';' || r == '\n'
}

func parseRow(data string) City {
	result := strings.FieldsFunc(data, Split) // this returns an array of strings

	name := result[0]

	value, err := strconv.ParseFloat(result[1], 32)

	if err != nil {
		fmt.Print("error parsing a row: ", err)
	}

	temp := float32(value)

	var city City

	city.name = name
	city.temp = temp

	return city
}

func processCity(city City, cities map[string]CityState) map[string]CityState {
	_, ok := cities[city.name]
	var new CityState

	if ok {
		new = cities[city.name]
		if city.temp > cities[city.name].max {
			new.max = city.temp
		}

		if city.temp < cities[city.name].min {
			new.min = city.temp
		}

		new.mean = cities[city.name].mean + city.temp / 2

		cities[city.name] = new

		return cities

	} else {
		new.name = city.name
		new.max, new.min, new.mean = city.temp, city.temp, city.temp
		cities[new.name] = new

		return cities
	}
}

func process(path string, cities map[string]CityState) error {
	fileHandle, err := os.Open(path)

	if err != nil {
		return err
	}

	defer fileHandle.Close()

	scanner := bufio.NewReader(fileHandle)

	for {
		textLine, err := scanner.ReadString('\n')

		if err == io.EOF {
			if len(textLine) != 0 {
				city := parseRow(textLine)

				processCity(city, cities)
			}

			break
		}

		if err != nil {
			return fmt.Errorf("error reading from file: %w", err)
		}

		city := parseRow(textLine)

		processCity(city, cities)
	}

	return nil
}

func main() {
	path := "/home/nvt-dev/projects/1brc/measurements.txt"

	cities := make(map[string]CityState) // use hash map for this

	process(path, cities)

	for k := range cities {
		fmt.Printf("In %s we have max: %.2f, min: %.2f, average: %.2f\n", cities[k].name, cities[k].max, cities[k].min, cities[k].mean)
	}
}