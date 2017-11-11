package profiles

import (
	"math"

	"github.com/Sovianum/turbocycle/utils/turbine/radial/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"gonum.org/v1/gonum/mat"
)

type BladeProfile interface {
	PSLine() Line
	SSLine() Line
	MeanLine() Line
	InletEdge() Line
	OutletEdge() Line
	Transform(t geom.Transformation)
}

func NewBladeProfile(
	inletPSPoint, outletPSPoint *mat.VecDense,
	inletSSPoint, outletSSPoint *mat.VecDense,
	inletMeanPoint, outletMeanPoint *mat.VecDense,
	inletPSAngle, outletPSAngle float64,
	inletSSAngle, outletSSAngle float64,
	inletMeanAngle, outletMeanAngle float64,
) BladeProfile {
	return &bladeProfile{
		psLine:     NewLine(inletPSPoint, outletPSPoint, inletPSAngle, outletPSAngle),
		ssLine:     NewLine(inletSSPoint, outletSSPoint, inletSSAngle, outletSSAngle),
		meanLine:   NewLine(inletMeanPoint, outletMeanPoint, inletMeanAngle, outletMeanAngle),
		inletEdge:  NewLine(inletPSPoint, inletSSPoint, -inletPSAngle, -inletSSAngle),
		outletEdge: NewLine(outletPSPoint, outletSSPoint, -outletPSAngle, -outletSSAngle),
	}
}

func NewBladeProfileWithRadiuses(
	inletMeanPoint, outletMeanPoint *mat.VecDense,
	inletPSAngle, outletPSAngle float64,
	inletSSAngle, outletSSAngle float64,
	inletMeanAngle, outletMeanAngle float64,
	unitInletRadius, unitOutletRadius float64,
) BladeProfile {
	var inletPSPoint = radialPoint(inletMeanPoint, inletPSAngle, unitInletRadius)
	var outletPSPoint = radialPoint(outletMeanPoint, outletPSAngle, unitOutletRadius)

	var inletSSPoint = radialPoint(inletMeanPoint, inletSSAngle, unitInletRadius)
	var outletSSPoint = radialPoint(outletMeanPoint, outletSSAngle, unitOutletRadius)

	return &bladeProfile{
		psLine:     NewLine(inletPSPoint, outletPSPoint, inletPSAngle, outletPSAngle),
		ssLine:     NewLine(inletSSPoint, outletSSPoint, inletSSAngle, outletSSAngle),
		meanLine:   NewLine(inletMeanPoint, outletMeanPoint, inletMeanAngle, outletMeanAngle),
		inletEdge:  NewLine(inletPSPoint, inletSSPoint, -inletPSAngle, -inletSSAngle),
		outletEdge: NewLine(outletPSPoint, outletSSPoint, -outletPSAngle, -outletSSAngle),
	}
}

func NewBladeProfileFromProfiler(hRel, unitInletRadius, unitOutletRadius float64, profiler profilers.Profiler) BladeProfile {
	var inletMeanPoint = mat.NewVecDense(2, []float64{0, 0})
	var outletMeanPoint = mat.NewVecDense(2, []float64{1, 0})
	return NewBladeProfileWithRadiuses(
		inletMeanPoint, outletMeanPoint,
		profilers.InletPSAngle(hRel, profiler),
		profilers.OutletPSAngle(hRel, profiler),
		profilers.InletSSAngle(hRel, profiler),
		profilers.OutletSSAngle(hRel, profiler),
		profiler.InletProfileAngle(hRel),
		profiler.OutletProfileAngle(hRel),
		unitInletRadius,
		unitOutletRadius,
	)
}

type bladeProfile struct {
	psLine     Line
	ssLine     Line
	meanLine   Line
	inletEdge  Line
	outletEdge Line
}

func (b *bladeProfile) Transform(t geom.Transformation) {
	b.psLine.Transform(t)
	b.ssLine.Transform(t)
	b.meanLine.Transform(t)
	b.inletEdge.Transform(t)
	b.outletEdge.Transform(t)
}

func (b *bladeProfile) MeanLine() Line {
	return b.meanLine
}

func (b *bladeProfile) SSLine() Line {
	return b.ssLine
}

func (b *bladeProfile) PSLine() Line {
	return b.psLine
}

func (b *bladeProfile) InletEdge() Line {
	return b.inletEdge
}

func (b *bladeProfile) OutletEdge() Line {
	return b.outletEdge
}

func radialPoint(startPoint *mat.VecDense, angle float64, radius float64) *mat.VecDense {
	var x = startPoint.At(0, 0) - radius*math.Sin(angle)
	var y = startPoint.At(1, 0) + radius*math.Cos(angle)
	return mat.NewVecDense(2, []float64{x, y})
}
