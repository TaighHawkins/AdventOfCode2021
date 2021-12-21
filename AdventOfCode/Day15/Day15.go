package main

import (
	"container/heap"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	Grid        map[Point]int
	InitialMaxX int
	InitialMaxY int
	MaxX        int
	MaxY        int
)

func init() {
}

func main() {
	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	PopulateMap(input)
	cost := Navigate(Point{0, 0}, Point{InitialMaxX, InitialMaxY})
	fmt.Printf("The minimum distance to get to the end is %d\n", cost)

	largerCost := Navigate(Point{0, 0}, Point{MaxX, MaxY})
	fmt.Printf("The minimum distance to get to the end of the larger map is %d\n", largerCost)

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func Navigate(start Point, end Point) int {

	pq := make(PriorityQueue, 0)
	score := make(map[Point]int)
	heap.Init(&pq)

	pq.Push(Node{Position: start, Risk: 0})
	coordinateModifiers := [2][4]int{{0, 0, 1, -1}, {1, -1, 0, 0}}

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(Node)

		for ii := 0; ii < 4; ii++ {
			neighbour := Point{
				current.Position.X + coordinateModifiers[0][ii],
				current.Position.Y + coordinateModifiers[1][ii],
			}

			if neighbour.OutOfBounds(end) {
				continue
			}
			risk := current.Risk + neighbour.GetRisk()
			curScore, present := score[neighbour]
			if present && curScore <= risk {
				continue
			} else {
				score[neighbour] = risk
			}

			pq.Push(Node{Position: neighbour, Risk: risk})
		}
	}

	return score[end]
}

func (p Point) GetRisk() int {
	x := p.X % (InitialMaxX + 1)
	y := p.Y % (InitialMaxY + 1)
	increment := (p.X / (InitialMaxX + 1)) + (p.Y / (InitialMaxY + 1))

	risk := Grid[Point{x, y}] + increment
	if risk > 9 {
		risk = risk - 9
	}
	return risk
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (p Point) OutOfBounds(end Point) bool {
	return p.X > end.X || p.X < 0 || p.Y > end.Y || p.Y < 0
}

func PopulateMap(input []string) {
	Grid = make(map[Point]int)
	InitialMaxY = len(input) - 1
	MaxY = 5*len(input) - 1
	for ii, rowRaw := range input {
		length := len(rowRaw)
		if length-1 > InitialMaxX {
			InitialMaxX = length - 1
			MaxX = length*5 - 1
		}
		for jj, r := range rowRaw {
			Grid[Point{jj, ii}] = int(r - '0')
		}
	}
}

// https://pkg.go.dev/container/heap#example-package-PriorityQueue
// A PriorityQueue implements heap.Interface and holds Items.
// We don't need a pointer here because we don't care to modify
// the nodes themselves
type PriorityQueue []Node

type Node struct {
	Position Point
	Risk     int
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int
}

type Point struct {
	X int
	Y int
}

func (q PriorityQueue) Len() int {
	return len(q)
}

func (q PriorityQueue) Less(i, j int) bool {
	// Kind of inverse of the normal 'priority queue'
	// We want to minimise risk to we want the lowest value
	return q[i].Risk < q[j].Risk
}

func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *PriorityQueue) Push(x interface{}) {
	n := q.Len()
	item := x.(Node)
	item.index = n
	*q = append(*q, item)
}

func (q *PriorityQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1
	*q = old[0 : n-1]
	return item
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
