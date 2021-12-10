package main

import (
	"errors"
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
	stringArray := strings.Split(string(rawinput), "\n\n")
	numbers := strings.Split(stringArray[0], ",")
	set := BingoBoardSet{}
	set.PopulateBoards(stringArray)
	set.CallNumbers(numbers)

	set.FirstWinningBoard.GetAndPrintBoardValue("First")
	set.FinalWinningBoard.GetAndPrintBoardValue("Last")

	timeTaken := time.Since(start)
	fmt.Printf("Process took %s", timeTaken)
	fmt.Scanf("h")
}

type BingoValue struct {
	Value  int
	Marked bool
}

type BingoBoard struct {
	Values        [][]*BingoValue
	HasWon        bool
	WinningNumber int
}

type BingoBoardSet struct {
	Boards               []*BingoBoard
	HasWinner            bool
	FirstWinningBoard    *BingoBoard
	FinalWinningBoard    *BingoBoard
	CountOfWinningBoards int
}

func (b *BingoBoardSet) CallNumbers(input []string) {
	foundFirst := false
	for _, v := range input {
		number, e := strconv.Atoi(v)
		Check(e)
		// fmt.Printf("Calling number %d\n", number)
		b.MarkBoardsWithValue(number)

		if b.FirstWinningBoard != nil && !foundFirst {
			fmt.Printf("\nWe have our first winner!\n\n")
			foundFirst = true
		}

		if b.FinalWinningBoard != nil {
			fmt.Printf("\nAll boards have won!\n\n")
			break
		}
	}
}

func (b *BingoBoardSet) PopulateBoards(input []string) {
	b.Boards = make([]*BingoBoard, len(input)-1)
	for ii := 1; ii < len(input); ii++ {
		board := BingoBoard{}
		boardRows := strings.Split(strings.TrimSpace(input[ii]), "\n")
		board.Populateboard(boardRows)
		b.Boards[ii-1] = &board
	}
}

func (b *BingoBoardSet) AllBoardsHaveWon() bool {
	return b.CountOfWinningBoards == len(b.Boards)
}

func (b *BingoBoardSet) MarkBoardsWithValue(number int) {
	for _, board := range b.Boards {
		if board.HasWon {
			continue
		}

		board.MarkValueIfPresent(number)
		if board.IsWinning() {
			if b.FirstWinningBoard == nil {
				b.FirstWinningBoard = board
			}
			if len(b.Boards)-b.CountOfWinningBoards == 1 {
				b.FinalWinningBoard = board
			}
			board.WinningNumber = number
			b.CountOfWinningBoards++
			b.HasWinner = true
		}
	}
}

func (b *BingoValue) IsMarked() bool {
	return b.Marked
}

func (b *BingoValue) Mark() {
	b.Marked = true
}

func (b *BingoBoard) GetAndPrintBoardValue(label string) {
	value, e := b.GetBoardValue()
	Check(e)
	b.PrintBoard(label)
	fmt.Printf("%s Winning board value is %d\n", label, value*b.WinningNumber)
}

func (b *BingoBoard) PrintBoard(label string) {
	fmt.Printf("The %s winning board is:\n", label)
	for _, row := range b.Values {
		for _, v := range row {
			fmt.Printf(" %02d ", v.Value)
		}
		fmt.Println("")
		for _, v := range row {
			if v.Marked {
				fmt.Print("  x ")
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func (b *BingoBoard) Populateboard(input []string) {
	values := make([][]*BingoValue, len(input))
	for ii, row := range input {
		splitRow := strings.FieldsFunc(row, SplitValues)
		rowValues := make([]*BingoValue, len(splitRow))
		for jj, v := range splitRow {
			number, e := strconv.Atoi(v)
			Check(e)
			rowValues[jj] = &BingoValue{number, false}
		}
		values[ii] = rowValues
	}
	b.Values = values
}

func (b *BingoBoard) GetValueFromBoard(number int) (*BingoValue, error) {
	for _, row := range b.Values {
		for _, v := range row {
			if v.Value == number {
				return v, nil
			}
		}
	}

	return nil, errors.New("value not on board")
}

func (b *BingoBoard) MarkValueIfPresent(number int) {
	v, e := b.GetValueFromBoard(number)
	if e == nil {
		v.Marked = true
	}
}

func (b *BingoBoard) IsWinning() bool {
	won := AllRowEntriesMarked(b.Values) || AllColumnEntriesMarked(b.Values)
	b.HasWon = won
	return b.HasWon
}

func (b *BingoBoard) GetBoardValue() (int, error) {
	if !b.IsWinning() {
		return 0, errors.New("the board has not won")
	}

	rv := 0
	for _, row := range b.Values {
		for _, v := range row {
			if !v.Marked {
				rv += v.Value
			}
		}
	}
	return rv, nil
}

func AllRowEntriesMarked(b [][]*BingoValue) bool {
	for _, row := range b {
		rv := true
		for _, v := range row {
			if !v.Marked {
				rv = false
			}
		}
		if rv {
			return true
		}
	}
	return false
}

func AllColumnEntriesMarked(b [][]*BingoValue) bool {
	for ix := 0; ix < len(b[0]); ix++ {
		rv := true
		for _, v := range b {
			if !v[ix].Marked {
				rv = false
			}
		}
		if rv {
			return true
		}
	}
	return false
}

func SplitValues(c rune) bool {
	return c == '\n' || c == ' '
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
