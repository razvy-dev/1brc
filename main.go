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

func processCity(city city, cities map[string]CityState) map[string]CityState {
	if cities[city.name] {
		if city.temp > cities[city.max] {
			cities[city.name] = city.temp
		}

		if city.temp < cities[city.min] {
			cities[city.name] = city.temp
		}

		cities[city.name].mean = cities[city.name].mean + city.temp / 2

		return cities

	} else {
		var new CityState
		new.name = city.name
		new.max, new.min, new.mean = new.temp
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

				fmt.Printf("Temperature in %s is %.2f\n", city.name, city.temp)
				process(city, cities)
			}

			break
		}

		if err != nil {
			return fmt.Errorf("error reading from file: %w", err)
		}

		city := parseRow(textLine)

		fmt.Printf("Temperature in %s is %.2f\n", city.name, city.temp)

		process(city, cities)
	}

	return nil
}

func main() {
	path := "/home/nvt-dev/projects/1brc/measurements.txt"

	cities := make(map[string]CityState) // use hash map for this

	readLinesFromFile(path, cities)
}