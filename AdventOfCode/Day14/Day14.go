package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

var (
	InsertionPairs map[string]rune
	InsertMap      map[string]int
	PolymerMap     map[rune]int
)

func init() {
	InsertionPairs = make(map[string]rune)
	InsertMap = make(map[string]int)
	PolymerMap = make(map[rune]int)
}

func main() {
	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	ExtractPolymerAndInsertionPairs(input)

	for ii := 0; ii < 10; ii++ {
		RunPairInsertion()
	}

	fmt.Printf("The difference between the most and least common elements after 10 steps is is %d\n", DetermineDifferenceBetweenMostAndLeastCommon())

	for ii := 0; ii < 30; ii++ {
		RunPairInsertion()
	}

	fmt.Printf("The difference between the most and least common elements after 40 steps is is %d\n", DetermineDifferenceBetweenMostAndLeastCommon())

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func DetermineDifferenceBetweenMostAndLeastCommon() int {
	max, min := 0, math.MaxInt

	for _, v := range PolymerMap {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max - min
}

func RunPairInsertion() {
	newMap := make(map[string]int)
	for k := range InsertMap {
		newMap[k] = 0
	}
	for k, v := range InsertMap {
		if v == 0 {
			continue
		}
		replacement := InsertionPairs[k]
		runes := []rune(k)
		s1 := string([]rune{runes[0], replacement})
		s2 := string([]rune{replacement, runes[1]})
		newMap[s1] += v
		newMap[s2] += v
		PolymerMap[replacement] += v
	}
	InsertMap = newMap
}

func SetInitialCounts(runes ...rune) {
	for ii, r := range runes {
		PolymerMap[r]++
		if ii < len(runes)-1 {
			UpdateInsertMap(r, runes[ii+1])
		}
	}
}

func UpdateInsertMap(r1, r2 rune) {
	s := string([]rune{r1, r2})
	InsertMap[s]++
}

func ExtractPolymerAndInsertionPairs(input []string) []rune {
	polymerTemplate := []rune(strings.TrimSpace(input[0]))

	for _, v := range strings.Split(input[1], "\n") {
		var key, insert string
		fmt.Sscanf(v, "%s -> %s", &key, &insert)
		char := []rune(insert)[0]
		InsertionPairs[key] = char
		InsertMap[key] = 0
		PolymerMap[char] = 0
	}
	SetInitialCounts(polymerTemplate...)

	return polymerTemplate
}

func GetStringFromInput(inputPath string) []string {

	rawinput, fileError := os.ReadFile(inputPath)
	Check(fileError)
	return strings.Split(strings.TrimSpace(string(rawinput)), "\n\n")
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
