package geom

import (
	"fmt"

	"github.com/Sovianum/turbocycle/common"
	"gonum.org/v1/gonum/mat"
)

func NewUnitSegment(curve Curve, tStartCurve, tEndCurve float64) Segment {
	return NewSegment(curve, tStartCurve, tEndCurve, 0, 1)
}

func NewSegment(curve Curve, tStartCurve, tEndCurve, tStartSegment, tEndSegment float64) Segment {
	return &segment{
		curve:         curve,
		tStartCurve:   tStartCurve,
		tEndCurve:     tEndCurve,
		tStartSegment: tStartSegment,
		tEndSegment:   tEndSegment,
	}
}

type Segment interface {
	Curve
	TSTart() float64
	TEnd() float64
}

// inner points are assumed to be asc sorted and be inside (0; 1)
func JoinToUnit(segments []Segment, innerPoints []float64) Segment {
	if len(segments)-1 != len(innerPoints) {
		panic("len(innerPoints must be equal (len(segments) - 1)")
	}

	var boundaryPoints = []float64{0}
	for _, t := range innerPoints {
		boundaryPoints = append(boundaryPoints, t)
	}
	boundaryPoints = append(boundaryPoints, 1)

	return &segmentJoin{
		segments:       segments,
		boundaryPoints: boundaryPoints,
	}
}

type segmentJoin struct {
	segments       []Segment
	boundaryPoints []float64
}

func (join *segmentJoin) Point(t float64) *mat.VecDense {
	var segment, tJoinStart, tJoinEnd = join.getSegment(t)
	var tSegment = common.LinScale(t, tJoinStart, tJoinEnd, segment.TSTart(), segment.TEnd())
	return segment.Point(tSegment)
}

func (join *segmentJoin) TSTart() float64 {
	return join.boundaryPoints[0]
}

func (join *segmentJoin) TEnd() float64 {
	return join.boundaryPoints[len(join.boundaryPoints)-1]
}

func (join *segmentJoin) getSegment(t float64) (Segment, float64, float64) {
	for i, tInner := range join.boundaryPoints[1:] {
		if t <= tInner {
			return join.segments[i], join.boundaryPoints[i], join.boundaryPoints[i+1]
		}
	}
	panic(fmt.Errorf(
		"out of range: t = %f, t_min = %f, t_max = %f",
		t,
		join.boundaryPoints[0],
		join.boundaryPoints[len(join.boundaryPoints)-1],
	))
}

type segment struct {
	curve         Curve
	tStartSegment float64
	tEndSegment   float64
	tStartCurve   float64
	tEndCurve     float64
}

func (s *segment) TSTart() float64 {
	return s.tStartSegment
}

func (s *segment) TEnd() float64 {
	return s.tEndSegment
}

func (s *segment) Point(t float64) *mat.VecDense {
	var tCurve = common.LinScale(t, s.tStartSegment, s.tEndSegment, s.tStartCurve, s.tEndCurve)
	return s.curve.Point(tCurve)
}
