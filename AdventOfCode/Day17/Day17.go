package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

type TargetZone struct {
	X MinMaxPair
	Y MinMaxPair
}

type MinMaxPair struct {
	Min int
	Max int
}

type Probe struct {
	Velocity *Velocity
	Location *Point
	MaxY     int
	Path     []Point
}

type Point struct {
	X int
	Y int
}

type Velocity struct {
	X int
	Y int
}

var ()

func init() {
}

func main() {
	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	fmt.Println("Beginning Trickshot Analysis")
	target := GetTargetFromInput(inputName)
	maxHeight, probeCount := target.LaunchProbesAndDetermineMaxYHeight()

	fmt.Printf("%d probes were launched and the max height reached by any probe that hit the target zone was: %dm\n", probeCount, maxHeight)
	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func (v *Velocity) Step() {
	v.Y -= 1

	if v.X > 0 {
		v.X -= 1
	} else if v.X < 0 {
		v.X += 1
	}
}

func (t TargetZone) LaunchProbesAndDetermineMaxYHeight() (int, int) {
	xClose, xFar := t.SortXByDistanceToZero()
	minX := GetMinimumXVelocity(xClose)
	fmt.Printf("Minimum x-component velocity is: %d\n", minX)
	// Due to directionality, after one step, being fired directly toward the lowest
	// y would enter the target
	minY := t.Y.Min
	maxX := xFar // same for x max
	yMax := 0
	if t.Y.Min < 0 {
		yMax = t.Y.Min * -1
	}
	fmt.Printf("Maximum x-component velocity is: %d\n", maxX)
	maxHeight := 0
	successfulProbes := 0
	validXVelocities, validYVelocities := make([]int, 0), make([]int, 0)
	probe := &Probe{}

	for ii := minX; ii <= maxX; ii++ {
		if t.HorizontalVelocityComponentEntersTargetZone(ii) {
			validXVelocities = append(validXVelocities, ii)
		}
	}

	for jj := minY; jj <= yMax; jj++ {
		if t.VerticalVelocityComponentEntersTargetZone(jj) {
			validYVelocities = append(validYVelocities, jj)
		}
	}

	for _, ii := range validXVelocities {
		for _, jj := range validYVelocities {

			//fmt.Printf("Launching probe with velocity: %d,%d\n", ii, jj)
			probe.Velocity = &Velocity{ii, jj}
			probe.Location = &Point{0, 0}
			probe.MaxY = math.MinInt
			probe.Path = []Point{{0, 0}}
			if probe.EntersTarget(t) {
				fmt.Printf("Probe launched at %d,%d entered target and reached a height of %d\n", ii, jj, probe.MaxY)
				if probe.MaxY > maxHeight {
					maxHeight = probe.MaxY
				}
				successfulProbes++
			}
			// fmt.Println(probe.Path)
		}
	}

	return maxHeight, successfulProbes
}

func (p *Probe) EntersTarget(t TargetZone) bool {
	for {
		if p.ExceededTarget(t) {
			return false
		}

		if p.WithinTargetZone(t) {
			return true
		}

		p.Step()
	}
}

func (p *Probe) WithinTargetZone(t TargetZone) bool {
	l := p.Location
	if l.X >= t.X.Min && l.X <= t.X.Max && l.Y >= t.Y.Min && l.Y <= t.Y.Max {
		return true
	}
	return false
}

// Get the minimum velocity to enter the target zone
func GetMinimumXVelocity(target int) int {
	rv := 0
	ii := 0
	invert := false

	if target < 0 {
		invert = true
		target *= -1
	}

	for ; rv < target; ii++ {
		rv += ii
	}
	ii--

	if invert {
		ii *= -1
	}

	return ii
}

// Determine the target x closest to 0
func (t TargetZone) SortXByDistanceToZero() (int, int) {
	if t.X.Max <= 0 {
		return t.X.Max, t.X.Min
	} else if t.X.Min >= 0 {
		return t.X.Min, t.X.Max
	}

	minDiffAbs := t.X.Min * -1

	if minDiffAbs > t.X.Max {
		return t.X.Max, t.X.Min
	}
	return t.X.Min, t.X.Max
}

// Determine if we've travelled beyond the target zone either horizontally
// or fallen below it vertically
func (p *Probe) ExceededTarget(t TargetZone) bool {
	if p.Location.Y < t.Y.Min {
		return true
	} else if p.Location.X > t.X.Max && p.Velocity.X >= 0 {
		return true
	} else if p.Location.X < t.X.Min && p.Velocity.X <= 0 {
		return true
	}

	return false
}

// Pre-launch check to ensure that the velocity with steps can actually enter
// the target zone x co-ordinate and won't fall short or end skip beyond
func (t TargetZone) HorizontalVelocityComponentEntersTargetZone(v int) bool {

	rv := 0
	increment := -1

	if v < 0 {
		increment = 1
	}

	for v != 0 {
		rv += v
		v += increment

		if rv >= t.X.Min && rv <= t.X.Max {
			// In target zone
			return true
		} else if rv < t.X.Min && increment == 1 {
			// If the increment is positive we have a negative x-component
			// and so if we're below the min we cannot get back
			return false
		} else if rv > t.X.Max && increment == -1 {
			// If the increment is negative we have a positive x-component
			// and so if we're above the max we cannot get back
			return false
		}
	}

	return false
}

// Pre-launch check to ensure that the velocity with steps can actually enter
// the target zone y co-ordinate and won't fall skip beyond
func (t TargetZone) VerticalVelocityComponentEntersTargetZone(v int) bool {
	rv := 0

	for {
		rv += v
		v--

		if rv >= t.Y.Min && rv <= t.Y.Max {
			// In target zone
			return true
		} else if rv < t.Y.Min {
			// If the increment is positive we have a negative x-component
			// and so if we're below the min we cannot get back
			return false
		}
	}
}

func (p *Probe) Step() {
	p.Location.X += p.Velocity.X
	p.Location.Y += p.Velocity.Y
	p.Path = append(p.Path, Point{p.Location.X, p.Location.Y})
	if p.Location.Y > p.MaxY {
		p.MaxY = p.Location.Y
	}
	p.Velocity.Step()
	//fmt.Printf("New location: %d,%d; New velocity:  %d,%d\n", p.Location.X, p.Velocity.Y, p.Velocity.X, p.Velocity.Y)
}

// Scans a string that looks like "target area: x=1..1, y=1..1" and extracts the numbers
func GetTargetFromInput(inputName *string) TargetZone {
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	var xMin, xMax, yMin, yMax int
	fmt.Sscanf(input, "target area: x=%d..%d, y=%d..%d",
		&xMin,
		&xMax,
		&yMin,
		&yMax)
	fmt.Printf("Target zone is: %d,%d to %d,%d\n", xMin, yMin, xMax, yMax)
	return TargetZone{
		X: MinMaxPair{Min: xMin, Max: xMax},
		Y: MinMaxPair{Min: yMin, Max: yMax},
	}
}

// Extracts the string and trims the trailing spaces
func GetStringFromInput(inputPath string) string {
	rawinput, fileError := os.ReadFile(inputPath)
	Check(fileError)
	return strings.TrimSpace(string(rawinput))
}

// Gets the string referenced by the pointer - if we have a null pointer return
// an empty string
func DereferenceStringPointer(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Lazy check and panic for errors
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
