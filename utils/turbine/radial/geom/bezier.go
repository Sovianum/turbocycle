package geom

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

func NewBezier(points []*mat.VecDense) Curve {
	return &bezierCurve{
		controlPoints: points,
	}
}

func NewBezier2FromOrientedPoints(inletPoint, outletPoint *mat.VecDense, inletAngle, outletAngle float64) Curve {
	return NewBezier([]*mat.VecDense{
		inletPoint,
		intersectionPoint(inletPoint, outletPoint, inletAngle, outletAngle),
		outletPoint,
	})
}

type bezierCurve struct {
	controlPoints []*mat.VecDense
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
	return curve.termFactor(i) * math.Pow(t, float64(i)) * math.Pow(1-t, float64(curve.order()-i))
}

func (curve *bezierCurve) termFactor(i int) float64 {
	var numer = float64(common.Factorial(curve.order()))
	var denom = float64(common.Factorial(i) * common.Factorial(curve.order()-i))
	return numer / denom
}

func (curve *bezierCurve) order() int {
	return len(curve.controlPoints) - 1
}

func newPoint() *mat.VecDense {
	return mat.NewVecDense(2, nil)
}

func intersectionPoint(inletPoint, outletPoint *mat.VecDense, inletAngle, outletAngle float64) *mat.VecDense {
	var xIn, yIn = inletPoint.At(0, 0), inletPoint.At(1, 0)
	var xOut, yOut = outletPoint.At(0, 0), outletPoint.At(1, 0)

	var x = (math.Tan(inletAngle)*xIn - math.Tan(outletAngle)*xOut + yOut - yIn) /
		(math.Tan(inletAngle) - math.Tan(outletAngle))
	var y = yIn + math.Tan(inletAngle)*(x-xIn)

	return mat.NewVecDense(2, []float64{x, y})
}
