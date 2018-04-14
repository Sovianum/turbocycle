package states

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressorVelocityTriangle(t *testing.T) {
	tc := []struct {
		u             float64
		ca            float64
		cu            float64
		expectedAlpha float64
		expectedBeta  float64
	}{
		{1, 1, 1, math.Pi / 4, math.Pi / 2},
		{1, 0.5, 0.5, math.Pi / 4, math.Pi / 4},
	}

	for i, c := range tc {
		triangle := NewCompressorVelocityTriangleFromProjections(c.cu, c.ca, c.u)
		assert.InDelta(t, c.expectedAlpha, triangle.Alpha(), 1e-9, "%d", i)
		assert.InDelta(t, c.expectedBeta, triangle.Beta(), 1e-9, "%d", i)
	}
}
