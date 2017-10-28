package geometry

import (
	"math"
	"testing"

	"github.com/Sovianum/turbocycle/common"
	"github.com/stretchr/testify/assert"
)

func TestAxialProfileLine(t *testing.T) {
	var x0 float64 = 10
	var d0 float64 = 10
	var angle = math.Pi / 4

	var line = NewAxialProfileLine(x0, d0, angle)
	var gotD = line.Diameter(11)
	assert.True(t, common.ApproxEqual(12, gotD, 0.001), "got %f", gotD)
}

func TestMeanLine(t *testing.T) {
	var x0 float64 = 0
	var x1 float64 = 10

	var d0In float64 = 10
	var angleIn = -math.Pi / 4
	var innerLine = NewAxialProfileLine(x0, d0In, angleIn)

	var d0Out float64 = 20
	var angleOut = math.Pi / 4
	var outerLine = NewAxialProfileLine(x0, d0Out, angleOut)

	var meanLine = MeanLine(innerLine, outerLine, 0.5)
	var dMeanIn = (innerLine.Diameter(x0) + outerLine.Diameter(x0)) / 2
	var dMeanOut = (innerLine.Diameter(x1) + outerLine.Diameter(x1)) / 2

	var delta = 0.0000001
	assert.InDelta(t, dMeanIn, meanLine.Diameter(x0), delta)
	assert.InDelta(t, dMeanOut, meanLine.Diameter(x1), delta)
}
