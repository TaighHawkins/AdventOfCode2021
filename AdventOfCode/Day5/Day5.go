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
	stringArray := strings.FieldsFunc(string(rawinput), SplitValues)
	plot := PlotGrid{}
	plot.InitialiseGrid(stringArray)
	plot.PlotPoints(VerticalAndHorizontalPair)
	simpleIntersections := plot.CountIntersections()
	// plot.PrintGrid()
	fmt.Printf("The number of intersections with only horizontal or vertical lines is %d\n", simpleIntersections)

	//plot.ClearGrid()
	plot.PlotPoints(DiagonalPair)

	allIntersections := plot.CountIntersections()
	// plot.PrintGrid()
	fmt.Printf("The total number of intersections is %d\n", allIntersections)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s", timeTaken)
	fmt.Scanf("h")
}

func VerticalAndHorizontalPair(p PlotPair) bool {
	return p.X1 == p.X2 || p.Y1 == p.Y2
}
func DiagonalPair(p PlotPair) bool {
	return !VerticalAndHorizontalPair(p)
}

func ConvertLineIntoPlotPair(input string) PlotPair {
	var x1, y1, x2, y2 int
	_, err := fmt.Sscanf(input, "%d,%d -> %d,%d", &x1, &y1, &x2, &y2)
	Check(err)
	var pair PlotPair
	if x1 <= x2 {
		pair = PlotPair{x1, x2, y1, y2}
	} else {
		pair = PlotPair{x2, x1, y2, y1}
	}
	return pair
}

func SplitStringIntoCoords(input string) (int, int) {
	split := strings.Split(input, ",")
	x, xerror := strconv.Atoi(split[0])
	y, yerror := strconv.Atoi(split[1])
	Check(xerror)
	Check(yerror)
	return x, y
}

type PlotGrid struct {
	Grid  [][]int
	Plots []PlotPair
	XMin  int
	XMax  int
	YMin  int
	YMax  int
}

type PlotPair struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

func (p *PlotGrid) InitialiseGrid(input []string) {
	p.Plots = make([]PlotPair, 0, len(input))
	for _, v := range input {
		pair := ConvertLineIntoPlotPair(v)
		p.UpdateMinAndMax(pair)
		p.Plots = append(p.Plots, pair)
	}

	p.Grid = make([][]int, p.XMax+1)
	for ii := range p.Grid {
		p.Grid[ii] = make([]int, p.YMax+1)
	}
}

func (p *PlotGrid) PlotPoints(condition func(PlotPair) bool) {
	for _, v := range p.Plots {
		if condition(v) {
			p.AddToGrid(v)
		}
	}
}

func (p *PlotGrid) PrintPoints() {
	for _, point := range p.Plots {
		fmt.Printf("%d , %d -> %d , %d\n", point.X1, point.Y1, point.X2, point.Y2)
	}
}

// Ill-advised to use this on the provided input, suitable for the test input though
func (p *PlotGrid) PrintGrid() {
	for _, row := range p.Grid {
		fmt.Println()
		for _, x := range row {
			fmt.Printf("%d ", x)
		}
	}
	fmt.Println("")
	fmt.Println("")
}

func (p *PlotGrid) CountIntersections() int {
	rv := 0
	for _, row := range p.Grid {
		for _, point := range row {
			if point > 1 {
				rv++
			}
		}
	}
	return rv
}

func (p *PlotGrid) UpdateMinAndMax(pair PlotPair) {
	if p.XMax < pair.X2 {
		p.XMax = pair.X2
	}
	if p.YMax < pair.MaxY() {
		p.YMax = pair.MaxY()
	}
}

func (p *PlotGrid) AddToGrid(plot PlotPair) {
	for _, point := range plot.RetrieveLineCoords() {
		p.Grid[point[0]][point[1]]++
	}
}

func (p *PlotGrid) ClearGrid() {
	for ii, row := range p.Grid {
		for jj := range row {
			p.Grid[ii][jj] = 0
		}
	}
}

func (p *PlotPair) RetrieveLineCoords() [][2]int {
	results := make([][2]int, 0)

	if p.X1 == p.X2 {
		for ii := p.MinY(); ii <= p.MaxY(); ii++ {
			results = append(results, [2]int{p.X1, ii})
		}
	} else if p.Y1 == p.Y2 {
		for ii := p.X1; ii <= p.X2; ii++ {
			results = append(results, [2]int{ii, p.Y1})
		}
	} else {
		increment := p.GetSlope()
		yStart := p.Y1
		for ii := p.X1; ii <= p.X2; ii++ {
			results = append(results, [2]int{ii, yStart})
			yStart += increment
		}
	}
	return results
}

func (p *PlotPair) GetSlope() int {
	if p.Y1 < p.Y2 {
		return 1
	} else {
		return -1
	}
}

func (p *PlotPair) MaxY() int {
	return GetMax(p.Y1, p.Y2)
}

func (p *PlotPair) MinY() int {
	return GetMin(p.Y1, p.Y2)
}

func GetMin(v1 int, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func GetMax(v1 int, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func SplitValues(c rune) bool {
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
