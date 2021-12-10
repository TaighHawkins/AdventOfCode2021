package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	data := CountFish(input)

	for ii := 0; ii < 256; ii++ {
		if ii == 80 {
			// Let this run in the background
			go fmt.Printf("There are %d Lanternfish after 80 days\n", SumArray(data))
		}
		TickFish(&data)
	}

	fmt.Printf("There are %d Lanternfish after 256 days\n", SumArray(data))

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func CountFish(input []string) [9]int {
	data := [9]int{}
	for _, v := range input {
		intValue, err := strconv.Atoi(v)
		Check(err)
		data[intValue]++
	}
	return data
}

func TickFish(data *[9]int) {
	countOfZeros := data[0]
	for ii := 0; ii < 8; ii++ {
		data[ii] = data[ii+1]
	}
	data[6] += countOfZeros
	data[8] = countOfZeros
}

func SumArray(input [9]int) int {
	rv := 0
	for _, v := range input {
		rv += v
	}
	return rv
}

func GetStringFromInput(inputPath string) []string {

	rawinput, fileError := os.ReadFile(inputPath)
	Check(fileError)
	return strings.Split(strings.TrimSpace(string(rawinput)), ",")
}

func RemoveEmptyValues(c rune) bool {
	return c == '\n'
}

func DereferenceStringPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
