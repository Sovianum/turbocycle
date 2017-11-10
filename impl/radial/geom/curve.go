package geom

import "gonum.org/v1/gonum/mat"

type ParametricCurve interface {
	Point(t float64) *mat.VecDense
	Points(tArr []float64) []*mat.VecDense
}
