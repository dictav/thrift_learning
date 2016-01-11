package main

import "testing"
import "github.com/dictav/thrift_learning/go/gen-go/tutorial"
import "math"

var tests = []struct {
	in  *tutorial.Work
	out int32
}{
	{&tutorial.Work{1, 2, tutorial.Operation_ADD, nil}, 3},
	{&tutorial.Work{2147483646, 1, tutorial.Operation_ADD, nil}, math.MaxInt32},
	{&tutorial.Work{1, 2, tutorial.Operation_SUBTRACT, nil}, -1},
	{&tutorial.Work{-2147483647, 1, tutorial.Operation_SUBTRACT, nil}, math.MinInt32},
	{&tutorial.Work{1, 2, tutorial.Operation_MULTIPLY, nil}, 2},
	{&tutorial.Work{-1, 2, tutorial.Operation_MULTIPLY, nil}, -2},
	{&tutorial.Work{1, 2, tutorial.Operation_DIVIDE, nil}, 0},
	{&tutorial.Work{2, 2, tutorial.Operation_DIVIDE, nil}, 1},
	{&tutorial.Work{-2, 1, tutorial.Operation_DIVIDE, nil}, -2},
}

func TestCalculate(t *testing.T) {
	c := NewCalculatorHandler()

	for i, tt := range tests {
		s, err := c.Calculate(int32(i), tt.in)
		if err != nil {
			t.Error(err)
		}

		if tt.out != s {
			t.Errorf("Calculate(%d, %q) => %d, want %d", i, tt.in, s, tt.out)
		}
	}
}
