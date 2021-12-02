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

	stringArray := strings.FieldsFunc(string(rawinput), RemoveEmptyValues)
	sub := Submarine{0, 0, 0, 0}
	for _, v := range stringArray {
		split := strings.Split(v, " ")
		magnitude, err := strconv.Atoi(split[1])
		Check(err)
		switch split[0] {
		case "forward":
			sub.Forward(magnitude)
		case "down":
			sub.Down(magnitude)
		case "up":
			sub.Up(magnitude)
		}
	}
	fmt.Printf("The product of the location scalars for the basic submarine is %d\n", sub.BasicProduct())
	fmt.Printf("The product of the location scalars for the aimed submarine is %d\n", sub.AimedProduct())

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s", timeTaken)
	fmt.Scanf("h")
}

type Submarine struct {
	Depth      int
	X          int
	AimedDepth int
	Aim        int
}

func (s *Submarine) Forward(i int) {
	s.X += i
	s.AimedDepth += i * s.Aim
}
func (s *Submarine) Up(i int) {
	s.Depth -= i
	s.Aim -= i
}
func (s *Submarine) Down(i int) {
	s.Depth += i
	s.Aim += i
}
func (s *Submarine) BasicProduct() int {
	return s.X * s.Depth
}
func (s *Submarine) AimedProduct() int {
	return s.X * s.AimedDepth
}

func RemoveEmptyValues(c rune) bool {
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
