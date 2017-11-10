package geom

import "gonum.org/v1/gonum/mat"

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
