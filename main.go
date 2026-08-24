package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type city struct {
	name string;
	max float32;
	min float32;
	mean float32;
}



func readLinesFromFile(path string) error {
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
				fmt.Print(textLine)
			}

			break
		}

		if err != nil {
			return fmt.Errorf("error reading from file: %w", err)
		}

		fmt.Print(textLine)
	}

	return nil
}

func main() {
	path := "/home/nvt-dev/projects/1brc/measurements.txt"

	readLinesFromFile(path)
}