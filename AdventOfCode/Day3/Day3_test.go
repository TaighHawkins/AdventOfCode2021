package main

import (
	"testing"
)

func Test_GetAllOneCounts(t *testing.T) {
	countOfOnes, halfArraySize := GetAllOneCounts(GetTestData())
	if halfArraySize != 6 {
		t.Fatalf("The array was not correctly halved, expected %d but got %d\n", 6, halfArraySize)
	}

	expectedOneCount := []int{7, 5, 8, 7, 5}
	if !SlicesEqual(countOfOnes, expectedOneCount) {
		t.Fatalf("The count of ones is incorrect, expected %d but got %d\n", countOfOnes, expectedOneCount)
	}
}

func Test_GetGammaAndEpsilon(t *testing.T) {
	countOfOnes, halfArraySize := GetAllOneCounts(GetTestData())
	gammaRateBits, epsilonBits := ExtractBits(countOfOnes, halfArraySize)
	gammaRate := ConvertStringBitsToInt(gammaRateBits)
	epsilon := ConvertStringBitsToInt(epsilonBits)
	expectedGammaRate := int64(22)
	if gammaRate != expectedGammaRate {
		t.Fatalf("The gamma rate is incorrect, expected %d but got %d\n", expectedGammaRate, gammaRate)
	}

	expectedEpsilon := int64(9)
	if epsilon != expectedEpsilon {
		t.Fatalf("The epsilon value is incorrect, expected %d but got %d\n", expectedEpsilon, epsilon)
	}
}

func SlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func GetTestData() []string {
	return []string{
		"00100",
		"11110",
		"10110",
		"10111",
		"10101",
		"01111",
		"00111",
		"11100",
		"10000",
		"11001",
		"00010",
		"01010",
	}
}
