package main

import (
	"testing"
)

func Test_YVelocityCheck(t *testing.T) {
	target := TargetZone{MinMaxPair{20, 30}, MinMaxPair{-10, -5}}
	test := 9
	if !target.VelocityCanEnterTargetY(test) {
		t.Errorf("Test y-velocity %d should have entered target zone", test)
	}

	t.Logf("Test y-velocity %d entered target zone", test)
}
