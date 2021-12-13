package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ()

func init() {
}

type Point struct {
	X int
	Y int
}
type Points struct {
	Points []*Point
	MaxX   int
	MaxY   int
}

func (p Points) Swap(i, j int) {
	p.Points[i], p.Points[j] = p.Points[j], p.Points[i]
}

func (p Points) Len() int {
	return len(p.Points)
}

func (p Points) Less(i, j int) bool {
	return p.Points[i].X < p.Points[j].X || (p.Points[i].X == p.Points[j].X && p.Points[i].Y < p.Points[j].Y)
}

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	points, instructions := ExtractCoordsAndInstructions(input)

	instructions = points.CompleteInstruction(instructions)

	fmt.Printf("Number of visible dots after one instruction is %d\n", len(points.Points))

	for len(instructions) > 0 {
		instructions = points.CompleteInstruction(instructions)
	}

	points.PrintCoords()

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func (p *Points) PrintCoords() {
	sort.Sort(p)

	grid := make([][]string, p.MaxY+1)
	for ii := range grid {
		grid[ii] = make([]string, p.MaxX+1)
	}

	for _, p := range p.Points {
		grid[p.Y][p.X] = "▓"
	}

	for _, row := range grid {
		for _, v := range row {
			if v == "" {
				fmt.Print(" ")
			} else {
				fmt.Print(v)
			}
		}
		fmt.Println("")
	}
}

func (p *Points) CompleteInstruction(instructions []string) []string {
	axis, foldLine := ExtractSpecifics(instructions[0])
	points := p.Points
	for ii := range points {
		if axis == "y" {
			p.MaxY = foldLine
			if points[ii].Y > foldLine {
				points[ii].Y = (2 * foldLine) - points[ii].Y
			}
		} else if axis == "x" {
			p.MaxX = foldLine
			if points[ii].X > foldLine {
				points[ii].X = (2 * foldLine) - points[ii].X
			}
		}
	}
	p.RemoveDuplicate()

	return instructions[1:]
}

func ExtractSpecifics(s string) (string, int) {
	var axis string
	var coord int
	fmt.Sscanf(strings.Replace(s, "=", " ", -1), "fold along %s %v", &axis, &coord)

	return axis, coord
}

func ExtractCoordsAndInstructions(input []string) (Points, []string) {
	split := strings.Split(input[0], "\n")
	points := Points{Points: make([]*Point, len(split))}
	for ii, v := range split {
		xy := strings.Split(v, ",")
		points.Points[ii] = &Point{GetInt(xy[0]), GetInt(xy[1])}
	}

	return points, strings.Split(input[1], "\n")
}

func (p *Points) RemoveDuplicate() {
	allKeys := make(map[Point]bool)
	list := make([]*Point, 0)
	for _, item := range p.Points {
		if _, value := allKeys[*item]; !value {
			allKeys[*item] = true
			list = append(list, item)
		}
	}
	p.Points = list
}

func GetInt(s string) int {
	intValue, err := strconv.Atoi(s)
	Check(err)
	return intValue
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
