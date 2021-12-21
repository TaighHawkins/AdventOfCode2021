package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

var (
	Caves      map[string]Cave
	Solutions2 [][]Cave
	Start      time.Time
)

type Cave struct {
	Name       string
	Big        bool
	VisitCount int
	Links      []string
}

func init() {
	Caves = make(map[string]Cave)
	Solutions2 = make([][]Cave, 0)
	Start = time.Now()
}

// https://www.reddit.com/r/adventofcode/comments/rehj2r/comment/hob126g/?utm_source=share&utm_medium=web2x&context=3
// Hints for optimisation

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	ReadMazeMap(input)
	TraverseCaves([]Cave{Caves["start"]}, map[string]Cave{"start": Caves["start"]}, IsNewPath)
	fmt.Printf("The number of paths through the maze is %d\nCalculated in %s\n", len(Solutions2), time.Since(Start))

	Solutions2 = make([][]Cave, 0)

	TraverseCaves([]Cave{Caves["start"]}, map[string]Cave{"start": Caves["start"]}, IsNewCaveOrOnlyOneSmallCaveHasASecondPass)
	fmt.Printf("The number of paths through the maze allowing small cave revisiting is %d\nCalculated in %s\n", len(Solutions2), time.Since(Start))

	timeTaken := time.Since(Start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func ReadMazeMap(input []string) {
	for _, v := range input {
		pair := strings.Split(v, "-")
		UpdateCave(pair[0], pair[1])
		UpdateCave(pair[1], pair[0])
	}
}

func UpdateCave(cave1, cave2 string) {
	_, onePresent := Caves[cave1]
	if onePresent {
		cave := Caves[cave1]
		cave.Links = append(cave.Links, cave2)
		Caves[cave1] = cave
	} else {
		Caves[cave1] = Cave{cave1, IsUpper(cave1), 0, []string{cave2}}
	}
}

func TraverseCaves(path []Cave, caveMap map[string]Cave, condition func([]Cave, map[string]Cave, string) bool) {
	if len(path) > 200 {
		fmt.Printf("Aborting run as depth of 200 hit:\n%v\n", path)
	}

	for _, v := range Caves[path[len(path)-1].Name].Links {
		if v == "end" {
			Solutions2 = append(Solutions2, append(path, Caves[v]))
			continue
		}

		if Caves[v].Big || condition(path, caveMap, v) {
			cave := Caves[v]
			caveLookup, present := caveMap[v]
			if present {
				caveLookup.VisitCount++
				caveMap[v] = caveLookup
			} else {
				caveMap[v] = cave
			}
			nextPath := append(path, cave)
			TraverseCaves(nextPath, caveMap, condition)
			if caveMap[v].VisitCount == 0 {
				delete(caveMap, v)
			} else {
				prevCave := caveMap[v]
				prevCave.VisitCount--
				caveMap[v] = prevCave
			}
		}
	}
}

func IsNewPath(p []Cave, cm map[string]Cave, s string) bool {
	_, present := cm[s]
	return !present
}

func IsNewCaveOrOnlyOneSmallCaveHasASecondPass(p []Cave, cm map[string]Cave, s string) bool {
	if s == "start" {
		return false
	}

	v, present := cm[s]

	if !present {
		return true
	} else if present && v.VisitCount >= 1 {
		return false
	} else if present {
		for _, value := range cm {
			if !value.Big && value.VisitCount == 1 {
				return false
			}
		}
	}

	return true
}

func IsNewPathOrIsLargeCave(path []string, s string) bool {
	return !IsAlreadyPresent(path, s) || IsUpper(s)
}

func IsAlreadyPresent(path []string, s string) bool {
	for _, v := range path {
		if s == v {
			return true
		}
	}
	return false
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
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
