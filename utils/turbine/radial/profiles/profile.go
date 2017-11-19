package profiles

import (
	"math"

	"github.com/Sovianum/turbocycle/utils/turbine/geom"
	"github.com/Sovianum/turbocycle/utils/turbine/radial/profilers"
	"gonum.org/v1/gonum/mat"
)

const (
	defaultN = 100
	positive = 1
	negative = -1
)

type BladeProfile interface {
	PSLine() geom.TransformableCurve
	SSLine() geom.TransformableCurve
	MeanLine() geom.TransformableCurve
	InletEdge() geom.TransformableCurve
	OutletEdge() geom.TransformableCurve
	Transform(t geom.Transformation)
}

func NewBladeProfileFromProfiler(
	hRel float64,
	unitInletRadius, unitOutletRadius float64,
	inletCurveFactor, outletCurveFactor float64,
	profiler profilers.Profiler,
) BladeProfile {
	var inletMeanPoint = mat.NewVecDense(2, []float64{0, 0})
	var outletMeanPoint = mat.NewVecDense(2, []float64{1, 0})
	return NewBladeProfileWithRadii(
		inletMeanPoint, outletMeanPoint,
		profilers.InletPSAngle(hRel, profiler),
		profilers.OutletPSAngle(hRel, profiler),
		profilers.InletSSAngle(hRel, profiler),
		profilers.OutletSSAngle(hRel, profiler),
		profiler.InletProfileAngle(hRel),
		profiler.OutletProfileAngle(hRel),
		unitInletRadius, unitOutletRadius,
		inletCurveFactor, outletCurveFactor,
	)
}

func NewBladeProfileWithRadii(
	inletMeanPoint, outletMeanPoint *mat.VecDense,
	inletPSAngle, outletPSAngle float64,
	inletSSAngle, outletSSAngle float64,
	inletMeanAngle, outletMeanAngle float64,
	unitInletRadius, unitOutletRadius float64,
	inletCurveFactor, outletCurveFactor float64,
) BladeProfile {
	var inletPSPoint = radialPoint(inletMeanPoint, inletPSAngle, unitInletRadius, negative)
	var outletPSPoint = radialPoint(outletMeanPoint, outletPSAngle, unitOutletRadius, negative)

	var inletSSPoint = radialPoint(inletMeanPoint, inletSSAngle, unitInletRadius, positive)
	var outletSSPoint = radialPoint(outletMeanPoint, outletSSAngle, unitOutletRadius, positive)

	return NewBladeProfile(
		inletPSPoint, outletPSPoint,
		inletSSPoint, outletSSPoint,
		inletMeanPoint, outletMeanPoint,
		inletPSAngle, outletPSAngle,
		inletSSAngle, outletSSAngle,
		inletMeanAngle, outletMeanAngle,
		inletCurveFactor, outletCurveFactor,
	)
}

func NewBladeProfile(
	inletPSPoint, outletPSPoint *mat.VecDense,
	inletSSPoint, outletSSPoint *mat.VecDense,
	inletMeanPoint, outletMeanPoint *mat.VecDense,
	inletPSAngle, outletPSAngle float64,
	inletSSAngle, outletSSAngle float64,
	inletMeanAngle, outletMeanAngle float64,
	inletCurveFactor, outletCurveFactor float64,
) BladeProfile {
	var psLine = geom.NewTransformableCurve(geom.NewBezier2FromOrientedPoints(
		inletPSPoint, outletPSPoint, inletPSAngle, math.Pi - outletPSAngle,
	))
	var ssLine = geom.NewTransformableCurve(geom.NewBezier2FromOrientedPoints(
		inletSSPoint, outletSSPoint, inletSSAngle, math.Pi - outletSSAngle,
	))
	var meanLine = geom.NewTransformableCurve(geom.NewBezier2FromOrientedPoints(
		inletMeanPoint, outletMeanPoint, inletMeanAngle, math.Pi - outletMeanAngle,
	))
	var inletEdge = geom.NewTransformableCurve(geom.NewBezier3FromOrientedPoints(
		inletPSPoint, inletSSPoint, math.Pi + inletPSAngle, inletSSAngle, inletCurveFactor,
	))
	var outletEdge = geom.NewTransformableCurve(geom.NewBezier3FromOrientedPoints(
		outletPSPoint, outletSSPoint, math.Pi - outletPSAngle, -outletSSAngle, outletCurveFactor,
	))

	return &bladeProfile{
		psLine:     psLine,
		ssLine:     ssLine,
		meanLine:   meanLine,
		inletEdge:  inletEdge,
		outletEdge: outletEdge,
	}
}

func Perimeter(profile BladeProfile) float64 {
	return geom.ApproxLength(profile.InletEdge(), 0, 1, defaultN) +
		geom.ApproxLength(profile.OutletEdge(), 0, 1, defaultN) +
		geom.ApproxLength(profile.PSLine(), 0, 1, defaultN) +
		geom.ApproxLength(profile.SSLine(), 0, 1, defaultN)
}

