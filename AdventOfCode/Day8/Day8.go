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

	patterns, fourdigitOutputs := ExtractPatternsAndOutputs(input)

	fmt.Printf("The digits 1, 4, 7 and 8 appear a total of %d times\n", CountUniqueDigits(fourdigitOutputs))

	ch := make(chan int)
	defer close(ch)

	for ii := 0; ii < len(patterns); ii++ {
		go ConvertSignalsAndInterpretValue(patterns[ii], fourdigitOutputs[ii], ch)
	}

	summed := 0
	for range patterns {
		summed += <-ch
	}
	fmt.Printf("The sum of all the digits is %d\n", summed)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func CountUniqueDigits(input [][4]string) int {
	count := 0
	for _, row := range input {
		for _, v := range row {
			length := len(v)
			if length == 2 || length == 3 || length == 4 || length == 7 {
				count++
			}
		}
	}
	return count
}

func ExtractPatternsAndOutputs(input []string) ([][10]string, [][4]string) {

	signalPatterns := make([][10]string, len(input))
	fourdigitOutputs := make([][4]string, len(input))

	for ii, v := range input {
		pairs := strings.Split(v, " | ")
		_, patternError := fmt.Sscanf(strings.TrimSpace(pairs[0]), "%s %s %s %s %s %s %s %s %s %s",
			&signalPatterns[ii][0],
			&signalPatterns[ii][1],
			&signalPatterns[ii][2],
			&signalPatterns[ii][3],
			&signalPatterns[ii][4],
			&signalPatterns[ii][5],
			&signalPatterns[ii][6],
			&signalPatterns[ii][7],
			&signalPatterns[ii][8],
			&signalPatterns[ii][9],
		)
		Check(patternError)
		_, outputError := fmt.Sscanf(strings.TrimSpace(pairs[1]), "%s %s %s %s",
			&fourdigitOutputs[ii][0],
			&fourdigitOutputs[ii][1],
			&fourdigitOutputs[ii][2],
			&fourdigitOutputs[ii][3],
		)
		Check(outputError)
	}
	return signalPatterns, fourdigitOutputs
}

func ConvertSignalsAndInterpretValue(signals [10]string, output [4]string, ch chan int) {
	dict := DecodeSignals(signals)
	number := ConvertOutputUsingMap(output, dict)
	ch <- number
}

func DecodeSignals(input [10]string) map[string]string {
	working := input[:]
	dict := make(map[string]string, 10)
	working, dict["1"] = ExtractValue(working, IsOne)
	working, dict["4"] = ExtractValue(working, IsFour)
	working, dict["7"] = ExtractValue(working, IsSeven)
	working, dict["8"] = ExtractValue(working, IsEight)
	working, dict["9"] = ExtractNine(working, dict["1"], dict["4"])
	working, dict["3"] = ExtractThree(working, dict["1"], dict["4"])
	working, dict["5"] = ExtractFive(working, dict["1"], dict["4"])
	working, dict["6"] = ExtractSix(working, dict["1"], dict["4"])
	working, dict["0"] = ExtractZero(working, dict["1"], dict["4"])
	dict["2"] = ExtractTwo(working, dict["1"], dict["4"])
	return dict
}

func ConvertOutputUsingMap(input [4]string, dict map[string]string) int {
	numberString := ""
	for _, v := range input {
		numberString += ExtractValueFromDict(dict, v)
	}
	rv, err := strconv.Atoi(numberString)
	Check(err)
	return rv
}

func ExtractValueFromDict(dict map[string]string, s string) string {
	for k, v := range dict {
		if AllRunesMatch(s, v) {
			return k
		}
	}
	panic("no matching key")
}

func AllRunesMatch(i, s string) bool {
	return len(i) == len(s) && AllCharactersArePresent(i, s)
}

func ExtractNine(input []string, one string, four string) ([]string, string) {
	for ii, v := range input {
		if len(v) == 6 && AllCharactersArePresent(v, four) && AllCharactersArePresent(v, four) {
			inputWithRemoved := RemoveElementAtIndex(input, ii)
			return inputWithRemoved, v
		}
	}
	panic("Nine not found in signal")
}

func ExtractSix(input []string, one string, four string) ([]string, string) {
	for ii, v := range input {
		if len(v) == 6 && AllExceptOneCharIsPresent(v, one) && AllExceptOneCharIsPresent(v, four) {
			inputWithRemoved := RemoveElementAtIndex(input, ii)
			return inputWithRemoved, v
		}
	}
	panic("Six not found in signal")
}

func ExtractZero(input []string, one string, four string) ([]string, string) {
	for ii, v := range input {
		if len(v) == 6 && AllCharactersArePresent(v, one) && AllExceptOneCharIsPresent(v, four) {
			inputWithRemoved := RemoveElementAtIndex(input, ii)
			return inputWithRemoved, v
		}
	}
	panic("Zero not found in signal")
}

func ExtractThree(input []string, one string, four string) ([]string, string) {
	for ii, v := range input {
		if len(v) == 5 && AllCharactersArePresent(v, one) && AllExceptOneCharIsPresent(v, four) {
			inputWithRemoved := RemoveElementAtIndex(input, ii)

			return inputWithRemoved, v
		}
	}
	panic("Three not found in signal")
}

func ExtractFive(input []string, one string, four string) ([]string, string) {
	for ii, v := range input {
		if len(v) == 5 && AllExceptOneCharIsPresent(v, one) && AllExceptOneCharIsPresent(v, four) {
			inputWithRemoved := RemoveElementAtIndex(input, ii)

			return inputWithRemoved, v
		}
	}
	panic("Five not found in signal")
}

func ExtractTwo(input []string, one string, four string) string {
	for _, v := range input {
		if len(v) == 5 && AllExceptOneCharIsPresent(v, one) && HasNCharactersInCommon(v, four, 2) {
			return v
		}
	}
	panic("Two not found in signal")
}

func HasNCharactersInCommon(input string, comparison string, number int) bool {
	inputArray := []rune(input)
	compareArray := []rune(comparison)
	count := 0
	for _, v := range inputArray {
		if count > number {
			return false
		}

		if CharInArray(v, compareArray) {
			count++
		}
	}
	return count == number
}

func AllCharactersArePresent(input string, comparer string) bool {
	return HasNCharactersInCommon(input, comparer, len(comparer))
}

func AllExceptOneCharIsPresent(input string, comparer string) bool {
	return HasNCharactersInCommon(input, comparer, len(comparer)-1)
}

func CharInArray(a rune, list []rune) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ExtractValue(input []string, condition func(string) bool) ([]string, string) {
	for ii, v := range input {
		if condition(v) {
			inputWithRemoved := RemoveElementAtIndex(input, ii)
			return inputWithRemoved, v
		}
	}
	panic("value not found in signal")
}

func IsOne(input string) bool {
	return len(input) == 2
}
func IsSeven(input string) bool {
	return len(input) == 3
}
func IsFour(input string) bool {
	return len(input) == 4
}
func IsEight(input string) bool {
	return len(input) == 7
}

func RemoveElementAtIndex(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func GetStringFromInput(inputPath string) []string {

	rawinput, fileError := os.ReadFile(inputPath)
	Check(fileError)
	return strings.Split(strings.TrimSpace(string(rawinput)), "\n")
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
