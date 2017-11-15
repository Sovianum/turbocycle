package profiles

import (
	"math"

	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/geom"
	"gonum.org/v1/gonum/mat"
)

func NewLine(inletPoint, outletPoint *mat.VecDense, inletAngle, outletAngle float64) Line {
	var result = &line{
		startPoint: inletPoint,
		endPoint:   outletPoint,
		interPoint: intersectionPoint(inletPoint, outletPoint, inletAngle, outletAngle),
		combinator: geom.NewCombinator(),
	}
	result.profile = geom.NewBezierCurve([]*mat.VecDense{result.startPoint, result.interPoint, result.endPoint})
	return result
}

type Line interface {
	GetPoints(pointNum int) []*mat.VecDense
	Transform(t geom.Transformation)
}

func ApproxLength(curve Line, segmentNum int) float64 {
	var points = curve.GetPoints(segmentNum+1)
	var delta = mat.NewVecDense(2, nil)
	var result float64

	for i := 1; i != segmentNum+1; i++ {
		delta.SubVec(points[i], points[i-1])
		result += mat.Norm(delta, 2)
	}
	return result
}

type line struct {
	startPoint *mat.VecDense
	endPoint   *mat.VecDense
	interPoint *mat.VecDense
	profile    geom.ParametricCurve
	combinator geom.Combinator
}

func (line *line) Transform(t geom.Transformation) {
	line.combinator.Combine(t)
}

func (line *line) GetPoints(pointNum int) []*mat.VecDense {
	var initialPoints = geom.Points(line.profile, common.Linspace(0, 1, pointNum))
	var points = geom.ApplyAll(line.combinator.Build(), initialPoints)
	return points
}

func intersectionPoint(inletPoint, outletPoint *mat.VecDense, inletAngle, outletAngle float64) *mat.VecDense {
	var xIn, yIn = inletPoint.At(0, 0), inletPoint.At(1, 0)
	var xOut, yOut = outletPoint.At(0, 0), outletPoint.At(1, 0)

	var x = (math.Tan(inletAngle)*xIn - math.Tan(outletAngle)*xOut + yOut - yIn) /
		(math.Tan(inletAngle) - math.Tan(outletAngle))
	var y = yIn + math.Tan(inletAngle)*(x-xIn)

	return mat.NewVecDense(2, []float64{x, y})
}
