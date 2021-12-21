package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	VersionSum int64
)

func init() {
}

func main() {
	inputName := flag.String("input", "input.txt", "The input file to work with")
	flag.Parse()
	start := time.Now()
	input := GetStringFromInput(DereferenceStringPointer(inputName))
	bitString := ConvertHexToBitString(input)
	_, v := DecodePacketStartingAt(bitString)

	fmt.Printf("Version sum of packet input is: %d\n", VersionSum)
	fmt.Printf("The value of the packet input is: %d\n", v)
	timeTaken := time.Since(start)
	fmt.Printf("Process took %s\n", timeTaken)
	fmt.Scanf("h")
}

func DecodePacketStartingAt(s string) (int, int64) {
	value := int64(0)
	soFar := 0
	length := 0
	version := ConvertSubStringToInt(s, 0, 3)
	VersionSum += version
	soFar += 3
	typeId := ConvertSubStringToInt(s, soFar, 3)
	// fmt.Printf("Extracted version and type : %d - %d\n", version, typeId)
	soFar += 3
	if typeId == 4 {
		value, length = ExtractNumberFromBits(s, soFar, 5)
		soFar += length
		// fmt.Printf("The number from the literal value is : %d\n", value)
	} else {
		lengthType := SubString(s, soFar, 1)
		// fmt.Printf("Operator Packet identified with lengthType %s\n", lengthType)
		var values []int64
		soFar++
		if lengthType == "0" {
			maxLength := ConvertSubStringToInt(s, soFar, 15)
			soFar += 15
			length, values = DecodePacketsUpToMaxLength(SubString(s, soFar, -1), int(maxLength))
		} else if lengthType == "1" {
			subPackets := ConvertSubStringToInt(s, soFar, 11)
			soFar += 11
			// fmt.Printf("Max number of subPackets %d\n", subPackets)
			length, values = DecodeNSubPackets(SubString(s, soFar, -1), int(subPackets))
		}
		soFar += length

		switch typeId {
		case 0:
			value += Sum(values)
		case 1:
			value += Product(values)
		case 2:
			value += Min(values)
		case 3:
			value += Max(values)
		case 5:
			value += GreaterThan(values)
		case 6:
			value += LessThan(values)
		case 7:
			value += EqualTo(values)
		default:
			panic("unavailable type detected!")
		}
	}
	// fmt.Printf("Extracted packet %s of %d length from %s\n", SubString(s, 0, soFar), soFar, s)
	return soFar, value
}

func Sum(values []int64) int64 {
	rv := int64(0)
	for _, v := range values {
		rv += v
	}
	return rv
}

func Product(values []int64) int64 {
	rv := int64(1)
	for _, v := range values {
		rv *= v
	}
	return rv
}

func Min(values []int64) int64 {
	rv := int64(math.MaxInt64)
	for _, v := range values {
		if v < rv {
			rv = v
		}
	}
	return rv
}

func Max(values []int64) int64 {
	rv := int64(0)
	for _, v := range values {
		if v > rv {
			rv = v
		}
	}
	return rv
}

func GreaterThan(values []int64) int64 {
	if len(values) != 2 {
		panic("wrong number of values for operator")
	}
	if values[0] > values[1] {
		return 1
	}
	return 0
}

func LessThan(values []int64) int64 {
	if len(values) != 2 {
		panic("wrong number of values for operator")
	}
	if values[0] < values[1] {
		return 1
	}
	return 0
}

func EqualTo(values []int64) int64 {
	if len(values) != 2 {
		panic("wrong number of values for operator")
	}
	if values[0] == values[1] {
		return 1
	}
	return 0
}
func DecodeNSubPackets(s string, n int) (int, []int64) {
	// fmt.Printf("Extracting %d sub-packets from string %s\n", n, s)
	lengthSoFar := 0
	values := make([]int64, n)
	for ii := 0; ii < n; ii++ {
		length, value := DecodePacketStartingAt(s)
		values[ii] = value
		lengthSoFar += length
		s = SubString(s, length, -1)
	}
	return lengthSoFar, values
}

func DecodePacketsUpToMaxLength(s string, maxLength int) (int, []int64) {
	// fmt.Printf("Extracting packets up to a maximum length of %d from %s\n", maxLength, s)
	lengthSoFar := 0
	values := make([]int64, 0)
	for lengthSoFar < maxLength {
		length, value := DecodePacketStartingAt(s)
		values = append(values, value)
		lengthSoFar += length
		s = SubString(s, length, -1)
	}
	return lengthSoFar, values
}

func ExtractNumberFromBits(s string, start int, length int) (int64, int) {
	final := false
	initialStart := start
	numberBits := ""
	for !final {
		ss := SubString(s, start, length)
		if SubString(ss, 0, 1) == "0" {
			final = true
		}
		numberBits += SubString(ss, 1, -1)
		start += length
	}
	return ConvertStringBitsToInt(numberBits), start - initialStart
}

func ConvertSubStringToInt(s string, start int, length int) int64 {
	ss := SubString(s, start, length)
	if ss == "" {
		panic("substring should not return an empty string")
	}
	d := ConvertStringBitsToInt(ss)
	// fmt.Printf("Substring being converted is: %s -> %d\n", ss, d)
	return d
}

func ConvertStringBitsToInt(input string) int64 {
	rv, err := strconv.ParseInt(input, 2, 64)
	Check(err)
	return rv
}

func SubString(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if length == -1 {
		return string(asRunes[start:])
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func ConvertHexToBitString(s string) string {
	return ConvertToBitString(ReadHexNumberToBytes(s))
}

func ReadHexNumberToBytes(s string) []byte {
	i, err := hex.DecodeString(s)
	Check(err)
	return i
}

func ConvertToBitString(bytes []byte) string {
	s := ""
	for _, v := range bytes {
		s += fmt.Sprintf("%08b", v)
	}
	return s
}

func GetBitSlice(val uint64, start int, end int) []uint16 {
	bits := AsBits(val)
	return bits[start:end]
}

func AsBits(val uint64) []uint16 {
	bits := []uint16{}
	for i := 0; i < 24; i++ {
		bits = append([]uint16{uint16(val & 0x1)}, bits...)
		val = val >> 1
	}
	return bits
}

func GetStringFromInput(inputPath string) string {

	rawinput, fileError := os.ReadFile(inputPath)
	Check(fileError)
	return strings.TrimSpace(string(rawinput))
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
