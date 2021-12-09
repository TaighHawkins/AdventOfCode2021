package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	subPositions, min, max := ExtractSubLocations(input)
	sort.Ints(subPositions)

	binaryCostChannel := make(chan int)
	binarylocationChannel := make(chan int)
	defer close(binaryCostChannel)
	defer close(binarylocationChannel)
	go BinarySearch(subPositions, min, max, 0, 0, binaryCostChannel, binarylocationChannel)

	// Migrating crabs with linear movement is minimised through by taking the median point
	usablelength := len(subPositions)
	if usablelength%2 != 0 {
		usablelength += 1
	}

	medianIx := usablelength / 2
	fmt.Printf("The median is %d and it takes %d fuel to get there \n",
		subPositions[medianIx],
		MoveAllSubsTo(subPositions, subPositions[medianIx], SimpleCost),
	)
	fmt.Printf("The binary search turned up: cost is %d and it will cost %d fuel across all subs\n",
		<-binaryCostChannel,
		<-binarylocationChannel)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func ExtractSubLocations(input []string) ([]int, int, int) {
	min := math.MaxInt
	max := math.MinInt
	data := make([]int, len(input))
	for ii, v := range input {
		intValue, err := strconv.Atoi(v)
		Check(err)
		data[ii] = intValue
		if intValue < min {
			min = intValue
		}
		if intValue > max {
			max = intValue
		}
	}
	return data, min, max
}

func BinarySearch(pos []int, min int, max int, minCost int, maxCost int, costChannel chan int, locationChannel chan int) {

	if minCost == 0 {
		minCost = MoveAllSubsTo(pos, min, ExpensiveCost)
	}
	if maxCost == 0 {
		maxCost = MoveAllSubsTo(pos, max, ExpensiveCost)
	}

	diff := max - min
	if diff == 1 {
		if minCost < maxCost {
			costChannel <- minCost
			locationChannel <- min
		} else {
			costChannel <- maxCost
			locationChannel <- max
		}
		return
	}

	if diff%2 != 0 {
		diff += 1
	}
	mid := min + (diff / 2)
	midCost := MoveAllSubsTo(pos, mid, ExpensiveCost)

	minMidDiff := float64(minCost - midCost)
	maxMidDiff := float64(maxCost - midCost)
	if minMidDiff > maxMidDiff {
		BinarySearch(pos, mid, max, midCost, maxCost, costChannel, locationChannel)
	} else {
		BinarySearch(pos, min, mid, minCost, midCost, costChannel, locationChannel)
	}
}

func MoveAllSubsTo(subPositions []int, position int, cost func(float64) float64) int {
	totalCost := float64(0)
	for _, v := range subPositions {
		totalCost += cost(float64(v - position))
	}
	return int(totalCost)
}

func SimpleCost(distance float64) float64 {
	return math.Abs(distance)
}

func ExpensiveCost(distance float64) float64 {
	abs := math.Abs(distance)
	return (abs * (abs + 1)) / 2
}

func GetStringFromInput(inputPath string) []string {

	rawinput, fileError := os.ReadFile(inputPath)
	Check(fileError)
	return strings.Split(strings.TrimSpace(string(rawinput)), ",")
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
