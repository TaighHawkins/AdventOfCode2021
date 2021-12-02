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
	rawinput, fileError := os.ReadFile(DereferenceStringPointer(inputName))
	Check(fileError)

	input, conversionError := ConvertInputToIntArray(strings.FieldsFunc(string(rawinput), RemoveEmptyValues))
	Check(conversionError)

	increases, countError := CountSequentialIncreases(input)
	Check(countError)
	fmt.Printf("There were %d direct increases\n", increases)

	summableIncreases, sumError := CountSummableIncreases(input, 3)
	Check(sumError)
	fmt.Printf("There were %d increases over a rolling period\n", summableIncreases)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s", timeTaken)
	fmt.Scanf("h")
}

func DereferenceStringPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func ConvertInputToIntArray(input []string) ([]int, error) {
	var intArray []int
	for _, v := range input {
		intValue, convError := strconv.Atoi(strings.TrimSpace(v))
		Check(convError)
		intArray = append(intArray, intValue)
	}
	return intArray, nil
}

func RemoveEmptyValues(c rune) bool {
	return c == '\n'
}

func CountSequentialIncreases(input []int) (int, error) {
	return CountSummableIncreases(input, 1)
}

func CountSummableIncreases(input []int, step int) (int, error) {

	var currentErr error
	currentSum := 0
	previousSum := 0

	increaseCount := 0
	for ii := range input {
		previousSum = currentSum
		currentSum, currentErr = SumSlice(input[ii : ii+step])
		Check(currentErr)

		if ii == 0 {
			continue
		} else if ii+step > len(input) {
			break
		} else if currentSum > previousSum {
			increaseCount++
		}
	}
	return increaseCount, nil
}

func SumSlice(input []int) (int, error) {
	sum := 0
	for _, v := range input {
		sum = sum + v
	}
	return sum, nil
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
