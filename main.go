package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
	"io"
	"time"
	"sync"
)

type Chunk struct {
	Start int64
	End int64
}

type City struct {
	name string;
	temp float32;
}

type CityState struct {
	name string;
	max float32;
	min float32;
	sum float32;
	count int;
}

func Split(r rune) bool {
	return r == ';' || r == '\n'
}

func worker(id int, jobs <-chan []string, wg *sync.WaitGroup, cities *sync.Map) {
	defer wg.Done()

	for batch := range jobs {
		for row := range batch {
			city := parseRow(batch[row])
			processCity(city, cities)
		}
	}
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

func processCity(city City, cities *sync.Map) *sync.Map {
	if new, ok := cities.Load(city.name); ok {
		new := new.(CityState)
		if city.temp > new.max {
			new.max = city.temp
		}

		if city.temp < new.min {
			new.min = city.temp
		}

		new.sum += city.temp
		new.count++

		cities.Store(new.name, new)

		return cities

	} else {
		var new CityState
		new.name = city.name
		new.max, new.min, new.sum = city.temp, city.temp, city.temp
		cities.Store(new.name, new)

		return cities
	}
}

func process_v2(path string, chunk Chunk, workerID int, wg *sync.WaitGroup, cities *sync.Map, workers int, size int) error {
	file, err := os.Open(path)

	if err != nil {
		fmt.Print("am I able to open this file?")
		return err
	}

	defer file.Close()

	_, err = file.Seek(chunk.Start, io.SeekStart)

	if err != nil {
		fmt.Printf("[Worker %d] Seek error: %v\n", workerID, err)
		return err
	}

	reader := bufio.NewReader(file)
	currentOffset := chunk.Start

	if chunk.Start != 0 {
		discardedBytes, err := reader.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return err
		}

		currentOffset += int64(len(discardedBytes))
	}

	jobs := make(chan []string, workers)
	var wg_subs sync.WaitGroup

	for w := 1; w <= workers; w++ {
		wg_subs.Add(1)
		go worker(w, jobs, &wg_subs, cities)
	}

	batch := make([]string, 0, size)

	for currentOffset < chunk.End {
		lineBytes, err := reader.ReadBytes('\n')
		lineLen := int64(len(lineBytes))

		currentOffset += lineLen

		if lineLen > 0 {
			line := string(lineBytes[:])
			line = strings.TrimRight(line, "\r\n")

			batch = append(batch, line)

			if len(batch) == size {
				jobs <- batch
				batch = make([]string, 0, size)
				fmt.Print("Starting to process next batch")
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}
	}

	if len(batch) > 0 {
		jobs <- batch
	}

	close(jobs)

	return nil
}

func process(path string, workers int, size int, cities *sync.Map) error {
	fileHandle, err := os.Open(path)

	if err != nil {
		return err
	}

	defer fileHandle.Close()

	jobs := make(chan []string, workers)
	var wg sync.WaitGroup

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg, cities)
	}

	scanner := bufio.NewReader(fileHandle)

	batch := make([]string, 0, size)

	for {

		line, err := scanner.ReadString('\n')

		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")

			batch = append(batch, line)

			if len(batch) == size {
				jobs <- batch
				batch = make([]string, 0, size)
				fmt.Print("Starting to process next batch")
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

	}

	if len(batch) > 0 {
		jobs <- batch
	}

	close(jobs)

	wg.Wait()

	return nil
}

func main() {
	path := "/home/nvt-dev/projects/1brc/measurements.txt"

	start := time.Now()

	var cities sync.Map // use hash map for this

	// process(path, 4, 50, &cities)

	workers := 4

	fileInfo, err := os.Stat(path)

	if err != nil {
		fmt.Print("something went wrong when getting the file stats")
		return
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(workers)

	var wg sync.WaitGroup

	fmt.Print("am i getting here?")

	for i := 0; i < workers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize

		if i == workers - 1 {
			end = fileSize // the last worker gets it in full, so no bytes are orphaned
		}

		chunk := Chunk{Start: start, End: end}

		wg.Add(1)

		fmt.Print("is the worker added?")
		
		go process_v2(path, chunk, i, &wg, &cities, 4, 50)

		fmt.Print("!!!")
	}

	cities.Range(func(key, value any) bool {
		city_name := key.(string)
		city_state := value.(CityState)

		fmt.Printf("City: %s, Max: %.2f, Min: %.2f, Mean: %.2f\n", city_name, city_state.max, city_state.min, city_state.sum / float32(city_state.count))

		return true
	})

	elapsed := time.Since(start)

	fmt.Printf("This took %s", elapsed)
}