package ode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	x0   = 0.
	xMax = 1.
	step = 1e-6
	y0   = 1.
)

func TestEulerSolution_Constant(t *testing.T) {
	var e = eulerSolution{
		step: step,
		x0:   x0,
		xMax: xMax,
		y0:   y0,
		derivative: func(x, y float64) float64 {
			return 0
		},
	}

	var xArr, yArr = e.Build()
	assert.Equal(t, len(xArr), len(yArr))
	for _, y := range yArr {
		assert.InDelta(t, y0, y, 1e-10)
	}
}

func TestEulerSolution_Linear(t *testing.T) {
	var k = 1.
	var e = eulerSolution{
		step: step,
		x0:   x0,
		xMax: xMax,
		y0:   y0,
		derivative: func(x, y float64) float64 {
			return k
		},
	}

	var xArr, yArr = e.Build()
	assert.Equal(t, len(xArr), len(yArr))
	for i := 0; i != len(xArr); i++ {
		var expected = y0 + k*(xArr[i]-x0)
		assert.InDelta(t, expected, yArr[i], 1e-8)
	}
}

func TestEulerSolution_Quadratic(t *testing.T) {
	var k = 1.

	var e = eulerSolution{
		step: step,
		x0:   x0,
		xMax: xMax,
		y0:   y0,
		derivative: func(x, y float64) float64 {
			return k * (x - x0)
		},
	}

	var xArr, yArr = e.Build()
	assert.Equal(t, len(xArr), len(yArr))
	for i := 0; i != len(xArr); i += 100 {
		var expected = y0 + k/2*(xArr[i]-x0)*(xArr[i]-x0)
		assert.InDelta(t, expected, yArr[i], 1e-6)
	}
}
