package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	syntaxScoreMap     = make(map[rune]int)
	completionScoreMap = make(map[rune]int)
	pairMap            = make(map[rune]rune)
)

func init() {
	syntaxScoreMap[')'] = 3
	syntaxScoreMap[']'] = 57
	syntaxScoreMap['}'] = 1197
	syntaxScoreMap['>'] = 25137

	completionScoreMap[')'] = 1
	completionScoreMap[']'] = 2
	completionScoreMap['}'] = 3
	completionScoreMap['>'] = 4

	pairMap['('] = ')'
	pairMap['['] = ']'
	pairMap['{'] = '}'
	pairMap['<'] = '>'
}

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))

	syntaxErrors, cleanLines := ExtractSyntaxErrors(input)
	score := ScoreErrors(syntaxErrors)
	fmt.Printf("The score for the syntax errors in the input is %d\n", score)

	ch := make(chan string)
	defer close(ch)
	completionLines := CompleteAppropriateLines(cleanLines, ch)
	completionScore := ScoreCompletions(completionLines)
	fmt.Printf("The score for the completions in the input is %d\n", completionScore)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func ScoreErrors(i []rune) int {
	score := 0
	for _, r := range i {
		score += syntaxScoreMap[r]
	}
	return score
}

func ExtractSyntaxErrors(input []string) ([]rune, []string) {
	errors := make([]rune, 0)
	cleanLines := make([]string, 0)
	var wasError bool
	for _, line := range input {
		errors, wasError = UpdateErrorWithSyntaxErrorIfNeeded(errors, line)
		if !wasError {
			cleanLines = append(cleanLines, line)
		}
	}
	return errors, cleanLines
}

func UpdateErrorWithSyntaxErrorIfNeeded(errors []rune, s string) ([]rune, bool) {
	awaitingClose := make([]rune, 0)
	error := false
	for _, c := range s {
		if c == '(' || c == '{' || c == '[' || c == '<' {
			awaitingClose = append(awaitingClose, c)
		} else if IsClosingBrace(awaitingClose[len(awaitingClose)-1], c) {
			awaitingClose = awaitingClose[:len(awaitingClose)-1]
		} else {
			errors = append(errors, c)
			error = true
			break
		}
	}
	return errors, error
}

func CompleteAppropriateLines(input []string, ch chan string) []string {
	completions := make([]string, 0)
	for _, line := range input {
		go UpdateCompletionStrings(line, ch)
	}
	for range input {
		completions = append(completions, <-ch)
	}
	return completions
}

func ScoreCompletions(i []string) int {
	scores := make([]int, len(i))

	for ii, s := range i {
		score := 0
		for _, r := range s {
			score *= 5
			score += completionScoreMap[r]
		}
		scores[ii] = score
	}

	sort.Ints(scores)
	return scores[(len(scores)-1)/2]
}

func UpdateCompletionStrings(s string, ch chan string) {
	awaitingClose := make([]rune, 0)
	completionString := make([]rune, 0)
	for _, c := range s {
		if IsOpeningBrace(c) {
			awaitingClose = append(awaitingClose, c)
		} else if IsClosingBrace(awaitingClose[len(awaitingClose)-1], c) {
			awaitingClose = awaitingClose[:len(awaitingClose)-1]
		} else {
			panic("this string has an error")
		}
	}

	for ii := len(awaitingClose) - 1; ii >= 0; ii-- {
		completionString = append(completionString, GetClosingBrace(awaitingClose[ii]))
	}

	ch <- string(completionString)
}

func IsOpeningBrace(r rune) bool {
	_, ok := GetClosingPair(r)
	return ok
}
func GetClosingBrace(r rune) rune {
	v, ok := GetClosingPair(r)
	if ok {
		return v
	}
	panic("no closing pair available")
}

func IsClosingBrace(r, x rune) bool {
	v, _ := GetClosingPair(r)
	return v == x
}

func GetClosingPair(r rune) (rune, bool) {
	v, ok := pairMap[r]
	return v, ok
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
