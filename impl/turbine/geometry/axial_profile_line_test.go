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
