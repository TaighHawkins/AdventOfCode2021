package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))

	hm := ExtractHeightMap(input)
	riskSum, lowPoints := hm.ExtractRiskLevelSum()

	fmt.Printf("The sum of the risk levels is %d\n", riskSum)

	hm.IdentifyBasins(lowPoints)
	sizes := hm.GetAllBasinSizes()
	largest := GetThreeLargest(sizes)
	prod := ProductOf(largest[:])
	fmt.Printf("The product of the largest three basins is %d\n", prod)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

type HeightMap struct {
	Grid   [][]int
	Basins []*Basin
}

type Basin struct {
	Coords [][2]int
}

func ExtractHeightMap(input []string) HeightMap {
	grid := make([][]int, len(input))
	for ii, row := range input {
		gridRow := make([]int, len(row))
		runes := []rune(row)
		for jj, char := range runes {
			gridRow[jj] = int(char - '0')
		}
		grid[ii] = gridRow
	}
	return HeightMap{grid, []*Basin{}}
}

func ProductOf(input []int) int {
	rv := 1
	for _, v := range input {
		rv *= v
	}
	return rv
}

func GetThreeLargest(input []int) [3]int {
	highest := [3]int{0, 0, 0}
	for _, v := range input {
		ix, lowest := GetLowest(highest)
		if v > lowest {
			highest[ix] = v
		}
	}
	return highest
}

func GetLowest(input [3]int) (int, int) {
	lowest := math.MaxInt
	ix := 0
	for ii, v := range input {
		if v < lowest {
			lowest = v
			ix = ii
		}
	}
	return ix, lowest
}

func IsHigherThanAValue(input int, array [3]int) (int, bool) {
	for ii, v := range array {
		if input > v {
			return ii, true
		}
	}
	return 0, false
}

func (h *HeightMap) ExtractRiskLevelSum() (int, [][2]int) {
	score := 0
	lowPoints := make([][2]int, 0)
	for ii, row := range h.Grid {
		for jj := range row {
			s := h.GetValueIfLowPoint(ii, jj)
			if s > 0 {
				score += s
				lowPoints = append(lowPoints, [2]int{ii, jj})
			}
		}
	}
	return score, lowPoints
}

func (h *HeightMap) GetValueIfLowPoint(y int, x int) int {
	chosen := h.Grid[y][x]
	if chosen == 9 {
		return 0
	}

	left := h.GetValueLeft(y, x)
	right := h.GetValueRight(y, x)
	up := h.GetValueAbove(y, x)
	down := h.GetValueBelow(y, x)
	if IsLowerThanAll(chosen, left, right, up, down) {
		// fmt.Printf("Low point found - score: %d\n", 1+chosen)
		// fmt.Printf(" %d \n%d%d%d\n %d\n", up, left, chosen, right, down)
		return 1 + chosen
	}
	return 0
}

func (h *HeightMap) IdentifyBasins(lowPoints [][2]int) {
	h.Basins = make([]*Basin, 0)
	var wg sync.WaitGroup
	for _, coords := range lowPoints {
		wg.Add(1)
		go h.AddToBasinIfAppropriate(coords[0], coords[1], &wg)
	}
	wg.Wait()
}

func (h *HeightMap) AddToBasinIfAppropriate(y int, x int, wg *sync.WaitGroup) {
	defer wg.Done()
	coords := make([][2]int, 0)
	newBasin := Basin{coords}
	newBasin.CrawlBasin(y, x, h.Grid)
	h.Basins = append(h.Basins, &newBasin)
}

func (h *HeightMap) GetAllBasinSizes() []int {
	sizes := make([]int, len(h.Basins))
	for ii, b := range h.Basins {
		sizes[ii] = len(b.Coords)
	}
	return sizes
}

func (b *Basin) CrawlBasin(y int, x int, grid [][]int) {
	if grid[y][x] == 9 {
		return
	}
	b.AddToBasin(y, x)
	if x > 0 && !b.IsAlreadyWithinBasin(y, x-1) {
		b.CrawlBasin(y, x-1, grid)
	}
	if x < len(grid[y])-1 && !b.IsAlreadyWithinBasin(y, x+1) {
		b.CrawlBasin(y, x+1, grid)
	}
	if y > 0 && !b.IsAlreadyWithinBasin(y-1, x) {
		b.CrawlBasin(y-1, x, grid)
	}
	if y < len(grid)-1 && !b.IsAlreadyWithinBasin(y+1, x) {
		b.CrawlBasin(y+1, x, grid)
	}
}

func (b *Basin) IsAlreadyWithinBasin(y, x int) bool {
	for _, v := range b.Coords {
		if v[0] == y && v[1] == x {
			return true
		}
	}
	return false
}

func (b *Basin) IsAdjacent(y, x int) bool {
	for _, v := range b.Coords {
		if IsBelow(v, y, x) || IsRight(v, y, x) {
			return true
		}
	}
	return false
}

func (b *Basin) AddToBasin(y, x int) {
	b.Coords = append(b.Coords, [2]int{y, x})
}

func IsRight(coords [2]int, y int, x int) bool {
	if x == 0 {
		return false
	} else if coords[1]+1 == x && coords[0] == y {
		return true
	}
	return false
}

func IsBelow(coords [2]int, y int, x int) bool {
	if y == 0 {
		return false
	} else if coords[0]+1 == y && coords[1] == x {
		return true
	}
	return false
}

func IsLowerThanAll(x int, compares ...int) bool {
	if len(compares) == 0 {
		return false
	} else {
		for _, i := range compares {
			if x >= i {
				return false
			}
		}
	}
	return true
}

func (g *HeightMap) GetValueAbove(y, x int) int {
	if y == 0 {
		return math.MaxInt
	} else {
		return g.Grid[y-1][x]
	}
}
func (g *HeightMap) GetValueBelow(y, x int) int {
	if y == len(g.Grid)-1 {
		return math.MaxInt
	} else {
		return g.Grid[y+1][x]
	}
}
func (g *HeightMap) GetValueLeft(y, x int) int {
	if x == 0 {
		return math.MaxInt
	} else {
		return g.Grid[y][x-1]
	}
}
func (g *HeightMap) GetValueRight(y, x int) int {
	if x == len(g.Grid[y])-1 {
		return math.MaxInt
	} else {
		return g.Grid[y][x+1]
	}
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