func PSSegment(profile BladeProfile, inletEdgeFraction float64, outletEdgeFraction float64) geom.Segment {
	var inletEdgeSegment = geom.NewUnitSegment(profile.InletEdge(), inletEdgeFraction, 0)
	var inletEdgeSegmentLength = geom.ApproxLength(inletEdgeSegment, 0, 1, defaultN)

	var psSegment = geom.NewUnitSegment(profile.PSLine(), 0, 1)
	var psSegmentLength = geom.ApproxLength(psSegment, 0, 1, defaultN)

	var outletEdgeSegment = geom.NewUnitSegment(profile.OutletEdge(), 0, outletEdgeFraction)
	var outletEdgeSegmentLength = geom.ApproxLength(outletEdgeSegment, 0, 1, defaultN)

	var totalLength = inletEdgeSegmentLength + psSegmentLength + outletEdgeSegmentLength

	return geom.JoinToUnit(
		[]geom.Segment{inletEdgeSegment, psSegment, outletEdgeSegment},
		[]float64{
			inletEdgeSegmentLength / totalLength,
			(inletEdgeSegmentLength + psSegmentLength) / totalLength,
		},
	)
}

func SSSegment(profile BladeProfile, inletEdgeFraction float64, outletEdgeFraction float64) geom.Segment {
	var inletEdgeSegment = geom.NewUnitSegment(profile.InletEdge(), inletEdgeFraction, 1)
	var inletEdgeSegmentLength = geom.ApproxLength(inletEdgeSegment, 0, 1, defaultN)

	var ssSegment = geom.NewUnitSegment(profile.SSLine(), 0, 1)
	var psSegmentLength = geom.ApproxLength(ssSegment, 0, 1, defaultN)

	var outletEdgeSegment = geom.NewUnitSegment(profile.OutletEdge(), 1, outletEdgeFraction)
	var outletEdgeSegmentLength = geom.ApproxLength(outletEdgeSegment, 0, 1, defaultN)

	var totalLength = inletEdgeSegmentLength + psSegmentLength + outletEdgeSegmentLength

	return geom.JoinToUnit(
		[]geom.Segment{inletEdgeSegment, ssSegment, outletEdgeSegment},
		[]float64{
			inletEdgeSegmentLength / totalLength,
			(inletEdgeSegmentLength + psSegmentLength) / totalLength,
		},
	)
}

func CircularSegment(profile BladeProfile) geom.Segment {
	var segments = []geom.Segment{
		geom.NewUnitSegment(profile.InletEdge(), 0, 1),
		geom.NewUnitSegment(profile.SSLine(), 0, 1),
		geom.NewUnitSegment(profile.OutletEdge(), 1, 0),
		geom.NewUnitSegment(profile.PSLine(), 1, 0),
	}
	var weights = []float64{5, 1, 5, 1}

	var weightedLength = make([]float64, len(segments))
	for i := 0; i != len(segments); i++ {
		weightedLength[i] = geom.ApproxLength(segments[i], 0, 1, defaultN) * weights[i]
	}

	var totalLength float64 = 0
	for _, num := range weightedLength {
		totalLength += num
	}

	var partialSum float64 = 0
	var partialSums = make([]float64, len(weightedLength) - 1)
	for i, length := range weightedLength[:len(weightedLength) - 1] {
		partialSum += length
		partialSums[i] = partialSum / totalLength
	}

	return geom.JoinToUnit(segments, partialSums)
}

type bladeProfile struct {
	psLine     geom.TransformableCurve
	ssLine     geom.TransformableCurve
	meanLine   geom.TransformableCurve
	inletEdge  geom.TransformableCurve
	outletEdge geom.TransformableCurve
}

func (b *bladeProfile) Transform(t geom.Transformation) {
	b.psLine.Transform(t)
	b.ssLine.Transform(t)
	b.meanLine.Transform(t)
	b.inletEdge.Transform(t)
	b.outletEdge.Transform(t)
}

func (b *bladeProfile) MeanLine() geom.TransformableCurve {
	return b.meanLine
}

func (b *bladeProfile) SSLine() geom.TransformableCurve {
	return b.ssLine
}

func (b *bladeProfile) PSLine() geom.TransformableCurve {
	return b.psLine
}

func (b *bladeProfile) InletEdge() geom.TransformableCurve {
	return b.inletEdge
}

func (b *bladeProfile) OutletEdge() geom.TransformableCurve {
	return b.outletEdge
}

func radialPoint(startPoint *mat.VecDense, angle float64, radius float64, direction float64) *mat.VecDense {
	var x = startPoint.At(0, 0) - direction * radius*math.Sin(angle)
	var y = startPoint.At(1, 0) + direction * radius*math.Cos(angle)
	return mat.NewVecDense(2, []float64{x, y})
}
