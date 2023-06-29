package main

import (
	"fmt"
	"testing"
)

func TestLenString(t *testing.T) {
	fName := "lenString"
	tests := []struct {
		in   string
		want int
	}{
		{in: "", want: 0},
		{in: " ", want: 1},
		{in: "a", want: 1},
		{in: "Aa", want: 2},
		{in: "\t", want: 2},
		{in: "\n", want: 2},
		{in: "Aa@", want: 3},
		{in: "\x00", want: 4},
		{in: " Jane Doe ", want: 10},
		// {in: "\x00\x10\x20\x30\x40\x50\x60\x70", want: 32},
	}
	for _, tt := range tests {
		if got := lenString(tt.in); got != tt.want {
			t.Errorf(errorString(fName, tt.in, got, tt.want))
		}
	}
}

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

func TestLenFloat64(t *testing.T) {
	fName := "lenFloat64"
	tests := []struct {
		in        float64
		precision int
		want      int
	}{
		{in: 0, precision: 5, want: 3},
		{in: 0.0, precision: 5, want: 3},
		{in: 0.75, precision: 5, want: 4},
		{in: .75, precision: 5, want: 4},
		{in: 80, precision: 5, want: 4},
		{in: 80.0, precision: 5, want: 4},
		{in: 3141567.98765456789, precision: 5, want: 13},
		{in: 3141567.98700456789, precision: 5, want: 11},
		{in: 3141567.98765456789, precision: 8, want: 16},
	}
	for _, tt := range tests {
		if got := lenFloat64(tt.in, tt.precision); got != tt.want {
			t.Errorf(
				"%s(%.*f, %d) = %d; want: %d\n",
				fName, tt.precision, tt.in, tt.precision, got, tt.want,
			)
		}
	}
}

func errorString(fName string, in, got, want interface{}) string {
	return fmt.Sprintf(
		"%s(%#v) = %#v; want: %#v\n", fName, in, got, want,
	)
}
