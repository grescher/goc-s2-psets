package main

import (
	"fmt"
	"testing"
)

func TestLenInt(t *testing.T) {
	fName := "lenInt"
	tests := []struct {
		in   int
		want int
	}{
		{in: 0, want: 1},
		{in: -0, want: 1},
		{in: 1, want: 1},
		{in: 9, want: 1},
		{in: -1, want: 2},
		{in: -9, want: 2},
		{in: 22, want: 2},
		{in: -22, want: 3},
		{in: 333, want: 3},
		{in: -999, want: 4},
		{in: 9999, want: 4},
	}
	for _, tt := range tests {
		if got := lenInt(tt.in); got != tt.want {
			t.Errorf(errorString(fName, tt.in, got, tt.want))
		}
	}
}

func errorString(fName string, in, got, want interface{}) string {
	return fmt.Sprintf(
		"x> %s(%#v) = %#v; want: %#v\n", fName, in, got, want,
	)
}
