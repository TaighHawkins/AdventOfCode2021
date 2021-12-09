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

func ReduceToFinalEntryBasedOnCriteria(input []string, criteria func([]string, []int, int) []string) string {
	workingCopy := input
	for ii := 0; len(workingCopy) > 1; ii++ {
		counts, _ := GetCountsOfOnes(workingCopy, ii, ii)
		workingCopy = criteria(workingCopy, counts, ii)
	}
	return workingCopy[0]
}

func ExtractMostCommonValues(input []string, counts []int, indexForComparison int) []string {
	mostCommon, _ := ExtractMostCommonAndLeastCommonValues(input, counts, indexForComparison)
	return mostCommon
}

func ExtractLeastCommonValues(input []string, counts []int, indexForComparison int) []string {
	_, leastCommon := ExtractMostCommonAndLeastCommonValues(input, counts, indexForComparison)
	return leastCommon
}

func ExtractMostCommonAndLeastCommonValues(input []string, counts []int, indexForComparison int) ([]string, []string) {
	var mostCommon []string
	var leastCommon []string
	halfInput := float64(len(input)) / 2

	for _, v := range input {
		for ii, c := range v {
			if ii == indexForComparison {
				if c == GetMostCommonChar(counts[indexForComparison], halfInput) {
					mostCommon = append(mostCommon, v)
				} else {
					leastCommon = append(leastCommon, v)
				}
			}
		}
	}
	return mostCommon, leastCommon
}

func GetMostCommonChar(oneCount int, halfInputCount float64) rune {
	if float64(oneCount) >= halfInputCount {
		return '1'
	} else {
		return '0'
	}
}

func GetAllOneCounts(input []string) ([]int, int) {
	return GetCountsOfOnes(input, 0, len(input)-1)
}

func GetCountsOfOnes(input []string, minIndex int, maxIndex int) ([]int, int) {
	count := make([]int, len(input[0]))
	halfInput := len(input) / 2
	for _, v := range input {
		for jj, c := range v {
			if jj >= minIndex && jj <= maxIndex {
				if c == '1' {
					count[jj] += 1
				}
			}
		}
	}
	return count, halfInput
}

func ExtractBits(counts []int, halfwayMarker int) (string, string) {
	mostCount := ""
	leastCount := ""

	for _, v := range counts {
		if v > halfwayMarker {
			mostCount += "1"
			leastCount += "0"
		} else {
			leastCount += "1"
			mostCount += "0"
		}
	}
	return mostCount, leastCount
}

func ConvertStringBitsToInt(input string) int64 {
	rv, err := strconv.ParseInt(input, 2, 64)
	Check(err)
	return rv
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
