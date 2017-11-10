package geom

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

func NewBezierCurve(points []*mat.VecDense) ParametricCurve {
	return &bezierCurve{
		controlPoints:points,
	}
}

type bezierCurve struct {
	controlPoints []*mat.VecDense
}

func (curve *bezierCurve) Points(tArr []float64) []*mat.VecDense {
	var result = make([]*mat.VecDense, 0)
	for _, t := range tArr {
		result = append(result, curve.point(t))
	}
	return result
}

func (curve *bezierCurve) Point(t float64) *mat.VecDense {
	return curve.point(t)
}

func (curve *bezierCurve) point(t float64) *mat.VecDense {
	var sum = newPoint()
	for i := range curve.controlPoints {
		sum.AddVec(sum, curve.term(t, i))
	}
	return sum
}

func (curve *bezierCurve) term(t float64, i int) *mat.VecDense {
	var result = newPoint()
	result.ScaleVec(curve.termBasis(t, i), curve.controlPoints[i])
	return result
}

func (curve *bezierCurve) termBasis(t float64, i int) float64 {
	return curve.termFactor(i) * math.Pow(t, float64(i)) * math.Pow(1 - t, float64(curve.order() - i))
}

func (curve *bezierCurve) termFactor(i int) float64 {
	var numer = float64(common.Factorial(curve.order()))
	var denom = float64(common.Factorial(i) * common.Factorial(curve.order() - i))
	return numer / denom
}

func (curve *bezierCurve) order() int {
	return len(curve.controlPoints) - 1
}

func newPoint() *mat.VecDense {
	return mat.NewVecDense(2, nil)
}

