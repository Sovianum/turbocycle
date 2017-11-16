package geom

import (
	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

type Curve interface {
	Point(t float64) *mat.VecDense
}

func ApproxLength(curve Curve, t0, t1 float64, segmentNum int) float64 {
	var points = LinPoints(curve, t0, t1, segmentNum+1)
	var delta = mat.NewVecDense(2, nil)
	var result float64

	for i := 1; i != segmentNum+1; i++ {
		delta.SubVec(points[i], points[i-1])
		result += mat.Norm(delta, 2)
	}
	return result
}

func LinPoints(curve Curve, t0, t1 float64, pointNum int) []*mat.VecDense {
	var tArr = common.Linspace(t0, t1, pointNum)
	return Points(curve, tArr)
}

func Points(curve Curve, tArr []float64) []*mat.VecDense {
	var result = make([]*mat.VecDense, len(tArr))
	for i, t := range tArr {
		result[i] = curve.Point(t)
	}
	return result
}
