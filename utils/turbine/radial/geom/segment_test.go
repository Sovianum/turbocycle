package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

func TestNewSegment(t *testing.T) {
	var curve = NewBezier([]*mat.VecDense{
		mat.NewVecDense(2, []float64{0, 0}),
		mat.NewVecDense(2, []float64{1, 1}),
	})
	var segment = NewSegment(curve, 0.5, 1, 0, 1)
	var point1 = segment.Point(0)
	assert.InDelta(t, 0.5, point1.At(0, 0), 1e-8)
	assert.InDelta(t, 0.5, point1.At(1, 0), 1e-8)

	var point2 = segment.Point(1)
	assert.InDelta(t, 1, point2.At(0, 0), 1e-8)
	assert.InDelta(t, 1, point2.At(1, 0), 1e-8)
}

func TestJoinToUnit(t *testing.T) {
	var curve1 = NewBezier([]*mat.VecDense{
		mat.NewVecDense(2, []float64{0, 0}),
		mat.NewVecDense(2, []float64{1, 1}),
	})
	var segment1 = NewUnitSegment(curve1, 0, 1)

	var curve2 = NewBezier([]*mat.VecDense{
		mat.NewVecDense(2, []float64{1, 1}),
		mat.NewVecDense(2, []float64{2, 0}),
	})
	var segment2 = NewUnitSegment(curve2, 0, 1)

	var join = JoinToUnit([]Segment{segment1, segment2}, []float64{0.5})

	var point1 = join.Point(0.25)
	assert.InDelta(t, 0.5, point1.At(0, 0), 1e-8)
	assert.InDelta(t, 0.5, point1.At(1, 0), 1e-8)

	var point2 = join.Point(0.5)
	assert.InDelta(t, 1, point2.At(0, 0), 1e-8)
	assert.InDelta(t, 1, point2.At(1, 0), 1e-8)

	var point3 = join.Point(0.75)
	assert.InDelta(t, 1.5, point3.At(0, 0), 1e-8)
	assert.InDelta(t, 0.5, point3.At(1, 0), 1e-8)
}
