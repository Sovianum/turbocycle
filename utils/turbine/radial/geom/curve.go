package geom

import (
	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

type ParametricCurve interface {
	Point(t float64) *mat.VecDense
}

func Points(curve ParametricCurve, tArr []float64) []*mat.VecDense {
	var result = make([]*mat.VecDense, len(tArr))
	for i, t := range tArr {
		result[i] = curve.Point(t)
	}
	return result
}

func ApproxLength(curve ParametricCurve, segmentNum int) float64 {
	var points = Points(curve, common.Linspace(0, 1, segmentNum+1))
	var delta = newPoint()
	var result float64

	for i := 1; i != segmentNum+1; i++ {
		delta.SubVec(points[i], points[i-1])
		result += mat.Norm(delta, 2)
	}
	return result
}
