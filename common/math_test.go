package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAverage_Const(t *testing.T) {
	var f = func(float64) float64 { return 1 }
	var ave = Average(f, 0, 100, 100)

	assert.Equal(t, ave, float64(1))
}

func TestAverage_Linear(t *testing.T) {
	var f = func(x float64) float64 { return x }
	var ave = Average(f, 0, 100, 100)

	assert.Equal(t, ave, float64(50))
}

func TestAngleTransforms(t *testing.T) {
	var x = 1.
	var xDeg = ToDegrees(x)
	var xRad = ToRadians(xDeg)

	assert.InDelta(t, x, xRad, 1e-8)
}

func TestRoundInt(t *testing.T) {
	var tc = []struct {
		x         float64
		expectedI int
	}{
		{-0.25, 0},
		{-0.75, -1},
		{0, 0},
		{1, 1},
		{-1, -1},
	}

	for i, tc := range tc {
		assert.Equal(t, tc.expectedI, RoundInt(tc.x), "Test case %d (%d)", i, tc.x)
	}
}

func TestFactorial(t *testing.T) {
	var tc = []struct {
		x         int
		expectedI int
	}{
		{0, 1},
		{-1, 1},
		{1, 1},
		{2, 2},
		{5, 120},
	}

	for i, tc := range tc {
		assert.Equal(t, tc.expectedI, Factorial(tc.x), "Test case %d (%f)", i, tc.x)
	}
}

func TestLinspace(t *testing.T) {
	var x1 = 0.
	var x2 = 1.
	var n = 3

	var linspace = LinSpace(x1, x2, n)

	assert.Equal(t, n, len(linspace))
	assert.InDelta(t, x1, linspace[0], 1e-8)
	assert.InDelta(t, x2, linspace[2], 1e-8)
	assert.InDelta(t, 0.5, linspace[1], 1e-8)
}
