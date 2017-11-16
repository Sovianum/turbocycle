package geom

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

type Curve interface {
	Point(t float64) *mat.VecDense
}

// works only for 2d curves
func CurvRadius2(curve Curve, t, dt float64) float64 {
	return 1 / Curvature2(curve, t, dt)
}

// works only for 2d curves
func Curvature2(curve Curve, t, dt float64) float64 {
	var der1 = ApproxDerivative1(curve, t, dt)
	var der2 = ApproxDerivative2(curve, t, dt)

	var enom = common.Cross2(der1, der2)
	var denom = math.Pow(mat.Norm(der1, 2), 3)

	return math.Abs(enom) / denom
}

func ApproxDerivative1(curve Curve, t, dt float64) *mat.VecDense {
	var leftPoint = curve.Point(t - dt / 2)
	var rightPoint = curve.Point(t + dt / 2)

	var result = mat.NewVecDense(leftPoint.Len(), nil)
	result.SubVec(rightPoint, leftPoint)
	result.ScaleVec(1 / dt, result)

	return result
}

func ApproxDerivative2(curve Curve, t, dt float64) *mat.VecDense {
	var leftPoint = curve.Point(t - dt / 2)
	var rightPoint = curve.Point(t + dt / 2)
	var centerPoint = curve.Point(t)

	var result = mat.NewVecDense(leftPoint.Len(), nil)
	result.AddVec(rightPoint, leftPoint)
	result.AddScaledVec(result, -2, centerPoint)
	result.ScaleVec(2 / (dt * dt), result)
	return result
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
	var tArr = common.LinSpace(t0, t1, pointNum)
	return Points(curve, tArr)
}

func Points(curve Curve, tArr []float64) []*mat.VecDense {
	var result = make([]*mat.VecDense, len(tArr))
	for i, t := range tArr {
		result[i] = curve.Point(t)
	}
	return result
}
