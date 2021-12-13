package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type DumboOctopus struct {
	Power      int
	HasFlashed bool
	Neighbours []*DumboOctopus
}

var (
	flashCount int
)

func init() {
	flashCount = 0
}

func main() {

	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))

	grid := CreateOctopusMap(input)
	var count, flashesAt100 int
	for ii := 0; !AllOctopoSynced(grid); ii++ {
		IncrementAndFlashOctopo(&grid)
		if ii == 99 {
			flashesAt100 = flashCount
		}
		// if ii > 275 {
		// 	PrintGrid(grid)
		// }
		count = ii
	}
	fmt.Printf("The total number of flashes at count 100 is is %d\n", flashesAt100)
	fmt.Printf("The octopo sync at count %d\n", count+1)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func PrintGrid(grid [][]*DumboOctopus) {
	for _, row := range grid {
		for _, o := range row {
			fmt.Printf("%d", o.Power)
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func AllOctopoSynced(grid [][]*DumboOctopus) bool {

	for _, row := range grid {
		for _, o := range row {
			if o.Power != 0 {
				return false
			}
		}
	}
	return true
}

func IncrementAndFlashOctopo(grid *[][]*DumboOctopus) {

	for _, row := range *grid {
		for _, o := range row {
			o.ClearFlashStatus()
			o.IncrementPower()
		}
	}

	// PrintGrid(*grid)

	for _, row := range *grid {
		for _, o := range row {
			o.FlashOctopusIfNeeded()
		}
	}
	// PrintGrid(*grid)
}

func (o *DumboOctopus) ClearFlashStatus() {
	o.HasFlashed = false
}

func (o *DumboOctopus) IncrementPower() {
	if !o.HasFlashed && o.Power < 10 {
		o.Power++
	}
}

func (o *DumboOctopus) FlashOctopusIfNeeded() {
	if o.Power == 10 {
		o.HasFlashed = true
		o.Power = 0
		flashCount++
		for _, v := range o.Neighbours {
			v.IncrementPower()
			v.FlashOctopusIfNeeded()
		}
	}
}

func CreateOctopusMap(input []string) [][]*DumboOctopus {
	octopoGrid := make([][]*DumboOctopus, len(input))
	for ii, row := range input {
		octopoRow := make([]*DumboOctopus, len(row))
		for jj, v := range row {
			octopoRow[jj] = &DumboOctopus{Power: int(v - '0'), HasFlashed: false, Neighbours: make([]*DumboOctopus, 0)}
		}
		octopoGrid[ii] = octopoRow
	}

	UpdateNeighbours(octopoGrid)
	return octopoGrid
}

func UpdateNeighbours(grid [][]*DumboOctopus) {
	for ii, row := range grid {
		for jj, oct := range row {

			var yMin, yLimit, xMin, xLimit int
			if ii > 0 {
				yMin = ii - 1
			} else {
				yMin = ii
			}

			if ii < len(grid)-1 {
				yLimit = ii + 1
			} else {
				yLimit = ii
			}

			if jj > 0 {
				xMin = jj - 1
			} else {
				xMin = jj
			}

			if jj < len(row)-1 {
				xLimit = jj + 1
			} else {
				xLimit = jj
			}

			for kk := yMin; kk <= yLimit; kk++ {
				for ll := xMin; ll <= xLimit; ll++ {
					if kk == ii && ll == jj {
						continue
					}
					oct.Neighbours = append(oct.Neighbours, grid[kk][ll])
				}
			}
		}
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
